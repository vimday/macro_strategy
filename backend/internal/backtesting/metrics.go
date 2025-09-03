package backtesting

import (
	"macro_strategy/internal/models"
	"math"
)

// calculatePerformanceMetrics calculates comprehensive performance metrics
func (be *BacktestEngine) calculatePerformanceMetrics(dailyReturns []models.DailyReturn, trades []models.Trade) models.PerformanceMetrics {
	if len(dailyReturns) == 0 {
		return models.PerformanceMetrics{}
	}

	initialValue := dailyReturns[0].PortfolioValue
	finalValue := dailyReturns[len(dailyReturns)-1].PortfolioValue

	// Basic return metrics - ensure we don't divide by zero
	totalReturn := 0.0
	if initialValue > 0 {
		totalReturn = (finalValue - initialValue) / initialValue
	}

	daysCount := len(dailyReturns)
	yearsCount := float64(daysCount) / 252.0 // Assuming 252 trading days per year

	// Calculate annualized return safely
	annualizedReturn := 0.0
	if yearsCount > 0 && totalReturn > -1 {
		annualizedReturn = math.Pow(1+totalReturn, 1/yearsCount) - 1
	}

	// Calculate daily returns for further analysis
	var returns []float64
	for i := 1; i < len(dailyReturns); i++ {
		returns = append(returns, dailyReturns[i].DailyReturn)
	}

	// Risk metrics
	maxDrawdown := be.calculateMaxDrawdown(dailyReturns)
	volatility := be.calculateVolatility(returns)

	// Risk-adjusted metrics
	riskFreeRate := 0.03 // Assume 3% annual risk-free rate
	sharpeRatio := be.calculateSharpeRatio(annualizedReturn, volatility, riskFreeRate)
	sortinoRatio := be.calculateSortinoRatio(returns, riskFreeRate)
	calmarRatio := be.calculateCalmarRatio(annualizedReturn, maxDrawdown)

	// Trade-based metrics
	tradeMetrics := be.calculateTradeMetrics(trades)

	// Drawdown periods
	maxDrawdownPeriod, recoveryPeriod := be.calculateDrawdownPeriods(dailyReturns)

	return models.PerformanceMetrics{
		TotalReturn:       totalReturn,
		AnnualizedReturn:  annualizedReturn,
		MaxDrawdown:       maxDrawdown,
		SharpeRatio:       sharpeRatio,
		SortinoRatio:      sortinoRatio,
		Volatility:        volatility,
		WinRate:           tradeMetrics.WinRate,
		ProfitFactor:      tradeMetrics.ProfitFactor,
		CalmarRatio:       calmarRatio,
		TotalTrades:       tradeMetrics.TotalTrades,
		WinningTrades:     tradeMetrics.WinningTrades,
		LosingTrades:      tradeMetrics.LosingTrades,
		AvgWinningTrade:   tradeMetrics.AvgWinningTrade,
		AvgLosingTrade:    tradeMetrics.AvgLosingTrade,
		MaxWinningTrade:   tradeMetrics.MaxWinningTrade,
		MaxLosingTrade:    tradeMetrics.MaxLosingTrade,
		MaxDrawdownPeriod: maxDrawdownPeriod,
		RecoveryPeriod:    recoveryPeriod,
	}
}

// calculateMaxDrawdown calculates the maximum drawdown
func (be *BacktestEngine) calculateMaxDrawdown(dailyReturns []models.DailyReturn) float64 {
	maxDrawdown := 0.0
	for _, dr := range dailyReturns {
		if dr.Drawdown > maxDrawdown {
			maxDrawdown = dr.Drawdown
		}
	}
	return maxDrawdown
}

// calculateVolatility calculates annualized volatility
func (be *BacktestEngine) calculateVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0.0
	}

	// Calculate mean return
	sum := 0.0
	for _, r := range returns {
		sum += r
	}
	mean := sum / float64(len(returns))

	// Calculate variance
	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(len(returns) - 1)

	// Annualize volatility (252 trading days)
	return math.Sqrt(variance) * math.Sqrt(252)
}

