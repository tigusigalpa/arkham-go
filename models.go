package arkham

import (
	"encoding/json"
)

// Address represents intelligence data for a blockchain address.
type Address struct {
	Address       string        `json:"address"`
	ArkhamEntity  *Entity       `json:"arkhamEntity,omitempty"`
	ArkhamLabel   *AddressLabel `json:"arkhamLabel,omitempty"`
	Chain         string        `json:"chain,omitempty"`
	Contract      bool          `json:"contract"`
	IsUserAddress bool          `json:"isUserAddress"`
}

// AddressLabel represents a label on an address.
type AddressLabel struct {
	Address   string `json:"address"`
	ChainType string `json:"chainType"`
	Name      string `json:"name"`
}

// Entity represents an Arkham entity.
type Entity struct {
	Crunchbase    string         `json:"crunchbase,omitempty"`
	Customized    bool           `json:"customized"`
	ID            string         `json:"id"`
	Instagram     string         `json:"instagram,omitempty"`
	LinkedIn      string         `json:"linkedin,omitempty"`
	LinkShareable bool           `json:"linkShareable,omitempty"`
	Name          string         `json:"name"`
	Note          string         `json:"note"`
	PopulatedTags []PopulatedTag `json:"populatedTags,omitempty"`
	Service       *bool          `json:"service"`
	Twitter       string         `json:"twitter,omitempty"`
	Type          string         `json:"type,omitempty"`
	Website       string         `json:"website,omitempty"`
}

// PopulatedTag represents a tag on an entity.
type PopulatedTag struct {
	Chain           string `json:"chain,omitempty"`
	DisablePage     bool   `json:"disablePage"`
	ExcludeEntities bool   `json:"excludeEntities"`
	ID              string `json:"id"`
	Label           string `json:"label"`
	Rank            int    `json:"rank,omitempty"`
	TagParams       string `json:"tagParams,omitempty"`
}

// EntitySummary represents summary statistics for an entity.
type EntitySummary struct {
	TotalBalance    float64 `json:"totalBalance"`
	TotalBalance24h float64 `json:"totalBalance24hAgo"`
	NumAddresses    int     `json:"numAddresses"`
	NumChains       int     `json:"numChains"`
}

// TokenBalance represents a token balance.
type TokenBalance struct {
	Balance           float64 `json:"balance"`
	BalanceExact      string  `json:"balanceExact"`
	EthereumAddress   string  `json:"ethereumAddress,omitempty"`
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	PriceChange24h    float64 `json:"priceChange24h"`
	PriceChange24hPct float64 `json:"priceChange24hPercent"`
	QuoteTime         string  `json:"quoteTime,omitempty"`
	Symbol            string  `json:"symbol"`
	USD               float64 `json:"usd"`
}

// AddressBalancesResponse represents the response from GET /balances/address/{address}.
type AddressBalancesResponse struct {
	Addresses       map[string]map[string]Address `json:"addresses"`
	Balances        map[string][]TokenBalance     `json:"balances"`
	TotalBalance    map[string]float64            `json:"totalBalance"`
	TotalBalance24h map[string]float64            `json:"totalBalance24hAgo"`
}

// EntityBalancesResponse represents the response from GET /balances/entity/{entity}.
type EntityBalancesResponse = AddressBalancesResponse

// Chain represents a supported blockchain.
type Chain struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ChainType   string `json:"chainType"`
	Explorer    string `json:"explorer,omitempty"`
	Logo        string `json:"logo,omitempty"`
	NativeToken string `json:"nativeToken,omitempty"`
}

// Transfer represents a blockchain transfer.
type Transfer struct {
	Chain       string        `json:"chain"`
	Token       string        `json:"token,omitempty"`
	Value       string        `json:"value"`
	ValueExact  string        `json:"valueExact,omitempty"`
	USD         float64       `json:"usd,omitempty"`
	From        string        `json:"from"`
	To          string        `json:"to"`
	Time        string        `json:"time"`
	TxHash      string        `json:"txHash,omitempty"`
	BlockNumber int64         `json:"blockNumber,omitempty"`
	FromEntity  *Entity       `json:"fromEntity,omitempty"`
	ToEntity    *Entity       `json:"toEntity,omitempty"`
	FromLabel   *AddressLabel `json:"fromLabel,omitempty"`
	ToLabel     *AddressLabel `json:"toLabel,omitempty"`
}

