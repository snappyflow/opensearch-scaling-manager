package fetchmetrics

import (
	"bytes"
	"context"
	"fmt"
	"log"
	elasticsearch "github.com/opensearch-project/opensearch-go"
	esapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	utils "fetchMetrics/utils"
)
//Input: opensearch client and context
//Description: Deletes documents that older than 72 hours 
func DeleteOldDocs(esClient *elasticsearch.Client, ctx context.Context) {
	var jsonQuery = []byte(`{
		"query": {
		  "bool": {
			  "filter": {
				  "range": {
						  "Timestamp": {
								  "lte":"now-72h"
						  }
					  }
				  }
			  }
		  }
	  }`)
	  var indexName = []string{utils.IndexName}
		  deleteByQueryRequest,err:=esapi.DeleteByQueryRequest{
				Index: indexName,
				Body: bytes.NewReader(jsonQuery),
		  }.Do(ctx,esClient)
		 if err!=nil{
			log.Fatal("Unable to execute request: ",err)
		 }
		 fmt.Println("Document deleted: ",deleteByQueryRequest)
}