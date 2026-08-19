package arkham

import (
	"context"
	"net/url"
)

// SwapsService provides access to swap endpoints.
type SwapsService struct {
	client *Client
}

// Swaps retrieves token swaps with filtering and pagination.
// Path: GET /swaps
func (s *SwapsService) Swaps(ctx context.Context, filter *TransferFilter) (*SwapsResponse, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out SwapsResponse
	meta, err := s.client.get(ctx, "/swaps", q, &out)
	return &out, meta, err
}
