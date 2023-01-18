package fetchmetrics

import (
	"bytes"
	"context"
	"fmt"
	"log"
	opensearch "github.com/opensearch-project/opensearch-go"
	osapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	"scaling_manager/utils"
)
//Input: opensearch client and context
//Description: Deletes documents that older than 72 hours 
func DeleteOldDocs(esClient *opensearch.Client, ctx context.Context) {
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
		  deleteByQueryRequest,err:=osapi.DeleteByQueryRequest{
				Index: indexName,
				Body: bytes.NewReader(jsonQuery),
		  }.Do(ctx,esClient)
		 if err!=nil{
			log.Fatal("Unable to execute request: ",err)
		 }
		 fmt.Println("Document deleted: ",deleteByQueryRequest)
}