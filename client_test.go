package arkham

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const testAPIKey = "test-key-12345"

func TestNewClient_EmptyKey(t *testing.T) {
	_, err := NewClient("")
	if !errors.Is(err, ErrMissingAPIKey) {
		t.Fatalf("expected ErrMissingAPIKey, got %v", err)
	}
}

func TestNewClient_Success(t *testing.T) {
	c, err := NewClient("key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.apiKey != "key" {
		t.Fatalf("apiKey mismatch")
	}
	if c.baseURL != DefaultBaseURL {
		t.Fatalf("baseURL mismatch")
	}
	if c.Analytics == nil || c.Intelligence == nil || c.Transfers == nil {
		t.Fatalf("services not initialized")
	}
}

func TestBuildURLNormalizesSlashes(t *testing.T) {
	c, err := NewClient(testAPIKey, WithBaseURL("https://api.example.test/"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	got := c.buildURL("/transfers", url.Values{"limit": {"10"}})
	if want := "https://api.example.test/transfers?limit=10"; got != want {
		t.Fatalf("buildURL() = %q, want %q", got, want)
	}
}

func TestAPIKeyHeader(t *testing.T) {
	var gotKey string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("API-Key")
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, _ = c.get(context.Background(), "/test", nil, nil)

	if gotKey != testAPIKey {
		t.Fatalf("API-Key header = %q, want %q", gotKey, testAPIKey)
	}
}

func TestUserAgentHeader(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, _ = c.get(context.Background(), "/test", nil, nil)

	if gotUA != DefaultUserAgent {
		t.Fatalf("User-Agent = %q, want %q", gotUA, DefaultUserAgent)
	}
}

func TestGETRequest(t *testing.T) {
	var gotMethod, gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"address":"0xabc","chain":"ethereum","contract":false,"isUserAddress":false}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	var out Address
	meta, err := c.get(context.Background(), "/intelligence/address/0xabc", nil, &out)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if gotMethod != "GET" {
		t.Fatalf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/intelligence/address/0xabc" {
		t.Fatalf("path = %q, want /intelligence/address/0xabc", gotPath)
	}
	if out.Address != "0xabc" {
		t.Fatalf("address = %q", out.Address)
	}
	if meta.StatusCode != 200 {
		t.Fatalf("status = %d", meta.StatusCode)
	}
}

func TestPOSTRequest(t *testing.T) {
	var gotMethod string
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"streamId":"test-stream-id","id":1}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	req := &CreateStreamV2Request{Base: []string{"binance"}, UsdGte: "500000"}
	var out StreamV2
	_, err := c.post(context.Background(), "/ws/v2/streams", req, &out)
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	if gotMethod != "POST" {
		t.Fatalf("method = %q, want POST", gotMethod)
	}
	if !strings.Contains(string(gotBody), "binance") {
		t.Fatalf("body doesn't contain base: %s", gotBody)
	}
	if out.StreamID != "test-stream-id" {
		t.Fatalf("streamId = %q", out.StreamID)
	}
}

func TestDELETERequest(t *testing.T) {
	var gotMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_, _ = fmt.Fprint(w, `{"success":true}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	var out DeleteStreamV2Response
	_, err := c.delete(context.Background(), "/ws/v2/streams/abc", &out)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Fatalf("method = %q, want DELETE", gotMethod)
	}
	if !out.Success {
		t.Fatalf("success = false")
	}
}

func TestError400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error":"invalid parameter"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Fatalf("status = %d", apiErr.StatusCode)
	}
	if apiErr.Message != "invalid parameter" {
		t.Fatalf("message = %q", apiErr.Message)
	}
}

func TestError401(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"error":"unauthorized"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestError403(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		fmt.Fprintf(w, `{"error":"forbidden"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestError404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintf(w, `{"error":"not found"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestError429(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(429)
		fmt.Fprintf(w, `{"error":"rate limited"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.RetryAfter != "1" {
			t.Fatalf("RetryAfter = %q", apiErr.RetryAfter)
		}
	}
}

func TestRetryDelayHTTPDate(t *testing.T) {
	c, err := NewClient(testAPIKey)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	delay := c.retryDelay(time.Now().Add(2*time.Second).UTC().Format(http.TimeFormat), 0)
	if delay <= 0 || delay > 2*time.Second {
		t.Fatalf("retry delay = %v, want a positive delay no greater than two seconds", delay)
	}
}

func TestError500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error":"internal server error"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if !errors.Is(err, ErrServerError) {
		t.Fatalf("expected ErrServerError, got %v", err)
	}
}

func TestErrorNonJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintf(w, "plain text error")
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError")
	}
	if !strings.Contains(apiErr.Message, "plain text") {
		t.Fatalf("message = %q", apiErr.Message)
	}
}

