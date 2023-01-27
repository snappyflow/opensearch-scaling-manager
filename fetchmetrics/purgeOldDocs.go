package fetchmetrics

import (
	"context"
	osutils "scaling_manager/opensearchUtils"
)

// Input: opensearch client and context
// Description: Deletes documents that older than 72 hours
func DeleteOldDocs(ctx context.Context) {
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
	deleteResp, err := osutils.DeleteWithQuery(jsonQuery, ctx)
	if err != nil {
		log.Panic.Println("Unable to execute request: ", err)
		panic(err)
	}
	log.Info.Println("Document deleted: ", deleteResp)
}
