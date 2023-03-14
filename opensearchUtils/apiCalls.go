package osutils

import (
	"bytes"
	"context"
	"strings"
	"time"

	osapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//	jsonDoc (byte): The request body in form of bytes
//
// Description:
//
//	Calls the osapi IndexRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func IndexMetrics(ctx context.Context, jsonDoc []byte) (*osapi.Response, error) {
	return osapi.IndexRequest{
		Index:        IndexName,
		DocumentType: "_doc",
		Body:         bytes.NewReader(jsonDoc),
		Refresh:      "wait_for",
	}.Do(ctx, osClient)
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterStatsRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetClusterStats(ctx context.Context) (*osapi.Response, error) {
	return osapi.ClusterStatsRequest{}.Do(ctx, osClient)
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterHealthRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetClusterHealth(ctx context.Context, WaitForShards *bool) (*osapi.Response, error) {
	return osapi.ClusterHealthRequest{
		WaitForNoInitializingShards: WaitForShards,
		WaitForNoRelocatingShards:   WaitForShards,
		Timeout:                     time.Duration(90 * time.Second),
	}.Do(ctx, osClient)
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterStateRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetClusterState(ctx context.Context) (*osapi.Response, error) {
	return osapi.ClusterStateRequest{}.Do(ctx, osClient)
}

// Input:
//
//	nodes ([]string): The list of nodes for which the stats needs to be fetched
//	metrics ([]string): The list of metrics that needs to be fetched for the specified node/s
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi NodeStatsRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetNodeStats(ctx context.Context, nodes []string, metrics []string) (*osapi.Response, error) {
	return osapi.NodesStatsRequest{
		Pretty: true,
		NodeID: nodes,
		Metric: metrics,
	}.Do(ctx, osClient)
}

// Input:
//
//	nodes ([]string): List of nodes for which the allocation needs to be fetched
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi CatAllocationRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func CatAllocation(ctx context.Context, nodes []string) (*osapi.Response, error) {
	return osapi.CatAllocationRequest{
		Pretty: true,
		NodeID: nodes,
	}.Do(ctx, osClient)
}

// Input:
//
//	jsonQuery ([]byte): The json query in bytes that needs to be queried from the index
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi SearchRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func SearchQuery(ctx context.Context, jsonQuery []byte) (*osapi.Response, error) {
	return osapi.SearchRequest{
		Index: []string{IndexName},
		Body:  bytes.NewReader(jsonQuery),
	}.Do(ctx, osClient)
}

// Input:
//
//	docId (string): The _id of the document which needs to be searched
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi GetRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func SearchDoc(ctx context.Context, docId string) (*osapi.Response, error) {
	return osapi.GetRequest{
		Index:      IndexName,
		DocumentID: docId,
	}.Do(ctx, osClient)
}

// Input:
//
//	docId (string): The _id of the document which needs to be updated
//	content (string): The body of the request that needs to be updated in the document
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi IndexRequest along with document ID to update and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func UpdateDoc(ctx context.Context, docId string, content string) (*osapi.Response, error) {
	return osapi.IndexRequest{
		Index:      IndexName,
		DocumentID: docId,
		Body:       strings.NewReader(content),
		Refresh:    "wait_for",
	}.Do(ctx, osClient)
}

// Input:
//
//	jsonQuery ([]byte): Query by which the deletion of documents is carried
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi DeleteByQueryRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func DeleteWithQuery(ctx context.Context, jsonQuery []byte) (*osapi.Response, error) {
	return osapi.DeleteByQueryRequest{
		Index: []string{IndexName},
		Body:  bytes.NewReader(jsonQuery),
	}.Do(ctx, osClient)
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterRerouteRequest with true and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func RerouteRetryFailed(ctx context.Context) (*osapi.Response, error) {
	retry := true
	return osapi.ClusterRerouteRequest{
		RetryFailed: &retry,
	}.Do(ctx, osClient)
}