func TestRetryOn429(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			fmt.Fprintf(w, `{"error":"rate limited"}`)
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"address":"0xabc","chain":"ethereum","contract":false,"isUserAddress":false}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(3), WithBaseDelay(10*time.Millisecond))
	var out Address
	_, err := c.get(context.Background(), "/test", nil, &out)
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	if out.Address != "0xabc" {
		t.Fatalf("address = %q", out.Address)
	}
}

func TestRetryExhausted(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(429)
		fmt.Fprintf(w, `{"error":"rate limited"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(2), WithBaseDelay(10*time.Millisecond))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls (1 + 2 retries), got %d", calls)
	}
}

func TestContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := c.get(ctx, "/test", nil, nil)
	if err == nil {
		t.Fatalf("expected context cancellation error")
	}
}

func TestIntelUsageHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Intel-Datapoints-Usage", "5000")
		w.Header().Set("X-Intel-Datapoints-Limit", "10000")
		w.Header().Set("X-Intel-Datapoints-Remaining", "5000")
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"totalCount":5000,"totalLimit":10000,"periodStart":"2024-01-01"}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	var out IntelUsage
	meta, err := c.get(context.Background(), "/subscription/intel-usage", nil, &out)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if meta.IntelDatapointsUsage != "5000" {
		t.Fatalf("usage = %q", meta.IntelDatapointsUsage)
	}
	if meta.IntelDatapointsLimit != "10000" {
		t.Fatalf("limit = %q", meta.IntelDatapointsLimit)
	}
	if meta.IntelDatapointsRemaining != "5000" {
		t.Fatalf("remaining = %q", meta.IntelDatapointsRemaining)
	}
	if out.TotalCount != 5000 {
		t.Fatalf("totalCount = %d", out.TotalCount)
	}
}

func TestTimeRangeValidation(t *testing.T) {
	tr := &TimeRange{TimeLast: "24h", TimeGte: "2024-01-01T00:00:00Z"}
	err := tr.Validate()
	if !errors.Is(err, ErrInvalidTimeRange) {
		t.Fatalf("expected ErrInvalidTimeRange, got %v", err)
	}

	tr2 := &TimeRange{TimeLast: "24h"}
	if err := tr2.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr3 := &TimeRange{TimeGte: "2024-01-01T00:00:00Z", TimeLte: "2024-06-01T00:00:00Z"}
	if err := tr3.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransferFilterValidation(t *testing.T) {
	filter := &TransferFilter{
		TimeRange: &TimeRange{TimeLast: "24h", TimeGte: "2024-01-01"},
	}
	err := filter.Validate()
	if !errors.Is(err, ErrInvalidTimeRange) {
		t.Fatalf("expected ErrInvalidTimeRange, got %v", err)
	}
}

func TestStreamFilterValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateStreamV2Request
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty", &CreateStreamV2Request{}, true},
		{"with base", &CreateStreamV2Request{Base: []string{"binance"}}, false},
		{"with from", &CreateStreamV2Request{From: []string{"0xabc"}}, false},
		{"with usdGte high", &CreateStreamV2Request{UsdGte: "300000"}, false},
		{"with usdGte low", &CreateStreamV2Request{UsdGte: "100000"}, true},
		{"with tokens", &CreateStreamV2Request{Tokens: []string{"ethereum"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStreamFilter(tt.req)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPagination(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		w.Header().Set("Content-Type", "application/json")
		if callCount <= 2 {
			fmt.Fprintf(w, `[{"chain":"ethereum","value":"1.0","from":"0xa","to":"0xb","time":"2024-01-01T00:00:00Z"}]`)
		} else {
			fmt.Fprintf(w, `[]`)
		}
		_ = limit
		_ = offset
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	p := NewPaginator(context.Background(), c, "/transfers", nil, 10, 0, 3)
	for p.HasNext() {
		var out []Transfer
		_, err := p.NextPage(&out)
		if err != nil {
			t.Fatalf("NextPage failed: %v", err)
		}
		if len(out) == 0 {
			break
		}
	}
	if callCount > 3 {
		t.Fatalf("expected max 3 calls, got %d", callCount)
	}
}

func TestPaginationMaxPages(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		fmt.Fprintf(w, `[{"chain":"ethereum","value":"1.0","from":"0xa","to":"0xb","time":"2024-01-01T00:00:00Z"}]`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	p := NewPaginator(context.Background(), c, "/transfers", nil, 10, 0, 2)
	for p.HasNext() {
		var out []Transfer
		_, err := p.NextPage(&out)
		if err != nil {
			t.Fatalf("NextPage failed: %v", err)
		}
	}
	if callCount != 2 {
		t.Fatalf("expected 2 calls, got %d", callCount)
	}
}

func TestPaginatorMaxItemsLimitsFinalRequest(t *testing.T) {
	var limits []string
	var offsets []string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limits = append(limits, r.URL.Query().Get("limit"))
		offsets = append(offsets, r.URL.Query().Get("offset"))
		_, _ = fmt.Fprint(w, `[]`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	p := NewPaginator(context.Background(), c, "/transfers", nil, 10, 23, 0)
	for p.HasNext() {
		var out []Transfer
		if _, err := p.NextPage(&out); err != nil {
			t.Fatalf("NextPage failed: %v", err)
		}
	}

	if got, want := strings.Join(limits, ","), "10,10,3"; got != want {
		t.Fatalf("limits = %q, want %q", got, want)
	}
	if got, want := strings.Join(offsets, ","), "0,10,20"; got != want {
		t.Fatalf("offsets = %q, want %q", got, want)
	}
}

func TestPaginatorDefaultsNonPositivePageSize(t *testing.T) {
	p := NewPaginator(context.Background(), nil, "/transfers", nil, 0, 0, 0)
	if p.pageSize != DefaultPageSize {
		t.Fatalf("pageSize = %d, want %d", p.pageSize, DefaultPageSize)
	}
}

func TestEmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	_, err := c.get(context.Background(), "/test", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMalformedJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{invalid json`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	var out Address
	_, err := c.get(context.Background(), "/test", nil, &out)
	if err == nil {
		t.Fatalf("expected decode error")
	}
}

func TestUnknownFieldsTolerated(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"address":"0xabc","chain":"ethereum","contract":false,"isUserAddress":false,"newField":"value","anotherNew":123}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	var out Address
	_, err := c.get(context.Background(), "/test", nil, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Address != "0xabc" {
		t.Fatalf("address = %q", out.Address)
	}
}

func TestIntelligenceAddress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/intelligence/address/0xabc" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"address":"0xabc","chain":"ethereum","contract":false,"isUserAddress":false}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Intelligence.Address(context.Background(), "0xabc", nil)
	if err != nil {
		t.Fatalf("Address failed: %v", err)
	}
	if out.Address != "0xabc" {
		t.Fatalf("address = %q", out.Address)
	}
}

func TestIntelligenceSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("query") != "binance" {
			t.Fatalf("query = %q", q.Get("query"))
		}
		if q.Get("arkhamEntities") != "5" {
			t.Fatalf("arkhamEntities = %q", q.Get("arkhamEntities"))
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"arkhamEntities":[{"id":"binance","name":"Binance","note":"","service":true,"type":"cex"}]}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	opts := &SearchOptions{ArkhamEntities: 5}
	out, _, err := c.Intelligence.Search(context.Background(), "binance", opts)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(out.ArkhamEntities) != 1 {
		t.Fatalf("entities = %d", len(out.ArkhamEntities))
	}
	if out.ArkhamEntities[0].ID != "binance" {
		t.Fatalf("id = %q", out.ArkhamEntities[0].ID)
	}
}

func TestTransfersWithFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("base") != "binance" {
			t.Fatalf("base = %q", q.Get("base"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("limit = %q", q.Get("limit"))
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"transfers":[{"chain":"ethereum","value":"1.0","from":"0xa","to":"0xb","time":"2024-01-01T00:00:00Z"}]}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	filter := &TransferFilter{Base: []string{"binance"}, Limit: 10}
	out, _, err := c.Transfers.Transfers(context.Background(), filter)
	if err != nil {
		t.Fatalf("Transfers failed: %v", err)
	}
	if len(out.Transfers) != 1 {
		t.Fatalf("transfers = %d", len(out.Transfers))
	}
}

func TestBalancesAddress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/balances/address/0xabc" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"balances":{"ethereum":[]},"totalBalance":{"ethereum":0},"addresses":{}}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Balances.Address(context.Background(), "0xabc", nil)
	if err != nil {
		t.Fatalf("Address failed: %v", err)
	}
	if out == nil {
		t.Fatalf("out is nil")
	}
}

