package arkham

import (
	"context"
)

// ChainsService provides access to chain endpoints.
// https://arkm.com/llms/get-chains.md
type ChainsService struct {
	client *Client
}

// List retrieves the list of supported blockchain chains.
// Path: GET /chains
func (s *ChainsService) List(ctx context.Context) ([]Chain, *ResponseMetadata, error) {
	var out []Chain
	meta, err := s.client.get(ctx, "/chains", nil, &out)
	return out, meta, err
}