// EnrichedTransfers represents the response from GET /transfers.
type EnrichedTransfers struct {
	Transfers []Transfer      `json:"transfers"`
	Count     int             `json:"count,omitempty"`
	Cursors   json.RawMessage `json:"cursors,omitempty"`
}

// UnenrichedTransfer represents a transfer without entity enrichment.
type UnenrichedTransfer struct {
	Chain       string `json:"chain"`
	Token       string `json:"token,omitempty"`
	Value       string `json:"value"`
	From        string `json:"from"`
	To          string `json:"to"`
	Time        string `json:"time"`
	TxHash      string `json:"txHash,omitempty"`
	BlockNumber int64  `json:"blockNumber,omitempty"`
}

// UnenrichedTransfersResponse represents the response from GET /transfers/unenriched.
type UnenrichedTransfersResponse struct {
	Transfers []UnenrichedTransfer `json:"transfers"`
	Count     int                  `json:"count,omitempty"`
}

// HistogramResponse represents a transfer histogram response.
type HistogramResponse struct {
	Bins []HistogramBin `json:"bins,omitempty"`
}

// HistogramBin represents a single histogram bin.
type HistogramBin struct {
	Time  string  `json:"time"`
	Count int     `json:"count,omitempty"`
	Value float64 `json:"value,omitempty"`
	USD   float64 `json:"usd,omitempty"`
}

// SimpleHistogramResponse represents a simple histogram response.
type SimpleHistogramResponse = HistogramResponse

