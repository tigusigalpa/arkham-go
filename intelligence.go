package arkham

import (
	"context"
	"net/url"
	"strconv"
)

// IntelligenceService provides access to intelligence endpoints.
type IntelligenceService struct {
	client *Client
}

// Address retrieves intelligence data for a blockchain address.
// Path: GET /intelligence/address/{address}
func (s *IntelligenceService) Address(ctx context.Context, address string, chains []string) (*Address, *ResponseMetadata, error) {
	q := url.Values{}
	if len(chains) > 0 {
		q.Set("chain", joinChains(chains))
	}
	var out Address
	meta, err := s.client.get(ctx, "/intelligence/address/"+pathEscape(address), q, &out)
	return &out, meta, err
}

// Entity retrieves intelligence data for an entity.
// Path: GET /intelligence/entity/{entity}
func (s *IntelligenceService) Entity(ctx context.Context, entity string) (*Entity, *ResponseMetadata, error) {
	var out Entity
	meta, err := s.client.get(ctx, "/intelligence/entity/"+pathEscape(entity), nil, &out)
	return &out, meta, err
}

// Search searches for entities, addresses, tokens, ENS names, and more.
// Path: GET /intelligence/search
func (s *IntelligenceService) Search(ctx context.Context, query string, opts *SearchOptions) (*SearchResults, *ResponseMetadata, error) {
	q := url.Values{}
	q.Set("query", query)
	if opts != nil {
		opts.ApplyToValues(q)
	}
	var out SearchResults
	meta, err := s.client.get(ctx, "/intelligence/search", q, &out)
	return &out, meta, err
}

// SearchOptions holds optional parameters for the search endpoint.
type SearchOptions struct {
	ArkhamEntities        int
	ArkhamAddresses       int
	UserEntities          int
	UserAddresses         int
	ENS                   int
	Types                 int
	Services              int
	Twitter               int
	Opensea               int
	Tokens                int
	Pools                 int
	Tags                  int
	PolymarketEvents      int
	ArkhamEntitiesOffset  int
	ArkhamAddressesOffset int
	UserEntitiesOffset    int
	UserAddressesOffset   int
	EnsOffset             int
	TypesOffset           int
	ServicesOffset        int
	TwitterOffset         int
	OpenseaOffset         int
	TokensOffset          int
	PoolsOffset           int
	TagsOffset            int
}

// ApplyToValues adds search parameters to the given url.Values.
func (o *SearchOptions) ApplyToValues(q url.Values) {
	if o == nil {
		return
	}
	setIfPositive := func(key string, val int) {
		if val > 0 {
			q.Set(key, strconv.Itoa(val))
		}
	}
	setIfPositive("arkhamEntities", o.ArkhamEntities)
	setIfPositive("arkhamAddresses", o.ArkhamAddresses)
	setIfPositive("userEntities", o.UserEntities)
	setIfPositive("userAddresses", o.UserAddresses)
	setIfPositive("ens", o.ENS)
	setIfPositive("types", o.Types)
	setIfPositive("services", o.Services)
	setIfPositive("twitter", o.Twitter)
	setIfPositive("opensea", o.Opensea)
	setIfPositive("tokens", o.Tokens)
	setIfPositive("pools", o.Pools)
	setIfPositive("tags", o.Tags)
	setIfPositive("polymarketEvents", o.PolymarketEvents)
	setIfPositive("arkhamEntitiesOffset", o.ArkhamEntitiesOffset)
	setIfPositive("arkhamAddressesOffset", o.ArkhamAddressesOffset)
	setIfPositive("userEntitiesOffset", o.UserEntitiesOffset)
	setIfPositive("userAddressesOffset", o.UserAddressesOffset)
	setIfPositive("ensOffset", o.EnsOffset)
	setIfPositive("typesOffset", o.TypesOffset)
	setIfPositive("servicesOffset", o.ServicesOffset)
	setIfPositive("twitterOffset", o.TwitterOffset)
	setIfPositive("openseaOffset", o.OpenseaOffset)
	setIfPositive("tokensOffset", o.TokensOffset)
	setIfPositive("poolsOffset", o.PoolsOffset)
	setIfPositive("tagsOffset", o.TagsOffset)
}

// Updates retrieves recent intelligence updates.
// Path: GET /intelligence/updates
func (s *IntelligenceService) Updates(ctx context.Context) ([]IntelligenceUpdate, *ResponseMetadata, error) {
	var out []IntelligenceUpdate
	meta, err := s.client.get(ctx, "/intelligence/updates", nil, &out)
	return out, meta, err
}

// Predictions retrieves entity predictions.
// Path: GET /intelligence/predictions
func (s *IntelligenceService) Predictions(ctx context.Context) ([]EntityPrediction, *ResponseMetadata, error) {
	var out []EntityPrediction
	meta, err := s.client.get(ctx, "/intelligence/predictions", nil, &out)
	return out, meta, err
}

// BalanceChanges retrieves entity balance changes.
// Path: GET /intelligence/balance-changes
func (s *IntelligenceService) BalanceChanges(ctx context.Context) ([]EntityBalanceChange, *ResponseMetadata, error) {
	var out []EntityBalanceChange
	meta, err := s.client.get(ctx, "/intelligence/balance-changes", nil, &out)
	return out, meta, err
}

// EntityTypes retrieves entity types.
// Path: GET /intelligence/entity-types
func (s *IntelligenceService) EntityTypes(ctx context.Context) ([]EntityType, *ResponseMetadata, error) {
	var out []EntityType
	meta, err := s.client.get(ctx, "/intelligence/entity-types", nil, &out)
	return out, meta, err
}
