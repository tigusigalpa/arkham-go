package arkham

import (
	"context"
	"net/url"
)

// MarketDataService provides access to market data endpoints.
type MarketDataService struct {
	client *Client
}

// AltcoinIndex retrieves the altcoin index.
// Path: GET /marketdata/altcoin-index
func (s *MarketDataService) AltcoinIndex(ctx context.Context) ([]AltcoinIndex, *ResponseMetadata, error) {
	var out []AltcoinIndex
	meta, err := s.client.get(ctx, "/marketdata/altcoin-index", nil, &out)
	return out, meta, err
}

// TokenTrending retrieves trending tokens.
// Path: GET /marketdata/tokens/trending
func (s *MarketDataService) TokenTrending(ctx context.Context) ([]TrendingToken, *ResponseMetadata, error) {
	var out []TrendingToken
	meta, err := s.client.get(ctx, "/marketdata/tokens/trending", nil, &out)
	return out, meta, err
}

// TokenTopFlow retrieves top flow data for a token.
// Path: GET /marketdata/tokens/top-flow/{token}
func (s *MarketDataService) TokenTopFlow(ctx context.Context, token string, filter *TransferFilter) ([]TokenTopFlow, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []TokenTopFlow
	meta, err := s.client.get(ctx, "/marketdata/tokens/top-flow/"+pathEscape(token), q, &out)
	return out, meta, err
}
