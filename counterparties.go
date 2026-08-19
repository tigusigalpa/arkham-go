package arkham

import (
	"context"
)

// CounterpartiesService provides access to counterparty endpoints.
// These are heavy endpoints (1 req/s rate limit).
type CounterpartiesService struct {
	client *Client
}

// Address retrieves top counterparties for an address.
// Path: GET /counterparties/address/{address}
func (s *CounterpartiesService) Address(ctx context.Context, address string) (*CounterpartiesResponse, *ResponseMetadata, error) {
	var out CounterpartiesResponse
	meta, err := s.client.get(ctx, "/counterparties/address/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// Entity retrieves top counterparties for an entity.
// Path: GET /counterparties/entity/{entity}
func (s *CounterpartiesService) Entity(ctx context.Context, entity string) (*CounterpartiesResponse, *ResponseMetadata, error) {
	var out CounterpartiesResponse
	meta, err := s.client.get(ctx, "/counterparties/entity/"+pathEscape(entity), nil, &out)
	return &out, meta, err
}
