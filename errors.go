package arkham

import (
	"errors"
	"fmt"
)

// APIError represents an error response from the Arkham API. It is
// returned when the API responds with a non-2xx status code. The error
// carries the status code, a safe response excerpt, retry metadata, and
// intel usage headers, but never includes the API key.
type APIError struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Message is a human-readable error message extracted from the
	// response body or a generated fallback.
	Message string
	// RawBody is a truncated copy of the response body (max 512 bytes).
	RawBody []byte
	// RetryAfter is the value of the Retry-After header, if present.
	RetryAfter string
	// IntelDatapointsUsage is from X-Intel-Datapoints-Usage header.
	IntelDatapointsUsage string
	// IntelDatapointsLimit is from X-Intel-Datapoints-Limit header.
	IntelDatapointsLimit string
	// IntelDatapointsRemaining is from X-Intel-Datapoints-Remaining header.
	IntelDatapointsRemaining string
	// Retried is the number of retry attempts before this error.
	Retried int
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("arkham: api error: status=%d message=%s", e.StatusCode, e.Message)
}

// IsRateLimit returns true if the error is a rate-limit error (429).
func (e *APIError) IsRateLimit() bool {
	return e.StatusCode == 429
}

// Sentinel errors for use with errors.Is.
var (
	// ErrMissingAPIKey is returned when no API key is provided.
	ErrMissingAPIKey = errors.New("arkham: missing API key — set ARKHAM_API_KEY or pass a key to NewClient")
	// ErrBadRequest is returned for HTTP 400.
	ErrBadRequest = errors.New("arkham: bad request — invalid parameters")
	// ErrUnauthorized is returned for HTTP 401.
	ErrUnauthorized = errors.New("arkham: unauthorized — check your API key")
	// ErrPaymentRequired is returned for HTTP 402.
	ErrPaymentRequired = errors.New("arkham: payment required")
	// ErrForbidden is returned for HTTP 403.
	ErrForbidden = errors.New("arkham: forbidden — endpoint not in plan")
	// ErrNotFound is returned for HTTP 404.
	ErrNotFound = errors.New("arkham: resource not found")
	// ErrRateLimited is returned for HTTP 429 after retries are exhausted.
	ErrRateLimited = errors.New("arkham: rate limit exceeded")
	// ErrServerError is returned for HTTP 5xx.
	ErrServerError = errors.New("arkham: server error")
	// ErrInvalidTimeRange is returned when timeLast is combined with
	// timeGte or timeLte.
	ErrInvalidTimeRange = errors.New("arkham: timeLast cannot be used with timeGte or timeLte")
	// ErrInvalidStreamFilter is returned when a WebSocket stream filter
	// does not meet the minimum requirement.
	ErrInvalidStreamFilter = errors.New("arkham: stream filter requires at least one of base, from, to, tokens, or usdGte >= 250000")
)

// wrappedAPIError wraps *APIError with a sentinel for errors.Is/As.
type wrappedAPIError struct {
	*APIError
	sentinel error
}

// Unwrap returns both the *APIError and the sentinel so that
// errors.Is and errors.As both work.
func (w *wrappedAPIError) Unwrap() []error {
	return []error{w.APIError, w.sentinel}
}

// newAPIError builds an error for the given status code, attaching the
// appropriate sentinel.
func newAPIError(statusCode int, message string, rawBody []byte, meta *ResponseMetadata, retried int) error {
	apiErr := &APIError{
		StatusCode: statusCode,
		Message:    message,
		RawBody:    truncateBody(rawBody),
		Retried:    retried,
	}
	if meta != nil {
		apiErr.RetryAfter = meta.RetryAfter
		apiErr.IntelDatapointsUsage = meta.IntelDatapointsUsage
		apiErr.IntelDatapointsLimit = meta.IntelDatapointsLimit
		apiErr.IntelDatapointsRemaining = meta.IntelDatapointsRemaining
	}
	switch statusCode {
	case 400:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrBadRequest}
	case 401:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrUnauthorized}
	case 402:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrPaymentRequired}
	case 403:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrForbidden}
	case 404:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrNotFound}
	case 429:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrRateLimited}
	default:
		if statusCode >= 500 {
			return &wrappedAPIError{APIError: apiErr, sentinel: ErrServerError}
		}
		return apiErr
	}
}

// truncateBody limits the raw body to 512 bytes.
func truncateBody(body []byte) []byte {
	if len(body) > 512 {
		return body[:512]
	}
	return body
}
