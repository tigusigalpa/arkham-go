package arkham

import (
	"context"
)

// TagService provides access to tag endpoints.
type TagService struct {
	client *Client
}

// Summary retrieves a tag summary.
// Path: GET /tag/{id}/summary
func (s *TagService) Summary(ctx context.Context, id string) (*TagSummary, *ResponseMetadata, error) {
	var out TagSummary
	meta, err := s.client.get(ctx, "/tag/"+pathEscape(id)+"/summary", nil, &out)
	return &out, meta, err
}

// Updates retrieves recent tag updates.
// Path: GET /tag/updates
func (s *TagService) Updates(ctx context.Context) ([]TagUpdate, *ResponseMetadata, error) {
	var out []TagUpdate
	meta, err := s.client.get(ctx, "/tag/updates", nil, &out)
	return out, meta, err
}

// AddressUpdates retrieves recent address-tag association updates.
// Path: GET /tag/address-updates
func (s *TagService) AddressUpdates(ctx context.Context) ([]AddressTagUpdate, *ResponseMetadata, error) {
	var out []AddressTagUpdate
	meta, err := s.client.get(ctx, "/tag/address-updates", nil, &out)
	return out, meta, err
}
