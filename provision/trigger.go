package provision

import (
	"context"
	"regexp"
	"scaling_manager/cluster"
	"scaling_manager/cluster_sim"
	"scaling_manager/config"
	"strconv"
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
			TriggerProvision(clusterCfg, usrCfg, state, numNodes, t, operation, recommendationQueue[0][task])
		} else {
			log.Warn.Println("Recommendation can not be provisioned as open search cluster is already in provisioning phase.")
		}
	}
}
