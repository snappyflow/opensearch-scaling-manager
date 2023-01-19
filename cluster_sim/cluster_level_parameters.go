// This package provide the data structure needed to get the metrics.
// There are two kind of metrics:
//
//	Cluster metrics: This data structure will provide cluster level metrics.
//	Node metrics: This data structure will provide node level metrics.
//
// The cluster metrics contains list of the node metrics collected over all the nodes present in a cluster.
// The package contains a struct called MetricStatsCluster which will calculate the statistics over a period of time.
// The package contains a struct called MetricViolatedCountCluster which will calculate the violated count over a period of time.
// The structs be used by recommendation module.
package cluster_sim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"scaling_manager/cluster"
	"scaling_manager/logger"
	"time"
)

var log logger.LOG

// Input:
//
// Description:
//
//	Initialize the logger module.
//
// Return:
func init() {
	log.Init("logger")
	log.Info.Println("Main module initialized")
}

// Input:
//
//		metricName: The Name of the metric for which the Cluster Average will be calculated(string).
//		decisionPeriod: The evaluation period for which the Average will be calculated.
//
// Description:
//
//		GetClusterAvg will use the opensearch query to find out the stats aggregation.
//		While getting stats aggregation it will pass the metricName and decisionPeriod as an input.
//		It will populate MetricStatsCluster struct and return it.
//
// Return:
//		Return populated MetricStatsCluster struct.

func GetClusterAvg(metricName string, decisionPeriod int) (cluster.MetricStats, []byte) {
	var metricStats cluster.MetricStats
	url := fmt.Sprintf("http://localhost:5000/stats/avg/%s/%d", metricName, decisionPeriod)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)

	if err != nil {
		log.Panic.Println(err)
		panic(err)
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			response, _ := ioutil.ReadAll(resp.Body)
			return metricStats, response
		}
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&metricStats)
	if err != nil {
		log.Panic.Println(err)
		panic(err)
	}
	log.Debug.Println(metricStats)

	return metricStats, nil
}

// Input:
//
//		metricName: The Name of the metric for which the Cluster Average will be calculated(string).
//		decisionPeriod: The evaluation period for which the Average will be calculated.(int)
//		limit: The limit for the particular metric for which the count is calculated.(float32)
//
// Description:
//
//		GetClusterCount will use the opensearch query to find out the stats aggregation.
//		While getting stats aggregation it will pass the metricName, decisionPeriod and limit as an input.
//		It will populate MetricViolatedCountCluster struct and return it.
//
// Return:
//		Return populated MetricViolatedCountCluster struct.

func GetClusterCount(metricName string, decisonPeriod int, limit float32) (cluster.MetricViolatedCount, []byte) {
	var metricViolatedCount cluster.MetricViolatedCount
	url := fmt.Sprintf("http://localhost:5000/stats/violated/%s/%d/%f", metricName, decisonPeriod, limit)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)

	if err != nil {
		log.Panic.Println(err)
		panic(err)
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			response, _ := ioutil.ReadAll(resp.Body)
			return metricViolatedCount, response
		}
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&metricViolatedCount)

	if err != nil {
		log.Panic.Println(err)
		panic(err)
	}
	log.Debug.Println(metricViolatedCount)
	return metricViolatedCount, nil
}

// Input:
// Description:
//
//		GetClusterCurrent will fetch the node level and cluster level metrics and fill in
//		ClusterDynamic, clusterStatic and Node struct using the given config file.
//		It will return the current cluster status.
//
// Return:
//		Return populated ClusterDynamic struct.

func GetClusterCurrent() cluster.ClusterDynamic {
	var clusterStats cluster.ClusterDynamic

	url := fmt.Sprintf("http://localhost:5000/stats/current")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Panic.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&clusterStats)
	if err != nil {
		log.Panic.Println(err)
		panic(err)
	}
	log.Debug.Println(clusterStats)
	return clusterStats
}
