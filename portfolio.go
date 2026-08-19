package arkham

import (
	"context"
	"net/url"
)

// PortfolioService provides access to portfolio history endpoints.
type PortfolioService struct {
	client *Client
}

// Address retrieves historical portfolio snapshots for an address.
// time is a Unix timestamp in milliseconds.
// Path: GET /portfolio/address/{address}
func (s *PortfolioService) Address(ctx context.Context, address string, time string, chains []string) (map[string]map[string]PortfolioSnapshot, *ResponseMetadata, error) {
	q := url.Values{}
	q.Set("time", time)
	if len(chains) > 0 {
		q.Set("chains", joinChains(chains))
	}
	var out map[string]map[string]PortfolioSnapshot
	meta, err := s.client.get(ctx, "/portfolio/address/"+pathEscape(address), q, &out)
	return out, meta, err
}

// Entity retrieves historical portfolio snapshots for an entity.
// time is a Unix timestamp in milliseconds.
// Path: GET /portfolio/entity/{entity}
func (s *PortfolioService) Entity(ctx context.Context, entity string, time string, chains []string) (map[string]map[string]PortfolioSnapshot, *ResponseMetadata, error) {
	q := url.Values{}
	q.Set("time", time)
	if len(chains) > 0 {
		q.Set("chains", joinChains(chains))
	}
	var out map[string]map[string]PortfolioSnapshot
	meta, err := s.client.get(ctx, "/portfolio/entity/"+pathEscape(entity), q, &out)
	return out, meta, err
}
