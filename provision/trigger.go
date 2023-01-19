package provision

import (
	"context"
	"regexp"
	"scaling_manager/cluster"
	"scaling_manager/config"
	"strconv"

	"github.com/opensearch-project/opensearch-go"
)

var ctx = context.Background()

// Input:
//
//	state: The current provisioning state of the system
//	recommendationQueue: Recommendations provided by the recommendation engine in the form of an array of strings
//
// Description:
//
//	GetRecommendation will fetch the recommendation from recommendation queue.
//	It will call the Provisioner with all the user defined configs.
//	Triggers the provisioning
//
// Return:
func GetRecommendation(state *State, recommendationQueue []map[string]string, osClient *opensearch.Client) {
	scaleRegexString := `(scale_up|scale_down)_by_([0-9]+)`
	scaleRegex := regexp.MustCompile(scaleRegexString)
	if len(recommendationQueue) > 0 {
		clusterCurrent := cluster.GetClusterCurrent(ctx, osClient)
		state.GetCurrentState()
		if clusterCurrent.ClusterStatus == "green" && state.CurrentState == "normal" {
			// Fill in the command struct with the recommendation queue and config file and trigger the recommendation.
			var subMatch []string
			var task string
			for task, _ = range recommendationQueue[0] {
				subMatch = scaleRegex.FindStringSubmatch(task)
			}

			numNodes, _ := strconv.Atoi(subMatch[2])
			operation := subMatch[1]
			configStruct, err := config.GetConfig("config.yaml")
			if err != nil {
				log.Warn.Println("Unable to get Config from GetConfig()", err)
				return
			}
			cfg := configStruct.ClusterDetails
			TriggerProvision(cfg, state, numNodes, osClient, operation, recommendationQueue[0][task])
		} else {
			log.Warn.Println("Recommendation can not be provisioned as open search cluster is already in provisioning phase or the cluster isn't healthy yet")
		}
	}
}
