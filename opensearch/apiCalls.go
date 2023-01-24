package osutils

import (
	"bytes"
	"context"
	"strings"

	osapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

func IndexMetrics(ctx context.Context, jsonDoc []byte) (*osapi.Response, error) {
	//Create an index request and pass the document to be indexed along with the index name
	return osapi.IndexRequest{
		Index:        IndexName,
		DocumentType: "_doc",
		Body:         bytes.NewReader(jsonDoc),
	}.Do(ctx, osClient)
}

func GetClusterStats(ctx context.Context) (*osapi.Response, error) {
	return osapi.ClusterStatsRequest{}.Do(ctx, osClient)
}

func GetClusterHealth(ctx context.Context) (*osapi.Response, error) {
	return osapi.ClusterHealthRequest{}.Do(ctx, osClient)
}

func GetClusterState(ctx context.Context) (*osapi.Response, error) {
	return osapi.ClusterStateRequest{}.Do(ctx, osClient)
}

func GetNodeStats(nodes []string, metrics []string, ctx context.Context) (*osapi.Response, error) {
	return osapi.NodesStatsRequest{
		Pretty: true,
		NodeID: nodes,
		Metric: metrics,
	}.Do(ctx, osClient)
}

func SearchQuery(jsonQuery []byte, ctx context.Context) (*osapi.Response, error) {
	return osapi.SearchRequest{
		Index: []string{IndexName},
		Body:  bytes.NewReader(jsonQuery),
	}.Do(ctx, osClient)
}

func SearchDoc(docId string, ctx context.Context) (*osapi.Response, error) {
	return osapi.GetRequest{
		Index:      IndexName,
		DocumentID: docId,
	}.Do(ctx, osClient)
}

func UpdateDoc(docId string, content string, ctx context.Context) (*osapi.Response, error) {
	return osapi.IndexRequest{
		Index:      IndexName,
		DocumentID: docId,
		Body:       strings.NewReader(content),
	}.Do(ctx, osClient)
}

func DeleteWithQuery(jsonQuery []byte, ctx context.Context) (*osapi.Response, error) {
	return osapi.DeleteByQueryRequest{
		Index: []string{IndexName},
		Body:  bytes.NewReader(jsonQuery),
	}.Do(ctx, osClient)
}
