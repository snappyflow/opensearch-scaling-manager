package fetchmetrics

import (
	"context"
	"fmt"
	"time"

	"scaling_manager/logger"
	utils "scaling_manager/utilities"
)

var ctx = context.Background()

var log = new(logger.LOG)

// Initializing logger module
func init() {
	log.Init("logger")
	log.Info.Println("FetchMetrics module initiated")
}

// Descriptions: Fetch metrics will create an opensearch client and fetches the node and cluster level details and
// indexes them into elasticesearch. It also deletes documents that are older than 72 hours
func FetchMetrics(pollingInterval int) {
	ticker := time.Tick(time.Duration(pollingInterval) * time.Second)
	for range ticker {
		//check if current node is the master node and update the cluster stats if it is master
		if utils.CheckIfMaster(ctx, "") {
			fmt.Println("It's the master node")
			IndexClusterHealth(ctx)
		}
		//Index the the node stats
		IndexNodeStats(ctx)
		//Purge documents from elasticsearch index that are older than 72 hours
		// DeleteOldDocs(ctx)
	}
}
