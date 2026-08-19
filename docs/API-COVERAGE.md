# Arkham Go SDK — API Coverage

## REST Endpoints

| Endpoint | Method | Path | Service | Status |
|----------|--------|------|---------|--------|
| Intelligence — Address | GET | `/intelligence/address/{address}` | Intelligence | ✅ |
| Intelligence — Entity | GET | `/intelligence/entity/{entity}` | Intelligence | ✅ |
| Intelligence — Search | GET | `/intelligence/search` | Intelligence | ✅ |
| Intelligence — Updates | GET | `/intelligence/updates` | Intelligence | ✅ |
| Intelligence — Predictions | GET | `/intelligence/predictions` | Intelligence | ✅ |
| Intelligence — Balance Changes | GET | `/intelligence/balance-changes` | Intelligence | ✅ |
| Intelligence — Entity Types | GET | `/intelligence/entity-types` | Intelligence | ✅ |
| Balances — Address | GET | `/balances/address/{address}` | Balances | ✅ |
| Balances — Entity | GET | `/balances/entity/{entity}` | Balances | ✅ |
| Balances — Solana Subaccounts (Address) | GET | `/balances/solana/subaccounts/address/{addresses}` | Balances | ✅ |
| Balances — Solana Subaccounts (Entity) | GET | `/balances/solana/subaccounts/entity/{entities}` | Balances | ✅ |
| Transfers — Enriched | GET | `/transfers` | Transfers | ✅ |
| Transfers — Unenriched | GET | `/transfers/unenriched` | Transfers | ✅ |
| Transfers — Histogram | GET | `/transfers/histogram` | Transfers | ✅ |
| Transfers — Simple Histogram | GET | `/transfers/histogram/simple` | Transfers | ✅ |
| Volume — Address | GET | `/volume/address/{address}` | Transfers | ✅ |
| Swaps | GET | `/swaps` | Swaps | ✅ |
| TX — By Hash | GET | `/tx/{hash}` | Tx | ✅ |
| Counterparties — Address | GET | `/counterparties/address/{address}` | Counterparties | ✅ |
| Counterparties — Entity | GET | `/counterparties/entity/{entity}` | Counterparties | ✅ |
| Flow — Address | GET | `/flow/address/{address}` | Flow | ✅ |
| Flow — Entity | GET | `/flow/entity/{entity}` | Flow | ✅ |
| History — Address | GET | `/history/address/{address}` | History | ✅ |
| History — Entity | GET | `/history/entity/{entity}` | History | ✅ |
| Portfolio — Address | GET | `/portfolio/address/{address}` | Portfolio | ✅ |
| Portfolio — Entity | GET | `/portfolio/entity/{entity}` | Portfolio | ✅ |
| Chains | GET | `/chains` | Chains | ✅ |
| Networks — Status | GET | `/networks/status` | Networks | ✅ |
| Networks — History | GET | `/networks/{chain}/history` | Networks | ✅ |
| Token — Info | GET | `/token/{id}` | Token | ✅ |
| Token — Market | GET | `/token/market/{id}` | Token | ✅ |
| Token — Price History | GET | `/token/price-history/{id}` | Token | ✅ |
| Token — Price Change | GET | `/token/price-change/{id}` | Token | ✅ |
| Token — Holders | GET | `/token/holders/{id}` | Token | ✅ |
| Token — Volume | GET | `/token/volume/{id}` | Token | ✅ |
| Token — Addresses | GET | `/token/addresses/{id}` | Token | ✅ |
| Market Data — Altcoin Index | GET | `/marketdata/altcoin-index` | MarketData | ✅ |
| Market Data — Trending Tokens | GET | `/marketdata/tokens/trending` | MarketData | ✅ |
| Market Data — Token Top Flow | GET | `/marketdata/tokens/top-flow/{token}` | MarketData | ✅ |
| ARKM — Circulating | GET | `/arkm/circulating` | Arkham | ✅ |
| Cluster — Summary | GET | `/cluster/{id}/summary` | Cluster | ✅ |
| Tag — Summary | GET | `/tag/{id}/summary` | Tag | ✅ |
| Tag — Updates | GET | `/tag/updates` | Tag | ✅ |
| Tag — Address Updates | GET | `/tag/address-updates` | Tag | ✅ |
| Loans — Address | GET | `/loans/address/{address}` | Loans | ✅ |
| Loans — Entity | GET | `/loans/entity/{entity}` | Loans | ✅ |
| Risk — Address | GET | `/risk/address/{address}` | Risk | ✅ |
| Risk — Entity | GET | `/risk/entity/{entity}` | Risk | ✅ |
| Risk — Paths (Address) | GET | `/risk/paths/address/{address}` | Risk | ✅ |
| Risk — Paths (Entity) | GET | `/risk/paths/entity/{entity}` | Risk | ✅ |
| Risk — Sources (Address) | GET | `/risk/sources/address/{address}` | Risk | ✅ |
| Risk — Sources (Entity) | GET | `/risk/sources/entity/{entity}` | Risk | ✅ |
| Risk — Summary | GET | `/risk/summary/address/{address}` | Risk | ✅ |
| Risk — Batch Addresses | POST | `/risk/address/batch` | Risk | ✅ |
| Risk — Batch Entities | POST | `/risk/entity/batch` | Risk | ✅ |
| Subscription — Intel Usage | GET | `/subscription/intel-usage` | Subscription | ✅ |
| Analytics — Credit Periods | GET | `/analytics/credit-periods` | Analytics | ✅ |
| Analytics — Endpoint Calls | GET | `/analytics/endpoint-calls` | Analytics | ✅ |
| User — Alerts | GET | `/user/alerts` | User | ✅ |
| User — Create Alert | POST | `/user/alerts` | User | ✅ |
| User — Update Alert | PUT | `/user/alerts/{id}` | User | ✅ |
| User — Delete Alert | DELETE | `/user/alerts/{id}` | User | ✅ |
| User — Entities | GET | `/user/entities` | User | ✅ |
| User — Labels | GET | `/user/labels` | User | ✅ |
| User — Create Labels | POST | `/user/labels` | User | ✅ |
| Polymarket — Events | GET | `/polymarket/events` | Polymarket | ✅ |
| Polymarket — Activity | GET | `/polymarket/activity/{address}` | Polymarket | ✅ |
| Polymarket — Positions | GET | `/polymarket/positions/{address}` | Polymarket | ✅ |
| Polymarket — OrderBook | GET | `/polymarket/orderbook/{conditionId}` | Polymarket | ✅ |
| Polymarket — Price History | GET | `/polymarket/price-history/{conditionId}` | Polymarket | ✅ |
| Polymarket — PnL Chart | GET | `/polymarket/pnl-chart/{address}` | Polymarket | ✅ |
| Polymarket — Wallet Summary | GET | `/polymarket/wallet/summary/{address}` | Polymarket | ✅ |
| Polymarket — Wallet Stats | GET | `/polymarket/wallet/stats/{address}` | Polymarket | ✅ |
| Polymarket — Wallet Tags | GET | `/polymarket/wallet/tags/{address}` | Polymarket | ✅ |
| Polymarket — Leaderboard | GET | `/polymarket/leaderboard` | Polymarket | ✅ |
| Polymarket — Stats | GET | `/polymarket/stats` | Polymarket | ✅ |
| Polymarket — Top Events | GET | `/polymarket/top-events` | Polymarket | ✅ |
| Polymarket — Top Event Breakdown | GET | `/polymarket/top-events/{eventId}/breakdown` | Polymarket | ✅ |
| Polymarket — Top Holders | GET | `/polymarket/top-holders/{conditionId}` | Polymarket | ✅ |
| Polymarket — Wallet Event History | GET | `/polymarket/wallet/event-history/{address}` | Polymarket | ✅ |
| Polymarket — Prediction History | GET | `/polymarket/prediction-history/{conditionId}` | Polymarket | ✅ |
| Hypercore — Markets | GET | `/hypercore/markets` | Hypercore | ✅ |
| Hypercore — Active | GET | `/hypercore/active/{address}` | Hypercore | ✅ |
| Hypercore — Summary | GET | `/hypercore/summary/{address}` | Hypercore | ✅ |
| Hypercore — Spot Balances | GET | `/hypercore/spot/balances/{address}` | Hypercore | ✅ |
| Hypercore — Perp Positions | GET | `/hypercore/perp/positions/{address}` | Hypercore | ✅ |
| Hypercore — Portfolio History | GET | `/hypercore/portfolio/history/{address}` | Hypercore | ✅ |
| Hypercore — Subaccounts | GET | `/hypercore/subaccounts/{address}` | Hypercore | ✅ |
| Hypercore — Trades | GET | `/hypercore/trades/{address}` | Hypercore | ✅ |
| Hypercore — Token Top Perps | GET | `/hypercore/token/top-perps/{token}` | Hypercore | ✅ |
| Hypercore — Token Top Spots | GET | `/hypercore/token/top-spots/{token}` | Hypercore | ✅ |
| Hypercore — Token Holders | GET | `/hypercore/token/holders/{token}` | Hypercore | ✅ |
| Hypercore — Token Trades | GET | `/hypercore/token/trades/{token}` | Hypercore | ✅ |

## WebSocket v2 Endpoints

| Endpoint | Method | Path | Service | Status |
|----------|--------|------|---------|--------|
| Create Stream | POST | `/ws/v2/streams` | Streams | ✅ |
| List Streams | GET | `/ws/v2/streams` | Streams | ✅ |
| Delete Stream | DELETE | `/ws/v2/streams/{id}` | Streams | ✅ |
| Connect | WS | `wss://api.arkm.com/ws/v2/transfers?stream_id={id}` | Streams | ✅ |

## Notes

- Risk Scoring endpoints require a paid beta add-on
- WebSocket v2: max 10 active streams, 2 credits per transfer delivered
- Heavy endpoints (counterparties, intelligence search) limited to 1 req/s
- Standard endpoints: 20 req/s