// calculateSharpeRatio calculates the Sharpe ratio
func (be *BacktestEngine) calculateSharpeRatio(annualizedReturn, volatility, riskFreeRate float64) float64 {
	if volatility == 0 {
		return 0.0
	}
	return (annualizedReturn - riskFreeRate) / volatility
}

// calculateSortinoRatio calculates the Sortino ratio
func (be *BacktestEngine) calculateSortinoRatio(returns []float64, riskFreeRate float64) float64 {
	if len(returns) == 0 {
		return 0.0
	}

	// Calculate mean return
	sum := 0.0
	for _, r := range returns {
		sum += r
	}
	meanReturn := sum / float64(len(returns))
	annualizedMeanReturn := meanReturn * 252

	// Calculate downside deviation
	downsideVariance := 0.0
	downsideCount := 0
	dailyRiskFreeRate := riskFreeRate / 252

	for _, r := range returns {
		if r < dailyRiskFreeRate {
			downsideVariance += math.Pow(r-dailyRiskFreeRate, 2)
			downsideCount++
		}
	}

	if downsideCount == 0 {
		return math.Inf(1) // No downside risk
	}

	downsideVariance /= float64(len(returns))
	downsideDeviation := math.Sqrt(downsideVariance) * math.Sqrt(252)

	if downsideDeviation == 0 {
		return math.Inf(1)
	}

	return (annualizedMeanReturn - riskFreeRate) / downsideDeviation
}

// calculateCalmarRatio calculates the Calmar ratio
func (be *BacktestEngine) calculateCalmarRatio(annualizedReturn, maxDrawdown float64) float64 {
	if maxDrawdown <= 0 {
		if annualizedReturn > 0 {
			return math.Inf(1) // Positive return with no drawdown
		}
		return 0 // No drawdown, no meaningful ratio
	}
	return annualizedReturn / maxDrawdown
}

// TradeMetrics holds trade-based performance metrics
type TradeMetrics struct {
	TotalTrades     int
	WinningTrades   int
	LosingTrades    int
	WinRate         float64
	ProfitFactor    float64
	AvgWinningTrade float64
	AvgLosingTrade  float64
	MaxWinningTrade float64
	MaxLosingTrade  float64
}

