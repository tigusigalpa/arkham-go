package arkham

import (
	"context"
	"net/url"
)

// FlowService provides access to historical USD flow endpoints.
type FlowService struct {
	client *Client
}

// Address retrieves historical USD flow for an address.
// Path: GET /flow/address/{address}
func (s *FlowService) Address(ctx context.Context, address string, filter *TransferFilter) ([]FlowPoint, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []FlowPoint
	meta, err := s.client.get(ctx, "/flow/address/"+pathEscape(address), q, &out)
	return out, meta, err
}

// Entity retrieves historical USD flow for an entity.
// Path: GET /flow/entity/{entity}
func (s *FlowService) Entity(ctx context.Context, entity string, filter *TransferFilter) ([]FlowPoint, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []FlowPoint
	meta, err := s.client.get(ctx, "/flow/entity/"+pathEscape(entity), q, &out)
	return out, meta, err
}
