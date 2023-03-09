package provision

import (
	"context"
	"encoding/json"
	"fmt"
	osutils "github.com/maplelabs/opensearch-scaling-manager/opensearchUtils"
	utils "github.com/maplelabs/opensearch-scaling-manager/utilities"
	"time"
)

// This struct contains the State of the opensearch scaling manager
// States can be of following types: (May change in real implementation with new stages identified)
// At any point, the state should have either "scaleup/scaledown" to identify the current operation happening
//
//   - normal : This is the state when the recommnedation will be provisioned.
//   - provisioning_scaleup/provisioning_scaledown : Once the provision module will start provisioning it will set this state.
//   - start_scaleup_process/start_scaledown_process : Indicates start of scaleup/scaledown process
//   - scaleup_triggered_spin_vm: Indicates trigger for spinning new vms while scaleup
//   - scaledown_node_identified: A state to identify node identification to scaledown
//   - provisioning_scaleup_completed/provisioning_scaledown_completed : Once the provision is completed then this state will be state.
//   - provisioning_scaleup_failed/provisioning_scaledown_failed: If the provision is failed then this state will be set.
//   - provisioned_scaleup_successfully/provisioned_scaledown_successfully: If the provision is completed and cluster state is green then this state will be set.
type State struct {
	// CurrentState indicate the current state of the scaling manager
	CurrentState string
	// PreviousState indicates the previous state of the scaling manager
	PreviousState string
	// Remark indicates the additional remarks for the state of the scaling manager
	Remark string
	// Last Provisioned time is when the last successful provision was completed
	LastProvisionedTime int64
	// Start time of the current provisioning in place
	ProvisionStartTime int64
	// Rule triggered for provisioning. i.e., scale_up/scale_down
	RuleTriggered string
	// Rule Responsible for provisioning. i.e., cpu, mem, heap, shard, disk.
	RulesResponsible string
	// Number of nodes being added(scale_up) / removed(scale_down) from the cluster due to current provision
	NumNodes int
	// Number of nodes remaining to be scaled up/scaled down
	RemainingNodes int
	// StatTag
	StatTag string
	// For snappyflow dashboard
	_documentType string
	// Timestamp
	Timestamp int64
        // Node Ip storage
        NodeIp   string
        // Node Name
        NodeName string

}

var state = new(State)

// A global variable which stores the document ID of the State document that will to stored and fetched frm Opensearch
var docId string

// Input:
//
// Description:
//
//	Creates a unique document ID for maintaining the state of the provisioning system and updates the global variable
//
// Return:
func InitializeDocId() {
	docId = fmt.Sprint(utils.Hash(utils.GetClusterId()))
}

// Input:
//
// Caller:
//      Object of type State
//
// Description:
//      GetCurrentState will update the state variable pointer such that it is in sync with the updated values.
//      Reads the document from Opensearch and updates the Struct
//
// Return:

func (s *State) GetCurrentState() {
	// Get the document.

	searchResponse, err := osutils.SearchDoc(context.Background(), docId)
	if err != nil {
		log.Panic.Println("failed to search document: ", err)
		panic(err)
	}
	defer searchResponse.Body.Close()
	var stateInterface map[string]interface{}
	log.Debug.Println("Get resp: ", searchResponse)
	if searchResponse.Status() == "404 Not Found" {
		//Setting the initial state
		s.CurrentState = "normal"
		s._documentType = "State"
		s.StatTag = "State"
		s.UpdateState()
		return
	}
	jsonErr := json.NewDecoder(searchResponse.Body).Decode(&stateInterface)
	if jsonErr != nil {
		log.Panic.Println("Unable to decode the response into interface: ", jsonErr)
		panic(jsonErr)
	}
	// convert map to json
	jsonString, errr := json.Marshal(stateInterface["_source"].(map[string]interface{}))
	if errr != nil {
		log.Panic.Println("Unable to unmarshal interface: ", errr)
		panic(errr)
	}

	// convert json to struct
	json.Unmarshal(jsonString, s)
}

// Input:
//
// Description:
//
//      Updates the opensearch document with the values in state Struct pointer.
//
// Return:

func (s *State) UpdateState() {
	// Update the document.

	s.Timestamp = time.Now().UnixMilli()

	state, err := json.Marshal(s)
	if err != nil {
		log.Panic.Println("json.Marshal ERROR: ", err)
		panic(err)
	}
	content := string(state)

	updateResponse, err := osutils.UpdateDoc(context.Background(), docId, content)
	if err != nil {
		log.Panic.Println("failed to update document: ", err)
		panic(err)
	}
	defer updateResponse.Body.Close()
	log.Debug.Println("Update resp: ", updateResponse)
}
