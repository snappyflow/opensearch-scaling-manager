package fetchmetrics

import (
	"context"
	osutils "scaling_manager/opensearchUtils"
)

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Deletes documents that older than 72 hours
//
// Return:
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
	deleteResp, err := osutils.DeleteWithQuery(ctx, jsonQuery)
	if err != nil {
		log.Panic.Println("Unable to execute request: ", err)
		panic(err)
	}
	defer deleteResp.Body.Close()
	log.Info.Println("Document deleted: ", deleteResp)
}
