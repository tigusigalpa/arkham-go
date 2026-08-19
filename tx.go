package arkham

import (
	"context"
)

// TxService provides access to transaction endpoints.
type TxService struct {
	client *Client
}

// Transaction retrieves a transaction by hash.
// Path: GET /tx/{hash}
func (s *TxService) Transaction(ctx context.Context, hash string) (*Transaction, *ResponseMetadata, error) {
	var out Transaction
	meta, err := s.client.get(ctx, "/tx/"+pathEscape(hash), nil, &out)
	return &out, meta, err
}