// Swap represents a token swap.
type Swap struct {
	Chain     string  `json:"chain"`
	TokenIn   string  `json:"tokenIn,omitempty"`
	TokenOut  string  `json:"tokenOut,omitempty"`
	AmountIn  string  `json:"amountIn"`
	AmountOut string  `json:"amountOut"`
	USD       float64 `json:"usd,omitempty"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Time      string  `json:"time"`
	TxHash    string  `json:"txHash,omitempty"`
}

// SwapsResponse represents the response from GET /swaps.
type SwapsResponse struct {
	Swaps []Swap `json:"swaps"`
	Count int    `json:"count,omitempty"`
}

// Transaction represents a blockchain transaction.
type Transaction struct {
	Hash     string `json:"hash"`
	Chain    string `json:"chain"`
	Block    int64  `json:"block"`
	Time     string `json:"time"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Value    string `json:"value,omitempty"`
	GasUsed  string `json:"gasUsed,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Status   int    `json:"status,omitempty"`
}

// Counterparty represents a top counterparty.
type Counterparty struct {
	Entity    *Entity `json:"entity,omitempty"`
	Address   string  `json:"address,omitempty"`
	TotalUSD  float64 `json:"totalUsd"`
	Incoming  float64 `json:"incomingUsd,omitempty"`
	Outgoing  float64 `json:"outgoingUsd,omitempty"`
	Transfers int     `json:"transfers,omitempty"`
}

// CounterpartiesResponse represents the response from counterparty endpoints.
type CounterpartiesResponse struct {
	Counterparties []Counterparty `json:"counterparties"`
}

// FlowPoint represents a historical USD flow data point.
type FlowPoint struct {
	Time     string  `json:"time"`
	USD      float64 `json:"usd,omitempty"`
	Incoming float64 `json:"incomingUsd,omitempty"`
	Outgoing float64 `json:"outgoingUsd,omitempty"`
}

// VolumePoint represents a volume data point.
type VolumePoint struct {
	Time    string         `json:"time"`
	Volumes []VolumeDetail `json:"volumes,omitempty"`
}

// VolumeDetail represents volume detail for a single exchange.
type VolumeDetail struct {
	DepositedUSD float64 `json:"depositedUSD"`
	ExchangeID   string  `json:"exchangeID"`
	ExchangeName string  `json:"exchangeName"`
	WithdrawnUSD float64 `json:"withdrawnUSD"`
}

// HistoryPoint represents a historical data point.
type HistoryPoint struct {
	Time  string  `json:"time"`
	USD   float64 `json:"usd,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// TokenMarketData represents current market data for a token.
type TokenMarketData struct {
	AllTimeHigh       float64 `json:"allTimeHigh"`
	AllTimeLow        float64 `json:"allTimeLow"`
	CirculatingSupply float64 `json:"circulatingSupply"`
	FullyDilutedValue float64 `json:"fullyDilutedValue"`
	MarketCap         float64 `json:"marketCap"`
	MaxPrice24h       float64 `json:"maxPrice24h"`
	MinPrice24h       float64 `json:"minPrice24h"`
	Price             float64 `json:"price"`
	Price180dAgo      float64 `json:"price180dAgo"`
	Price24hAgo       float64 `json:"price24hAgo"`
	Price30dAgo       float64 `json:"price30dAgo"`
	Price7dAgo        float64 `json:"price7dAgo"`
	TotalSupply       float64 `json:"totalSupply"`
	TotalVolume       float64 `json:"totalVolume"`
}

// TokenPriceHistory represents a token price history point.
type TokenPriceHistory struct {
	Time  string  `json:"time"`
	Price float64 `json:"price"`
}

// TokenPriceChange represents a token price change response.
type TokenPriceChange struct {
	Price         float64 `json:"price"`
	PriceAtTime   float64 `json:"priceAtTime"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

// TokenHolder represents a top token holder.
type TokenHolder struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	USD     float64 `json:"usd,omitempty"`
	Entity  *Entity `json:"entity,omitempty"`
}

// TokenInfo represents basic token information.
type TokenInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// TokenTopFlow represents top flow data for a token.
type TokenTopFlow struct {
	Address   string  `json:"address"`
	USD       float64 `json:"usd"`
	Direction string  `json:"direction,omitempty"`
}

// TokenVolume represents token volume data.
type TokenVolume struct {
	Time   string  `json:"time"`
	Volume float64 `json:"volume"`
	USD    float64 `json:"usd,omitempty"`
}

// TokenAddress represents a chain address for a token.
type TokenAddress struct {
	Chain   string `json:"chain"`
	Address string `json:"address"`
}

// TrendingToken represents a trending token.
type TrendingToken struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,omitempty"`
}

// PortfolioSnapshot represents a portfolio snapshot for a token.
type PortfolioSnapshot struct {
	Balance float64 `json:"balance"`
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Symbol  string  `json:"symbol"`
	USD     float64 `json:"usd"`
}

// NetworkStatus represents the status of a blockchain network.
type NetworkStatus struct {
	Active            bool            `json:"active"`
	AsOf              json.RawMessage `json:"asOf,omitempty"`
	Chain             string          `json:"chain"`
	GasFee            *float64        `json:"gasFee"`
	MarketCap         float64         `json:"marketCap"`
	Price             float64         `json:"price"`
	PriceChange24h    float64         `json:"priceChange24h"`
	PriceChange24hPct float64         `json:"priceChange24hPercent"`
	SatFee            *int            `json:"satFee"`
	Tip               int64           `json:"tip"`
	TPS               float64         `json:"tps"`
	Transfers         int64           `json:"transfers"`
	Volume            float64         `json:"volume"`
}

// NetworkHistoryPoint represents a historical data point for a network.
type NetworkHistoryPoint struct {
	Time   string  `json:"time"`
	Price  float64 `json:"price,omitempty"`
	Tip    int64   `json:"tip,omitempty"`
	Volume float64 `json:"volume,omitempty"`
}

// ARKMCirculating represents the ARKM circulating supply response.
type ARKMCirculating struct {
	Circulating float64 `json:"circulating"`
	Total       float64 `json:"total,omitempty"`
}

// IntelUsage represents the intel usage response.
type IntelUsage struct {
	TotalCount  int            `json:"totalCount"`
	TotalLimit  int            `json:"totalLimit"`
	PeriodStart string         `json:"periodStart"`
	Chains      map[string]int `json:"chains,omitempty"`
}

// CreditUsagePeriod represents a billing period's credit usage.
type CreditUsagePeriod struct {
	PeriodStart  string `json:"periodStart"`
	PeriodEnd    string `json:"periodEnd"`
	PlanCredits  int    `json:"planCredits"`
	ExtraCredits int    `json:"extraCredits"`
	TotalCredits int    `json:"totalCredits"`
}

// EndpointCallAnalytics represents endpoint call analytics.
type EndpointCallAnalytics struct {
	Endpoint string `json:"endpoint"`
	Calls    int    `json:"calls"`
	Credits  int    `json:"credits,omitempty"`
}

// AltcoinIndex represents the altcoin index response.
type AltcoinIndex struct {
	Time  string  `json:"time"`
	Index float64 `json:"index"`
}

// ClusterSummary represents a cluster summary.
type ClusterSummary struct {
	TotalBalance float64 `json:"totalBalance"`
	NumAddresses int     `json:"numAddresses"`
	NumChains    int     `json:"numChains"`
}

// TagParams represents tag parameters.
type TagParams struct {
	Params string `json:"params,omitempty"`
}

// TagSummary represents tag summary statistics.
type TagSummary struct {
	TotalBalance float64 `json:"totalBalance,omitempty"`
	NumAddresses int     `json:"numAddresses,omitempty"`
}

// Label represents a user label.
type Label struct {
	Address   string `json:"address"`
	ChainType string `json:"chainType"`
	Name      string `json:"name"`
	Note      string `json:"note,omitempty"`
}

// Alert represents an alert.
type Alert struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Enabled       bool     `json:"enabled"`
	Base          []string `json:"base,omitempty"`
	From          []string `json:"from,omitempty"`
	To            []string `json:"to,omitempty"`
	Tokens        []string `json:"tokens,omitempty"`
	Chains        []string `json:"chains,omitempty"`
	UsdGte        *float64 `json:"usdGte,omitempty"`
	UsdLte        *float64 `json:"usdLte,omitempty"`
	ValueGte      *float64 `json:"valueGte,omitempty"`
	ValueLte      *float64 `json:"valueLte,omitempty"`
	AlertMethodID int      `json:"alertMethodId"`
	Description   string   `json:"description,omitempty"`
}

