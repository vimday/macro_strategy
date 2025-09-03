package models

import "time"

// MarketType represents different market types
type MarketType string

const (
	MarketTypeAShare MarketType = "a_share"
	MarketTypeCrypto MarketType = "crypto"
	MarketTypeHKUS   MarketType = "hk_us"
	MarketTypeETF    MarketType = "etf"
	MarketTypeBond   MarketType = "bond"
	MarketTypeFuture MarketType = "future"
	MarketTypeOption MarketType = "option"
)

// AssetClass represents different asset classes
type AssetClass string

const (
	AssetClassEquity     AssetClass = "equity"
	AssetClassCrypto     AssetClass = "crypto"
	AssetClassCommodity  AssetClass = "commodity"
	AssetClassBond       AssetClass = "bond"
	AssetClassDerivative AssetClass = "derivative"
	AssetClassIndex      AssetClass = "index"
)

// Currency represents different currencies
type Currency string

const (
	CurrencyCNY Currency = "CNY"
	CurrencyUSD Currency = "USD"
	CurrencyHKD Currency = "HKD"
	CurrencyBTC Currency = "BTC"
	CurrencyETH Currency = "ETH"
)

// Index represents a market index or tradable asset
type Index struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Symbol       string                 `json:"symbol"`
	MarketType   MarketType             `json:"market_type"`
	AssetClass   AssetClass             `json:"asset_class"`
	Currency     Currency               `json:"currency"`
	Description  string                 `json:"description"`
	TradingHours *TradingHours          `json:"trading_hours,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// TradingHours represents trading hours for different markets
type TradingHours struct {
	Timezone    string `json:"timezone"`
	OpenTime    string `json:"open_time"`
	CloseTime   string `json:"close_time"`
	WeekendDays []int  `json:"weekend_days"` // 0=Sunday, 1=Monday, etc.
}

// OHLCV represents Open, High, Low, Close, Volume data with enhanced metadata
type OHLCV struct {
	Date     time.Time `json:"date"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`
	Amount   float64   `json:"amount,omitempty"`   // 成交额
	Turnover float64   `json:"turnover,omitempty"` // 换手率
	PctChg   float64   `json:"pct_chg,omitempty"`  // 涨跌幅
}

// MarketData represents historical market data for an asset with enhanced metadata
type MarketData struct {
	AssetID    string                 `json:"asset_id"`
	Symbol     string                 `json:"symbol"`
	MarketType MarketType             `json:"market_type"`
	AssetClass AssetClass             `json:"asset_class"`
	Currency   Currency               `json:"currency"`
	Data       []OHLCV                `json:"data"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	LastUpdate time.Time              `json:"last_update"`
}

// StrategyType represents different strategy types
type StrategyType string

const (
	StrategyTypeMonthlyRotation StrategyType = "monthly_rotation"
	StrategyTypeBuyAndHold      StrategyType = "buy_and_hold"
	StrategyTypeGridTrading     StrategyType = "grid_trading"
	StrategyTypeMeanReversion   StrategyType = "mean_reversion"
	StrategyTypeMomentum        StrategyType = "momentum"
	// Future strategy types can be added here
)

