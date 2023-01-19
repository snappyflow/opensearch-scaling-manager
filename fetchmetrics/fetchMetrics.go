package fetchmetrics

import (
	"context"
	"fmt"

	"scaling_manager/utils"

	opensearch "github.com/opensearch-project/opensearch-go"
)

var ctx = context.Background()

// Descriptions: Fetch metrics will create an opensearch client and fetches the node and cluster level details and
// indexes them into elasticesearch. It also deletes documents that are older than 72 hours
func FetchMetrics(osClient *opensearch.Client) {
	//check if current node is the master node and update the cluster stats if it is master
	if utils.CheckIfMaster(osClient, ctx) {
		fmt.Println("It's the master node")
		IndexClusterHealth(osClient, ctx)
	}
	//Index the the node stats
	//IndexNodeStats(osClient, ctx)
	//Purge documents from elasticsearch index that are older than 72 hours
	DeleteOldDocs(osClient, ctx)
}