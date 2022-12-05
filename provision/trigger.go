package provision

import (
	"fmt"
	"regexp"
	"scaling_manager/cluster"
	"scaling_manager/config"
	"strconv"
)

// Input:
// Description:
//
//	GetRecommendation will fetch the recommendation from recommendation queue and clear the queue.
//	It will populate the command queue which contains all the details to scale out the cluster.
//
// Return:
func GetRecommendation(state *State, recommendation_queue []string) {
	scaleRegexString := `(scale_up|down)_by_([0-9]+)`
	scaleRegex := regexp.MustCompile(scaleRegexString)
	if len(recommendation_queue) > 0 {
		clusterCurrent := cluster.GetClusterCurrent()
		// Remove later
		clusterCurrent.ClusterDynamic.ClusterStatus = "green"
		current_state := state.GetCurrentState()
		if clusterCurrent.ClusterDynamic.ClusterStatus == "green" && current_state == "normal" {
			var command Command
			// Fill in the command struct with the recommendation queue and config file and trigger the recommendation.
			subMatch := scaleRegex.FindStringSubmatch(recommendation_queue[0])
			command.NumNodes, _ = strconv.Atoi(subMatch[2])
			command.Operation = subMatch[1]
			configStruct := config.GetConfig("config.yaml")
			command.ClusterDetails = configStruct.ClusterDetails
			command.triggerRecommendation(state)
		} else {
			fmt.Println("Recommendation can not be provisioned as open search cluster is already in provisioning phase or the cluster isn't healthy yet")
		}
	}
}

// Input:
// Description:
//
//	triggerRecommendation will get the status of the provisioner
//	and cluster and trigger the provisioning.
//
// Return:
func (c *Command) triggerRecommendation(state *State) {
	clusterCurrent := cluster.GetClusterCurrent()
	// Remove later
	clusterCurrent.ClusterDynamic.ClusterStatus = "green"
	current_state := state.GetCurrentState()
	if clusterCurrent.ClusterDynamic.ClusterStatus == "green" && current_state == "normal" {
		if c.Operation == "scale_up" {
			state.SetState("provisioning_scaleup", current_state)
		} else if c.Operation == "scale_down" {
			state.SetState("provisioning_scaledown", current_state)
		}
		go c.Provision(state)
	} else {
		fmt.Println("Recommendation can not be provisioned as open search cluster is already in provisioning phase or the cluster isn't healthy yet")
	}
}

// Input:
// Description:
//
//	GetCurrentState will get the current state of provisioning system of the scaling manager.
//
// Return:
//
//	Returns a string which contains the current state.
func (s *State) GetCurrentState() string {
	if s.CurrentState == "" {
		s.CurrentState = "normal"
	}
	return s.CurrentState
}

// Input:
//
//	currentState(string): The current state for the provisioner.
//	previousState(string): The previous state for the provisioner.
//
// Description:
//
//	SetState will set the state of provisioning system of the scaling manager.
//
// Return:
func (s *State) SetState(currentState string, previousState string) {
	// set the state for the opensearch scaling manager
	// This state can be either pushed to OS or else kept locally.
	s.CurrentState = currentState
	s.PreviousState = previousState
}
