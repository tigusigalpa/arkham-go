# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-01

### Added
- Initial release of arkham-go SDK
- Complete coverage of Arkham Intel REST API (80+ endpoints)
- WebSocket v2 transfer stream lifecycle (create, connect, receive, reconnect, delete)
- 24 typed service structs: Analytics, Arkham, Balances, Chains, Cluster, Counterparties,
  Flow, History, Hypercore, Intelligence, Loans, MarketData, Networks, Polymarket,
  Portfolio, Risk, Subscription, Swaps, Tag, Token, Transfers, Tx, User, Streams
- Functional options for client configuration
- Automatic retry with exponential backoff and jitter on 429/5xx
- Typed errors with `errors.Is`/`errors.As` support
- Response metadata with intel datapoint headers
- Offset-based paginator with max items/pages limits
- Context-aware requests with cancellation support
- Zero external dependencies (standard library only)
- 45+ unit tests with httptest mock server
- 4 examples: intelligence, transfers, websocket, pagination
- GitHub Actions CI workflow
- API coverage documentation
