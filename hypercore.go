package arkham

import (
	"context"
	"net/url"
)

// HypercoreService provides access to HyperCore endpoints.
type HypercoreService struct {
	client *Client
}

// Markets lists all HyperCore spot and perpetual markets.
// Path: GET /hypercore/markets
func (s *HypercoreService) Markets(ctx context.Context) (*HypercoreMarketsResponse, *ResponseMetadata, error) {
	var out HypercoreMarketsResponse
	meta, err := s.client.get(ctx, "/hypercore/markets", nil, &out)
	return &out, meta, err
}

// Active checks if an address is active on HyperCore.
// Path: GET /hypercore/active/{address}
func (s *HypercoreService) Active(ctx context.Context, address string) (*HypercoreActiveResponse, *ResponseMetadata, error) {
	var out HypercoreActiveResponse
	meta, err := s.client.get(ctx, "/hypercore/active/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// Summary retrieves a HyperCore account summary.
// Path: GET /hypercore/summary/{address}
func (s *HypercoreService) Summary(ctx context.Context, address string) (*HypercoreSummary, *ResponseMetadata, error) {
	var out HypercoreSummary
	meta, err := s.client.get(ctx, "/hypercore/summary/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// SpotBalances retrieves HyperCore spot balances for an address.
// Path: GET /hypercore/spot/balances/{address}
func (s *HypercoreService) SpotBalances(ctx context.Context, address string) ([]HypercoreSpotBalance, *ResponseMetadata, error) {
	var out []HypercoreSpotBalance
	meta, err := s.client.get(ctx, "/hypercore/spot/balances/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// PerpPositions retrieves HyperCore perpetual positions for an address.
// Path: GET /hypercore/perp/positions/{address}
func (s *HypercoreService) PerpPositions(ctx context.Context, address string) ([]HypercorePerpPosition, *ResponseMetadata, error) {
	var out []HypercorePerpPosition
	meta, err := s.client.get(ctx, "/hypercore/perp/positions/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// PortfolioHistory retrieves HyperCore portfolio history for an address.
// Path: GET /hypercore/portfolio/history/{address}
func (s *HypercoreService) PortfolioHistory(ctx context.Context, address string) ([]HypercorePortfolioHistory, *ResponseMetadata, error) {
	var out []HypercorePortfolioHistory
	meta, err := s.client.get(ctx, "/hypercore/portfolio/history/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// Subaccounts retrieves HyperCore subaccounts for an address.
// Path: GET /hypercore/subaccounts/{address}
func (s *HypercoreService) Subaccounts(ctx context.Context, address string) ([]HypercoreSubaccount, *ResponseMetadata, error) {
	var out []HypercoreSubaccount
	meta, err := s.client.get(ctx, "/hypercore/subaccounts/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// Trades retrieves HyperCore trades for an address.
// Path: GET /hypercore/trades/{address}
func (s *HypercoreService) Trades(ctx context.Context, address string) ([]HypercoreTrade, *ResponseMetadata, error) {
	var out []HypercoreTrade
	meta, err := s.client.get(ctx, "/hypercore/trades/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// TokenTopPerps retrieves top perp positions for a token.
// Path: GET /hypercore/token/top-perps/{token}
func (s *HypercoreService) TokenTopPerps(ctx context.Context, token string) ([]HypercoreTokenPosition, *ResponseMetadata, error) {
	var out []HypercoreTokenPosition
	meta, err := s.client.get(ctx, "/hypercore/token/top-perps/"+pathEscape(token), nil, &out)
	return out, meta, err
}

// TokenTopSpots retrieves top spot balances for a token.
// Path: GET /hypercore/token/top-spots/{token}
func (s *HypercoreService) TokenTopSpots(ctx context.Context, token string) ([]HypercoreTokenPosition, *ResponseMetadata, error) {
	var out []HypercoreTokenPosition
	meta, err := s.client.get(ctx, "/hypercore/token/top-spots/"+pathEscape(token), nil, &out)
	return out, meta, err
}

// TokenHolders retrieves top token holders on HyperCore.
// Path: GET /hypercore/token/holders/{token}
func (s *HypercoreService) TokenHolders(ctx context.Context, token string) ([]TokenHolder, *ResponseMetadata, error) {
	var out []TokenHolder
	meta, err := s.client.get(ctx, "/hypercore/token/holders/"+pathEscape(token), nil, &out)
	return out, meta, err
}

// TokenTrades retrieves trades for a token on HyperCore.
// Path: GET /hypercore/token/trades/{token}
func (s *HypercoreService) TokenTrades(ctx context.Context, token string, filter *TransferFilter) ([]HypercoreTrade, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []HypercoreTrade
	meta, err := s.client.get(ctx, "/hypercore/token/trades/"+pathEscape(token), q, &out)
	return out, meta, err
}
