package models

import "time"

// MarketType represents different market types
type MarketType string

const (
	MarketTypeAShareIndex MarketType = "a_share_index" // A股指数
	MarketTypeAShareStock MarketType = "a_share_stock" // A股个股
	MarketTypeUSIndex     MarketType = "us_index"      // 美股指数
	MarketTypeUSStock     MarketType = "us_stock"      // 美股个股
	MarketTypeCrypto      MarketType = "crypto"        // 数字货币
	MarketTypeHKIndex     MarketType = "hk_index"      // 港股指数
	MarketTypeHKStock     MarketType = "hk_stock"      // 港股个股
	MarketTypeETF         MarketType = "etf"           // ETF基金
	MarketTypeBond        MarketType = "bond"          // 债券
	MarketTypeFuture      MarketType = "future"        // 期货
	MarketTypeOption      MarketType = "option"        // 期权
	MarketTypeCommodity   MarketType = "commodity"     // 大宗商品
)

// AssetClass represents different asset classes
type AssetClass string

const (
	AssetClassEquity     AssetClass = "equity"     // 股票
	AssetClassIndex      AssetClass = "index"      // 指数
	AssetClassCrypto     AssetClass = "crypto"     // 数字货币
	AssetClassETF        AssetClass = "etf"        // ETF基金
	AssetClassCommodity  AssetClass = "commodity"  // 大宗商品
	AssetClassBond       AssetClass = "bond"       // 债券
	AssetClassDerivative AssetClass = "derivative" // 衍生品
	AssetClassREITs      AssetClass = "reits"      // 房地产信托基金
)

// Currency represents different currencies
type Currency string

