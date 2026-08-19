// Package arkham provides a Go SDK for the Arkham Intel API v1.1.0
// (https://api.arkm.com). It exposes service structs grouped by resource
// family, accessible as fields on Client.
//
// Authentication uses the API-Key header on every request. All
// network-facing methods accept context.Context as their first parameter
// and return typed responses alongside response metadata.
//
// Quick start:
//
//	client, err := arkham.NewClient(os.Getenv("ARKHAM_API_KEY"),
//	    arkham.WithTimeout(20*time.Second),
//	    arkham.WithMaxRetries(4),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	result, meta, err := client.Intelligence.Address(ctx, "0x28C6c06298d514Db089934071355E5743bf21d60", nil)
//
// WebSocket v2 streams are supported via client.Streams. Create a stream
// with a transfer filter, then connect to receive real-time transfers.
//
// Risk Scoring endpoints are beta and require a paid add-on. See
// https://arkm.com/llms/guides/risk-scoring-beta.md for details.
//
// Arkham's x402 pay-per-request protocol is a separate USDC payment
// interface and is not implemented in this SDK. See
// https://arkm.com/llms/guides/x402-pay-per-request-for-agents.md for
// information about that protocol.
package arkham
