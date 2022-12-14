package provision

import (
	"regexp"
	"scaling_manager/cluster"
	"scaling_manager/config"
	"strconv"

	log "scaling_manager/logger"
)

// Input:
// Description:
//
//	GetRecommendation will fetch the recommendation from recommendation queue and clear the queue.
//	It will populate the command queue which contains all the details to scale out the cluster.
//
// Return:
func GetRecommendation(state *State, recommendation_queue []string) {
	scaleRegexString := `(scale_up|scale_down)_by_([0-9]+)`
	scaleRegex := regexp.MustCompile(scaleRegexString)
	if len(recommendation_queue) > 0 {
		clusterCurrent := cluster.GetClusterCurrent()
		state.GetCurrentState()
		if clusterCurrent.ClusterDynamic.ClusterStatus == "green" && state.CurrentState == "normal" {
			var command Command
			// Fill in the command struct with the recommendation queue and config file and trigger the recommendation.
			subMatch := scaleRegex.FindStringSubmatch(recommendation_queue[0])
			command.NumNodes, _ = strconv.Atoi(subMatch[2])
			command.Operation = subMatch[1]
			configStruct, err := config.GetConfig("config.yaml")
			if err != nil {
				log.Warn(log.ProvisionerWarn, "Unable to get Config from GetConfig()")
				return
			}
			command.ClusterDetails = configStruct.ClusterDetails
			go command.TriggerProvision(state)
		} else {
			log.Warn(log.ProvisionerWarn, "Recommendation can not be provisioned as open search cluster is already in provisioning phase or the cluster isn't healthy yet")
		}
	}
}
