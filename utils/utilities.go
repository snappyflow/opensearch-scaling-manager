package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	opensearch "github.com/opensearch-project/opensearch-go"
	osapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// Index that holds the node and cluster level metrics
const (
	IndexName string = "monitor-stats-1"
)

// Input: context, opensearch client and the json document that is to be indexed to the elasticsearch
// Description: Indexes the json document to the elasticsearch
// Output: Document indexed to elasticsearch index
func IndexMetrics(ctx context.Context, esClient *opensearch.Client, jsonDoc []byte) {
	//Create an index request and pass the document to be indexed along with the index name
	indexMetricsRequest, err := osapi.IndexRequest{
		Index:        IndexName,
		DocumentType: "_doc",
		Body:         bytes.NewReader(jsonDoc),
	}.Do(ctx, esClient)
	if err != nil {
		log.Fatal("Error indexing document: ", err)
	}
	fmt.Println("Document index successfull!: ", indexMetricsRequest)
}

// Input: opensearch client and context
// Description:The function checks if index exists, if it exists it does nothing and returns. If it does not exists
// It creates the index and returns
// Output: Cretes a new index if does not exists
func CheckIfIndexExists(esClient *opensearch.Client, ctx context.Context) {
	//Read the mappings file to create index with mappings if index is not present
	mappings, err := ioutil.ReadFile("mappings.json")
	if err != nil {
		log.Fatal("Unable to find mappings: ", err)
	}

	var indexName = []string{IndexName}

	//Create a index exists request to fetch if index is already present or not
	exist, err := osapi.IndicesExistsRequest{
		Index: indexName,
	}.Do(ctx, esClient)
	if err != nil {
		log.Fatal("Check index exists request error: ", err)
	}
	//If status code == 200 then index exists, print index exists, return
	if exist.StatusCode == 200 {
		fmt.Println("Index Exists!")
		return
	}
	//If status code is not 200 then index does not exist, so crete a new Index via index create request API,
	// pass mappings and index name.
	indexCreateRequest, err := osapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  bytes.NewReader(mappings),
	}.Do(ctx, esClient)
	if err != nil {
		log.Fatal("Index create request error: ", err)
	}
	fmt.Println("Created!: ", indexCreateRequest)
}

// Input: opensearch client and context
// The function checks if the current node is having the data node role.
// Output: Returns a boolean value, true if current node has data role
func CheckIfData(esClient *opensearch.Client, ctx context.Context) bool {
	//interface to dump the node stats response
	var nodeStatsInterface map[string]interface{}

	nodes := []string{"_local"}

	//Creating node stats request and fetching the node stats for the current node
	nodeStatReq, err := osapi.NodesStatsRequest{
		Pretty: true,
		NodeID: nodes,
	}.Do(ctx, esClient)
	if err != nil {
		log.Fatal("Node stat fetch error: ", err)
	}

	//Decoding the response and dumping the node stats in the interface
	nodeDecodeErr := json.NewDecoder(nodeStatReq.Body).Decode(&nodeStatsInterface)
	if nodeDecodeErr != nil {
		log.Fatal("decode Error: ", nodeDecodeErr)
	}

	//Parsing for the node id of the current node
	nodeId := GetNodeId(nodeStatsInterface["nodes"].(map[string]interface{}))

	//get the roles for current node
	roles := nodeStatsInterface["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["roles"].([]interface{})

	//check if the role of the node is data, if so return true
	for _, v := range roles {
		if v == "data" {
			return true
		}
	}
	return false
}

// Inputs: map[string]interface which holds the node stats response
// Description: Fetches the node ID of the current node
// Output: string which describes node ID of the current node
func GetNodeId(m map[string]interface{}) string {
	for k := range m {
		return k
	}
	return ""
}

// Input: opensearch client and context
// Description: Checks if the current node is the master node.
// Output: A boolean value, true if current node is master, false if it is not.
func CheckIfMaster(esClient *opensearch.Client, ctx context.Context) bool {
	var clusterStateInterface map[string]interface{} //To store the cluster state info and parse for master node ID
	var nodeStatsInterface map[string]interface{}    //To store current node stats and parse for current node ID

	//Create cluster state request and fetch cluster state
	clusterState, err := osapi.ClusterStateRequest{}.Do(ctx, esClient)
	if err != nil {
		panic(err)
	}

	//Decoding the response and dumping in the cluster state interface
	decodeErr := json.NewDecoder(clusterState.Body).Decode(&clusterStateInterface)
	if decodeErr != nil {
		log.Fatal("decode Error: ", decodeErr)
	}

	//Parsing interface to get the id of the master node
	masterNode := clusterStateInterface["master_node"].(string)

	nodes := []string{"_local"}
	//Creating node stats request and fetching the node stats for the current node
	nodeStatReq, err := osapi.NodesStatsRequest{
		Pretty: true,
		NodeID: nodes,
	}.Do(ctx, esClient)
	if err != nil {
		log.Fatal("Node stat fetch error: ", err)
	}

	//Decoding the response and dumping the node stats in the interface
	nodeDecodeErr := json.NewDecoder(nodeStatReq.Body).Decode(&nodeStatsInterface)
	if nodeDecodeErr != nil {
		log.Fatal("decode Error: ", nodeDecodeErr)
	}

	//Parsing for the node id of the current node
	currentNode := GetNodeId(nodeStatsInterface["nodes"].(map[string]interface{}))
	return masterNode == currentNode
}

func CreateOsClient(OsUsername, OsPassword string) *opensearch.Client {
	//create a configuration that is to be passed while creating the client
	cfg := opensearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Username: OsUsername,
		Password: OsPassword,
	}

	//create the client using the configuration
	osClient, err := opensearch.NewClient(cfg)
	if err != nil {
		fmt.Println("Opensearch connection error:", err)
	}
	res, err := osClient.Info()
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	}
	fmt.Println("Response: ", res)
	return osClient
}