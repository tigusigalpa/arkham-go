package arkham

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestFiltersEncodeConfiguredValues(t *testing.T) {
	values := url.Values{}
	filter := &TransferFilter{
		Base:           []string{"base"},
		Chains:         []string{"ethereum"},
		Flow:           FlowIn,
		From:           []string{"from"},
		To:             []string{"to"},
		Counterparties: []string{"counterparty"},
		Tokens:         []string{"token"},
		TimeRange:      &TimeRange{TimeLast: "24h", TimeGte: "1", TimeLte: "2"},
		ValueGte:       "1",
		ValueLte:       "2",
		UsdGte:         "3",
		UsdLte:         "4",
		SortKey:        SortKeyUSD,
		SortDir:        SortDirDesc,
		Limit:          10,
		Offset:         1,
	}
	filter.ApplyToValues(values)
	if got, want := len(values), 18; got != want {
		t.Fatalf("encoded values = %d, want %d: %v", got, want, values)
	}
	if err := filter.Validate(); !errors.Is(err, ErrInvalidTimeRange) {
		t.Fatalf("Validate() error = %v, want ErrInvalidTimeRange", err)
	}

	search := SearchOptions{}
	populateEndpointValue(reflect.ValueOf(&search).Elem())
	values = url.Values{}
	search.ApplyToValues(values)
	if got, want := len(values), 25; got != want {
		t.Fatalf("search values = %d, want %d", got, want)
	}

	active, groupGames := true, false
	polymarket := &PolymarketEventOptions{
		Tag: "tag", ExcludeTag: "excluded", Active: &active, Search: "search",
		GroupGames: &groupGames, SortBy: "volume", Order: "desc", Limit: 10, Offset: 1,
	}
	values = url.Values{}
	polymarket.ApplyToValues(values)
	if got, want := len(values), 9; got != want {
		t.Fatalf("Polymarket values = %d, want %d", got, want)
	}

	values = url.Values{}
	(&ChainsFilter{Chains: []string{"ethereum", "bitcoin"}}).ApplyToValues(values)
	if got, want := values.Get("chains"), "ethereum,bitcoin"; got != want {
		t.Fatalf("chains = %q, want %q", got, want)
	}
	if got, want := UnixMilli(time.UnixMilli(123)), "123"; got != want {
		t.Fatalf("UnixMilli() = %q, want %q", got, want)
	}
}

func TestWebSocketTransport(t *testing.T) {
	conn := &memoryConn{}
	frames := append(serverFrame(0x9, []byte("ping")), serverFrame(0x1, []byte(`{"id":1}`))...)
	ws := &wsConn{conn: conn, reader: bufio.NewReader(bytes.NewReader(frames))}
	message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage() error = %v", err)
	}
	if got, want := string(message), `{"id":1}`; got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
	if got := conn.Bytes(); len(got) == 0 || got[0] != 0x8A {
		t.Fatalf("pong frame = %x, want masked pong", got)
	}

	for _, payload := range [][]byte{bytes.Repeat([]byte("a"), 126), bytes.Repeat([]byte("b"), 127)} {
		frame := extendedServerFrame(payload)
		ws = &wsConn{reader: bufio.NewReader(bytes.NewReader(frame))}
		parsed, err := ws.readFrame()
		if err != nil {
			t.Fatalf("readFrame() error = %v", err)
		}
		if got, want := len(parsed.data), len(payload); got != want {
			t.Fatalf("frame length = %d, want %d", got, want)
		}
	}

	if err := (&wsConn{conn: &memoryConn{}}).Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !contains("websocket", "socket") || searchString("abc", "z") {
		t.Fatal("string helpers returned an unexpected result")
	}
}

func TestWSDialAndConnectionLifecycle(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil || line == "\r\n" {
				break
			}
		}
		_, _ = io.WriteString(conn, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n")
	}()

	connection, err := wsDial(context.Background(), "ws://"+listener.Addr().String()+"/ws", testAPIKey)
	if err != nil {
		t.Fatalf("wsDial() error = %v", err)
	}
	if err := connection.Close(); err != nil {
		t.Fatalf("connection.Close() error = %v", err)
	}
}

type memoryConn struct {
	bytes.Buffer
	closed bool
}

func (c *memoryConn) Read([]byte) (int, error)         { return 0, io.EOF }
func (c *memoryConn) Close() error                     { c.closed = true; return nil }
func (c *memoryConn) LocalAddr() net.Addr              { return nil }
func (c *memoryConn) RemoteAddr() net.Addr             { return nil }
func (c *memoryConn) SetDeadline(time.Time) error      { return nil }
func (c *memoryConn) SetReadDeadline(time.Time) error  { return nil }
func (c *memoryConn) SetWriteDeadline(time.Time) error { return nil }

func serverFrame(opcode byte, payload []byte) []byte {
	return append([]byte{0x80 | opcode, byte(len(payload))}, payload...)
}

func extendedServerFrame(payload []byte) []byte {
	if len(payload) == 126 {
		frame := []byte{0x81, 126, 0, 126}
		return append(frame, payload...)
	}
	frame := []byte{0x81, 127, 0, 0, 0, 0, 0, 0, 0, byte(len(payload))}
	return append(frame, payload...)
}

var _ net.Conn = (*memoryConn)(nil)
