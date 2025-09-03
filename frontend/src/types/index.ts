// Market Types
export type MarketType = 'a_share' | 'crypto' | 'hk_us' | 'etf' | 'bond' | 'future' | 'option';

// Asset Classes
export type AssetClass = 'equity' | 'crypto' | 'commodity' | 'bond' | 'derivative' | 'index';

// Currencies
export type Currency = 'CNY' | 'USD' | 'HKD' | 'BTC' | 'ETH';

// Trading Hours
export interface TradingHours {
  timezone: string;
  open_time: string;
  close_time: string;
  weekend_days: number[];
}

// Index interface with enhanced metadata
export interface Index {
  id: string;
  name: string;
  symbol: string;
  market_type: MarketType;
  asset_class: AssetClass;
  currency: Currency;
  description: string;
  trading_hours?: TradingHours;
  metadata?: Record<string, any>;
}

// OHLCV data point with enhanced fields
export interface OHLCV {
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  amount?: number;   // 成交额
  turnover?: number; // 换手率
  pct_chg?: number;  // 涨跌幅
}

// Market data with enhanced metadata
export interface MarketData {
  asset_id: string;
  symbol: string;
  market_type: MarketType;
  asset_class: AssetClass;
  currency: Currency;
  data: OHLCV[];
  metadata?: Record<string, any>;
  last_update: string;
}

// Strategy types with enhanced options
export type StrategyType = 'monthly_rotation' | 'buy_and_hold' | 'grid_trading' | 'mean_reversion' | 'momentum';

// Risk management configuration
export interface RiskManagementConfig {
  max_position_size?: number;
  stop_loss?: number;
  take_profit?: number;
  max_drawdown?: number;
  commission_rate?: number;
  slippage_rate?: number;
}

// Strategy configuration with enhanced flexibility
export interface StrategyConfig {
  type: StrategyType;
  parameters: Record<string, any>;
  description: string;
  risk_management?: RiskManagementConfig;
  metadata?: Record<string, any>;
}

// Monthly rotation strategy parameters
export interface MonthlyRotationParams {
  buy_days_before_month_end: number;
  sell_days_after_month_start: number;
}

// Backtest request with enhanced configuration
export interface BacktestRequest {
  asset_id?: string;  // 新字段
  index_id?: string;  // 向后兼容
  strategy: StrategyConfig;
  start_date: string;
  end_date: string;
  initial_cash: number;
  benchmark?: string;
  rebalance_freq?: string;
  data_source?: string;
  metadata?: Record<string, any>;
}

// Trade record
export interface Trade {
  date: string;
  action: 'buy' | 'sell';
  price: number;
  quantity: number;
  amount: number;
  commission: number;
}

// Position
export interface Position {
  quantity: number;
  avg_price: number;
  market_value: number;
  unrealized_pl: number;
}

// Performance metrics
export interface PerformanceMetrics {
  total_return: number;
  annualized_return: number;
  max_drawdown: number;
  sharpe_ratio: number;
  sortino_ratio: number;
  volatility: number;
  win_rate: number;
  profit_factor: number;
  calmar_ratio: number;
  total_trades: number;
  winning_trades: number;
  losing_trades: number;
  avg_winning_trade: number;
  avg_losing_trade: number;
  max_winning_trade: number;
  max_losing_trade: number;
  max_drawdown_period: number;
  recovery_period: number;
}

// Daily return
export interface DailyReturn {
  date: string;
  portfolio_value: number;
  daily_return: number;
  cumulative_return: number;
  drawdown: number;
  cash: number;
  position: Position;
}

// Backtest result
export interface BacktestResult {
  id: string;
  request: BacktestRequest;
  trades: Trade[];
  daily_returns: DailyReturn[];
  performance_metrics: PerformanceMetrics;
  created_at: string;
  duration: number;
}

// API Response types
export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface ErrorResponse {
  error: string;
  message: string;
  code: number;
}

// Form data types for frontend
export interface BacktestFormData {
  indexId: string;
  strategyType: StrategyType;
  buyDaysBeforeMonthEnd: number;
  sellDaysAfterMonthStart: number;
  startDate: string;
  endDate: string;
  initialCash: number;
}

// Chart data types
export interface ChartDataPoint {
  date: string;
  value: number;
  [key: string]: any;
}

export interface PerformanceChartData {
  portfolioValue: ChartDataPoint[];
  cumulativeReturns: ChartDataPoint[];
  drawdown: ChartDataPoint[];
}