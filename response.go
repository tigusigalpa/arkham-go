package arkham

import (
	"time"
)

// ResponseMetadata holds non-body information from an HTTP response.
// It is returned alongside the decoded payload so applications can
// inspect rate-limit headers, intel usage headers, and timing.
type ResponseMetadata struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// RetryAfter is the value of the Retry-After header, if present.
	RetryAfter string
	// IntelDatapointsUsage is from X-Intel-Datapoints-Usage header.
	IntelDatapointsUsage string
	// IntelDatapointsLimit is from X-Intel-Datapoints-Limit header.
	IntelDatapointsLimit string
	// IntelDatapointsRemaining is from X-Intel-Datapoints-Remaining header.
	IntelDatapointsRemaining string
	// Duration is the time taken for the HTTP round-trip.
	Duration time.Duration
	// FinalURL is the redacted final URL of the request.
	FinalURL string
}

// Logger is a minimal logging interface for request diagnostics.
// Implementations must redact secrets; the SDK never passes the API
// key to the logger.
type Logger interface {
	Logf(format string, args ...interface{})
}

// nopLogger is a no-op logger that discards all output.
type nopLogger struct{}

func (nopLogger) Logf(string, ...interface{}) {}
