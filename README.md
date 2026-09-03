# Arkham Golang Client/SDK/Library

![Arkham Intel Golang library](https://i.postimg.cc/yYPnLb5s/arkham-golang-github.jpg)

[![CI](https://github.com/tigusigalpa/arkham-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/arkham-go/actions/workflows/ci.yml)
[![Tests](https://github.com/tigusigalpa/arkham-go/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/arkham-go/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![CodeQL](https://github.com/tigusigalpa/arkham-go/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/arkham-go/actions/workflows/codeql.yml)
[![Codecov](https://codecov.io/gh/tigusigalpa/arkham-go/graph/badge.svg)](https://codecov.io/gh/tigusigalpa/arkham-go)
[![GitHub Release](https://img.shields.io/github/v/release/tigusigalpa/arkham-go?style=flat-square)](https://github.com/tigusigalpa/arkham-go/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://pkg.go.dev/github.com/tigusigalpa/arkham-go)

A friendly, production-ready Go SDK for the [Arkham Intel API](https://arkm.com).

It handles the REST endpoints and WebSocket v2 streams for you, so you can focus on building with on-chain intelligence instead of wrestling with HTTP plumbing.

## What you get

- **Full API coverage** — every documented REST endpoint and WebSocket v2 transfer stream, ready to use.
- **Idiomatic Go** — strongly typed structs for every request and response, so your IDE can help you.
- **Context support** — every method accepts `context.Context` for cancellation and deadlines.
- **Flexible configuration** — tune the base URL, timeout, retries, HTTP client, and logger with functional options.
- **Smart retries** — exponential backoff with jitter on `429` and `5xx`, and it respects `Retry-After`.
- **Typed errors** — sentinel errors like `ErrBadRequest`, `ErrUnauthorized`, and `ErrRateLimited`, all compatible with `errors.Is` and `errors.As`.
- **Credit tracking** — response metadata keeps the `X-Intel-Datapoints-*` headers close at hand.
- **Pagination that behaves** — caller-controlled offset pagination with limits on items and pages.
- **WebSocket v2 made simple** — create, connect, receive, reconnect, and delete streams without fuss.
- **No external dependencies** — everything is built on the Go standard library.

## Installation

Requirements: Go 1.21 or newer and an Arkham API key. Create and manage a key
in your Arkham account, then keep it out of source control:

```bash
export ARKHAM_API_KEY="your-api-key"
```

On PowerShell:

```powershell
$env:ARKHAM_API_KEY = "your-api-key"
```

Install the module in your application:

```bash
go get github.com/tigusigalpa/arkham-go
```

The SDK has no runtime dependencies outside the Go standard library.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/tigusigalpa/arkham-go"
)

func main() {
    client, err := arkham.NewClient(os.Getenv("ARKHAM_API_KEY"))
    if err != nil {
        log.Fatal(err)
    }

    // Look up intelligence for an address
    addr, meta, err := client.Intelligence.Address(
        context.Background(),
        "0x28C6c06298d514Db089934071355E5743bf21d60",
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }

    if addr.ArkhamEntity != nil {
        fmt.Printf("Entity: %s\n", addr.ArkhamEntity.Name)
    }
    fmt.Printf("Intel:  %s/%s datapoints\n", meta.IntelDatapointsUsage, meta.IntelDatapointsLimit)
}
```

`NewClientFromEnv` is useful when the application follows the conventional
`ARKHAM_API_KEY` environment variable:

```go
client, err := arkham.NewClientFromEnv()
if err != nil {
    return err
}
```

Every network method takes a `context.Context`. Use a deadline on work that
must not outlive a request window:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

address, _, err := client.Intelligence.Address(ctx, "0x28C6c06298d514Db089934071355E5743bf21d60", nil)
```

## Configuration

Fine-tune the client when you create it:

```go
client, err := arkham.NewClient(
    apiKey,
    arkham.WithBaseURL("https://api.arkm.com"),
    arkham.WithTimeout(30 * time.Second),
    arkham.WithMaxRetries(3),
    arkham.WithBaseDelay(500 * time.Millisecond),
    arkham.WithUserAgent("my-app/1.0"),
)
```

`WithMaxRetries` is the number of additional attempts after the first request.
GET requests retry `429` and `5xx` responses by default. The retry delay
honours `Retry-After` (seconds or an HTTP date); otherwise it uses exponential
backoff with jitter. Mutating requests do not retry unless the SDK knows the
endpoint is safe to retry.

When supplying `WithHTTPClient`, configure transport behaviour there; a later
`WithTimeout` option sets that client's `Timeout` field.

## Services

| Service | What it covers |
|---------|----------------|
| `client.Analytics` | Credit periods and endpoint calls |
| `client.Arkham` | ARKM circulating supply |
| `client.Balances` | Address/entity balances and Solana subaccounts |
| `client.Chains` | Supported chains |
| `client.Cluster` | Cluster summaries |
| `client.Counterparties` | Address/entity counterparties |
| `client.Flow` | Historical USD flow |
| `client.History` | Historical balance snapshots |
| `client.Hypercore` | HyperCore markets, positions, and trades |
| `client.Intelligence` | Address/entity intelligence, search, and updates |
| `client.Loans` | Loan and borrow positions |
| `client.MarketData` | Altcoin index and trending tokens |
| `client.Networks` | Network status and history |
| `client.Polymarket` | Events, positions, orderbook, and leaderboard |
| `client.Portfolio` | Historical portfolio snapshots |
| `client.Risk` | Risk scores (beta), batch checks, paths, and sources |
| `client.Subscription` | Intel usage |
| `client.Swaps` | Token swaps |
| `client.Tag` | Tag summaries and updates |
| `client.Token` | Market data, price history, holders, and volume |
| `client.Transfers` | Enriched/unenriched transfers, histograms, and volume |
| `client.Tx` | Transaction details |
| `client.User` | Alerts, entities, and labels |
| `client.Streams` | WebSocket v2 stream management and connections |

For the endpoint-by-endpoint list, see [API coverage](docs/API-COVERAGE.md).

## Fetching transfers

Filters are expressed with typed option structs and converted to query
parameters by the SDK. Values such as `UsdGte` are strings to preserve the
decimal representation expected by the API.

```go
filter := &arkham.TransferFilter{
    Base:    []string{"binance"},
    Chains:  []string{"ethereum", "bitcoin"},
    Flow:    arkham.FlowOut,
    UsdGte:  "100000",
    SortKey: arkham.SortKeyTime,
    SortDir: arkham.SortDirDesc,
    Limit:   25,
    TimeRange: &arkham.TimeRange{
        TimeLast: "24h",
    },
}

transfers, meta, err := client.Transfers.Transfers(ctx, filter)
if err != nil {
    return err
}
for _, transfer := range transfers.Transfers {
    fmt.Printf("%s → %s: $%s\n", transfer.From, transfer.To, transfer.USD)
}
fmt.Println("Datapoints remaining:", meta.IntelDatapointsRemaining)
```

`TimeLast` cannot be combined with `TimeGte` or `TimeLte`. Call `Validate` if
you construct filters before issuing a request, or rely on service methods that
validate their relevant options.

## Pagination

`Paginator` is useful for offset-based list endpoints. It fetches one page at
a time; it does not infer when the remote list is exhausted, so stop when the
decoded response is empty. `maxItems` caps the requested total, and the final
request's `limit` is reduced to the remaining number of items.

```go
query := url.Values{"from": {"0x28C6c06298d514Db089934071355E5743bf21d60"}}
pages := arkham.NewPaginator(ctx, client, "/transfers", query,
    100,  // items per page
    500,  // at most 500 items requested (0 means unlimited)
    10,   // at most 10 requests (0 means unlimited)
)

for pages.HasNext() {
    var page []arkham.Transfer
    if _, err := pages.NextPage(&page); err != nil {
        return err
    }
    if len(page) == 0 {
        break
    }
    // Process page before requesting the next one.
}
```

A non-positive page size uses `arkham.DefaultPageSize` (100), preventing
repeated requests for offset zero.

## WebSocket v2

```go
ctx := context.Background()

// 1. Create a stream
stream, _, err := client.Streams.Create(ctx, &arkham.CreateStreamV2Request{
    Base:   []string{"binance"},
    UsdGte: "500000",
})

// 2. Connect
conn, err := client.Streams.Connect(ctx, stream.StreamID)
if err != nil {
    return err
}
defer conn.Close()

// 3. Receive transfers
for {
    msg, err := conn.ReceiveTyped()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(msg.Type, string(msg.Payload))
}

// 4. Delete when you are done
client.Streams.Delete(ctx, stream.StreamID)
```

The stream filter must include `Base`, `From`, `To`, or `Tokens`; alternatively
it may use `UsdGte` of at least `250000`. Delete streams you no longer use.
If a connection drops, call `conn.Reconnect(ctx)` within the API's reactivation
window and continue receiving with `Receive` or `ReceiveTyped`.

Runnable programs are available in
[`examples/`](examples): `intelligence`, `transfers`, `pagination`, and
`websocket`. For example:

```bash
go run ./examples/transfers
```

## Response metadata and errors

Each successful service call returns `*arkham.ResponseMetadata` alongside its
decoded value. It contains the HTTP status, request duration, final URL, and
the `X-Intel-Datapoints-*` headers. This makes it straightforward to record
credit consumption without parsing headers yourself.

## Error handling

```go
addr, _, err := client.Intelligence.Address(ctx, "0xabc", nil)
if err != nil {
    if errors.Is(err, arkham.ErrRateLimited) {
        // Handle rate limit
    }
    var apiErr *arkham.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("Status: %d, Message: %s\n", apiErr.StatusCode, apiErr.Message)
    }
}
```

`APIError` also exposes `StatusCode`, a short response-body excerpt,
`RetryAfter`, and the datapoint headers. Transport failures and cancelled
contexts are wrapped with `%w`, so standard `errors.Is` / `errors.As` handling
works throughout.

## Development

Run the test suite and static checks from the module root:

```bash
go test ./...
go vet ./...
```

## License

MIT