// AlertRequest represents the request body for creating/updating alerts.
type AlertRequest struct {
	Name          string   `json:"name"`
	Enabled       bool     `json:"enabled"`
	Base          []string `json:"base,omitempty"`
	From          []string `json:"from,omitempty"`
	To            []string `json:"to,omitempty"`
	Tokens        []string `json:"tokens,omitempty"`
	Chains        []string `json:"chains,omitempty"`
	UsdGte        *float64 `json:"usdGte,omitempty"`
	UsdLte        *float64 `json:"usdLte,omitempty"`
	ValueGte      *float64 `json:"valueGte,omitempty"`
	ValueLte      *float64 `json:"valueLte,omitempty"`
	AlertMethodID int      `json:"alertMethodId"`
	Description   string   `json:"description,omitempty"`
}

// UserEntity represents a private entity.
type UserEntity struct {
	Addresses       map[string][]string `json:"addresses,omitempty"`
	CreatedAt       string              `json:"createdAt,omitempty"`
	Crunchbase      string              `json:"crunchbase,omitempty"`
	CustomImageSlug string              `json:"customImageSlug,omitempty"`
	Customized      bool                `json:"customized"`
	ID              string              `json:"id"`
	Instagram       string              `json:"instagram,omitempty"`
	LinkShareable   bool                `json:"linkShareable,omitempty"`
	LinkedIn        string              `json:"linkedin,omitempty"`
	Name            string              `json:"name"`
	Note            string              `json:"note"`
	OwnerUID        string              `json:"ownerUID,omitempty"`
	PopulatedTags   []PopulatedTag      `json:"populatedTags,omitempty"`
	Service         *bool               `json:"service"`
	Twitter         string              `json:"twitter,omitempty"`
	Type            string              `json:"type,omitempty"`
	UpdatedAt       string              `json:"updatedAt,omitempty"`
	Website         string              `json:"website,omitempty"`
}

