package arkham

import (
	"net/url"
	"strconv"
	"time"
)

// SortKey represents a sort field for list endpoints.
type SortKey string

const (
	SortKeyTime  SortKey = "time"
	SortKeyValue SortKey = "value"
	SortKeyUSD   SortKey = "usd"
)

// SortDir represents a sort direction.
type SortDir string

const (
	SortDirDesc SortDir = "desc"
	SortDirAsc  SortDir = "asc"
)

// FlowDirection represents the flow direction filter.
type FlowDirection string

const (
	FlowIn   FlowDirection = "in"
	FlowOut  FlowDirection = "out"
	FlowSelf FlowDirection = "self"
	FlowAll  FlowDirection = "all"
)

// TimeRange holds mutually exclusive time filter options. TimeLast
// cannot be used together with TimeGte or TimeLte.
type TimeRange struct {
	// TimeLast is a duration string like "24h", "7d", "30d", "1M", "1y".
	TimeLast string
	// TimeGte is an absolute lower bound (Unix ms or ISO-8601).
	TimeGte string
	// TimeLte is an absolute upper bound (Unix ms or ISO-8601).
	TimeLte string
}

// Validate returns an error if the time range is invalid.
func (t *TimeRange) Validate() error {
	if t == nil {
		return nil
	}
	if t.TimeLast != "" && (t.TimeGte != "" || t.TimeLte != "") {
		return ErrInvalidTimeRange
	}
	return nil
}

// ApplyToValues adds time parameters to the given url.Values.
func (t *TimeRange) ApplyToValues(q url.Values) {
	if t == nil {
		return
	}
	if t.TimeLast != "" {
		q.Set("timeLast", t.TimeLast)
	}
	if t.TimeGte != "" {
		q.Set("timeGte", t.TimeGte)
	}
	if t.TimeLte != "" {
		q.Set("timeLte", t.TimeLte)
	}
}

// TransferFilter holds the common filter parameters for transfer,
// swap, and WebSocket stream endpoints.
type TransferFilter struct {
	// Base filters from or to any of the listed entities or addresses.
	Base []string
	// Chains filters by blockchain.
	Chains []string
	// Flow filters by transfer direction.
	Flow FlowDirection
	// From filters by sender addresses, entities, or deposit services.
	From []string
	// To filters by receiver addresses, entities, or deposit services.
	To []string
	// Counterparties filters to only base <-> counterparty transfers.
	Counterparties []string
	// Tokens filters by token addresses or IDs.
	Tokens []string
	// TimeRange holds the time filter.
	TimeRange *TimeRange
	// ValueGte filters above a minimum token value.
	ValueGte string
	// ValueLte filters below a maximum token value.
	ValueLte string
	// UsdGte filters above a minimum USD value.
	UsdGte string
	// UsdLte filters below a maximum USD value.
	UsdLte string
	// SortKey is the sort field.
	SortKey SortKey
	// SortDir is the sort direction.
	SortDir SortDir
	// Limit is the max results per page.
	Limit int
	// Offset is the pagination offset.
	Offset int
}

// ApplyToValues adds filter parameters to the given url.Values.
func (f *TransferFilter) ApplyToValues(q url.Values) {
	if f == nil {
		return
	}
	if len(f.Base) > 0 {
		q.Set("base", joinChains(f.Base))
	}
	if len(f.Chains) > 0 {
		q.Set("chains", joinChains(f.Chains))
	}
	if f.Flow != "" {
		q.Set("flow", string(f.Flow))
	}
	if len(f.From) > 0 {
		q.Set("from", joinChains(f.From))
	}
	if len(f.To) > 0 {
		q.Set("to", joinChains(f.To))
	}
	if len(f.Counterparties) > 0 {
		q.Set("counterparties", joinChains(f.Counterparties))
	}
	if len(f.Tokens) > 0 {
		q.Set("tokens", joinChains(f.Tokens))
	}
	if f.TimeRange != nil {
		f.TimeRange.ApplyToValues(q)
	}
	if f.ValueGte != "" {
		q.Set("valueGte", f.ValueGte)
	}
	if f.ValueLte != "" {
		q.Set("valueLte", f.ValueLte)
	}
	if f.UsdGte != "" {
		q.Set("usdGte", f.UsdGte)
	}
	if f.UsdLte != "" {
		q.Set("usdLte", f.UsdLte)
	}
	if f.SortKey != "" {
		q.Set("sortKey", string(f.SortKey))
	}
	if f.SortDir != "" {
		q.Set("sortDir", string(f.SortDir))
	}
	if f.Limit > 0 {
		q.Set("limit", strconv.Itoa(f.Limit))
	}
	if f.Offset > 0 {
		q.Set("offset", strconv.Itoa(f.Offset))
	}
}

// Validate checks the transfer filter for valid time range.
func (f *TransferFilter) Validate() error {
	if f == nil {
		return nil
	}
	return f.TimeRange.Validate()
}

// ChainsFilter holds chain filter parameters for balance and portfolio
// endpoints.
type ChainsFilter struct {
	// Chains filters by blockchain.
	Chains []string
}

// ApplyToValues adds chain filter parameters to the given url.Values.
func (c *ChainsFilter) ApplyToValues(q url.Values) {
	if c == nil {
		return
	}
	if len(c.Chains) > 0 {
		q.Set("chains", joinChains(c.Chains))
	}
}

// UnixMilli returns a Unix millisecond timestamp string from a time.Time.
func UnixMilli(t time.Time) string {
	return strconv.FormatInt(t.UnixMilli(), 10)
}
