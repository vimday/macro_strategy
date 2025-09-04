import { PerformanceMetrics } from '@/types';

/**
 * Format percentage value
 */
export function formatPercentage(value: number): string {
  return `${(value * 100).toFixed(2)}%`;
}

/**
 * Format currency value
 */
export function formatCurrency(value: number): string {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value);
}

/**
 * Get value color based on performance
 */
export function getValueColor(value: number): string {
  if (value > 0) return 'text-green-600';
  if (value < 0) return 'text-red-600';
  return 'text-gray-600';
}

/**
 * Calculate performance metrics from trades and returns
 */
export function calculatePerformanceMetrics(
  initialCash: number,
  finalValue: number,
  dailyReturns: { date: string; daily_return: number; cumulative_return: number; drawdown: number }[],
  trades: { action: string; amount: number; commission: number }[]
): PerformanceMetrics {
  const totalReturn = (finalValue - initialCash) / initialCash;
  
  // Calculate annualized return (assuming daily returns)
  const numDays = dailyReturns.length;
  const annualizedReturn = Math.pow(1 + totalReturn, 365 / numDays) - 1;
  
  // Calculate volatility (standard deviation of daily returns)
  const avgDailyReturn = dailyReturns.reduce((sum, dr) => sum + dr.daily_return, 0) / numDays;
  const variance = dailyReturns.reduce((sum, dr) => sum + Math.pow(dr.daily_return - avgDailyReturn, 2), 0) / numDays;
  const volatility = Math.sqrt(variance) * Math.sqrt(365); // Annualized volatility
  
  // Calculate max drawdown
  const maxDrawdown = Math.min(...dailyReturns.map(dr => dr.drawdown));
  
  // Calculate Sharpe ratio (assuming risk-free rate of 0 for simplicity)
  const sharpeRatio = volatility > 0 ? annualizedReturn / volatility : 0;
  
  // Calculate Sortino ratio (using downside deviation)
  const negativeReturns = dailyReturns.filter(dr => dr.daily_return < 0);
  const downsideVariance = negativeReturns.reduce((sum, dr) => sum + Math.pow(dr.daily_return - avgDailyReturn, 2), 0) / numDays;
  const downsideDeviation = Math.sqrt(downsideVariance) * Math.sqrt(365);
  const sortinoRatio = downsideDeviation > 0 ? annualizedReturn / downsideDeviation : 0;
  
  // Calculate Calmar ratio
  const calmarRatio = maxDrawdown < 0 ? Math.abs(annualizedReturn / maxDrawdown) : 0;
  
  // Calculate trade statistics
  const totalTrades = trades.length;
  const winningTrades = trades.filter(trade => trade.action === 'sell' && trade.amount > 0).length;
  const losingTrades = totalTrades - winningTrades;
  const winRate = totalTrades > 0 ? winningTrades / totalTrades : 0;
  
  // Calculate profit factor
  const grossProfit = trades.filter(trade => trade.action === 'sell' && trade.amount > 0)
    .reduce((sum, trade) => sum + trade.amount, 0);
  const grossLoss = Math.abs(trades.filter(trade => trade.action === 'sell' && trade.amount < 0)
    .reduce((sum, trade) => sum + trade.amount, 0));
  const profitFactor = grossLoss > 0 ? grossProfit / grossLoss : grossProfit > 0 ? Infinity : 0;
  
  // Calculate average winning/losing trades
  const winningTradeValues = trades.filter(trade => trade.action === 'sell' && trade.amount > 0)
    .map(trade => trade.amount / initialCash);
  const losingTradeValues = trades.filter(trade => trade.action === 'sell' && trade.amount < 0)
    .map(trade => trade.amount / initialCash);
  
  const avgWinningTrade = winningTradeValues.length > 0 
    ? winningTradeValues.reduce((sum, val) => sum + val, 0) / winningTradeValues.length 
    : 0;
  
  const avgLosingTrade = losingTradeValues.length > 0 
    ? losingTradeValues.reduce((sum, val) => sum + val, 0) / losingTradeValues.length 
    : 0;
  
  // Find max winning/losing trades
  const maxWinningTrade = winningTradeValues.length > 0 
    ? Math.max(...winningTradeValues) 
    : 0;
  
  const maxLosingTrade = losingTradeValues.length > 0 
    ? Math.min(...losingTradeValues) 
    : 0;
  
  // Calculate max drawdown period and recovery period
  let maxDrawdownPeriod = 0;
  let recoveryPeriod = 0;
  
  // Simple implementation - in a real scenario, you'd want more sophisticated calculation
  if (dailyReturns.length > 0) {
    maxDrawdownPeriod = dailyReturns.filter(dr => dr.drawdown < 0).length;
    recoveryPeriod = dailyReturns.filter(dr => dr.cumulative_return > 0).length;
  }
  
  return {
    total_return: totalReturn,
    annualized_return: annualizedReturn,
    max_drawdown: maxDrawdown,
    sharpe_ratio: sharpeRatio,
    sortino_ratio: sortinoRatio,
    volatility: volatility,
    win_rate: winRate,
    profit_factor: profitFactor,
    calmar_ratio: calmarRatio,
    total_trades: totalTrades,
    winning_trades: winningTrades,
    losing_trades: losingTrades,
    avg_winning_trade: avgWinningTrade,
    avg_losing_trade: avgLosingTrade,
    max_winning_trade: maxWinningTrade,
    max_losing_trade: maxLosingTrade,
    max_drawdown_period: maxDrawdownPeriod,
    recovery_period: recoveryPeriod,
  };
}

/**
 * Debounce function
 */
export function debounce<T extends (...args: unknown[]) => unknown>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

/**
 * Deep clone object
 */
export function deepClone<T>(obj: T): T {
  if (obj === null || typeof obj !== 'object') return obj;
  if (obj instanceof Date) return new Date(obj.getTime()) as T;
  if (obj instanceof Array) return obj.map(deepClone) as T;
  if (typeof obj === 'object') {
    const cloned = {} as T;
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        cloned[key] = deepClone(obj[key]);
      }
    }
    return cloned;
  }
  return obj;
}

/**
 * Generate random ID
 */
export function generateId(): string {
  return Math.random().toString(36).substring(2) + Date.now().toString(36);
}

/**
 * Validate date range
 */
export function validateDateRange(startDate: string, endDate: string): boolean {
  const start = new Date(startDate);
  const end = new Date(endDate);
  const today = new Date();
  
  return start < end && end <= today;
}

/**
 * Calculate Sharpe ratio color
 */
export function getSharpeRatioColor(ratio: number): string {
  if (ratio >= 2) return 'text-green-700';
  if (ratio >= 1) return 'text-green-600';
  if (ratio >= 0.5) return 'text-yellow-600';
  if (ratio >= 0) return 'text-orange-600';
  return 'text-red-600';
}

/**
 * Calculate max drawdown color
 */
export function getDrawdownColor(drawdown: number): string {
  const absDrawdown = Math.abs(drawdown);
  if (absDrawdown <= 0.05) return 'text-green-600';
  if (absDrawdown <= 0.10) return 'text-yellow-600';
  if (absDrawdown <= 0.20) return 'text-orange-600';
  return 'text-red-600';
}

/**
 * Convert array to chart data format
 */
export function convertToChartData(data: Record<string, unknown>[], xKey: string, yKey: string) {
  return data.map(item => ({
    x: item[xKey],
    y: item[yKey],
  }));
}

/**
 * Safe division
 */
export function safeDivide(numerator: number, denominator: number, fallback: number = 0): number {
  return denominator === 0 ? fallback : numerator / denominator;
}