// SearchResults represents the response from GET /intelligence/search.
type SearchResults struct {
	ArkhamAddresses []SearchAddress `json:"arkhamAddresses,omitempty"`
	ArkhamEntities  []Entity        `json:"arkhamEntities,omitempty"`
	ENS             []ENSResult     `json:"ens,omitempty"`
	Opensea         []OpenSeaResult `json:"opensea,omitempty"`
	Pools           []PoolResult    `json:"pools,omitempty"`
	Services        []Entity        `json:"services,omitempty"`
	Tokens          []SearchToken   `json:"tokens,omitempty"`
	Twitter         []Entity        `json:"twitter,omitempty"`
}

// SearchAddress represents an address search result.
type SearchAddress struct {
	Address     string        `json:"address"`
	ArkhamLabel *AddressLabel `json:"arkhamLabel,omitempty"`
	Chain       string        `json:"chain,omitempty"`
}

// ENSResult represents an ENS name search result.
type ENSResult struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

// OpenSeaResult represents an OpenSea username search result.
type OpenSeaResult struct {
	Address      string  `json:"address"`
	ProfileImage *string `json:"profileImage"`
	Username     string  `json:"username"`
}

// PoolResult represents a pool search result.
type PoolResult struct {
	LiquidityUSD   float64 `json:"liquidityUsd"`
	MarketCapUSD   float64 `json:"marketCapUsd"`
	PoolAddress    string  `json:"poolAddress"`
	PriceUSD       float64 `json:"priceUsd"`
	PricingID      string  `json:"pricingID"`
	ProgramAddress string  `json:"programAddress"`
	TokenAddress   string  `json:"tokenAddress"`
	TokenName      string  `json:"tokenName"`
	TokenSymbol    string  `json:"tokenSymbol"`
	Volume1h       float64 `json:"volume1h"`
}

// SearchToken represents a token search result.
type SearchToken struct {
	Deployments []TokenDeployment `json:"deployments,omitempty"`
	Identifier  TokenIdentifier   `json:"identifier"`
	Name        string            `json:"name"`
	Price       float64           `json:"price,omitempty"`
	Price24hAgo float64           `json:"price24hAgo,omitempty"`
	Symbol      string            `json:"symbol"`
}

// TokenDeployment represents a token deployment.
type TokenDeployment struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
}

// TokenIdentifier represents a token identifier.
type TokenIdentifier struct {
	PricingID string `json:"pricingID"`
}

