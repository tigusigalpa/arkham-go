package arkham

import (
	"context"
)

// ClusterService provides access to cluster endpoints.
// https://arkm.com/llms/get-cluster-id-summary.md
type ClusterService struct {
	client *Client
}

// Summary retrieves a cluster summary.
// Path: GET /cluster/{id}/summary
func (s *ClusterService) Summary(ctx context.Context, id string) (*ClusterSummary, *ResponseMetadata, error) {
	var out ClusterSummary
	meta, err := s.client.get(ctx, "/cluster/"+pathEscape(id)+"/summary", nil, &out)
	return &out, meta, err
}
