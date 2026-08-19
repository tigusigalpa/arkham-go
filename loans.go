package arkham

import (
	"context"
)

// LoansService provides access to loan/borrow position endpoints.
type LoansService struct {
	client *Client
}

// Address retrieves loan positions for an address.
// Path: GET /loans/address/{address}
func (s *LoansService) Address(ctx context.Context, address string) (*LoansResponse, *ResponseMetadata, error) {
	var out LoansResponse
	meta, err := s.client.get(ctx, "/loans/address/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// Entity retrieves loan positions for an entity.
// Path: GET /loans/entity/{entity}
func (s *LoansService) Entity(ctx context.Context, entity string) (*LoansResponse, *ResponseMetadata, error) {
	var out LoansResponse
	meta, err := s.client.get(ctx, "/loans/entity/"+pathEscape(entity), nil, &out)
	return &out, meta, err
}