// EntityPrediction represents a prediction for an entity.
type EntityPrediction struct {
	EntityID   string  `json:"entityId,omitempty"`
	Prediction string  `json:"prediction,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

// EntityBalanceChange represents an entity balance change.
type EntityBalanceChange struct {
	EntityID string  `json:"entityId,omitempty"`
	Chain    string  `json:"chain,omitempty"`
	Token    string  `json:"token,omitempty"`
	Change   float64 `json:"change,omitempty"`
	USD      float64 `json:"usd,omitempty"`
	Time     string  `json:"time,omitempty"`
}

// EntityType represents an entity type.
type EntityType struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

// IntelligenceUpdate represents an intelligence update entry.
type IntelligenceUpdate struct {
	Address string `json:"address,omitempty"`
	Entity  string `json:"entity,omitempty"`
	Chain   string `json:"chain,omitempty"`
	Label   string `json:"label,omitempty"`
	Time    string `json:"time,omitempty"`
}

// TagUpdate represents a tag definition update.
type TagUpdate struct {
	ID    string `json:"id"`
	Label string `json:"label,omitempty"`
}

// AddressTagUpdate represents an address-tag association update.
type AddressTagUpdate struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
	TagID   string `json:"tagId,omitempty"`
	Time    string `json:"time,omitempty"`
}

// RiskScoreResponse represents a risk score response.
type RiskScoreResponse struct {
	ChainType                    string       `json:"chain_type"`
	Address                      string       `json:"address"`
	RiskLevel                    string       `json:"risk_level"`
	HackerScore                  int          `json:"hacker_score,omitempty"`
	SanctionsScore               int          `json:"sanctions_score,omitempty"`
	MaxScore                     int          `json:"max_score,omitempty"`
	Sanctioned1hopScore          int          `json:"sanctioned_1hop_score,omitempty"`
	GreatestRiskCategory         string       `json:"greatest_risk_category,omitempty"`
	MaxScoreForward              int          `json:"max_score_forward,omitempty"`
	GreatestRiskCategoryForward  string       `json:"greatest_risk_category_forward,omitempty"`
	MaxScoreBackward             int          `json:"max_score_backward,omitempty"`
	GreatestRiskCategoryBackward string       `json:"greatest_risk_category_backward,omitempty"`
	IsSeed                       bool         `json:"is_seed"`
	HopDistance                  int          `json:"hop_distance,omitempty"`
	MaxHopReached                int          `json:"max_hop_reached,omitempty"`
	RiskWeightedIncomingUSD      float64      `json:"risk_weighted_incoming_usd,omitempty"`
	RiskWeightedOutgoingUSD      float64      `json:"risk_weighted_outgoing_usd,omitempty"`
	TopSources                   []RiskSource `json:"top_sources,omitempty"`
	UpdatedAt                    string       `json:"updated_at,omitempty"`
}

// RiskSource represents a risk source.
type RiskSource struct {
	SeedAddress     string  `json:"seed_address"`
	RiskCategory    string  `json:"risk_category"`
	Direction       string  `json:"direction"`
	ContributionPct float64 `json:"contribution_pct"`
	ContributionUSD float64 `json:"contribution_usd"`
	HopDistance     int     `json:"hop_distance"`
	FirstTS         string  `json:"first_ts"`
	LastTS          string  `json:"last_ts"`
}

// RiskPathsResponse represents risk paths response.
type RiskPathsResponse struct {
	Paths json.RawMessage `json:"paths"`
}

// RiskBatchRequest represents a batch risk score request.
type RiskBatchRequest struct {
	Addresses []string `json:"addresses"`
}

// RiskEntityBatchRequest represents a batch entity risk request.
type RiskEntityBatchRequest struct {
	Entities []string `json:"entities"`
}

// StreamV2 represents a WebSocket v2 stream.
type StreamV2 struct {
	StreamID  string `json:"streamId"`
	ID        int    `json:"id"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateStreamV2Request represents the request body for creating a v2 stream.
type CreateStreamV2Request struct {
	Base           []string `json:"base,omitempty"`
	Chains         []string `json:"chains,omitempty"`
	Flow           string   `json:"flow,omitempty"`
	From           []string `json:"from,omitempty"`
	To             []string `json:"to,omitempty"`
	Counterparties []string `json:"counterparties,omitempty"`
	Tokens         []string `json:"tokens,omitempty"`
	UsdGte         string   `json:"usdGte,omitempty"`
	UsdLte         string   `json:"usdLte,omitempty"`
	ValueGte       string   `json:"valueGte,omitempty"`
	ValueLte       string   `json:"valueLte,omitempty"`
}

// DeleteStreamV2Response represents the delete stream response.
type DeleteStreamV2Response struct {
	Success bool `json:"success"`
}

// WSMessage represents a WebSocket v2 message envelope.
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// WSTransferPayload represents the payload of a transfer message.
type WSTransferPayload struct {
	Transfer Transfer `json:"transfer"`
	AlertID  int      `json:"alertId,omitempty"`
}

// WSErrorPayload represents the payload of an error message.
type WSErrorPayload struct {
	Code      string `json:"code"`
	Message   string `json:"message,omitempty"`
	LimitType string `json:"limitType,omitempty"`
	ResetIn   int    `json:"resetIn,omitempty"`
}

// PolymarketEvent represents a Polymarket event.
type PolymarketEvent struct {
	ID          string             `json:"id"`
	Slug        string             `json:"slug,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Active      bool               `json:"active"`
	Closed      bool               `json:"closed,omitempty"`
	EndDate     string             `json:"endDate,omitempty"`
	Volume      float64            `json:"volume,omitempty"`
	Liquidity   float64            `json:"liquidity,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
	Markets     []PolymarketMarket `json:"markets,omitempty"`
}

// PolymarketMarket represents a Polymarket market within an event.
type PolymarketMarket struct {
	ID       string  `json:"id"`
	Question string  `json:"question,omitempty"`
	Outcome  string  `json:"outcome,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Volume   float64 `json:"volume,omitempty"`
}

// PolymarketEventsResponse represents the response from GET /polymarket/events.
type PolymarketEventsResponse struct {
	Events []PolymarketEvent `json:"events"`
	Count  int               `json:"count,omitempty"`
}

// PolymarketActivity represents a Polymarket activity entry.
type PolymarketActivity struct {
	Type    string  `json:"type"`
	Address string  `json:"address,omitempty"`
	EventID string  `json:"eventId,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
	USD     float64 `json:"usd,omitempty"`
	Time    string  `json:"time,omitempty"`
}

// PolymarketPosition represents a Polymarket user position.
type PolymarketPosition struct {
	ConditionID string  `json:"conditionId"`
	Outcome     string  `json:"outcome,omitempty"`
	Size        float64 `json:"size"`
	Price       float64 `json:"price,omitempty"`
	USD         float64 `json:"usd,omitempty"`
}

// PolymarketOrderBookEntry represents an order book entry.
type PolymarketOrderBookEntry struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// PolymarketOrderBook represents an order book.
type PolymarketOrderBook struct {
	Bids []PolymarketOrderBookEntry `json:"bids,omitempty"`
	Asks []PolymarketOrderBookEntry `json:"asks,omitempty"`
}

// PolymarketLeaderboardEntry represents a leaderboard entry.
type PolymarketLeaderboardEntry struct {
	Address string  `json:"address"`
	PnL     float64 `json:"pnl"`
	Rank    int     `json:"rank,omitempty"`
}

// PolymarketStats represents Polymarket platform stats.
type PolymarketStats struct {
	TotalVolume   float64 `json:"totalVolume,omitempty"`
	TotalUsers    int     `json:"totalUsers,omitempty"`
	ActiveMarkets int     `json:"activeMarkets,omitempty"`
}

// PolymarketPriceHistory represents a price history point.
type PolymarketPriceHistory struct {
	Time  string  `json:"time"`
	Price float64 `json:"price"`
}

// PolymarketPnLChart represents a PnL chart point.
type PolymarketPnLChart struct {
	Time string  `json:"time"`
	PnL  float64 `json:"pnl"`
}

// PolymarketWalletSummary represents a wallet summary.
type PolymarketWalletSummary struct {
	Value float64 `json:"value,omitempty"`
}

// PolymarketWalletStats represents wallet trading stats.
type PolymarketWalletStats struct {
	TotalVolume float64 `json:"totalVolume,omitempty"`
	TotalTrades int     `json:"totalTrades,omitempty"`
	WinRate     float64 `json:"winRate,omitempty"`
}

// PolymarketWalletTags represents wallet tags.
type PolymarketWalletTags struct {
	Tags []string `json:"tags,omitempty"`
}

// PolymarketTopEvent represents a top event.
type PolymarketTopEvent struct {
	EventID string  `json:"eventId"`
	Title   string  `json:"title,omitempty"`
	Volume  float64 `json:"volume,omitempty"`
}

// PolymarketTopEventBreakdown represents a top event breakdown.
type PolymarketTopEventBreakdown struct {
	Address string  `json:"address"`
	Volume  float64 `json:"volume,omitempty"`
	PnL     float64 `json:"pnl,omitempty"`
}

// PolymarketTopHolder represents a top holder.
type PolymarketTopHolder struct {
	Address string  `json:"address"`
	Size    float64 `json:"size"`
	USD     float64 `json:"usd,omitempty"`
}

// PolymarketWalletEventHistory represents a wallet event history entry.
type PolymarketWalletEventHistory struct {
	EventID string  `json:"eventId,omitempty"`
	Type    string  `json:"type"`
	Time    string  `json:"time,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
}

// PolymarketPredictionHistory represents a prediction history entry.
type PolymarketPredictionHistory struct {
	ConditionID string  `json:"conditionId,omitempty"`
	Outcome     string  `json:"outcome,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Size        float64 `json:"size,omitempty"`
	Time        string  `json:"time,omitempty"`
}

// HypercoreMarket represents a HyperCore market.
type HypercoreMarket struct {
	Name      string `json:"name"`
	PricingID string `json:"pricingId"`
	Symbol    string `json:"symbol"`
}

// HypercoreMarketsResponse represents the response from GET /hypercore/markets.
type HypercoreMarketsResponse struct {
	Markets []HypercoreMarket `json:"markets"`
}

// HypercoreActiveResponse represents the active check response.
type HypercoreActiveResponse struct {
	Active bool `json:"active"`
}

// HypercoreSummary represents a HyperCore account summary.
type HypercoreSummary struct {
	AccountValue float64 `json:"accountValue,omitempty"`
	PnL          float64 `json:"pnl,omitempty"`
}

// HypercoreSpotBalance represents a HyperCore spot balance.
type HypercoreSpotBalance struct {
	Token   string  `json:"token,omitempty"`
	Balance float64 `json:"balance"`
	USD     float64 `json:"usd,omitempty"`
}

// HypercorePerpPosition represents a HyperCore perp position.
type HypercorePerpPosition struct {
	Market     string  `json:"market,omitempty"`
	Size       float64 `json:"size"`
	EntryPrice float64 `json:"entryPrice,omitempty"`
	USD        float64 `json:"usd,omitempty"`
	PnL        float64 `json:"pnl,omitempty"`
}

// HypercorePortfolioHistory represents a portfolio history point.
type HypercorePortfolioHistory struct {
	Time         string  `json:"time"`
	AccountValue float64 `json:"accountValue,omitempty"`
	PnL          float64 `json:"pnl,omitempty"`
}

// HypercoreSubaccount represents a HyperCore subaccount.
type HypercoreSubaccount struct {
	Name string `json:"name,omitempty"`
}

// HypercoreTrade represents a HyperCore trade.
type HypercoreTrade struct {
	Market string  `json:"market,omitempty"`
	Side   string  `json:"side,omitempty"`
	Size   float64 `json:"size"`
	Price  float64 `json:"price"`
	Time   string  `json:"time,omitempty"`
	USD    float64 `json:"usd,omitempty"`
}

// HypercoreTokenPosition represents a top perp position for a token.
type HypercoreTokenPosition struct {
	Address string  `json:"address,omitempty"`
	Size    float64 `json:"size"`
	USD     float64 `json:"usd,omitempty"`
	PnL     float64 `json:"pnl,omitempty"`
}

// LoanPosition represents a loan/borrow position.
type LoanPosition struct {
	Chain    string  `json:"chain,omitempty"`
	Protocol string  `json:"protocol,omitempty"`
	Token    string  `json:"token,omitempty"`
	Borrowed float64 `json:"borrowed,omitempty"`
	Supplied float64 `json:"supplied,omitempty"`
	USD      float64 `json:"usd,omitempty"`
}

// LoansResponse represents the response from loans endpoints.
type LoansResponse struct {
	Loans []LoanPosition `json:"loans,omitempty"`
}
