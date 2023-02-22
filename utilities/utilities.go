package utilities

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/maplelabs/opensearch-scaling-manager/logger"
	"github.com/maplelabs/opensearch-scaling-manager/config"
	osutils "github.com/maplelabs/opensearch-scaling-manager/opensearchUtils"
	"hash/fnv"
	"os"
)

// A global logger variable used across the package for logging.
var log = new(logger.LOG)

func init() {
	log.Init("logger")
	log.Info.Println("Utilities module initiated")
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//	nodeId (string): The node id which needs to be checked if a master
//
// Description:
//
//	Checks if the current node or the node with the passed nodeId is a master or not.
//
// Output:
//
//	(bool): A boolean value, true if current node is master, false if it is not.
func CheckIfMaster(ctx context.Context, nodeId string) bool {
	var clusterStateInterface map[string]interface{} //To store the cluster state info and parse for master node ID
	var nodeStatsInterface map[string]interface{}    //To store current node stats and parse for current node ID

	//Create cluster state request and fetch cluster state
	clusterState, err := osutils.GetClusterState(ctx)
	if err != nil {
		panic(err)
	}
	defer clusterState.Body.Close()

	//Decoding the response and dumping in the cluster state interface
	decodeErr := json.NewDecoder(clusterState.Body).Decode(&clusterStateInterface)
	if decodeErr != nil {
		log.Panic.Println("decode Error: ", decodeErr)
		panic(err)
	}

	//Parsing interface to get the id of the master node
	masterNode := clusterStateInterface["master_node"].(string)

	if nodeId != "" {
		return masterNode == nodeId
	}

	nodes := []string{"_local"}

	//Creating node stats request and fetching the node stats for the current node
	nodeStatReq, err := osutils.GetNodeStats(ctx, nodes, nil)
	if err != nil {
		log.Panic.Println("Node stat fetch error: ", err)
		panic(err)
	}
	defer nodeStatReq.Body.Close()

	//Decoding the response and dumping the node stats in the interface
	nodeDecodeErr := json.NewDecoder(nodeStatReq.Body).Decode(&nodeStatsInterface)
	if nodeDecodeErr != nil {
		log.Panic.Println("decode Error: ", nodeDecodeErr)
		panic(err)
	}

	//Parsing for the node id of the current node
	var currentNode string
	for node := range nodeStatsInterface["nodes"].(map[string]interface{}) {
		currentNode = node
		break
	}
	return masterNode == currentNode
}

// Input:
//
// Description:
//
//	Read the cluster UUID from the Cluster Stats api response and return the id.
//
// Output:
//
//	(string): Return the cluster UUID of the cluster
func GetClusterId() string {
	var clusterStatsInterface map[string]interface{}
	resp, err := osutils.GetClusterStats(context.Background())
	if err != nil {
		log.Error.Println("cluster Stats fetch ERROR:", err)
	}
	defer resp.Body.Close()

	decodeErr := json.NewDecoder(resp.Body).Decode(&clusterStatsInterface)
	if decodeErr != nil {
		log.Error.Println("decode Error: ", decodeErr)
	}

	return clusterStatsInterface["cluster_uuid"].(string)
}

// Input:
//
// Description:
//
//	Return the map which contains the node details of each node present in the cluster
//
// Output:
//
//	(string): Return the map which contains node details
func GetNodes() map[string]interface{} {
	var nodeStatsInterface map[string]interface{}

	nodes := []string{"_all"}
	metrics := []string{}

	nodeStatResp, err := osutils.GetNodeStats(context.Background(), nodes, metrics)
	if err != nil {
		log.Error.Println("Node stat fetch error: ", err)
	}
	defer nodeStatResp.Body.Close()

	decodeErr := json.NewDecoder(nodeStatResp.Body).Decode(&nodeStatsInterface)
	if decodeErr != nil {
		log.Error.Println("decode Error: ", decodeErr)
	}

	nodeMap := make(map[string]interface{}, 0)

	for node, nodeInfo := range nodeStatsInterface["nodes"].(map[string]interface{}) {
		nodeInfoMap := nodeInfo.(map[string]interface{})
		nodeMap[node] = map[string]string{"name": nodeInfoMap["name"].(string), "hostIp": nodeInfoMap["host"].(string)}
	}

	return nodeMap

}

// Input:
//	s (string): String to be hashed
//
// Description:
//	Returns a hashed value of the string passed as input
//
// Output:
//	(uint32): Hashed value of string

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Input:
//
//	mapp (map[string]interface{}): A map of length 1 which contains the local node details
//
// Description:
//
//	Parses the first element in the map and returns the first matched nodeID
//
// Output:
//
//	(string): Node Id as string
func ParseNodeId(mapp map[string]interface{}) string {
	for node := range mapp {
		return node
	}
	return ""
}

func HostsWithCurrentNodes(fileName string, clusterCfg config.ClusterDetails) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal.Println(err)
		panic(err)
	}
	defer f.Close()
	nodes := GetNodes()
	dataWriter := bufio.NewWriter(f)
	dataWriter.WriteString("[current_nodes]\n")
	for _, nodeIdMap := range nodes {
		_, writeErr := dataWriter.WriteString(nodeIdMap.(map[string]string)["name"] + " ansible_user=" + clusterCfg.SshUser + " roles=master,data,ingest ansible_private_host=" + nodeIdMap.(map[string]string)["hostIp"] + " ansible_ssh_private_key_file=" + clusterCfg.CloudCredentials.PemFilePath + "\n")
		if writeErr != nil {
			log.Error.Println("Error writing the node data into hosts file", writeErr)
			panic(err)
		}
	}
	dataWriter.Flush()
}
