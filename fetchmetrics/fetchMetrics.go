package fetchmetrics

import (
	"context"
	"fmt"
	"log"

	elasticsearch "github.com/opensearch-project/opensearch-go"
	utils "fetchMetrics/utils"
)

//Descriptions: Fetch metrics will create an opensearch client and fetches the node and cluster level details and 
//indexes them into elasticesearch. It also deletes documents that are older than 72 hours
func FetchMetrics() {
	ctx := context.Background()

	//create a configuration that is to be passed while creating the client 
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	//create the client using the configuration
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	}
	fmt.Println("Response: ", res)

	//check if current node is the master node and update the cluster stats if it is master 
	if utils.CheckIfMaster(esClient, ctx) {  
		fmt.Println("It's the master node")
		IndexClusterHealth(esClient, ctx)
	}
	//Index the the node stats
	IndexNodeStats(esClient, ctx)
	//Purge documents from elasticsearch index that are older than 72 hours
	DeleteOldDocs(esClient,ctx) 
}