// StrategyConfig represents strategy configuration with enhanced flexibility
type StrategyConfig struct {
	Type           StrategyType           `json:"type"`
	Parameters     map[string]interface{} `json:"parameters"`
	Description    string                 `json:"description"`
	RiskManagement *RiskManagementConfig  `json:"risk_management,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// RiskManagementConfig represents risk management settings
type RiskManagementConfig struct {
	MaxPositionSize float64 `json:"max_position_size,omitempty"` // 最大仓位比例
	StopLoss        float64 `json:"stop_loss,omitempty"`         // 止损比例
	TakeProfit      float64 `json:"take_profit,omitempty"`       // 止盈比例
	MaxDrawdown     float64 `json:"max_drawdown,omitempty"`      // 最大回撤限制
	CommissionRate  float64 `json:"commission_rate,omitempty"`   // 手续费率
	SlippageRate    float64 `json:"slippage_rate,omitempty"`     // 滑点率
}

// MonthlyRotationParams represents parameters for monthly rotation strategy
type MonthlyRotationParams struct {
	BuyDaysBeforeMonthEnd   int `json:"buy_days_before_month_end"`
	SellDaysAfterMonthStart int `json:"sell_days_after_month_start"`
}

// BacktestRequest represents a backtest request with enhanced configuration
type BacktestRequest struct {
	AssetID       string                 `json:"asset_id"`           // 兼容原 IndexID
	IndexID       string                 `json:"index_id,omitempty"` // 向后兼容
	Strategy      StrategyConfig         `json:"strategy"`
	StartDate     time.Time              `json:"start_date"`
	EndDate       time.Time              `json:"end_date"`
	InitialCash   float64                `json:"initial_cash"`
	Benchmark     string                 `json:"benchmark,omitempty"`      // 基准指数
	RebalanceFreq string                 `json:"rebalance_freq,omitempty"` // 再平衡频率
	DataSource    string                 `json:"data_source,omitempty"`    // 数据源
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Trade represents a single trade
type Trade struct {
	Date       time.Time `json:"date"`
	Action     string    `json:"action"` // "buy" or "sell"
	Price      float64   `json:"price"`
	Quantity   float64   `json:"quantity"`
	Amount     float64   `json:"amount"`
	Commission float64   `json:"commission"`
}

// Position represents current position
type Position struct {
	Quantity     float64 `json:"quantity"`
	AvgPrice     float64 `json:"avg_price"`
	MarketValue  float64 `json:"market_value"`
	UnrealizedPL float64 `json:"unrealized_pl"`
}

// PerformanceMetrics represents strategy performance metrics
type PerformanceMetrics struct {
	TotalReturn       float64 `json:"total_return"`
	AnnualizedReturn  float64 `json:"annualized_return"`
	MaxDrawdown       float64 `json:"max_drawdown"`
	SharpeRatio       float64 `json:"sharpe_ratio"`
	SortinoRatio      float64 `json:"sortino_ratio"`
	Volatility        float64 `json:"volatility"`
	WinRate           float64 `json:"win_rate"`
	ProfitFactor      float64 `json:"profit_factor"`
	CalmarRatio       float64 `json:"calmar_ratio"`
	TotalTrades       int     `json:"total_trades"`
	WinningTrades     int     `json:"winning_trades"`
	LosingTrades      int     `json:"losing_trades"`
	AvgWinningTrade   float64 `json:"avg_winning_trade"`
	AvgLosingTrade    float64 `json:"avg_losing_trade"`
	MaxWinningTrade   float64 `json:"max_winning_trade"`
	MaxLosingTrade    float64 `json:"max_losing_trade"`
	MaxDrawdownPeriod int     `json:"max_drawdown_period"`
	RecoveryPeriod    int     `json:"recovery_period"`
}

// BacktestResult represents the complete backtest result
type BacktestResult struct {
	ID                 string             `json:"id"`
	Request            BacktestRequest    `json:"request"`
	Trades             []Trade            `json:"trades"`
	DailyReturns       []DailyReturn      `json:"daily_returns"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
	CreatedAt          time.Time          `json:"created_at"`
	Duration           time.Duration      `json:"duration"`
}

// DailyReturn represents daily portfolio value and returns
type DailyReturn struct {
	Date             time.Time `json:"date"`
	PortfolioValue   float64   `json:"portfolio_value"`
	DailyReturn      float64   `json:"daily_return"`
	CumulativeReturn float64   `json:"cumulative_return"`
	Drawdown         float64   `json:"drawdown"`
	Cash             float64   `json:"cash"`
	Position         Position  `json:"position"`
}

// DataSourceConfig represents data source configuration
type DataSourceConfig struct {
	Provider   string                 `json:"provider"`
	APIKey     string                 `json:"api_key,omitempty"`
	BaseURL    string                 `json:"base_url,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
