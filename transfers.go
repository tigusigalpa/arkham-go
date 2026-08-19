package arkham

import (
	"context"
	"net/url"
)

// TransfersService provides access to transfer endpoints.
type TransfersService struct {
	client *Client
}

// Transfers retrieves enriched transfers with filtering and pagination.
// Path: GET /transfers
func (s *TransfersService) Transfers(ctx context.Context, filter *TransferFilter) (*EnrichedTransfers, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out EnrichedTransfers
	meta, err := s.client.get(ctx, "/transfers", q, &out)
	return &out, meta, err
}

// Unenriched retrieves unenriched transfers with filtering and pagination.
// Path: GET /transfers/unenriched
func (s *TransfersService) Unenriched(ctx context.Context, filter *TransferFilter) (*UnenrichedTransfersResponse, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out UnenrichedTransfersResponse
	meta, err := s.client.get(ctx, "/transfers/unenriched", q, &out)
	return &out, meta, err
}

// Histogram retrieves a transfer histogram.
// Path: GET /transfers/histogram
func (s *TransfersService) Histogram(ctx context.Context, filter *TransferFilter) (*HistogramResponse, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out HistogramResponse
	meta, err := s.client.get(ctx, "/transfers/histogram", q, &out)
	return &out, meta, err
}

// SimpleHistogram retrieves a simple transfer histogram.
// Path: GET /transfers/histogram/simple
func (s *TransfersService) SimpleHistogram(ctx context.Context, filter *TransferFilter) (*SimpleHistogramResponse, *ResponseMetadata, error) {
	if err := filter.Validate(); err != nil {
		return nil, nil, err
	}
	q := url.Values{}
	filter.ApplyToValues(q)
	var out SimpleHistogramResponse
	meta, err := s.client.get(ctx, "/transfers/histogram/simple", q, &out)
	return &out, meta, err
}

// Volume retrieves transfer volume for an address.
// Path: GET /volume/address/{address}
func (s *TransfersService) Volume(ctx context.Context, address string, filter *ChainsFilter) (map[string][]VolumePoint, *ResponseMetadata, error) {
	q := url.Values{}
	filter.ApplyToValues(q)
	var out map[string][]VolumePoint
	meta, err := s.client.get(ctx, "/volume/address/"+pathEscape(address), q, &out)
	return out, meta, err
}
