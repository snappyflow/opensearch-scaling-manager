// This package will fetch the recommendation from the recommendation Queue and provision the scale in/out
// based on command.
package provision

import (
	"fmt"
	"scaling_manager/cluster"
	"scaling_manager/config"
	"time"
)

var counter uint8 = 1

// This struct contains the State of the opensearch scaling manager
// States can be of following types:
//  1. normal : This is the state when the recommnedation will be provisioned.
//  2. provision_scaleup/provision_scaledown : Once the trigger module will call provision it will set this state.
//  3. provisioning_scaleup/provisioning_scaledown : Once the provision module will start provisioning it will set this state.
//  4. provisioning_scaleup_completed/provisioning_scaledown_completed : Once the provision is completed then this state will be state.
//  5. provisioning_scaleup_failed/provisioning_scaledown_failed: If the provision is failed then this state will be set.
//  6. provisioned_successfully: If the provision is completed and cluster state is green then
//     this state will be set.
//  7. provisioned_failed: If the provision is completed and the cluster state is not green after
//     certain retries then this state will be set.
type State struct {
	// CurrentState indicate the current state of the scaling manager
	CurrentState string
	// PreviousState indicates the previous state of the scaling manager
	PreviousState string
	// Remark indicates the additional remarks for the state of the scaling manager
	Remark string
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
//		Think about the scenario where event based scaling needs to be performed.
//		Morning need to scale up and evening need to scale down.
//		If in morning the scale up was not successful then we should not perform the scale down.
//		May be we can keep a concept of minimum number of nodes as a configuration input.
//
// Return:
func (c *Command) Provision(state *State) {
	current_state := state.GetCurrentState()
	if c.Operation == "scale_up" {
		state.SetState("provisioning_scaleup", current_state)
		isScaledUp := c.ScaleOut(1, state)
		if isScaledUp {
			fmt.Println("Scaleup successful")
		} else {
			current_state = state.GetCurrentState()
			// Add a retry mechanism
			state.SetState("provisioning_scaleup_failed", current_state)
		}
	} else if c.Operation == "scale_down" {
		state.SetState("provisioning_scaledown", current_state)
		isScaledDown := c.ScaleIn(1, state)
		if isScaledDown {
			fmt.Println("Scaledown successful")
		} else {
			current_state = state.GetCurrentState()
			// Add a retry mechanism
			state.SetState("provisioning_scaledown_failed", current_state)
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
//		ScaleOut will scale out the cluster with the number of nodes.
//		This function will invoke commands to create a VM based on cloud type.
//	 	Then it will configure the opensearch on newly created nodes.
//
// Return:
//
//	Return the status of scale out of the nodes.
func (c *Command) ScaleOut(numNodes int, state *State) bool {
	// Read the current state of scaleup process and proceed with next step
	// If no stage was already set. The function returns an empty string. Then, start the scaleup process
	if state.GetCurrentState() == "provisioning_scaleup" {
		state.SetState("start_scaleup_process", "provision_scaleup")
		fmt.Println("Starting scaleUp process")
	}
	// Spin new VMs based on number of nodes and cloud type
	if state.GetCurrentState() == "start_scaleup_process" {
		fmt.Println("Spin new vms based on the cloud type")
		state.SetState("scaleup_triggered_spin_vm", "start_scaleup_process")
		fmt.Println("Spinning new vms")
		time.Sleep(5 * time.Second)
	}
	// Add the newly added VM to the list of VMs
	// Configure OS on newly created VM
	if state.GetCurrentState() == "scaleup_triggered_spin_vm" {
		fmt.Println("Check if the vm creation is complete and wait till done")
		fmt.Println("Adding the spinned nodes into the list of vms")
		fmt.Println("Configure ES")
		state.SetState("provisioning_scaleup_completed", "scaleup_triggered_spin_vm")
		fmt.Println("Configuring in progress")
		time.Sleep(5 * time.Second)
	}
	// Check cluster status after the configuration
	if state.GetCurrentState() == "provisioning_scaleup_completed" {
		fmt.Println("Wait for the cluster health and return status")
		fmt.Println("Waiting for the cluster to become healthy")
		time.Sleep(5 * time.Second)
		cluster := cluster.GetClusterCurrent()
		// Remove later
		cluster.ClusterDynamic.ClusterStatus = "yellow"
		for i := 0; i <= 12; i++ {
			if cluster.ClusterDynamic.ClusterStatus == "green" {
				current_state := state.GetCurrentState()
				state.SetState("provisioned_successfully", current_state)
				fmt.Println("Provisioned successfully")
				break
			}
			time.Sleep(10 * time.Second)
			cluster.ClusterDynamic.ClusterStatus = "green"
		}
		current_state := state.GetCurrentState()
		if current_state != "provisioned_successfully" {
			state.SetState("provisioned_failed", current_state)
			fmt.Println("Failed to provision")
		}
		time.Sleep(5 * time.Second)
		state.SetState("normal", state.GetCurrentState())
		fmt.Println("State set back to normal")
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
func (c *Command) ScaleIn(numNodes int, state *State) bool {
	// Read the current state of scaledown process and proceed with next step
	// If no stage was already set. The function returns an empty string. Then, start the scaledown process
	if state.GetCurrentState() == "provisioning_scaledown" {
		state.SetState("start_scaledown_process", "provision_scaledown")
		fmt.Println("Staring scaleDown process")
	}

	// Identify the node which can be removed from the cluster.
	if state.GetCurrentState() == "start_scaledown_process" {
		fmt.Println("Identify the node to remove from the cluster and store the node_ip")
		time.Sleep(2 * time.Second)
		state.SetState("scaledown_node_identified", "start_scaledown_process")
	}
	// Configure OS to tell master node that the present node is going to be removed
	if state.GetCurrentState() == "scaledown_node_identified" {
		fmt.Println("Configure ES to remove the node ip from cluster")
		state.SetState("provisioning_scaledown_completed", "scaledown_node_identified")
		time.Sleep(5 * time.Second)
		fmt.Println("Node removed from ES configuration")
	}
	// Wait for cluster to be in stable state(Shard rebalance)
	// Shut down the node
	if state.GetCurrentState() == "provisioning_scaledown_completed" {
		fmt.Println("Wait for the cluster to become healthy (in a loop of 5*12 minutes) and then proceed")
		cluster := cluster.GetClusterCurrent()
		// Remove later
		cluster.ClusterDynamic.ClusterStatus = "yellow"
		for i := 0; i <= 12; i++ {
			if cluster.ClusterDynamic.ClusterStatus == "green" {
				current_state := state.GetCurrentState()
				state.SetState("provisioned_successfully", current_state)
				break
			}
			// Remove later
			cluster.ClusterDynamic.ClusterStatus = "green"
			time.Sleep(300 * time.Second)
		}
		current_state := state.GetCurrentState()
		if current_state != "provisioned_successfully" {
			state.SetState("provisioned_failed", current_state)
			fmt.Println("Cluster hasn't come back to healthy state. Returning false")
		}
		time.Sleep(5 * time.Second)
		fmt.Println("Shutdown the node")
		state.SetState("normal", state.GetCurrentState())
	}
	return true
}

// Input:
// Description:
//
//		CheckClusterHealth will check the current cluster health and also check if there are any relocating
//	 	shards. If the cluster status is green and there are no relocating shard then we will update the status
//	 	to provisioned_successfully. Else, we will wait for 3 minutes and perform this check again for 3 times.
//
// Return:
func CheckClusterHealth(state *State) {
	cluster := cluster.GetClusterCurrent()
	if cluster.ClusterDynamic.ClusterStatus == "green" {
		current_state := state.GetCurrentState()
		state.SetState("provisioned_successfully", current_state)
	} else if counter >= 3 {
		time.Sleep(180 * time.Second)
		CheckClusterHealth(state)
	} else {
		current_state := state.GetCurrentState()
		state.SetState("provisioned_failed", current_state)
	}
	// We should wait for buffer period after provisioned_successfully state to stablize the cluster.
	// After that buffer period we should change the state to normal, which can tell trigger module to trigger
	// the recommendation.
}

// Input:
// Description:
//		Read the current stage that the provisioning process is in from Elasticsearch or any centralized DB which will be updated after each stage.
// Return: Stage returned from ES

func readStageFromEs() string {
	return ""
}

