package arkham

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"
)

// StreamsService provides access to WebSocket v2 stream management endpoints
// and real-time transfer streaming.
//
// WebSocket v2 lifecycle:
// 1. Create a stream with POST /ws/v2/streams (filter set at creation time)
// 2. Connect via wss://api.arkm.com/ws/v2/transfers?stream_id=<id>
// 3. Receive transfers automatically (no subscribe/unsubscribe needed)
// 4. Reconnect within 10 minutes to reactivate and receive missed transfers
// 5. Delete with DELETE /ws/v2/streams/{id}
//
// Limits: max 10 active streams per user.
// Billing: 2 credits per transfer delivered. Stream management is free.
type StreamsService struct {
	client *Client
}

// Create creates a new WebSocket v2 stream with the given filter.
// Path: POST /ws/v2/streams
func (s *StreamsService) Create(ctx context.Context, req *CreateStreamV2Request) (*StreamV2, *ResponseMetadata, error) {
	if err := validateStreamFilter(req); err != nil {
		return nil, nil, err
	}
	var out StreamV2
	meta, err := s.client.post(ctx, "/ws/v2/streams", req, &out)
	return &out, meta, err
}

// List retrieves all WebSocket v2 streams for the authenticated user.
// Path: GET /ws/v2/streams
func (s *StreamsService) List(ctx context.Context) ([]StreamV2, *ResponseMetadata, error) {
	var out []StreamV2
	meta, err := s.client.get(ctx, "/ws/v2/streams", nil, &out)
	return out, meta, err
}

// Delete deletes a WebSocket v2 stream by streamId.
// Path: DELETE /ws/v2/streams/{id}
func (s *StreamsService) Delete(ctx context.Context, streamID string) (*DeleteStreamV2Response, *ResponseMetadata, error) {
	var out DeleteStreamV2Response
	meta, err := s.client.delete(ctx, "/ws/v2/streams/"+pathEscape(streamID), &out)
	return &out, meta, err
}

// validateStreamFilter checks that the stream filter meets minimum requirements.
// At least one of base, from, to, tokens must be set, or usdGte >= 250000.
func validateStreamFilter(req *CreateStreamV2Request) error {
	if req == nil {
		return ErrInvalidStreamFilter
	}
	hasIdentifying := len(req.Base) > 0 || len(req.From) > 0 || len(req.To) > 0 || len(req.Tokens) > 0
	if hasIdentifying {
		return nil
	}
	if req.UsdGte != "" {
		if val, err := strconv.ParseFloat(req.UsdGte, 64); err == nil && val >= 250000 {
			return nil
		}
	}
	return ErrInvalidStreamFilter
}

// WSConn represents a WebSocket v2 connection for receiving transfers.
type WSConn struct {
	streamID  string
	apiKey    string
	wsBaseURL string
	conn      wsConnAdapter
	mu        sync.Mutex
	closed    bool
}

// Connect establishes a WebSocket connection to the v2 transfer stream.
// The connection receives transfers automatically; no subscribe message
// is needed. Use Receive() to read messages.
func (s *StreamsService) Connect(ctx context.Context, streamID string) (*WSConn, error) {
	wsURL := s.client.wsBaseURL + "/ws/v2/transfers?stream_id=" + url.QueryEscape(streamID)

	conn, err := dialWebSocket(ctx, wsURL, s.client.apiKey)
	if err != nil {
		return nil, fmt.Errorf("arkham: websocket connection failed: %w", err)
	}

	return &WSConn{
		streamID:  streamID,
		apiKey:    s.client.apiKey,
		wsBaseURL: s.client.wsBaseURL,
		conn:      conn,
	}, nil
}

// Receive reads the next message from the WebSocket connection. It blocks
// until a message is available or the connection is closed/canceled.
// Returns the raw JSON message.
func (w *WSConn) Receive() (json.RawMessage, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return nil, fmt.Errorf("arkham: websocket connection closed")
	}
	return w.conn.ReadMessage()
}

// ReceiveTyped reads the next message and decodes it as a WSMessage.
func (w *WSConn) ReceiveTyped() (*WSMessage, error) {
	raw, err := w.Receive()
	if err != nil {
		return nil, err
	}
	var msg WSMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, fmt.Errorf("arkham: failed to decode websocket message: %w", err)
	}
	return &msg, nil
}

// Close closes the WebSocket connection.
func (w *WSConn) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return nil
	}
	w.closed = true
	return w.conn.Close()
}

// Reconnect attempts to reconnect to the stream. Streams expire after
// 10 minutes of inactivity; reconnecting within that window reactivates
// the stream and delivers any missed transfers.
func (w *WSConn) Reconnect(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.closed {
		_ = w.conn.Close()
	}

	wsURL := w.wsBaseURL + "/ws/v2/transfers?stream_id=" + url.QueryEscape(w.streamID)
	conn, err := dialWebSocket(ctx, wsURL, w.apiKey)
	if err != nil {
		return fmt.Errorf("arkham: websocket reconnection failed: %w", err)
	}
	w.conn = conn
	w.closed = false
	return nil
}

// wsConnAdapter is an interface for WebSocket connections, allowing
// the underlying implementation to be swapped (e.g. for testing).
type wsConnAdapter interface {
	ReadMessage() (json.RawMessage, error)
	Close() error
}

// dialWebSocket establishes a WebSocket connection with the API-Key header.
// Uses a minimal standard-library WebSocket client implementation.
func dialWebSocket(ctx context.Context, wsURL, apiKey string) (wsConnAdapter, error) {
	return wsDial(ctx, wsURL, apiKey)
}
