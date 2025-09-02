// Market Types
export type MarketType = 'a_share' | 'crypto' | 'hk_us';

// Index interface
export interface Index {
  id: string;
  name: string;
  symbol: string;
  market_type: MarketType;
  description: string;
}

// OHLCV data point
export interface OHLCV {
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}

// Market data
export interface MarketData {
  index_id: string;
  data: OHLCV[];
}

// Strategy types
export type StrategyType = 'monthly_rotation';

// Strategy configuration
export interface StrategyConfig {
  type: StrategyType;
  parameters: Record<string, any>;
  description: string;
}

// Monthly rotation strategy parameters
export interface MonthlyRotationParams {
  buy_days_before_month_end: number;
  sell_days_after_month_start: number;
}

// Backtest request
export interface BacktestRequest {
  index_id: string;
  strategy: StrategyConfig;
  start_date: string;
  end_date: string;
  initial_cash: number;
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