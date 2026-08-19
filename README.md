# arkham-go

A production-quality Go SDK for the [Arkham Intel API](https://arkm.com).

## Features

- **Complete API coverage** — All documented REST endpoints and WebSocket v2 transfer streams
- **Strongly typed** — Idiomatic Go structs for every request and response
- **Context-aware** — All methods accept `context.Context` for cancellation and deadlines
- **Functional options** — Configure base URL, timeout, retries, HTTP client, logger
- **Automatic retry** — Exponential backoff with jitter on 429 and 5xx, respects `Retry-After`
- **Typed errors** — Sentinel errors (`ErrBadRequest`, `ErrUnauthorized`, `ErrRateLimited`, etc.) with `errors.Is`/`errors.As` support
- **Credit safety** — Response metadata preserves `X-Intel-Datapoints-*` headers
- **Offset pagination** — Caller-controlled paginator with max items/pages limits
- **WebSocket v2** — Full stream lifecycle: create, connect, receive, reconnect, delete
- **Zero dependencies** — Standard library only

## Installation

```bash
go get github.com/tigusigalpa/arkham-go
```

## Quick Start

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

    // Look up address intelligence
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

| Service | Endpoints |
|---------|-----------|
| `client.Analytics` | Credit periods, endpoint calls |
| `client.Arkham` | ARKM circulating supply |
| `client.Balances` | Address/entity balances, Solana subaccounts |
| `client.Chains` | Supported chains |
| `client.Cluster` | Cluster summaries |
| `client.Counterparties` | Address/entity counterparties |
| `client.Flow` | Historical USD flow |
| `client.History` | Historical balance snapshots |
| `client.Hypercore` | HyperCore markets, positions, trades |
| `client.Intelligence` | Address/entity intelligence, search, updates |
| `client.Loans` | Loan/borrow positions |
| `client.MarketData` | Altcoin index, trending tokens |
| `client.Networks` | Network status and history |
| `client.Polymarket` | Events, positions, orderbook, leaderboard |
| `client.Portfolio` | Historical portfolio snapshots |
| `client.Risk` | Risk scores (beta), batch, paths, sources |
| `client.Subscription` | Intel usage |
| `client.Swaps` | Token swaps |
| `client.Tag` | Tag summaries and updates |
| `client.Token` | Market data, price history, holders, volume |
| `client.Transfers` | Enriched/unenriched transfers, histograms, volume |
| `client.Tx` | Transaction details |
| `client.User` | Alerts, entities, labels |
| `client.Streams` | WebSocket v2 stream management + connection |

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

// 4. Delete when done
client.Streams.Delete(ctx, stream.StreamID)
```

## Error Handling

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
