package provision

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"scaling_manager/cluster"
	log "scaling_manager/logger"
	"strings"

	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

var docId = fmt.Sprint(hash(cluster.GetClusterId()))

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

const IndexName = "monitor-stats-1"

var client *opensearch.Client

func init() {

	var err error
	client, err = opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Fatal(log.ProvisionerError, err)
		os.Exit(1)
	}

	mappingFile, err := os.ReadFile("provision/mappings.json") // just pass the file name
	if err != nil {
		log.Error(log.ProvisionerError, err)
	}
	mapping := string(mappingFile)

	createNewIndexWithMappings(mapping)
}

func createNewIndexWithMappings(mapping string) {
	ctx := context.Background()
	createReq := opensearchapi.IndicesCreateRequest{}
	createReq.Index = IndexName
	createReq.Body = strings.NewReader(mapping)
	req := opensearchapi.IndicesExistsRequest{}
	req.Index = []string{IndexName}
	resp, err := req.Do(ctx, client)
	if err != nil {
		log.Fatal(log.ProvisionerError, fmt.Sprintf("Index exists check failed: %v", err))
	}
	log.Info(log.ProvisionerInfo, "Index already exists")
	if resp.Status() != "200 OK" {
		res, err := createReq.Do(ctx, client)
		if err != nil {
			log.Info(log.ProvisionerInfo, fmt.Sprintf("Create Index request error: %v ", err))
		}
		log.Info(fmt.Sprintf("Index create Response: %v", res))
	}
}

// Input:
// Description:
//
//      GetCurrentState will get the current state of provisioning system of the scaling manager.
//
// Return:
//
//      Returns a string which contains the current state.

func (s *State) GetCurrentState() {
	// Get the document.

	search := opensearchapi.GetRequest{
		Index:      IndexName,
		DocumentID: fmt.Sprint(docId),
	}

	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		log.Error(log.ProvisionerError, fmt.Sprintf("failed to search document: %v ", err))
		os.Exit(1)
	}
	var stateInterface map[string]interface{}
	log.Info(log.ProvisionerInfo, fmt.Sprintf("Get resp: %v ", searchResponse))
	if searchResponse.Status() == "404 Not Found" {
		//Setting the initial state
		s.CurrentState = "normal"
		s.UpdateState()
	}
	jsonErr := json.NewDecoder(searchResponse.Body).Decode(&stateInterface)
	if jsonErr != nil {
		log.Error(log.ProvisionerError, fmt.Sprintf("Unable to decode the response into interface: %v", jsonErr))
		return
	}
	// convert map to json
	jsonString, errr := json.Marshal(stateInterface["_source"].(map[string]interface{}))
	if errr != nil {
		log.Fatal(log.ProvisionerFatal, fmt.Sprintf("Unable to unmarshal interface: %v", errr))
	}

	// convert json to struct
	json.Unmarshal(jsonString, s)
}

// Input:
//
//      currentState(string): The current state for the provisioner.
//      previousState(string): The previous state for the provisioner.
//
// Description:
//
//      SetState will set the state of provisioning system of the scaling manager.
//
// Return:

func (s *State) UpdateState() {
	// Update the document.

	state, err := json.Marshal(s)
	if err != nil {
		log.Fatal(log.ProvisionerError, fmt.Sprintf("json.Marshal ERROR: %v", err))
	}
	content := string(state)

	updateReq := opensearchapi.IndexRequest{
		Index:      IndexName,
		DocumentID: fmt.Sprint(docId),
		Body:       strings.NewReader(content),
	}

	updateResponse, err := updateReq.Do(context.Background(), client)
	if err != nil {
		log.Fatal(log.ProvisionerError, fmt.Sprintf("failed to update document: %v ", err))
	}
	log.Info(log.ProvisionerInfo, fmt.Sprintf("Update resp: %v ", updateResponse))
}