func TestChainsList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `[{"id":"ethereum","name":"Ethereum","chainType":"evm"}]`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Chains.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("chains = %d", len(out))
	}
	if out[0].ID != "ethereum" {
		t.Fatalf("id = %q", out[0].ID)
	}
}

func TestUserAlerts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `[{"id":1,"name":"test","enabled":true,"alertMethodId":1}]`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.User.ListAlerts(context.Background())
	if err != nil {
		t.Fatalf("ListAlerts failed: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("alerts = %d", len(out))
	}
	if out[0].ID != 1 {
		t.Fatalf("id = %d", out[0].ID)
	}
}

func TestUserCreateAlert(t *testing.T) {
	var gotMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"id":1,"name":"test","enabled":true,"alertMethodId":1}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	req := &AlertRequest{Name: "test", Enabled: true, AlertMethodID: 1}
	out, _, err := c.User.CreateAlert(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateAlert failed: %v", err)
	}
	if gotMethod != "POST" {
		t.Fatalf("method = %q", gotMethod)
	}
	if out.ID != 1 {
		t.Fatalf("id = %d", out.ID)
	}
}

func TestStreamCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"streamId":"abc-123","id":1}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	req := &CreateStreamV2Request{Base: []string{"binance"}}
	out, _, err := c.Streams.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if out.StreamID != "abc-123" {
		t.Fatalf("streamId = %q", out.StreamID)
	}
}

func TestStreamCreateInvalidFilter(t *testing.T) {
	c, _ := NewClient(testAPIKey, WithMaxRetries(0))
	req := &CreateStreamV2Request{}
	_, _, err := c.Streams.Create(context.Background(), req)
	if !errors.Is(err, ErrInvalidStreamFilter) {
		t.Fatalf("expected ErrInvalidStreamFilter, got %v", err)
	}
}

func TestStreamList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `[{"streamId":"abc","id":1}]`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Streams.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("streams = %d", len(out))
	}
}

func TestStreamDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"success":true}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Streams.Delete(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !out.Success {
		t.Fatalf("success = false")
	}
}

func TestRiskAddress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"chain_type":"ethereum","address":"0xabc","risk_level":"low","is_seed":false}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Risk.Address(context.Background(), "0xabc")
	if err != nil {
		t.Fatalf("Address failed: %v", err)
	}
	if out.RiskLevel != "low" {
		t.Fatalf("risk_level = %q", out.RiskLevel)
	}
}

func TestPolymarketEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("limit") != "50" {
			t.Fatalf("limit = %q", q.Get("limit"))
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"events":[{"id":"1","active":true}],"count":1}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	opts := &PolymarketEventOptions{Limit: 50}
	out, _, err := c.Polymarket.Events(context.Background(), opts)
	if err != nil {
		t.Fatalf("Events failed: %v", err)
	}
	if len(out.Events) != 1 {
		t.Fatalf("events = %d", len(out.Events))
	}
}

func TestHypercoreMarkets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"markets":[{"name":"Bitcoin","pricingId":"bitcoin","symbol":"BTC"}]}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Hypercore.Markets(context.Background())
	if err != nil {
		t.Fatalf("Markets failed: %v", err)
	}
	if len(out.Markets) != 1 {
		t.Fatalf("markets = %d", len(out.Markets))
	}
	if out.Markets[0].Symbol != "BTC" {
		t.Fatalf("symbol = %q", out.Markets[0].Symbol)
	}
}

func TestArkhamCirculating(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"circulating":1000000.0}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	out, _, err := c.Arkham.Circulating(context.Background())
	if err != nil {
		t.Fatalf("Circulating failed: %v", err)
	}
	if out.Circulating != 1000000.0 {
		t.Fatalf("circulating = %f", out.Circulating)
	}
}

