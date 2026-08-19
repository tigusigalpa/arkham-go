package arkham

import (
	"context"
)

// AnalyticsService provides access to analytics endpoints.
// https://arkm.com/llms/get-analytics-credit-periods.md
// https://arkm.com/llms/get-analytics-endpoint-calls.md
type AnalyticsService struct {
	client *Client
}

// CreditPeriods retrieves historical credit usage by billing period.
// Only available to organization admins and individual API subscribers.
// Path: GET /analytics/credit-periods
func (s *AnalyticsService) CreditPeriods(ctx context.Context) ([]CreditUsagePeriod, *ResponseMetadata, error) {
	var out []CreditUsagePeriod
	meta, err := s.client.get(ctx, "/analytics/credit-periods", nil, &out)
	return out, meta, err
}

// EndpointCalls retrieves API endpoint call analytics.
// Path: GET /analytics/endpoint-calls
func (s *AnalyticsService) EndpointCalls(ctx context.Context) ([]EndpointCallAnalytics, *ResponseMetadata, error) {
	var out []EndpointCallAnalytics
	meta, err := s.client.get(ctx, "/analytics/endpoint-calls", nil, &out)
	return out, meta, err
}
