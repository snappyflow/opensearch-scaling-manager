package provision

import (
	"context"
	"encoding/json"
	"regexp"
	"scaling_manager/cluster"
	"scaling_manager/cluster_sim"
	"scaling_manager/config"
	osutils "scaling_manager/opensearchUtils"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

// Input:
//
//	state (*State): A pointer to the state struct which is state maintained in OS document
//	recommendationQueue ([]map[string]string): Recommendations provided by the recommendation engine in the form of an array of strings
//	clusterCfg (config.ClusterDetails): Cluster Level config details
//	usrCfg (config.UserConfig): User defined config for applicatio behavior
//
// Description:
//
//	GetRecommendation will fetch the recommendation from recommendation queue.
//	It will call the Provisioner with all the user defined configs.
//	Triggers the provisioning
//
// Return:
func GetRecommendation(state *State, recommendationQueue []map[string]string, clusterCfg config.ClusterDetails, usrCfg config.UserConfig, t *time.Time) {
	var clusterCurrent cluster.ClusterDynamic
	scaleRegexString := `(scale_up|scale_down)_by_([0-9]+)`
	scaleRegex := regexp.MustCompile(scaleRegexString)
	if len(recommendationQueue) > 0 {
		if usrCfg.MonitorWithSimulator {
			clusterCurrent = cluster_sim.GetClusterCurrent(usrCfg.IsAccelerated)
		} else {
			clusterCurrent = cluster.GetClusterCurrent()
		}

		state.GetCurrentState()
		if state.CurrentState == "normal" {
			var subMatch []string
			var task string
			for task, _ = range recommendationQueue[0] {
				subMatch = scaleRegex.FindStringSubmatch(task)
			}

			numNodes, _ := strconv.Atoi(subMatch[2])
			operation := subMatch[1]
			// Call scale down provisioning only when the cluster status is green. No recommended to scale down when cluster is in yellow or red state
			if operation == "scale_down" && clusterCurrent.ClusterStatus != "green" {
				log.Warn.Println("Recommendation can not be provisioned as open search cluster is unhealthy for a scale_down. \n Discarding this recommendation")
				return
			}

			ruleResponsible := recommendationQueue[0][task]
			// Split the rules if more than one rule is responsible for recommendation
			splitRules := strings.Split(ruleResponsible, "_and_")
			var largestDecisionPeriod int

			// Find the largest decision period among the rules responsiblle
			for _, rule := range splitRules {
				if decisionPeriod, err := strconv.Atoi(rule[strings.LastIndex(rule, "-")+1:]); err == nil && decisionPeriod > largestDecisionPeriod {
					largestDecisionPeriod = decisionPeriod
				} else if err != nil {
					log.Error.Println("Invalid decision period:", err)
					return
				}
			}

			// Get the latest document of successful provision happened
			resp, err := osutils.SearchQuery(context.Background(), []byte(getLatestProvisionQuery()))
			if err != nil {
				log.Error.Println("Error querying the last provision document frm Opensearch", err)
			}
			defer resp.Body.Close()

			var respInterface map[string]interface{}

			decodeErr := json.NewDecoder(resp.Body).Decode(&respInterface)
			if decodeErr != nil {
				log.Error.Println("decode Error: ", decodeErr)
				return
			}

			respHits := respInterface["hits"].(map[string]interface{})["hits"].([]interface{})

			var lastProvisionTime time.Time

			// Get the last successful provision time
			for _, doc := range respHits {
				provisionEndTime := doc.(map[string]interface{})["_source"].(map[string]interface{})["ProvisionEndTime"].(float64)
				lastProvisionTime = time.UnixMilli(int64(provisionEndTime))
			}

			duration, dErr := time.ParseDuration(strconv.Itoa(largestDecisionPeriod) + "m")
			if dErr != nil {
				log.Error.Println("Error converting string to time.Duration", dErr)
			}

			// If the last provision has occured in the range of the largest decision period and now. Discard the current recommendation
			diff := time.Now().Sub(lastProvisionTime)
			if diff < duration {
				log.Warn.Println("During the current recommendation's decision time, there was already a successful provision. Therefore, discarding this recommendation until next polling interval.")
				// Warning message for huge decision periods.
				if duration > time.Duration(12)*time.Hour {
					log.Warn.Println("The current wait time until next provision if recommended is ", duration-diff)
					log.Warn.Println("If you believe the delay is too long, please consider reducing the decision period of your rule.")
				}
				return
			}

			TriggerProvision(clusterCfg, usrCfg, state, numNodes, t, operation, ruleResponsible)
		} else {
			log.Warn.Println("Recommendation can not be provisioned as open search cluster is already in provisioning phase.")
		}
	}
}

// Input:
//
// Description:
//
//	Generates the query string to get the latest document of successful Provision
//
// Return:
//
//	(string): Returns the query string that can be given as an OS query api parameter.
func getLatestProvisionQuery() string {
	return `{
                  "size": 1,
                  "sort": {
                    "Timestamp": "desc"
                  },
                  "query": {
                    "bool": {
                      "must": [
                        {
                          "match": {
                            "StatTag": "ProvisionStats"
                          }
                        },
                        {
                          "match": {
                            "Status": "Success"
                          }
                        }
                      ]
                    }
                  }
                }`
}
