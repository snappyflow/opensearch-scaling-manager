package fetchmetrics

import (
	"context"
	"encoding/json"
	utils "fetchMetrics/utils"
	"log"
	"os/exec"
	"strconv"
	"strings"

	elasticsearch "github.com/opensearch-project/opensearch-go"
	esapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// Description: NodeMetrics struct holds the node level metrics that are to be populated and indexed to elasticsearch
type NodeMetrics struct {
	utils.Node
	Timestamp int64
	StatTag   string
}

// Input:map[string]interface which holds the node stats and nodeId which used to parse the node stats response
// Description: The function calculates and returns disk utilization.
// The disk utilization is fetched from the node stats response from elasticsearch, there is no difference in terms
// of output when we fetch from linux or elasticsearch. And to fetch from the linux we need to make a call to elasticsearch
// to find where the binary of es is installed, when we make this call we get the disk utilization along with the mount path
// hence the overhead of calculating from linux is minimized by making use of response from elasticsearch.
// Output: Returns the disk utilization
func getDiskUtil(m map[string]interface{}, nodeId string) float32 {
	//Parse the node stats interface for required info
	list := m["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["fs"].(map[string]interface{})["data"].([]interface{})
	listJson, err := json.MarshalIndent(list[0], "", " ")
	if err != nil {
		log.Fatal("Cannot marshall: ", err)
	}
	var listInterface map[string]interface{}
	err1 := json.Unmarshal(listJson, &listInterface)
	if err1 != nil {
		log.Fatal("Unmarshal Error: ", err1)
	}
	// disk utilization = ((total space - available space) / total space) *100
	return float32(((listInterface["total_in_bytes"].(float64) - listInterface["available_in_bytes"].(float64)) / listInterface["total_in_bytes"].(float64)) * 100)
}

// Description: The functions fetchs the CPU utilization directly from the system through linux commands
// Output: Returns the CPU utilization
func getCpuUtil() float32 {
	//Executing the top command to fetch the CPU utilization
	cmd := exec.Command("bash", "-c", "top -bn2 | grep '%Cpu' | tail -1 | awk '{print 100-$8}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("error: ", out)
	}

	//Converting the result to float
	cpuFloat, err := strconv.ParseFloat(strings.TrimSuffix(string(out), "\n"), 8)
	if err != nil {
		log.Fatal("Unable to convert to float: ", err)
	}
	return float32(cpuFloat)
}

// Description:The functions fetchs the Memory utilization directly from the system through linux commands
// Output: Returns the memory utilization of the system.
func getRamUtil() float32 {
	//Executing the top command to fetch the memory utilization
	cmd := exec.Command("bash", "-c", "top -bn2 | grep 'KiB Mem' | tail -1 | awk '{print ($8/$4)*100}'")
	memout, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	//Converting the result to float
	memFloat, err := strconv.ParseFloat(strings.TrimSuffix(string(memout), "\n"), 8)
	if err != nil {
		log.Fatal("Unable to convert to float: ", err)
	}

	return float32(memFloat)
}

// Input: opensearch client and context
// Description: The function fetches and indexes the node stats
func IndexNodeStats(esClient *elasticsearch.Client, ctx context.Context) {

	nodeMetrics := new(NodeMetrics)

	//creating a node stats requests with filter to reduce the response to requirement
	nodes := []string{"_local"}
	metric := []string{"jvm", "os", "fs", "indices"}
	nodeStatReq, err := esapi.NodesStatsRequest{
		Pretty: true,
		NodeID: nodes,
		Metric: metric,
	}.Do(ctx, esClient)
	if err != nil {
		log.Fatal("Node stat fetch error: ", err)
	}

	//A map to dump the values from node stats response
	var nodeStatsInterface map[string]interface{}

	//Decoding the response into the the interface
	decodeErr := json.NewDecoder(nodeStatReq.Body).Decode(&nodeStatsInterface)
	if decodeErr != nil {
		log.Fatal("decode Error: ", decodeErr)
	}

	//parsing the interface and populating the node stats structure
	nodeId := utils.GetNodeId(nodeStatsInterface["nodes"].(map[string]interface{}))
	nodeMetrics.NodeId = nodeId
	nodeMetrics.NodeName = nodeStatsInterface["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["name"].(string)
	nodeMetrics.Timestamp = int64(nodeStatsInterface["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["timestamp"].(float64))
	nodeMetrics.HostIp = nodeStatsInterface["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["host"].(string)
	nodeMetrics.IsMaster = utils.CheckIfMaster(esClient, ctx)
	nodeMetrics.IsData = utils.CheckIfData(esClient, ctx)
	nodeMetrics.CpuUtil = getCpuUtil()
	nodeMetrics.RamUtil = getRamUtil()
	nodeMetrics.HeapUtil = float32(nodeStatsInterface["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["jvm"].(map[string]interface{})["mem"].(map[string]interface{})["heap_used_percent"].(float64))
	nodeMetrics.DiskUtil = getDiskUtil(nodeStatsInterface, nodeId)
	nodeMetrics.NumShards = int(nodeStatsInterface["nodes"].(map[string]interface{})[nodeId].(map[string]interface{})["indices"].(map[string]interface{})["shard_stats"].(map[string]interface{})["total_count"].(float64))
	nodeMetrics.StatTag = "NodesStats"

	//marshall the node metrics, to index into the elasticsearch
	nodeMetricsJson, err := json.MarshalIndent(nodeMetrics, "", "\t")
	if err != nil {
		log.Fatal("Error converting struct to Json: ", err)
	}

	utils.CheckIfIndexExists(esClient, ctx)
	utils.IndexMetrics(ctx, esClient, nodeMetricsJson)
}
