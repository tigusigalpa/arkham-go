package arkham

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"
)

// wsConn implements wsConnAdapter using a raw TCP/TLS connection
// with the WebSocket protocol (RFC 6455). This is a minimal
// read-only client implementation using only the standard library.
type wsConn struct {
	conn   net.Conn
	reader *bufio.Reader
	closed bool
}

// wsDial establishes a WebSocket connection to the given URL with
// the API-Key header. It performs the HTTP upgrade handshake and
// returns a connection that can read text frames.
func wsDial(ctx context.Context, wsURL, apiKey string) (wsConnAdapter, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("arkham: invalid websocket URL: %w", err)
	}

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		if u.Scheme == "wss" {
			port = "443"
		} else {
			port = "80"
		}
	}
	addr := net.JoinHostPort(host, port)

	var conn net.Conn
	dialer := &net.Dialer{Timeout: 30 * time.Second}

	if u.Scheme == "wss" {
		tlsDialer := &tls.Dialer{
			NetDialer: dialer,
		}
		conn, err = tlsDialer.DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, fmt.Errorf("arkham: websocket dial failed: %w", err)
	}

	// Generate Sec-WebSocket-Key
	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		conn.Close()
		return nil, fmt.Errorf("arkham: failed to generate websocket key: %w", err)
	}
	wsKey := base64.StdEncoding.EncodeToString(keyBytes)

	// Build the HTTP upgrade request
	path := u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}
	if path == "" {
		path = "/"
	}

	reqStr := fmt.Sprintf(
		"GET %s HTTP/1.1\r\n"+
			"Host: %s\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Key: %s\r\n"+
			"Sec-WebSocket-Version: 13\r\n"+
			"API-Key: %s\r\n"+
			"\r\n",
		path, u.Host, wsKey, apiKey,
	)

	_, err = conn.Write([]byte(reqStr))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("arkham: failed to send websocket upgrade: %w", err)
	}

	reader := bufio.NewReader(conn)

	// Read the HTTP upgrade response
	respLine, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("arkham: failed to read websocket response: %w", err)
	}

	if !contains(respLine, "101") {
		conn.Close()
		return nil, fmt.Errorf("arkham: websocket upgrade failed: %s", respLine)
	}

	// Read and discard headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("arkham: failed to read websocket headers: %w", err)
		}
		if line == "\r\n" || line == "\n" {
			break
		}
	}

	return &wsConn{conn: conn, reader: reader}, nil
}

// ReadMessage reads the next WebSocket text frame and returns its payload.
// It handles fragmentation (continuation frames) and skips non-text frames.
func (w *wsConn) ReadMessage() (json.RawMessage, error) {
	var payload []byte
	for {
		frame, err := w.readFrame()
		if err != nil {
			return nil, err
		}

		if frame.opcode == 0x8 { // Close frame
			return nil, fmt.Errorf("arkham: websocket closed by server")
		}
		if frame.opcode == 0x9 { // Ping frame — send pong
			w.writePong(frame.data)
			continue
		}
		if frame.opcode == 0xA { // Pong frame — skip
			continue
		}

		payload = append(payload, frame.data...)

		if frame.fin {
			break
		}
		// Continuation frames follow
	}

	return json.RawMessage(payload), nil
}

// wsFrame represents a parsed WebSocket frame.
type wsFrame struct {
	fin    bool
	opcode byte
	data   []byte
}

// readFrame reads a single WebSocket frame from the connection.
func (w *wsConn) readFrame() (*wsFrame, error) {
	if w.closed {
		return nil, fmt.Errorf("arkham: connection closed")
	}

	header := make([]byte, 2)
	if _, err := io.ReadFull(w.reader, header); err != nil {
		return nil, fmt.Errorf("arkham: failed to read frame header: %w", err)
	}

	fin := header[0]&0x80 != 0
	opcode := header[0] & 0x0F
	masked := header[1]&0x80 != 0
	length := int64(header[1] & 0x7F)

	if length == 126 {
		extLen := make([]byte, 2)
		if _, err := io.ReadFull(w.reader, extLen); err != nil {
			return nil, fmt.Errorf("arkham: failed to read extended length: %w", err)
		}
		length = int64(binary.BigEndian.Uint16(extLen))
	} else if length == 127 {
		extLen := make([]byte, 8)
		if _, err := io.ReadFull(w.reader, extLen); err != nil {
			return nil, fmt.Errorf("arkham: failed to read extended length: %w", err)
		}
		length = int64(binary.BigEndian.Uint64(extLen))
	}

	var maskKey []byte
	if masked {
		maskKey = make([]byte, 4)
		if _, err := io.ReadFull(w.reader, maskKey); err != nil {
			return nil, fmt.Errorf("arkham: failed to read mask key: %w", err)
		}
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(w.reader, data); err != nil {
		return nil, fmt.Errorf("arkham: failed to read frame data: %w", err)
	}

	if masked {
		for i := range data {
			data[i] ^= maskKey[i%4]
		}
	}

	return &wsFrame{fin: fin, opcode: opcode, data: data}, nil
}

// writePong sends a masked pong frame in response to a ping. Per RFC 6455
// §5.1, all frames sent from client to server MUST be masked.
func (w *wsConn) writePong(data []byte) {
	w.conn.Write(maskedFrame(0x0A, data))
}

// Close closes the WebSocket connection, sending a masked close frame
// as required by RFC 6455 for client-originated frames.
func (w *wsConn) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true
	w.conn.Write(maskedFrame(0x08, nil))
	return w.conn.Close()
}

// maskedFrame builds a single-frame, masked WebSocket control frame
// (opcode + payload) as required for all client-to-server frames.
func maskedFrame(opcode byte, payload []byte) []byte {
	maskKey := make([]byte, 4)
	_, _ = rand.Read(maskKey)

	masked := make([]byte, len(payload))
	for i, b := range payload {
		masked[i] = b ^ maskKey[i%4]
	}

	frame := make([]byte, 0, 2+4+len(masked))
	frame = append(frame, 0x80|opcode)             // FIN + opcode
	frame = append(frame, 0x80|byte(len(payload))) // MASK bit + length (<=125)
	frame = append(frame, maskKey...)
	frame = append(frame, masked...)
	return frame
}

// contains checks if s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
