package fetchmetrics

import (
	"context"
	"encoding/json"
	"log"

	"scaling_manager/utils"
	"scaling_manager/cluster"
	opensearch "github.com/opensearch-project/opensearch-go"
	osapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// Description: ClusterMetrics holds the cluster level information that are to be populated and indexed into elasticsearch
type ClusterMetrics struct {
	cluster.ClusterDynamic
	Timestamp   int64
	StatTag     string
	ClusterName string
}

// Input: opensearch client and context.
// Description: Fetches cluster level info and populates ClusterMetrics struct
// Output: Returns the populated ClusterMetrics struct
func FetchClusterHealthMetrics(esClient *opensearch.Client, ctx context.Context) ClusterMetrics {

	//Create an interface to capture the response from cluster health and cluster stats API
	var clusterStatsInterface map[string]interface{}
	var clusterHealthInterface map[string]interface{}

	clusterStats := new(ClusterMetrics)

	//Create a cluster stats request and fetch the response
	clusterStatsRequest, err := osapi.ClusterStatsRequest{}.Do(ctx, esClient)
	if err != nil {
		log.Fatalf("cluster stats fetch ERROR:", err)
	}

	//decode and dump the cluster stats response into interface
	decodeErr := json.NewDecoder(clusterStatsRequest.Body).Decode(&clusterStatsInterface)
	if decodeErr != nil {
		log.Fatal("decode Error: ", decodeErr)
	}

	//Parse the interface and populate required fields in cluster stats
	clusterStats.NumActiveDataNodes = int(clusterStatsInterface["nodes"].(map[string]interface{})["count"].(map[string]interface{})["data"].(float64))
	clusterStats.NumMasterNodes = int(clusterStatsInterface["nodes"].(map[string]interface{})["count"].(map[string]interface{})["master"].(float64))
	clusterStats.Timestamp = int64(clusterStatsInterface["timestamp"].(float64))

	//create a cluster health request and fetch cluster health
	clusterHealthRequest, err := osapi.ClusterHealthRequest{}.Do(ctx, esClient)
	if err != nil {
		log.Fatalf("cluster Health fetch ERROR:", err)
	}

	//Decode the response and dump the response into the cluster health interface
	decodeErr2 := json.NewDecoder(clusterHealthRequest.Body).Decode(&clusterHealthInterface)
	if decodeErr2 != nil {
		log.Fatal("decode Error: ", decodeErr)
	}

	//Parse the interface and populate required fields in cluster stats
	clusterStats.NumNodes = int(clusterHealthInterface["number_of_nodes"].(float64))
	clusterStats.ClusterStatus = clusterStatsInterface["status"].(string)
	clusterStats.NumActiveShards = int(clusterHealthInterface["active_shards"].(float64))
	clusterStats.NumActivePrimaryShards = int(clusterHealthInterface["active_primary_shards"].(float64))
	clusterStats.NumInitializingShards = int(clusterHealthInterface["initializing_shards"].(float64))
	clusterStats.NumUnassignedShards = int(clusterHealthInterface["unassigned_shards"].(float64))
	clusterStats.NumRelocatingShards = int(clusterHealthInterface["relocating_shards"].(float64))
	clusterStats.StatTag = "ClusterStats"
	clusterStats.ClusterName = clusterHealthInterface["cluster_name"].(string)
	return *clusterStats
}

// Input: opensearch client and context
// Description: Fetches the cluster level info and indexes into the elasticsearch
func IndexClusterHealth(esClient *opensearch.Client, ctx context.Context) {
	var clusterHealth = ClusterMetrics{}

	//fetch the cluster stats
	clusterHealth = FetchClusterHealthMetrics(esClient, ctx)

	//Convert the cluster stats struct into Json
	clusterHealthJson, jsonErr := json.MarshalIndent(clusterHealth, "", "\t")
	if jsonErr != nil {
		log.Fatal("Error converting struct to json: ", jsonErr)
	}

	//Check and index the Json document into elasticsearch
	utils.CheckIfIndexExists(esClient, ctx)
	utils.IndexMetrics(ctx, esClient, clusterHealthJson)
}