const (
	// 法定货币
	CurrencyCNY Currency = "CNY" // 人民币
	CurrencyUSD Currency = "USD" // 美元
	CurrencyHKD Currency = "HKD" // 港币
	CurrencyEUR Currency = "EUR" // 欧元
	CurrencyJPY Currency = "JPY" // 日元
	CurrencyGBP Currency = "GBP" // 英镑
	// 数字货币
	CurrencyBTC  Currency = "BTC"  // 比特币
	CurrencyETH  Currency = "ETH"  // 以太坊
	CurrencyUSDT Currency = "USDT" // 泰达币
	CurrencyBNB  Currency = "BNB"  // 币安币
	CurrencyADA  Currency = "ADA"  // 艾达币
	CurrencySOL  Currency = "SOL"  // Solana
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
	// 基础策略
	StrategyTypeMonthlyRotation StrategyType = "monthly_rotation" // 月末轮动策略
	StrategyTypeBuyAndHold      StrategyType = "buy_and_hold"     // 买入持有策略
	// 技术策略
	StrategyTypeGridTrading   StrategyType = "grid_trading"   // 网格交易策略
	StrategyTypeMeanReversion StrategyType = "mean_reversion" // 均值回归策略
	StrategyTypeMomentum      StrategyType = "momentum"       // 动量策略
	StrategyTypeBreakout      StrategyType = "breakout"       // 突破策略
	// 高频策略
	StrategyTypeDCA       StrategyType = "dca"       // 定投策略 (Dollar Cost Averaging)
	StrategyTypeRebalance StrategyType = "rebalance" // 再平衡策略
	StrategyTypePairs     StrategyType = "pairs"     // 配对交易策略
	StrategyTypeArbitrage StrategyType = "arbitrage" // 套利策略
	// 多因子策略
	StrategyTypeMultiFactor StrategyType = "multi_factor" // 多因子策略
	StrategyTypePortfolio   StrategyType = "portfolio"    // 组合策略
	StrategyTypeRiskParity  StrategyType = "risk_parity"  // 风险平价策略
	StrategyTypeMinVariance StrategyType = "min_variance" // 最小方差策略
	// 机器学习策略
	StrategyTypeML            StrategyType = "ml"            // 机器学习策略
	StrategyTypeReinforcement StrategyType = "reinforcement" // 强化学习策略
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

// BuyAndHoldParams represents parameters for buy and hold strategy
type BuyAndHoldParams struct {
	RebalanceFrequency string  `json:"rebalance_frequency,omitempty"` // "monthly", "quarterly", "yearly", "never"
	DividendReinvest   bool    `json:"dividend_reinvest,omitempty"`   // 股息再投资
	TargetAllocation   float64 `json:"target_allocation,omitempty"`   // 目标仓位比例
}

// GridTradingParams represents parameters for grid trading strategy
type GridTradingParams struct {
	GridCount    int     `json:"grid_count"`    // 网格数量
	GridSpacing  float64 `json:"grid_spacing"`  // 网格间距(百分比)
	BasePrice    float64 `json:"base_price"`    // 基准价格
	MaxPosition  float64 `json:"max_position"`  // 最大仓位
	ProfitTarget float64 `json:"profit_target"` // 止盈目标
}

// MomentumParams represents parameters for momentum strategy
type MomentumParams struct {
	LookbackPeriod    int     `json:"lookback_period"`    // 回望周期
	MomentumThreshold float64 `json:"momentum_threshold"` // 动量阈值
	HoldingPeriod     int     `json:"holding_period"`     // 持有周期
	RebalanceFreq     string  `json:"rebalance_freq"`     // 再平衡频率
}

// DCAParams represents parameters for Dollar Cost Averaging strategy
type DCAParams struct {
	InvestmentAmount float64 `json:"investment_amount"` // 每次投资金额
	Frequency        string  `json:"frequency"`         // 投资频率: "daily", "weekly", "monthly"
	DurationMonths   int     `json:"duration_months"`   // 投资期限(月)
	StartDelay       int     `json:"start_delay"`       // 开始延迟(天)
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
	Type      string                 `json:"type"` // "akshare", "yahoo", "binance", "mock"
	APIKey    string                 `json:"api_key,omitempty"`
	APISecret string                 `json:"api_secret,omitempty"`
	BaseURL   string                 `json:"base_url,omitempty"`
	RateLimit int                    `json:"rate_limit,omitempty"`
	Timeout   int                    `json:"timeout,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MultiStrategyBacktestRequest represents a request to compare multiple strategies
type MultiStrategyBacktestRequest struct {
	AssetID       string                 `json:"asset_id"`
	Strategies    []StrategyConfig       `json:"strategies"` // 多个策略配置
	StartDate     time.Time              `json:"start_date"`
	EndDate       time.Time              `json:"end_date"`
	InitialCash   float64                `json:"initial_cash"`
	Benchmark     string                 `json:"benchmark,omitempty"`      // 基准指数
	DataSource    string                 `json:"data_source,omitempty"`    // 数据源
	ComparisonOpt *ComparisonOptions     `json:"comparison_opt,omitempty"` // 对比选项
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ComparisonOptions represents options for strategy comparison
type ComparisonOptions struct {
	ShowBenchmark      bool     `json:"show_benchmark"`       // 显示基准
	NormalizeReturns   bool     `json:"normalize_returns"`    // 归一化收益
	ShowDrawdown       bool     `json:"show_drawdown"`        // 显示回撤
	ShowRollingMetrics bool     `json:"show_rolling_metrics"` // 显示滚动指标
	RollingWindow      int      `json:"rolling_window"`       // 滚动窗口
	Metrics            []string `json:"metrics"`              // 需要对比的指标
}

// MultiStrategyBacktestResult represents the result of multiple strategy comparison
type MultiStrategyBacktestResult struct {
	ID              string                       `json:"id"`
	Request         MultiStrategyBacktestRequest `json:"request"`
	Results         []BacktestResult             `json:"results"`                    // 各策略结果
	Comparison      *StrategyComparison          `json:"comparison"`                 // 对比结果
	BenchmarkResult *BacktestResult              `json:"benchmark_result,omitempty"` // 基准结果
	CreatedAt       time.Time                    `json:"created_at"`
	Duration        time.Duration                `json:"duration"`
}

// StrategyComparison represents comparison results between strategies
type StrategyComparison struct {
	MetricsComparison map[string][]float64 `json:"metrics_comparison"` // 指标对比
	Rankings          map[string][]int     `json:"rankings"`           // 排名
	CorrelationMatrix [][]float64          `json:"correlation_matrix"` // 相关性矩阵
	BestStrategy      string               `json:"best_strategy"`      // 最佳策略
	WorstStrategy     string               `json:"worst_strategy"`     // 最差策略
	Summary           string               `json:"summary"`            // 对比总结
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
