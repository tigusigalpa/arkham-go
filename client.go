package arkham

import (
	"net/http"
	"time"
)

// DefaultBaseURL is the default Arkham API base URL used by NewClient
// unless overridden with WithBaseURL.
const DefaultBaseURL = "https://api.arkm.com"

// DefaultWebSocketURL is the default Arkham WebSocket base URL.
const DefaultWebSocketURL = "wss://api.arkm.com"

// DefaultTimeout is the default HTTP client timeout.
const DefaultTimeout = 30 * time.Second

// DefaultMaxRetries is the default number of retry attempts for
// retryable responses (429, 5xx) before giving up.
const DefaultMaxRetries = 3

// DefaultBaseDelay is the default initial backoff delay for retries.
const DefaultBaseDelay = 500 * time.Millisecond

// DefaultUserAgent is the default User-Agent header value.
const DefaultUserAgent = "arkham-go/1.0.0"

// Client is the entry point of the Arkham Go SDK. It holds HTTP
// configuration and exposes one service struct per API resource family.
// A Client is safe for concurrent use by multiple goroutines.
type Client struct {
	apiKey     string
	baseURL    string
	wsBaseURL  string
	httpClient *http.Client
	maxRetries int
	baseDelay  time.Duration
	userAgent  string
	logger     Logger

	// Analytics provides access to analytics endpoints.
	Analytics *AnalyticsService
	// Arkham provides access to ARKM token endpoints.
	Arkham *ArkhamService
	// Balances provides access to balance endpoints.
	Balances *BalancesService
	// Chains provides access to chain endpoints.
	Chains *ChainsService
	// Cluster provides access to cluster endpoints.
	Cluster *ClusterService
	// Counterparties provides access to counterparty endpoints.
	Counterparties *CounterpartiesService
	// Flow provides access to flow endpoints.
	Flow *FlowService
	// History provides access to history endpoints.
	History *HistoryService
	// Hypercore provides access to HyperCore endpoints.
	Hypercore *HypercoreService
	// Intelligence provides access to intelligence endpoints.
	Intelligence *IntelligenceService
	// Loans provides access to loan endpoints.
	Loans *LoansService
	// MarketData provides access to market data endpoints.
	MarketData *MarketDataService
	// Networks provides access to network endpoints.
	Networks *NetworksService
	// Polymarket provides access to Polymarket endpoints.
	Polymarket *PolymarketService
	// Portfolio provides access to portfolio endpoints.
	Portfolio *PortfolioService
	// Risk provides access to risk scoring endpoints (beta add-on).
	Risk *RiskService
	// Subscription provides access to subscription endpoints.
	Subscription *SubscriptionService
	// Swaps provides access to swap endpoints.
	Swaps *SwapsService
	// Tag provides access to tag endpoints.
	Tag *TagService
	// Token provides access to token endpoints.
	Token *TokenService
	// Transfers provides access to transfer endpoints.
	Transfers *TransfersService
	// Tx provides access to transaction endpoints.
	Tx *TxService
	// User provides access to user endpoints (alerts, entities, labels).
	User *UserService
	// Streams provides access to WebSocket v2 stream endpoints.
	Streams *StreamsService
}

// NewClient creates a new Arkham API client authenticated with the given
// apiKey. The key is sent in the API-Key header on every request.
// Behaviour can be customized via Option values such as WithBaseURL,
// WithHTTPClient, WithTimeout, WithMaxRetries, and WithUserAgent.
//
// Returns an error if apiKey is empty.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}
	c := &Client{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		wsBaseURL:  DefaultWebSocketURL,
		httpClient: &http.Client{Timeout: DefaultTimeout},
		maxRetries: DefaultMaxRetries,
		baseDelay:  DefaultBaseDelay,
		userAgent:  DefaultUserAgent,
		logger:     nopLogger{},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.initServices()
	return c, nil
}

// NewClientFromEnv reads the API key from the ARKHAM_API_KEY environment
// variable and returns a new Arkham client.
func NewClientFromEnv(opts ...Option) (*Client, error) {
	apiKey := osLookupEnv("ARKHAM_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}
	return NewClient(apiKey, opts...)
}

// initServices creates all service structs and wires them to the client.
func (c *Client) initServices() {
	c.Analytics = &AnalyticsService{client: c}
	c.Arkham = &ArkhamService{client: c}
	c.Balances = &BalancesService{client: c}
	c.Chains = &ChainsService{client: c}
	c.Cluster = &ClusterService{client: c}
	c.Counterparties = &CounterpartiesService{client: c}
	c.Flow = &FlowService{client: c}
	c.History = &HistoryService{client: c}
	c.Hypercore = &HypercoreService{client: c}
	c.Intelligence = &IntelligenceService{client: c}
	c.Loans = &LoansService{client: c}
	c.MarketData = &MarketDataService{client: c}
	c.Networks = &NetworksService{client: c}
	c.Polymarket = &PolymarketService{client: c}
	c.Portfolio = &PortfolioService{client: c}
	c.Risk = &RiskService{client: c}
	c.Subscription = &SubscriptionService{client: c}
	c.Swaps = &SwapsService{client: c}
	c.Tag = &TagService{client: c}
	c.Token = &TokenService{client: c}
	c.Transfers = &TransfersService{client: c}
	c.Tx = &TxService{client: c}
	c.User = &UserService{client: c}
	c.Streams = &StreamsService{client: c}
}
