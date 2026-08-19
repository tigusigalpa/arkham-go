package arkham

import (
	"context"
	"net/url"
	"strconv"
)

// PolymarketService provides access to Polymarket endpoints.
type PolymarketService struct {
	client *Client
}

// Events retrieves a filtered, paginated list of Polymarket events.
// Path: GET /polymarket/events
func (s *PolymarketService) Events(ctx context.Context, opts *PolymarketEventOptions) (*PolymarketEventsResponse, *ResponseMetadata, error) {
	q := url.Values{}
	if opts != nil {
		opts.ApplyToValues(q)
	}
	var out PolymarketEventsResponse
	meta, err := s.client.get(ctx, "/polymarket/events", q, &out)
	return &out, meta, err
}

// PolymarketEventOptions holds optional parameters for the events endpoint.
type PolymarketEventOptions struct {
	Tag        string
	ExcludeTag string
	Active     *bool
	Search     string
	GroupGames *bool
	SortBy     string
	Order      string
	Limit      int
	Offset     int
}

// ApplyToValues adds event parameters to the given url.Values.
func (o *PolymarketEventOptions) ApplyToValues(q url.Values) {
	if o == nil {
		return
	}
	if o.Tag != "" {
		q.Set("tag", o.Tag)
	}
	if o.ExcludeTag != "" {
		q.Set("excludeTag", o.ExcludeTag)
	}
	if o.Active != nil {
		q.Set("active", strconv.FormatBool(*o.Active))
	}
	if o.Search != "" {
		q.Set("search", o.Search)
	}
	if o.GroupGames != nil {
		q.Set("groupGames", strconv.FormatBool(*o.GroupGames))
	}
	if o.SortBy != "" {
		q.Set("sortBy", o.SortBy)
	}
	if o.Order != "" {
		q.Set("order", o.Order)
	}
	if o.Limit > 0 {
		q.Set("limit", strconv.Itoa(o.Limit))
	}
	if o.Offset > 0 {
		q.Set("offset", strconv.Itoa(o.Offset))
	}
}

