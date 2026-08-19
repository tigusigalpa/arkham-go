package arkham

import (
	"context"
	"encoding/json"
	"net/url"
)

// BalancesService provides access to balance endpoints.
type BalancesService struct {
	client *Client
}

// AddressBalances retrieves token balances for an address.
// Path: GET /balances/address/{address}
func (s *BalancesService) Address(ctx context.Context, address string, filter *ChainsFilter) (*AddressBalancesResponse, *ResponseMetadata, error) {
	q := url.Values{}
	filter.ApplyToValues(q)
	var out AddressBalancesResponse
	meta, err := s.client.get(ctx, "/balances/address/"+pathEscape(address), q, &out)
	return &out, meta, err
}

// EntityBalances retrieves token balances for an entity.
// Path: GET /balances/entity/{entity}
func (s *BalancesService) Entity(ctx context.Context, entity string, filter *ChainsFilter) (*EntityBalancesResponse, *ResponseMetadata, error) {
	q := url.Values{}
	filter.ApplyToValues(q)
	var out EntityBalancesResponse
	meta, err := s.client.get(ctx, "/balances/entity/"+pathEscape(entity), q, &out)
	return &out, meta, err
}

// SolanaSubaccountsAddress retrieves Solana subaccount balances for addresses.
// Path: GET /balances/solana/subaccounts/address/{addresses}
func (s *BalancesService) SolanaSubaccountsAddress(ctx context.Context, addresses string) (json.RawMessage, *ResponseMetadata, error) {
	var out json.RawMessage
	meta, err := s.client.get(ctx, "/balances/solana/subaccounts/address/"+pathEscape(addresses), nil, &out)
	return out, meta, err
}

// SolanaSubaccountsEntity retrieves Solana subaccount balances for entities.
// Path: GET /balances/solana/subaccounts/entity/{entities}
func (s *BalancesService) SolanaSubaccountsEntity(ctx context.Context, entities string) (json.RawMessage, *ResponseMetadata, error) {
	var out json.RawMessage
	meta, err := s.client.get(ctx, "/balances/solana/subaccounts/entity/"+pathEscape(entities), nil, &out)
	return out, meta, err
}