// calculateTradeMetrics calculates trade-based metrics
func (be *BacktestEngine) calculateTradeMetrics(trades []models.Trade) TradeMetrics {
	if len(trades) == 0 {
		return TradeMetrics{}
	}

	// Group trades into round trips (buy-sell pairs)
	var roundTrips []float64
	var buyTrade *models.Trade

	for i, trade := range trades {
		currentTrade := trade // Create a copy to avoid pointer issues
		if currentTrade.Action == "buy" {
			buyTrade = &currentTrade
		} else if currentTrade.Action == "sell" && buyTrade != nil {
			// Calculate P&L for this round trip
			grossPnL := (currentTrade.Price - buyTrade.Price) * currentTrade.Quantity
			totalCommission := currentTrade.Commission + buyTrade.Commission
			netPnL := grossPnL - totalCommission
			
			// Calculate percentage return based on initial investment
			initialInvestment := buyTrade.Price * buyTrade.Quantity + buyTrade.Commission
			percentageReturn := netPnL / initialInvestment
			
			// Store percentage return for analysis
			roundTrips = append(roundTrips, percentageReturn)
			buyTrade = nil
		}
		_ = i // Suppress unused variable warning
	}

	if len(roundTrips) == 0 {
		return TradeMetrics{}
	}

	// Analyze round trips
	winningTrades := 0
	losingTrades := 0
	totalWinningReturn := 0.0
	totalLosingReturn := 0.0
	maxWinningTrade := math.Inf(-1)
	maxLosingTrade := math.Inf(1)

	for _, returnPct := range roundTrips {
		if returnPct > 0 {
			winningTrades++
			totalWinningReturn += returnPct
			if returnPct > maxWinningTrade {
				maxWinningTrade = returnPct
			}
		} else if returnPct < 0 {
			losingTrades++
			totalLosingReturn += math.Abs(returnPct)
			if returnPct < maxLosingTrade {
				maxLosingTrade = returnPct
			}
		}
		// Note: returnPct == 0 (breakeven trades) are not counted as wins or losses
	}

	winRate := 0.0
	if len(roundTrips) > 0 {
		winRate = float64(winningTrades) / float64(len(roundTrips))
	}

	profitFactor := 0.0
	if totalLosingReturn > 0 {
		profitFactor = totalWinningReturn / totalLosingReturn
	} else if totalWinningReturn > 0 {
		// All trades are winning, set a high profit factor
		profitFactor = math.Inf(1)
	}

	avgWinningTrade := 0.0
	if winningTrades > 0 {
		avgWinningTrade = totalWinningReturn / float64(winningTrades)
	}

	avgLosingTrade := 0.0
	if losingTrades > 0 {
		avgLosingTrade = -totalLosingReturn / float64(losingTrades)
	}

	if math.IsInf(maxWinningTrade, -1) {
		maxWinningTrade = 0
	}
	if math.IsInf(maxLosingTrade, 1) {
		maxLosingTrade = 0
	}

	return TradeMetrics{
		TotalTrades:     len(roundTrips),
		WinningTrades:   winningTrades,
		LosingTrades:    losingTrades,
		WinRate:         winRate,
		ProfitFactor:    profitFactor,
		AvgWinningTrade: avgWinningTrade,
		AvgLosingTrade:  avgLosingTrade,
		MaxWinningTrade: maxWinningTrade,
		MaxLosingTrade:  maxLosingTrade,
	}
}

// calculateDrawdownPeriods calculates maximum drawdown period and recovery period
func (be *BacktestEngine) calculateDrawdownPeriods(dailyReturns []models.DailyReturn) (int, int) {
	if len(dailyReturns) == 0 {
		return 0, 0
	}

	maxDrawdownPeriod := 0
	recoveryPeriod := 0
	currentDrawdownPeriod := 0
	maxDrawdownStartIndex := -1
	maxDrawdown := 0.0

	for i, dr := range dailyReturns {
		if dr.Drawdown > 0 {
			if currentDrawdownPeriod == 0 {
				// Start of new drawdown period
				currentDrawdownPeriod = 1
			} else {
				currentDrawdownPeriod++
			}

			// Track maximum drawdown
			if dr.Drawdown > maxDrawdown {
				maxDrawdown = dr.Drawdown
				maxDrawdownStartIndex = i - currentDrawdownPeriod + 1
			}
		} else {
			// End of drawdown period
			if currentDrawdownPeriod > maxDrawdownPeriod {
				maxDrawdownPeriod = currentDrawdownPeriod
			}
			currentDrawdownPeriod = 0
		}
	}

	// Handle case where drawdown period extends to the end
	if currentDrawdownPeriod > maxDrawdownPeriod {
		maxDrawdownPeriod = currentDrawdownPeriod
	}

	// Calculate recovery period (from max drawdown to recovery)
	if maxDrawdownStartIndex >= 0 {
		for i := maxDrawdownStartIndex; i < len(dailyReturns); i++ {
			if dailyReturns[i].Drawdown == 0 {
				recoveryPeriod = i - maxDrawdownStartIndex
				break
			}
		}

		// If no recovery found, set to remaining period
		if recoveryPeriod == 0 {
			recoveryPeriod = len(dailyReturns) - maxDrawdownStartIndex
		}
	}

	return maxDrawdownPeriod, recoveryPeriod
}
