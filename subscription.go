package arkham

import (
	"context"
)

// SubscriptionService provides access to subscription endpoints.
type SubscriptionService struct {
	client *Client
}

// IntelUsage retrieves intel data usage for the current billing period.
// Path: GET /subscription/intel-usage
func (s *SubscriptionService) IntelUsage(ctx context.Context) (*IntelUsage, *ResponseMetadata, error) {
	var out IntelUsage
	meta, err := s.client.get(ctx, "/subscription/intel-usage", nil, &out)
	return &out, meta, err
}
