package arkham

import (
	"context"
	"net/url"
)

// DefaultPageSize is used when NewPaginator receives a non-positive page
// size. It keeps offsets advancing and prevents accidental repeated requests
// for the first page.
const DefaultPageSize = 100

// Paginator provides offset-based pagination for list endpoints.
// It fetches one page at a time and never eagerly traverses all pages.
type Paginator struct {
	client       *Client
	ctx          context.Context
	path         string
	query        url.Values
	pageSize     int
	maxItems     int
	maxPages     int
	offset       int
	totalFetched int
	pageCount    int
}

// NewPaginator creates a new Paginator for the given path and query.
// pageSize controls the limit parameter, maxItems caps total items
// fetched (0 = unlimited), maxPages caps the number of page requests
// (0 = unlimited).
func NewPaginator(ctx context.Context, client *Client, path string, query url.Values, pageSize, maxItems, maxPages int) *Paginator {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	return &Paginator{
		client:   client,
		ctx:      ctx,
		path:     path,
		query:    query,
		pageSize: pageSize,
		maxItems: maxItems,
		maxPages: maxPages,
	}
}

// HasNext returns true if more pages may be available.
func (p *Paginator) HasNext() bool {
	if p.maxPages > 0 && p.pageCount >= p.maxPages {
		return false
	}
	if p.maxItems > 0 && p.totalFetched >= p.maxItems {
		return false
	}
	return true
}

// NextOffset returns the offset for the next page request.
func (p *Paginator) NextOffset() int {
	return p.offset
}

// NextPage fetches the next page of results into out. Returns the
// response metadata and any error. The caller should inspect the number
// of items returned to decide whether to continue.
func (p *Paginator) NextPage(out interface{}) (*ResponseMetadata, error) {
	if err := p.ctx.Err(); err != nil {
		return nil, err
	}
	if !p.HasNext() {
		return nil, nil
	}

	q := url.Values{}
	for k, v := range p.query {
		q[k] = v
	}
	limit := p.pageSize
	if p.maxItems > 0 {
		remaining := p.maxItems - p.totalFetched
		if remaining < limit {
			limit = remaining
		}
	}
	q.Set("limit", intToString(limit))
	q.Set("offset", intToString(p.offset))

	meta, err := p.client.get(p.ctx, p.path, q, out)
	if err != nil {
		return meta, err
	}

	p.pageCount++
	p.totalFetched += limit
	p.offset += limit
	return meta, nil
}

// intToString converts an int to a string without importing strconv
// elsewhere in this file.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
