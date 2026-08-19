package arkham

import (
	"net/http"
	"os"
	"time"
)

// Option configures a Client. Options are applied in order.
type Option func(*Client)

// WithBaseURL overrides the default Arkham API base URL
// (https://api.arkm.com). Primarily useful for testing.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.baseURL = url
		}
	}
}

// WithWebSocketURL overrides the default WebSocket base URL
// (wss://api.arkm.com).
func WithWebSocketURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.wsBaseURL = url
		}
	}
}

// WithHTTPClient sets a custom *http.Client for all requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithTimeout sets the timeout on the underlying HTTP client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.httpClient.Timeout = d
		}
	}
}

// WithMaxRetries sets the maximum number of retry attempts for
// retryable responses (429, 5xx). Set to 0 to disable retries.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		c.maxRetries = n
	}
}

// WithBaseDelay sets the initial backoff delay for retry attempts.
func WithBaseDelay(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.baseDelay = d
		}
	}
}

// WithUserAgent overrides the default User-Agent header value.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// WithLogger sets a logger for request/response diagnostics. The logger
// must redact secrets; the SDK never passes the API key to it.
func WithLogger(l Logger) Option {
	return func(c *Client) {
		if l != nil {
			c.logger = l
		}
	}
}

// osLookupEnv wraps os.LookupEnv for testability.
var osLookupEnv = func(key string) string {
	v, _ := os.LookupEnv(key)
	return v
}
