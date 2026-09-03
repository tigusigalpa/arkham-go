package arkham

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// get performs an HTTP GET request and decodes the JSON response.
func (c *Client) get(ctx context.Context, path string, query url.Values, out interface{}) (*ResponseMetadata, error) {
	return c.do(ctx, http.MethodGet, path, query, nil, out, true, 0)
}

// post performs an HTTP POST with a JSON body.
func (c *Client) post(ctx context.Context, path string, body interface{}, out interface{}) (*ResponseMetadata, error) {
	return c.doWithBody(ctx, http.MethodPost, path, nil, body, out, false, 0)
}

// postWithRetry performs an HTTP POST with retry enabled (use only
// for idempotent POST endpoints).
func (c *Client) postWithRetry(ctx context.Context, path string, body interface{}, out interface{}) (*ResponseMetadata, error) {
	return c.doWithBody(ctx, http.MethodPost, path, nil, body, out, true, 0)
}

// put performs an HTTP PUT with a JSON body. No retry by default.
func (c *Client) put(ctx context.Context, path string, body interface{}, out interface{}) (*ResponseMetadata, error) {
	return c.doWithBody(ctx, http.MethodPut, path, nil, body, out, false, 0)
}

// delete performs an HTTP DELETE. No retry by default.
func (c *Client) delete(ctx context.Context, path string, out interface{}) (*ResponseMetadata, error) {
	return c.do(ctx, http.MethodDelete, path, nil, nil, out, false, 0)
}

// do is the core request method for requests without a body.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, bodyBytes []byte, out interface{}, allowRetry bool, retried int) (*ResponseMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	fullURL := c.buildURL(path, query)

	// A fresh reader is created from bodyBytes on every attempt (including
	// retries) since an io.Reader is drained after a single request and
	// cannot be safely reused.
	var reqBody io.Reader
	if bodyBytes != nil {
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("arkham: failed to build request: %w", err)
	}

	c.applyHeaders(req)

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("arkham: request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	meta := c.captureMetadata(resp, fullURL, elapsed)

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return meta, fmt.Errorf("arkham: failed to read response body: %w", readErr)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if out != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, out); err != nil {
				return meta, fmt.Errorf("arkham: failed to decode response: %w", err)
			}
		}
		return meta, nil
	}

	if allowRetry && isRetryable(resp.StatusCode) && retried < c.maxRetries {
		delay := c.retryDelay(meta.RetryAfter, retried)
		c.logger.Logf("arkham: retrying %s %s (status %d, attempt %d, delay %v)", method, path, resp.StatusCode, retried+1, delay)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return meta, ctx.Err()
		case <-timer.C:
		}
		return c.do(ctx, method, path, query, bodyBytes, out, allowRetry, retried+1)
	}

	return meta, newAPIError(resp.StatusCode, extractMessage(respBody), respBody, meta, retried)
}

// doWithBody handles requests with a JSON body.
func (c *Client) doWithBody(ctx context.Context, method, path string, query url.Values, body interface{}, out interface{}, allowRetry bool, retried int) (*ResponseMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var jsonBytes []byte
	if body != nil {
		var err error
		jsonBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("arkham: failed to marshal request body: %w", err)
		}
	}

	return c.do(ctx, method, path, query, jsonBytes, out, allowRetry, retried)
}

// buildURL constructs the full request URL with query parameters.
func (c *Client) buildURL(path string, query url.Values) string {
	fullURL := strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(path, "/")
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}
	return fullURL
}

// applyHeaders sets authentication and content headers on the request.
func (c *Client) applyHeaders(req *http.Request) {
	req.Header.Set("API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
}

// captureMetadata builds a ResponseMetadata from the HTTP response.
func (c *Client) captureMetadata(resp *http.Response, fullURL string, elapsed time.Duration) *ResponseMetadata {
	return &ResponseMetadata{
		StatusCode:               resp.StatusCode,
		RetryAfter:               resp.Header.Get("Retry-After"),
		IntelDatapointsUsage:     resp.Header.Get("X-Intel-Datapoints-Usage"),
		IntelDatapointsLimit:     resp.Header.Get("X-Intel-Datapoints-Limit"),
		IntelDatapointsRemaining: resp.Header.Get("X-Intel-Datapoints-Remaining"),
		Duration:                 elapsed,
		FinalURL:                 fullURL,
	}
}

// isRetryable returns true for status codes that warrant a retry.
func isRetryable(statusCode int) bool {
	return statusCode == 429 || statusCode >= 500
}

// retryDelay computes the delay before the next retry. It honors Retry-After
// values expressed as either seconds or an HTTP date; otherwise it uses
// exponential backoff with jitter.
func (c *Client) retryDelay(retryAfter string, attempt int) time.Duration {
	if retryAfter = strings.TrimSpace(retryAfter); retryAfter != "" {
		if secs, err := strconv.Atoi(retryAfter); err == nil && secs >= 0 {
			return time.Duration(secs) * time.Second
		}
		if retryAt, err := http.ParseTime(retryAfter); err == nil {
			if delay := time.Until(retryAt); delay > 0 {
				return delay
			}
			return 0
		}
	}
	base := c.baseDelay
	if base <= 0 {
		base = DefaultBaseDelay
	}
	backoff := base * time.Duration(1<<uint(attempt))
	jitter := time.Duration(rand.Int63n(int64(base)))
	return backoff + jitter
}

// extractMessage attempts to pull a human-readable message from a JSON
// error response body. Arkham uses {"error": "description"} format.
func extractMessage(body []byte) string {
	var parsed struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &parsed); err == nil {
		if parsed.Error != "" {
			return parsed.Error
		}
		if parsed.Message != "" {
			return parsed.Message
		}
	}
	if len(body) == 0 {
		return "unknown error"
	}
	snippet := string(body)
	if len(snippet) > 200 {
		snippet = snippet[:200]
	}
	return snippet
}

// pathEscape escapes a single path segment.
func pathEscape(s string) string {
	return url.PathEscape(s)
}

// joinChains joins chain values with commas for query parameters.
func joinChains(chains []string) string {
	return strings.Join(chains, ",")
}
