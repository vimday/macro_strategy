package models
package models

import "time"

// MarketType represents different market types
type MarketType string

const (
	MarketTypeAShare MarketType = "a_share"
	MarketTypeCrypto MarketType = "crypto"
	MarketTypeHKUS   MarketType = "hk_us"
)

// Index represents a market index
type Index struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Symbol      string     `json:"symbol"`
	MarketType  MarketType `json:"market_type"`
	Description string     `json:"description"`
}

// OHLCV represents Open, High, Low, Close, Volume data
type OHLCV struct {
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int64     `json:"volume"`
}

// MarketData represents historical market data for an index
type MarketData struct {
	IndexID string   `json:"index_id"`
	Data    []OHLCV  `json:"data"`
}

// StrategyType represents different strategy types
type StrategyType string

const (
	StrategyTypeMonthlyRotation StrategyType = "monthly_rotation"
	// Future strategy types can be added here
)

// StrategyConfig represents strategy configuration
type StrategyConfig struct {
	Type        StrategyType           `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Description string                 `json:"description"`
}

// MonthlyRotationParams represents parameters for monthly rotation strategy
type MonthlyRotationParams struct {
	BuyDaysBeforeMonthEnd  int `json:"buy_days_before_month_end"`
	SellDaysAfterMonthStart int `json:"sell_days_after_month_start"`
}

// BacktestRequest represents a backtest request
type BacktestRequest struct {
	IndexID     string         `json:"index_id"`
	Strategy    StrategyConfig `json:"strategy"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	InitialCash float64        `json:"initial_cash"`
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
	Quantity    float64 `json:"quantity"`
	AvgPrice    float64 `json:"avg_price"`
	MarketValue float64 `json:"market_value"`
	UnrealizedPL float64 `json:"unrealized_pl"`
}

// PerformanceMetrics represents strategy performance metrics
type PerformanceMetrics struct {
	TotalReturn         float64 `json:"total_return"`
	AnnualizedReturn    float64 `json:"annualized_return"`
	MaxDrawdown         float64 `json:"max_drawdown"`
	SharpeRatio         float64 `json:"sharpe_ratio"`
	SortinoRatio        float64 `json:"sortino_ratio"`
	Volatility          float64 `json:"volatility"`
	WinRate             float64 `json:"win_rate"`
	ProfitFactor        float64 `json:"profit_factor"`
	CalmarRatio         float64 `json:"calmar_ratio"`
	TotalTrades         int     `json:"total_trades"`
	WinningTrades       int     `json:"winning_trades"`
	LosingTrades        int     `json:"losing_trades"`
	AvgWinningTrade     float64 `json:"avg_winning_trade"`
	AvgLosingTrade      float64 `json:"avg_losing_trade"`
	MaxWinningTrade     float64 `json:"max_winning_trade"`
	MaxLosingTrade      float64 `json:"max_losing_trade"`
	MaxDrawdownPeriod   int     `json:"max_drawdown_period"`
	RecoveryPeriod      int     `json:"recovery_period"`
}

// BacktestResult represents the complete backtest result
type BacktestResult struct {
	ID                 string              `json:"id"`
	Request            BacktestRequest     `json:"request"`
	Trades             []Trade             `json:"trades"`
	DailyReturns       []DailyReturn       `json:"daily_returns"`
	PerformanceMetrics PerformanceMetrics  `json:"performance_metrics"`
	CreatedAt          time.Time           `json:"created_at"`
	Duration           time.Duration       `json:"duration"`
}

// DailyReturn represents daily portfolio value and returns
type DailyReturn struct {
	Date            time.Time `json:"date"`
	PortfolioValue  float64   `json:"portfolio_value"`
	DailyReturn     float64   `json:"daily_return"`
	CumulativeReturn float64  `json:"cumulative_return"`
	Drawdown        float64   `json:"drawdown"`
	Cash            float64   `json:"cash"`
	Position        Position  `json:"position"`
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