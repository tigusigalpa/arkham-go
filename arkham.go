package arkham

import (
	"context"
)

// ArkhamService provides access to ARKM token endpoints.
// https://arkm.com/llms/get-arkm-circulating.md
type ArkhamService struct {
	client *Client
}

// Circulating retrieves the current ARKM circulating supply.
// Path: GET /arkm/circulating
func (s *ArkhamService) Circulating(ctx context.Context) (*ARKMCirculating, *ResponseMetadata, error) {
	var out ARKMCirculating
	meta, err := s.client.get(ctx, "/arkm/circulating", nil, &out)
	return &out, meta, err
}
