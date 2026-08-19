package arkham

import (
	"context"
	"net/url"
)

// NetworksService provides access to network status and history endpoints.
type NetworksService struct {
	client *Client
}

// Status retrieves current status for all blockchain networks.
// Path: GET /networks/status
func (s *NetworksService) Status(ctx context.Context) (map[string]NetworkStatus, *ResponseMetadata, error) {
	var out map[string]NetworkStatus
	meta, err := s.client.get(ctx, "/networks/status", nil, &out)
	return out, meta, err
}

// History retrieves historical data for a specific network.
// Path: GET /networks/{chain}/history
func (s *NetworksService) History(ctx context.Context, chain string, filter *TransferFilter) ([]NetworkHistoryPoint, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out []NetworkHistoryPoint
	meta, err := s.client.get(ctx, "/networks/"+pathEscape(chain)+"/history", q, &out)
	return out, meta, err
}
