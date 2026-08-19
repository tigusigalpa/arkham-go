package arkham

import (
	"context"
	"net/url"
)

// TokenService provides access to token endpoints.
type TokenService struct {
	client *Client
}

// Market retrieves current market data for a token.
// Path: GET /token/market/{id}
func (s *TokenService) Market(ctx context.Context, id string) (*TokenMarketData, *ResponseMetadata, error) {
	var out TokenMarketData
	meta, err := s.client.get(ctx, "/token/market/"+pathEscape(id), nil, &out)
	return &out, meta, err
}

// PriceHistory retrieves price history for a token.
// Path: GET /token/price-history/{id}
func (s *TokenService) PriceHistory(ctx context.Context, id string, filter *TransferFilter) ([]TokenPriceHistory, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []TokenPriceHistory
	meta, err := s.client.get(ctx, "/token/price-history/"+pathEscape(id), q, &out)
	return out, meta, err
}

// PriceChange retrieves price change for a token.
// Path: GET /token/price-change/{id}
func (s *TokenService) PriceChange(ctx context.Context, id string, timeGte string, timeLte string) (*TokenPriceChange, *ResponseMetadata, error) {
	q := url.Values{}
	if timeGte != "" {
		q.Set("timeGte", timeGte)
	}
	if timeLte != "" {
		q.Set("timeLte", timeLte)
	}
	var out TokenPriceChange
	meta, err := s.client.get(ctx, "/token/price-change/"+pathEscape(id), q, &out)
	return &out, meta, err
}

// Holders retrieves top token holders.
// Path: GET /token/holders/{id}
func (s *TokenService) Holders(ctx context.Context, id string, filter *TransferFilter) ([]TokenHolder, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []TokenHolder
	meta, err := s.client.get(ctx, "/token/holders/"+pathEscape(id), q, &out)
	return out, meta, err
}

// Volume retrieves token volume data.
// Path: GET /token/volume/{id}
func (s *TokenService) Volume(ctx context.Context, id string, filter *TransferFilter) ([]TokenVolume, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []TokenVolume
	meta, err := s.client.get(ctx, "/token/volume/"+pathEscape(id), q, &out)
	return out, meta, err
}

// Addresses retrieves chain addresses for a token.
// Path: GET /token/addresses/{id}
func (s *TokenService) Addresses(ctx context.Context, id string) ([]TokenAddress, *ResponseMetadata, error) {
	var out []TokenAddress
	meta, err := s.client.get(ctx, "/token/addresses/"+pathEscape(id), nil, &out)
	return out, meta, err
}

// Info retrieves basic token information.
// Path: GET /token/{id}
func (s *TokenService) Info(ctx context.Context, id string) (*TokenInfo, *ResponseMetadata, error) {
	var out TokenInfo
	meta, err := s.client.get(ctx, "/token/"+pathEscape(id), nil, &out)
	return &out, meta, err
}
