package arkham

import (
	"context"
	"net/url"
)

// HistoryService provides access to historical balance endpoints.
type HistoryService struct {
	client *Client
}

// Address retrieves historical balance snapshots for an address.
// Path: GET /history/address/{address}
func (s *HistoryService) Address(ctx context.Context, address string, filter *ChainsFilter) ([]HistoryPoint, *ResponseMetadata, error) {
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []HistoryPoint
	meta, err := s.client.get(ctx, "/history/address/"+pathEscape(address), q, &out)
	return out, meta, err
}

// Entity retrieves historical balance snapshots for an entity.
// Path: GET /history/entity/{entity}
func (s *HistoryService) Entity(ctx context.Context, entity string, filter *ChainsFilter) ([]HistoryPoint, *ResponseMetadata, error) {
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []HistoryPoint
	meta, err := s.client.get(ctx, "/history/entity/"+pathEscape(entity), q, &out)
	return out, meta, err
}
