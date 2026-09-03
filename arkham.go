package arkham

import (
	"context"
)

// Service provides access to ARKM token endpoints.
// https://arkm.com/llms/get-arkm-circulating.md
type Service struct {
	client *Client
}

// Circulating retrieves the current ARKM circulating supply.
// Path: GET /arkm/circulating
func (s *Service) Circulating(ctx context.Context) (*ARKMCirculating, *ResponseMetadata, error) {
	var out ARKMCirculating
	meta, err := s.client.get(ctx, "/arkm/circulating", nil, &out)
	return &out, meta, err
}