func TestQueryParamsEncoding(t *testing.T) {
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"transfers":[]}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	filter := &TransferFilter{
		Base:   []string{"binance", "!wintermute"},
		Chains: []string{"ethereum", "bsc"},
		Flow:   FlowIn,
		Limit:  10,
		Offset: 0,
	}
	_, _, err := c.Transfers.Transfers(context.Background(), filter)
	if err != nil {
		t.Fatalf("Transfers failed: %v", err)
	}
	if !strings.Contains(gotQuery, "base=binance") {
		t.Fatalf("query missing base: %s", gotQuery)
	}
	if !strings.Contains(gotQuery, "flow=in") {
		t.Fatalf("query missing flow: %s", gotQuery)
	}
	if !strings.Contains(gotQuery, "limit=10") {
		t.Fatalf("query missing limit: %s", gotQuery)
	}
}

func TestWithTimeout(t *testing.T) {
	c, _ := NewClient(testAPIKey, WithTimeout(5*time.Second))
	if c.httpClient.Timeout != 5*time.Second {
		t.Fatalf("timeout = %v", c.httpClient.Timeout)
	}
}

func TestWithUserAgent(t *testing.T) {
	c, _ := NewClient(testAPIKey, WithUserAgent("custom/1.0"))
	if c.userAgent != "custom/1.0" {
		t.Fatalf("userAgent = %q", c.userAgent)
	}
}

func TestWithBaseURL(t *testing.T) {
	c, _ := NewClient(testAPIKey, WithBaseURL("https://custom.example.com"))
	if c.baseURL != "https://custom.example.com" {
		t.Fatalf("baseURL = %q", c.baseURL)
	}
}

// TestPOSTRetryPreservesBody is a regression test for a bug where retrying
// a request with a body would resend an already-drained io.Reader,
// resulting in an empty body on the retried attempt.
func TestPOSTRetryPreservesBody(t *testing.T) {
	calls := 0
	var gotBodies []string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		b, _ := io.ReadAll(r.Body)
		gotBodies = append(gotBodies, string(b))
		if calls < 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			fmt.Fprintf(w, `{"error":"rate limited"}`)
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"streamId":"abc","id":1}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(3), WithBaseDelay(5*time.Millisecond))
	req := &CreateStreamV2Request{Base: []string{"binance"}}
	var out StreamV2
	// postWithRetry allows retry; exercises the fixed retry path.
	_, err := c.postWithRetry(context.Background(), "/ws/v2/streams", req, &out)
	if err != nil {
		t.Fatalf("postWithRetry failed: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
	for i, b := range gotBodies {
		if !strings.Contains(b, "binance") {
			t.Fatalf("call %d: body missing data, got empty/wrong body: %q", i, b)
		}
	}
}

// TestMaskedFrameIsMasked is a regression test ensuring client-to-server
// WebSocket control frames (pong/close) are masked per RFC 6455 §5.1.
func TestMaskedFrameIsMasked(t *testing.T) {
	payload := []byte("ping-data")
	frame := maskedFrame(0x0A, payload)

	if len(frame) < 2 {
		t.Fatalf("frame too short: %d bytes", len(frame))
	}
	// Byte 1 bit 7 must be the MASK bit (1 = masked).
	if frame[1]&0x80 == 0 {
		t.Fatalf("MASK bit not set on client frame: %08b", frame[1])
	}
	// FIN bit set, opcode preserved.
	if frame[0]&0x80 == 0 {
		t.Fatalf("FIN bit not set: %08b", frame[0])
	}
	if frame[0]&0x0F != 0x0A {
		t.Fatalf("opcode mismatch: got %x, want 0xA", frame[0]&0x0F)
	}
	// Verify payload round-trips through the mask key.
	length := int(frame[1] & 0x7F)
	if length != len(payload) {
		t.Fatalf("length mismatch: got %d, want %d", length, len(payload))
	}
	maskKey := frame[2:6]
	masked := frame[6:]
	unmasked := make([]byte, len(masked))
	for i, b := range masked {
		unmasked[i] = b ^ maskKey[i%4]
	}
	if string(unmasked) != string(payload) {
		t.Fatalf("payload did not round-trip: got %q, want %q", unmasked, payload)
	}
}

func TestJSONRawMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"key":"value","nested":{"a":1}}`)
	}))
	defer ts.Close()

	c, _ := NewClient(testAPIKey, WithBaseURL(ts.URL), WithMaxRetries(0))
	var out json.RawMessage
	_, err := c.get(context.Background(), "/test", nil, &out)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if !strings.Contains(string(out), "value") {
		t.Fatalf("raw message = %s", out)
	}
}
