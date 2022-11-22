package utils

import (
	"context"
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
	"bytes"
	elasticsearch "github.com/opensearch-project/opensearch-go"
	esapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// This struct will contain the dynamic metrics of the cluster.
type ClusterDynamic struct {
	// NumNodes indicates the number of nodes present in the OpenSearch cluster at any time.
	NumNodes int
	// ClusterStatus indicates the present state of a cluster.
	// 	red: One or more primary shards are unassigned, so some data is unavailable.
	//		This can occur briefly during cluster startup as primary shards are assigned.
	//  yellow: All primary shards are assigned, but one or more replica shards are unassigned.
	//		If a node in the cluster fails, some data could be unavailable until that node is repaired.
	//  green: All shards are assigned.
	ClusterStatus string
	// NumActiveShards indicates the total number of active primary and replica shards.
	NumActiveShards int
	// NumActivePrimaryShards indicates the number of active primary shards.
	NumActivePrimaryShards int
	// NumInitializingShards indicates the number of shards that are under initialization.
	NumInitializingShards int
	// NumUnassignedShards indicats the number of shards that are not allocated.
	NumUnassignedShards int
	// NumRelocatingShards indicates the number of shards that are under relocation.
	NumRelocatingShards int
	// NumMasterNodes indicates the number of master eligible nodes present in the cluster.
	NumMasterNodes int
	// NumActiveDataNodes indicates the number of active data nodes present in the cluster.
	NumActiveDataNodes int
}

// This struct will contain node metrics for a node in the OpenSearch cluster.
type Node struct {
	// NodeId indicates a unique ID of the node given by OpenSearch.
	NodeId string
	// NodeName indicates human-readable identifier for a particular instance of OpenSearch which is a configurable input.
	NodeName string
	// HostIp indicates the IP address of the node.
	HostIp string
	// IsMater indicates if the node is a master node.
	IsMaster bool
	// IsData indicates if the node is a data node.
	IsData bool
	// CpuUtil indicates the overall CPU Utilization in percentage for a node.
	CpuUtil float32
	// MemUtil indicates the overall Memory Utilization in percentage for a node.
	RamUtil float32
	// HeapUtil indicates the overall Java Heap Utilization in percentage for a node.
	HeapUtil float32
	// DiskUtil indicates the overall Disk Utilization in percentage for a node.
	DiskUtil float32
	// NumShards Number of shards present on a node.
	NumShards int
}

//Index that holds the node and cluster level metrics
const (
	IndexName string = "monitor-stats-1"
)

//Input: context, opensearch client and the json document that is to be indexed to the elasticsearch
//Description: Indexes the json document to the elasticsearch
//Output: Document indexed to elasticsearch index
func IndexMetrics(ctx context.Context, esClient *elasticsearch.Client, jsonDoc []byte) {
	//Create an index request and pass the document to be indexed along with the index name 
	indexMetricsRequest,err:=esapi.IndexRequest{
		Index: IndexName,
		DocumentType: "_doc",
		Body: bytes.NewReader(jsonDoc),
	}.Do(ctx,esClient)
	if err != nil {
		log.Fatal("Error indexing document: ", err)
	}
	fmt.Println("Document index successfull!: ",indexMetricsRequest)
}

//Input: opensearch client and context 
//Description:The function checks if index exists, if it exists it does nothing and returns. If it does not exists 
//It creates the index and returns
//Output: Cretes a new index if does not exists
func CheckIfIndexExists(esClient *elasticsearch.Client, ctx context.Context) {
	//Read the mappings file to create index with mappings if index is not present
	mappings, err := ioutil.ReadFile("mappings.json")
	if err != nil {
		log.Fatal("Unable to find mappings: ", err)
	}

	var indexName = []string{IndexName}

	//Create a index exists request to fetch if index is already present or not
	exist, err := esapi.IndicesExistsRequest{
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
	indexCreateRequest, err := esapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  bytes.NewReader(mappings),
	}.Do(ctx, esClient)
	if err != nil {
		log.Fatal("Index create request error: ", err)
	}
	fmt.Println("Created!: ", indexCreateRequest)
}

//Input: opensearch client and context
//The function checks if the current node is having the data node role.
//Output: Returns a boolean value, true if current node has data role
func CheckIfData(esClient *elasticsearch.Client, ctx context.Context) bool {
	//interface to dump the node stats response 
	var nodeStatsInterface map[string]interface{}

	nodes := []string{"_local"}

	//Creating node stats request and fetching the node stats for the current node
	nodeStatReq, err := esapi.NodesStatsRequest{
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
	nodeId:=GetNodeId(nodeStatsInterface["nodes"].(map[string]interface{}))

	//get the roles for current node
	roles:=nodeStatsInterface["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["roles"].([]interface{})

	//check if the role of the node is data, if so return true
	for _,v:=range roles{
		if v == "data"{
			return true
		}
	}
	return false
}

//Inputs: map[string]interface which holds the node stats response 
//Description: Fetches the node ID of the current node 
//Output: string which describes node ID of the current node 
func GetNodeId(m map[string]interface{}) string {
	for k := range m {
		return k
	}
	return ""
}

//Input: opensearch client and context 
//Description: Checks if the current node is the master node.
//Output: A boolean value, true if current node is master, false if it is not.
func CheckIfMaster(esClient *elasticsearch.Client, ctx context.Context) bool {
	var clusterStateInterface map[string]interface{} //To store the cluster state info and parse for master node ID
	var nodeStatsInterface map[string]interface{} //To store current node stats and parse for current node ID 
	
	//Create cluster state request and fetch cluster state
	clusterState, err := esapi.ClusterStateRequest{}.Do(ctx,esClient)
	if err != nil {
		panic(err)
	}
	
	//Decoding the response and dumping in the cluster state interface 
	decodeErr := json.NewDecoder(clusterState.Body).Decode(&clusterStateInterface)
	if decodeErr != nil {
		log.Fatal("decode Error: ", decodeErr)
	}

	//Parsing interface to get the id of the master node
	masterNode:= clusterStateInterface["master_node"].(string)
	
	nodes := []string{"_local"}
	//Creating node stats request and fetching the node stats for the current node
	nodeStatReq, err := esapi.NodesStatsRequest{
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
	currentNode:=GetNodeId(nodeStatsInterface["nodes"].(map[string]interface{}))
	return masterNode == currentNode
}