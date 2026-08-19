package arkham

import (
	"context"
	"net/url"
)

// RiskService provides access to Risk Scoring endpoints (beta add-on).
// Risk Scoring requires a paid add-on. See:
// https://arkm.com/llms/guides/risk-scoring-beta.md
type RiskService struct {
	client *Client
}

// Address retrieves the compliance risk score for an address.
// Path: GET /risk/address/{address}
func (s *RiskService) Address(ctx context.Context, address string) (*RiskScoreResponse, *ResponseMetadata, error) {
	var out RiskScoreResponse
	meta, err := s.client.get(ctx, "/risk/address/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// Entity retrieves the compliance risk score for an entity.
// Path: GET /risk/entity/{entity}
func (s *RiskService) Entity(ctx context.Context, entity string) (*RiskScoreResponse, *ResponseMetadata, error) {
	var out RiskScoreResponse
	meta, err := s.client.get(ctx, "/risk/entity/"+pathEscape(entity), nil, &out)
	return &out, meta, err
}

// Paths retrieves risk paths for an address.
// Path: GET /risk/paths/address/{address}
func (s *RiskService) Paths(ctx context.Context, address string) (*RiskPathsResponse, *ResponseMetadata, error) {
	var out RiskPathsResponse
	meta, err := s.client.get(ctx, "/risk/paths/address/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// BatchAddresses retrieves risk scores for multiple addresses.
// Path: POST /risk/address/batch
func (s *RiskService) BatchAddresses(ctx context.Context, addresses []string) ([]RiskScoreResponse, *ResponseMetadata, error) {
	body := RiskBatchRequest{Addresses: addresses}
	var out []RiskScoreResponse
	meta, err := s.client.post(ctx, "/risk/address/batch", body, &out)
	return out, meta, err
}

// BatchEntities retrieves risk scores for multiple entities.
// Path: POST /risk/entity/batch
func (s *RiskService) BatchEntities(ctx context.Context, entities []string) ([]RiskScoreResponse, *ResponseMetadata, error) {
	body := RiskEntityBatchRequest{Entities: entities}
	var out []RiskScoreResponse
	meta, err := s.client.post(ctx, "/risk/entity/batch", body, &out)
	return out, meta, err
}

// RiskSources retrieves risk sources for an address.
// Path: GET /risk/sources/address/{address}
func (s *RiskService) Sources(ctx context.Context, address string) ([]RiskSource, *ResponseMetadata, error) {
	var out []RiskSource
	meta, err := s.client.get(ctx, "/risk/sources/address/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// RiskEntitySources retrieves risk sources for an entity.
// Path: GET /risk/sources/entity/{entity}
func (s *RiskService) EntitySources(ctx context.Context, entity string) ([]RiskSource, *ResponseMetadata, error) {
	var out []RiskSource
	meta, err := s.client.get(ctx, "/risk/sources/entity/"+pathEscape(entity), nil, &out)
	return out, meta, err
}

// RiskEntityPaths retrieves risk paths for an entity.
// Path: GET /risk/paths/entity/{entity}
func (s *RiskService) EntityPaths(ctx context.Context, entity string) (*RiskPathsResponse, *ResponseMetadata, error) {
	var out RiskPathsResponse
	meta, err := s.client.get(ctx, "/risk/paths/entity/"+pathEscape(entity), nil, &out)
	return &out, meta, err
}

// RiskSummary retrieves a risk summary for an address.
// Path: GET /risk/summary/address/{address}
func (s *RiskService) Summary(ctx context.Context, address string) (*RiskScoreResponse, *ResponseMetadata, error) {
	q := url.Values{}
	var out RiskScoreResponse
	meta, err := s.client.get(ctx, "/risk/summary/address/"+pathEscape(address), q, &out)
	return &out, meta, err
}
