# Arkham Golang library

![Arkham Intel Golang library](https://i.postimg.cc/yYPnLb5s/arkham-golang-github.jpg)

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

```bash
go get github.com/tigusigalpa/arkham-go
```

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

    fmt.Printf("Entity: %s\n", addr.ArkhamEntity.Name)
    fmt.Printf("Intel:  %s/%s datapoints\n", meta.IntelDatapointsUsage, meta.IntelDatapointsLimit)
}
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

## WebSocket v2

```go
// 1. Create a stream
stream, _, err := client.Streams.Create(ctx, &arkham.CreateStreamV2Request{
    Base:   []string{"binance"},
    UsdGte: "500000",
})

// 2. Connect
conn, err := client.Streams.Connect(ctx, stream.StreamID)
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

## License

MIT
