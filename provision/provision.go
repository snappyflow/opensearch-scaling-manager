// This package will fetch the recommendation from the recommendation Queue and provision the scale in/out
// based on command.
package provision

import (
	"bytes"
	"fmt"
	"net/http"
	"scaling_manager/cluster"
	"scaling_manager/config"
	"strings"
	"time"

	log "scaling_manager/logger"
)

var counter uint8 = 1

// This struct contains the State of the opensearch scaling manager
// States can be of following types: (May change in real implementation with new stages identified)
// At any point, the state should have either "scaleup/scaledown" to identify the current operation happening
//
//	normal : This is the state when the recommnedation will be provisioned.
//	provisioning_scaleup/provisioning_scaledown : Once the provision module will start provisioning it will set this state.
//	start_scaleup_process/start_scaledown_process : Indicates start of scaleup/scaledown process
//	scaleup_triggered_spin_vm: Indicates trigger for spinning new vms while scaleup
//	scaledown_node_identified: A state to identify node identification to scaledown
//	provisioning_scaleup_completed/provisioning_scaledown_completed : Once the provision is completed then this state will be state.
//	provisioning_scaleup_failed/provisioning_scaledown_failed: If the provision is failed then this state will be set.
//	provisioned_scaleup_successfully/provisioned_scaledown_successfully: If the provision is completed and cluster state is green then
//	   this state will be set.
type State struct {
	// CurrentState indicate the current state of the scaling manager
	CurrentState string
	// PreviousState indicates the previous state of the scaling manager
	PreviousState string
	// Remark indicates the additional remarks for the state of the scaling manager
	Remark              string
	LastProvisionedTime time.Time
	ProvisionStartTime  time.Time
	RuleTriggered       string
	NumNodes            int
	RemainingNodes      int
}

// This struct contains the operation and details to scale the cluster
type Command struct {
	// Operation indicates the operation will be performed by the provisioner.
	// As of now two operations can be performed by the provisioner:
	//  1) scale_up
	//  2) scale_down
	Operation string
	// NumNodes indicates the number of nodes need to be scaled in or out.
	NumNodes int
	config.ClusterDetails
}

// Input:
// Caller: Object of Command
// Description:
//
//	Provision will scale in/out the cluster based on the operation.
//	ToDo:
//	        Think about the scenario where event based scaling needs to be performed.
//	        Morning need to scale up and evening need to scale down.
//	        If in morning the scale up was not successful then we should not perform the scale down.
//	        May be we can keep a concept of minimum number of nodes as a configuration input.
//
// Return:
func (c *Command) TriggerProvision(state *State) {
	state.GetCurrentState()
	if c.Operation == "scale_up" {
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaleup"
		state.UpdateState()
		isScaledUp := c.ScaleOut(state)
		if isScaledUp {
			log.Info(log.ProvisionerInfo, "Scaleup successful")
		} else {
			state.GetCurrentState()
			// Add a retry mechanism
			state.PreviousState = state.CurrentState
			state.CurrentState = "provisioning_scaleup_failed"
			state.UpdateState()
		}
	} else if c.Operation == "scale_down" {
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaledown"
		state.UpdateState()
		isScaledDown := c.ScaleIn(state)
		if isScaledDown {
			log.Info(log.ProvisionerInfo, "Scaledown successful")
		} else {
			state.GetCurrentState()
			// Add a retry mechanism
			state.PreviousState = state.CurrentState
			state.CurrentState = "provisioning_scaledown_failed"
			state.UpdateState()
		}
	}
}

// Input:
//
//	numNodes(int): Number of nodes to scale out.
//
// Caller: Object of Command
// Description:
//
//	ScaleOut will scale out the cluster with the number of nodes.
//	This function will invoke commands to create a VM based on cloud type.
//	Then it will configure the opensearch on newly created nodes.
//
// Return:
//
//	Return the status of scale out of the nodes.
func (c *Command) ScaleOut(state *State) bool {
	// Read the current state of scaleup process and proceed with next step
	// If no stage was already set. The function returns an empty string. Then, start the scaleup process
	state.GetCurrentState()
	if state.CurrentState == "provisioning_scaleup" {
		log.Info(log.ProvisionerInfo, "Starting scaleUp process")
		state.PreviousState = state.CurrentState
		state.CurrentState = "start_scaleup_process"
		state.ProvisionStartTime = time.Now()
		state.RuleTriggered = "scale_up"
		state.NumNodes = c.NumNodes
		state.UpdateState()
	}
	// Spin new VMs based on number of nodes and cloud type
	if state.CurrentState == "start_scaleup_process" {
		log.Info(log.ProvisionerInfo, "Spin new vms based on the cloud type")
		log.Info(log.ProvisionerInfo, "Spinning new vms")
		time.Sleep(5 * time.Second)
		state.PreviousState = state.CurrentState
		state.CurrentState = "scaleup_triggered_spin_vm"
		state.UpdateState()
	}
	// Add the newly added VM to the list of VMs
	// Configure OS on newly created VM
	if state.CurrentState == "scaleup_triggered_spin_vm" {
		log.Info(log.ProvisionerInfo, "Check if the vm creation is complete and wait till done")
		log.Info(log.ProvisionerInfo, "Adding the spinned nodes into the list of vms")
		log.Info(log.ProvisionerInfo, "Configure ES")
		log.Info(log.ProvisionerInfo, "Configuring in progress")
		time.Sleep(5 * time.Second)
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaleup_completed"
		state.UpdateState()
	}
	// Check cluster status after the configuration
	if state.CurrentState == "provisioning_scaleup_completed" {
		SimulateSharRebalancing()
		log.Info(log.ProvisionerInfo, "Wait for the cluster health and return status")
		log.Info(log.ProvisionerInfo, "Waiting for the cluster to become healthy")
		time.Sleep(5 * time.Second)
		CheckClusterHealth(state)
		state.LastProvisionedTime = time.Now()
		state.ProvisionStartTime = time.Time{}
		state.PreviousState = state.CurrentState
		state.CurrentState = "normal"
		state.RuleTriggered = ""
		state.UpdateState()
		time.Sleep(5 * time.Second)
		log.Info(log.ProvisionerInfo, "State set back to normal")
	}
	return true
}