// Activity retrieves Polymarket activity for an address.
// Path: GET /polymarket/activity/{address}
func (s *PolymarketService) Activity(ctx context.Context, address string) ([]PolymarketActivity, *ResponseMetadata, error) {
	var out []PolymarketActivity
	meta, err := s.client.get(ctx, "/polymarket/activity/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// Positions retrieves Polymarket positions for an address.
// Path: GET /polymarket/positions/{address}
func (s *PolymarketService) Positions(ctx context.Context, address string) ([]PolymarketPosition, *ResponseMetadata, error) {
	var out []PolymarketPosition
	meta, err := s.client.get(ctx, "/polymarket/positions/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// OrderBook retrieves the order book for a condition.
// Path: GET /polymarket/orderbook/{conditionId}
func (s *PolymarketService) OrderBook(ctx context.Context, conditionID string) (*PolymarketOrderBook, *ResponseMetadata, error) {
	var out PolymarketOrderBook
	meta, err := s.client.get(ctx, "/polymarket/orderbook/"+pathEscape(conditionID), nil, &out)
	return &out, meta, err
}

// PriceHistory retrieves price history for a condition.
// Path: GET /polymarket/price-history/{conditionId}
func (s *PolymarketService) PriceHistory(ctx context.Context, conditionID string) ([]PolymarketPriceHistory, *ResponseMetadata, error) {
	var out []PolymarketPriceHistory
	meta, err := s.client.get(ctx, "/polymarket/price-history/"+pathEscape(conditionID), nil, &out)
	return out, meta, err
}

// PnLChart retrieves PnL chart data for an address.
// Path: GET /polymarket/pnl-chart/{address}
func (s *PolymarketService) PnLChart(ctx context.Context, address string) ([]PolymarketPnLChart, *ResponseMetadata, error) {
	var out []PolymarketPnLChart
	meta, err := s.client.get(ctx, "/polymarket/pnl-chart/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// WalletSummary retrieves a wallet summary for an address.
// Path: GET /polymarket/wallet/summary/{address}
func (s *PolymarketService) WalletSummary(ctx context.Context, address string) (*PolymarketWalletSummary, *ResponseMetadata, error) {
	var out PolymarketWalletSummary
	meta, err := s.client.get(ctx, "/polymarket/wallet/summary/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// WalletStats retrieves wallet trading stats for an address.
// Path: GET /polymarket/wallet/stats/{address}
func (s *PolymarketService) WalletStats(ctx context.Context, address string) (*PolymarketWalletStats, *ResponseMetadata, error) {
	var out PolymarketWalletStats
	meta, err := s.client.get(ctx, "/polymarket/wallet/stats/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// WalletTags retrieves wallet tags for an address.
// Path: GET /polymarket/wallet/tags/{address}
func (s *PolymarketService) WalletTags(ctx context.Context, address string) (*PolymarketWalletTags, *ResponseMetadata, error) {
	var out PolymarketWalletTags
	meta, err := s.client.get(ctx, "/polymarket/wallet/tags/"+pathEscape(address), nil, &out)
	return &out, meta, err
}

// Leaderboard retrieves the Polymarket leaderboard.
// Path: GET /polymarket/leaderboard
func (s *PolymarketService) Leaderboard(ctx context.Context) ([]PolymarketLeaderboardEntry, *ResponseMetadata, error) {
	var out []PolymarketLeaderboardEntry
	meta, err := s.client.get(ctx, "/polymarket/leaderboard", nil, &out)
	return out, meta, err
}

// Stats retrieves Polymarket platform stats.
// Path: GET /polymarket/stats
func (s *PolymarketService) Stats(ctx context.Context) (*PolymarketStats, *ResponseMetadata, error) {
	var out PolymarketStats
	meta, err := s.client.get(ctx, "/polymarket/stats", nil, &out)
	return &out, meta, err
}

// TopEvents retrieves top Polymarket events.
// Path: GET /polymarket/top-events
func (s *PolymarketService) TopEvents(ctx context.Context) ([]PolymarketTopEvent, *ResponseMetadata, error) {
	var out []PolymarketTopEvent
	meta, err := s.client.get(ctx, "/polymarket/top-events", nil, &out)
	return out, meta, err
}

// TopEventBreakdown retrieves a breakdown of a top event.
// Path: GET /polymarket/top-events/{eventId}/breakdown
func (s *PolymarketService) TopEventBreakdown(ctx context.Context, eventID string) ([]PolymarketTopEventBreakdown, *ResponseMetadata, error) {
	var out []PolymarketTopEventBreakdown
	meta, err := s.client.get(ctx, "/polymarket/top-events/"+pathEscape(eventID)+"/breakdown", nil, &out)
	return out, meta, err
}

// TopHolders retrieves top holders for a condition.
// Path: GET /polymarket/top-holders/{conditionId}
func (s *PolymarketService) TopHolders(ctx context.Context, conditionID string) ([]PolymarketTopHolder, *ResponseMetadata, error) {
	var out []PolymarketTopHolder
	meta, err := s.client.get(ctx, "/polymarket/top-holders/"+pathEscape(conditionID), nil, &out)
	return out, meta, err
}

// WalletEventHistory retrieves wallet event history for an address.
// Path: GET /polymarket/wallet/event-history/{address}
func (s *PolymarketService) WalletEventHistory(ctx context.Context, address string) ([]PolymarketWalletEventHistory, *ResponseMetadata, error) {
	var out []PolymarketWalletEventHistory
	meta, err := s.client.get(ctx, "/polymarket/wallet/event-history/"+pathEscape(address), nil, &out)
	return out, meta, err
}

// PredictionHistory retrieves prediction history for a condition.
// Path: GET /polymarket/prediction-history/{conditionId}
func (s *PolymarketService) PredictionHistory(ctx context.Context, conditionID string) ([]PolymarketPredictionHistory, *ResponseMetadata, error) {
	var out []PolymarketPredictionHistory
	meta, err := s.client.get(ctx, "/polymarket/prediction-history/"+pathEscape(conditionID), nil, &out)
	return out, meta, err
}