// Input:
//
//	numNodes(int): Number of nodes to scale in.
//
// Caller: Object of Command
// Description:
//
//	ScaleIn will scale in the cluster with the number of nodes.
//	This function will invoke commands to remove a node from opensearch cluster.
//
// Return:
//
//	Return the status of scale in of the nodes.
func (c *Command) ScaleIn(state *State) bool {
	// Read the current state of scaledown process and proceed with next step
	// If no stage was already set. The function returns an empty string. Then, start the scaledown process
	state.GetCurrentState()
	if state.CurrentState == "provisioning_scaledown" {
		log.Info(log.ProvisionerInfo, "Staring scaleDown process")
		state.PreviousState = state.CurrentState
		state.CurrentState = "start_scaledown_process"
		state.ProvisionStartTime = time.Now()
		state.RuleTriggered = "scale_down"
		state.NumNodes = c.NumNodes
		state.UpdateState()
	}

	// Identify the node which can be removed from the cluster.
	if state.CurrentState == "start_scaledown_process" {
		log.Info(log.ProvisionerInfo, "Identify the node to remove from the cluster and store the node_ip")
		time.Sleep(5 * time.Second)
		state.PreviousState = state.CurrentState
		state.CurrentState = "scaledown_node_identified"
		state.UpdateState()
	}
	// Configure OS to tell master node that the present node is going to be removed
	if state.CurrentState == "scaledown_node_identified" {
		log.Info(log.ProvisionerInfo, "Configure ES to remove the node ip from cluster")
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaledown_completed"
		state.UpdateState()
		time.Sleep(5 * time.Second)
		log.Info(log.ProvisionerInfo, "Node removed from ES configuration")
	}
	// Wait for cluster to be in stable state(Shard rebalance)
	// Shut down the node
	if state.CurrentState == "provisioning_scaledown_completed" {
		SimulateSharRebalancing()
		log.Info(log.ProvisionerInfo, "Wait for the cluster to become healthy (in a loop of 5*12 minutes) and then proceed")
		CheckClusterHealth(state)
		log.Info(log.ProvisionerInfo, "Shutdown the node")
		time.Sleep(5 * time.Second)
		state.LastProvisionedTime = time.Now()
		state.ProvisionStartTime = time.Time{}
		state.RuleTriggered = ""
		state.PreviousState = state.CurrentState
		state.CurrentState = "normal"
		state.UpdateState()
		log.Info(log.ProvisionerInfo, "State set back to normal")
	}
	return true
}

// Input:
// Description:
//
//	CheckClusterHealth will check the current cluster health and also check if there are any relocating
//	shards. If the cluster status is green and there are no relocating shard then we will update the status
//	to provisioned_successfully. Else, we will wait for 3 minutes and perform this check again for 3 times.
//
// Return:
func CheckClusterHealth(state *State) {
	for i := 0; i <= 12; i++ {
		cluster := cluster.GetClusterCurrent()
		log.Info(cluster.ClusterDynamic.ClusterStatus)
		if cluster.ClusterDynamic.ClusterStatus == "green" {
			state.GetCurrentState()
			state.PreviousState = state.CurrentState
			if strings.Contains(state.PreviousState, "scaleup") {
				state.CurrentState = "provisioned_scaleup_successfully"
			} else {
				state.CurrentState = "provisioned_scaledown_successfully"
			}
			state.UpdateState()
			break
		}
		log.Info(log.ProvisionerInfo, "Waiting for cluster to be healthy.......")
		time.Sleep(15 * time.Second)
	}
	state.GetCurrentState()
	if !(strings.Contains(state.CurrentState, "success")) {
		state.PreviousState = state.CurrentState
		if strings.Contains(state.PreviousState, "scaleup") {
			state.CurrentState = "provisioning_scaleup_failed"
		} else {
			state.CurrentState = "provisioning_scaledown_failed"
		}
		state.UpdateState()
		log.Warn(log.ProvisionerWarn, "Cluster hasn't come back to healthy state.")
	}
	// We should wait for buffer period after provisioned_successfully state to stablize the cluster.
	// After that buffer period we should change the state to normal, which can tell trigger module to trigger
	// the recommendation.
}

// Input:
// Description:
//              Read the current stage that the provisioning process is in from Elasticsearch or any centralized DB which will be updated after each stage.
// Return: Stage returned from ES

func readStageFromEs() string {
	return ""
}

func SimulateSharRebalancing() {
	// Add logic to call the simulator's end point
	var jsonStr = []byte(`{"nodes":1}`)
	urlLink := fmt.Sprintf("http://localhost:5000/provision/addnode")
	req, err := http.NewRequest("POST", urlLink, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(log.ProvisionerFatal, err)
	}

	defer resp.Body.Close()
}
