package backtesting

import (
	"fmt"
	"macro_strategy/internal/models"
	"math"
	"sort"
	"time"
)

// BacktestEngine handles strategy backtesting
type BacktestEngine struct {
	commissionRate float64 // Commission rate per trade
}

// NewBacktestEngine creates a new backtest engine
func NewBacktestEngine() *BacktestEngine {
	return &BacktestEngine{
		commissionRate: 0.0003, // 0.03% commission rate
	}
}

// RunBacktest executes a backtest for given request
func (be *BacktestEngine) RunBacktest(request models.BacktestRequest, marketData *models.MarketData) (*models.BacktestResult, error) {
	startTime := time.Now()

	// Validate request
	if err := be.validateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Filter market data to backtest period
	filteredData := be.filterDataByDateRange(marketData.Data, request.StartDate, request.EndDate)
	if len(filteredData) == 0 {
		return nil, fmt.Errorf("no market data available for the specified period")
	}

	// Execute strategy
	trades, dailyReturns, err := be.executeStrategy(request, filteredData)
	if err != nil {
		return nil, fmt.Errorf("strategy execution failed: %w", err)
	}

	// Calculate performance metrics
	metrics := be.calculatePerformanceMetrics(dailyReturns, trades)

	// Create result
	result := &models.BacktestResult{
		ID:                 generateBacktestID(),
		Request:            request,
		Trades:             trades,
		DailyReturns:       dailyReturns,
		PerformanceMetrics: metrics,
		CreatedAt:          startTime,
		Duration:           time.Since(startTime),
	}

	return result, nil
}

// validateRequest validates the backtest request
func (be *BacktestEngine) validateRequest(request models.BacktestRequest) error {
	if request.IndexID == "" {
		return fmt.Errorf("index_id is required")
	}
	if request.InitialCash <= 0 {
		return fmt.Errorf("initial_cash must be positive")
	}
	if request.StartDate.After(request.EndDate) {
		return fmt.Errorf("start_date must be before end_date")
	}
	if request.Strategy.Type == "" {
		return fmt.Errorf("strategy type is required")
	}
	return nil
}

// filterDataByDateRange filters market data by date range
func (be *BacktestEngine) filterDataByDateRange(data []models.OHLCV, startDate, endDate time.Time) []models.OHLCV {
	var filtered []models.OHLCV
	for _, point := range data {
		if (point.Date.Equal(startDate) || point.Date.After(startDate)) &&
			(point.Date.Equal(endDate) || point.Date.Before(endDate)) {
			filtered = append(filtered, point)
		}
	}
	return filtered
}

// executeStrategy executes the trading strategy
func (be *BacktestEngine) executeStrategy(request models.BacktestRequest, data []models.OHLCV) ([]models.Trade, []models.DailyReturn, error) {
	switch request.Strategy.Type {
	case models.StrategyTypeMonthlyRotation:
		return be.executeMonthlyRotationStrategy(request, data)
	case models.StrategyTypeBuyAndHold:
		return be.executeBuyAndHoldStrategy(request, data)
	default:
		return nil, nil, fmt.Errorf("unsupported strategy type: %s", request.Strategy.Type)
	}
}

// executeMonthlyRotationStrategy executes monthly rotation strategy
func (be *BacktestEngine) executeMonthlyRotationStrategy(request models.BacktestRequest, data []models.OHLCV) ([]models.Trade, []models.DailyReturn, error) {
	// Parse strategy parameters
	params, err := be.parseMonthlyRotationParams(request.Strategy.Parameters)
	if err != nil {
		return nil, nil, err
	}

	var trades []models.Trade
	var dailyReturns []models.DailyReturn

	cash := request.InitialCash
	position := models.Position{}

	// Store initial portfolio value for calculations
	initialPortfolioValue := request.InitialCash

	for i, dataPoint := range data {
		currentDate := dataPoint.Date
		currentPrice := dataPoint.Close

		// Check if we should buy
		if be.shouldBuy(currentDate, params.BuyDaysBeforeMonthEnd, data, i) && position.Quantity == 0 {
			// Buy with all available cash
			quantity := math.Floor(cash / currentPrice)
			if quantity > 0 {
				amount := quantity * currentPrice
				commission := amount * be.commissionRate
				totalCost := amount + commission

				if totalCost <= cash {
					trade := models.Trade{
						Date:       currentDate,
						Action:     "buy",
						Price:      currentPrice,
						Quantity:   quantity,
						Amount:     amount,
						Commission: commission,
					}
					trades = append(trades, trade)

					cash -= totalCost
					position.Quantity = quantity
					position.AvgPrice = currentPrice
				}
			}
		}

		// Check if we should sell
		if be.shouldSell(currentDate, params.SellDaysAfterMonthStart, data, i) && position.Quantity > 0 {
			// Sell all positions
			amount := position.Quantity * currentPrice
			commission := amount * be.commissionRate
			netAmount := amount - commission

			trade := models.Trade{
				Date:       currentDate,
				Action:     "sell",
				Price:      currentPrice,
				Quantity:   position.Quantity,
				Amount:     amount,
				Commission: commission,
			}
			trades = append(trades, trade)

			cash += netAmount
			position = models.Position{} // Reset position
		}

		// Update position market value if holding
		if position.Quantity > 0 {
			position.MarketValue = position.Quantity * currentPrice
			position.UnrealizedPL = position.MarketValue - (position.Quantity * position.AvgPrice)
		}

		// Calculate daily portfolio value
		portfolioValue := cash + position.MarketValue
		dailyReturn := 0.0
		cumulativeReturn := (portfolioValue - initialPortfolioValue) / initialPortfolioValue

		// Calculate daily return correctly: only after first day
		if i > 0 {
			prevValue := dailyReturns[i-1].PortfolioValue
			if prevValue > 0 {
				dailyReturn = (portfolioValue - prevValue) / prevValue
			}
		}

		dailyReturns = append(dailyReturns, models.DailyReturn{
			Date:             currentDate,
			PortfolioValue:   portfolioValue,
			DailyReturn:      dailyReturn,
			CumulativeReturn: cumulativeReturn,
			Cash:             cash,
			Position:         position,
		})
	}

	// Force sell any remaining positions at the end of backtest period
	if position.Quantity > 0 {
		lastDay := data[len(data)-1]
		lastPrice := lastDay.Close
		lastDate := lastDay.Date

		// Sell all remaining positions
		amount := position.Quantity * lastPrice
		commission := amount * be.commissionRate
		netAmount := amount - commission

		trade := models.Trade{
			Date:       lastDate,
			Action:     "sell",
			Price:      lastPrice,
			Quantity:   position.Quantity,
			Amount:     amount,
			Commission: commission,
		}
		trades = append(trades, trade)

		// Update final portfolio state
		lastIndex := len(dailyReturns) - 1
		dailyReturns[lastIndex].Cash += netAmount
		dailyReturns[lastIndex].Position = models.Position{}
		dailyReturns[lastIndex].PortfolioValue = dailyReturns[lastIndex].Cash
		dailyReturns[lastIndex].CumulativeReturn = (dailyReturns[lastIndex].PortfolioValue - initialPortfolioValue) / initialPortfolioValue
	}

	// Calculate drawdown for each day
	be.calculateDrawdown(dailyReturns)

	return trades, dailyReturns, nil
}

// parseMonthlyRotationParams parses monthly rotation strategy parameters
func (be *BacktestEngine) parseMonthlyRotationParams(parameters map[string]interface{}) (*models.MonthlyRotationParams, error) {
	params := &models.MonthlyRotationParams{}

	if buyDays, ok := parameters["buy_days_before_month_end"]; ok {
		if buyDaysFloat, ok := buyDays.(float64); ok {
			params.BuyDaysBeforeMonthEnd = int(buyDaysFloat)
		} else {
			return nil, fmt.Errorf("invalid buy_days_before_month_end parameter")
		}
	} else {
		params.BuyDaysBeforeMonthEnd = 1 // Default value
	}

	if sellDays, ok := parameters["sell_days_after_month_start"]; ok {
		if sellDaysFloat, ok := sellDays.(float64); ok {
			params.SellDaysAfterMonthStart = int(sellDaysFloat)
		} else {
			return nil, fmt.Errorf("invalid sell_days_after_month_start parameter")
		}
	} else {
		params.SellDaysAfterMonthStart = 1 // Default value
	}

	return params, nil
}

// shouldBuy determines if we should buy on a given date
func (be *BacktestEngine) shouldBuy(currentDate time.Time, buyDaysBeforeMonthEnd int, data []models.OHLCV, currentIndex int) bool {
	// Find the last trading day of the month for the current date's month
	year, month, _ := currentDate.Date()

	// Get all trading days in this month from the data
	var monthTradingDays []time.Time
	for _, d := range data {
		if d.Date.Year() == year && d.Date.Month() == month {
			monthTradingDays = append(monthTradingDays, d.Date)
		}
	}

	if len(monthTradingDays) < buyDaysBeforeMonthEnd {
		return false
	}

	// Sort trading days to ensure correct order
	sort.Slice(monthTradingDays, func(i, j int) bool {
		return monthTradingDays[i].Before(monthTradingDays[j])
	})

	// Find the target buy date (buyDaysBeforeMonthEnd from the end)
	// If buyDaysBeforeMonthEnd = 1, we want the last trading day
	// If buyDaysBeforeMonthEnd = 2, we want the second-to-last trading day
	targetIndex := len(monthTradingDays) - buyDaysBeforeMonthEnd
	if targetIndex < 0 {
		targetIndex = 0
	}

	targetDate := monthTradingDays[targetIndex]
	return currentDate.Equal(targetDate)
}

// shouldSell determines if we should sell on a given date
func (be *BacktestEngine) shouldSell(currentDate time.Time, sellDaysAfterMonthStart int, data []models.OHLCV, currentIndex int) bool {
	// Find the first trading day of the month for the current date's month
	year, month, _ := currentDate.Date()

	// Get all trading days in this month from the data
	var monthTradingDays []time.Time
	for _, d := range data {
		if d.Date.Year() == year && d.Date.Month() == month {
			monthTradingDays = append(monthTradingDays, d.Date)
		}
	}

	if len(monthTradingDays) < sellDaysAfterMonthStart {
		return false
	}

	// Sort trading days to ensure correct order
	sort.Slice(monthTradingDays, func(i, j int) bool {
		return monthTradingDays[i].Before(monthTradingDays[j])
	})

	// Find the target sell date (sellDaysAfterMonthStart from the beginning)
	// If sellDaysAfterMonthStart = 1, we want the first trading day
	// If sellDaysAfterMonthStart = 2, we want the second trading day
	targetIndex := sellDaysAfterMonthStart - 1
	if targetIndex >= len(monthTradingDays) {
		targetIndex = len(monthTradingDays) - 1
	}

	targetDate := monthTradingDays[targetIndex]
	return currentDate.Equal(targetDate)
}

// calculateDrawdown calculates drawdown for each day
func (be *BacktestEngine) calculateDrawdown(dailyReturns []models.DailyReturn) {
	peak := dailyReturns[0].PortfolioValue

	for i := range dailyReturns {
		if dailyReturns[i].PortfolioValue > peak {
			peak = dailyReturns[i].PortfolioValue
		}

		drawdown := (peak - dailyReturns[i].PortfolioValue) / peak
		dailyReturns[i].Drawdown = drawdown
	}
}

// generateBacktestID generates a unique ID for the backtest
func generateBacktestID() string {
	return fmt.Sprintf("bt_%d", time.Now().UnixNano())
}

// executeBuyAndHoldStrategy executes buy and hold strategy
func (be *BacktestEngine) executeBuyAndHoldStrategy(request models.BacktestRequest, data []models.OHLCV) ([]models.Trade, []models.DailyReturn, error) {
	// Parse strategy parameters
	params, err := be.parseBuyAndHoldParams(request.Strategy.Parameters)
	if err != nil {
		return nil, nil, err
	}

	var trades []models.Trade
	var dailyReturns []models.DailyReturn

	cash := request.InitialCash
	position := models.Position{}
	initialPurchase := false
	lastRebalanceMonth := -1

	// Store initial portfolio value for calculations
	initialPortfolioValue := request.InitialCash

	for i, dataPoint := range data {
		currentDate := dataPoint.Date
		currentPrice := dataPoint.Close

		// Initial purchase on first day
		if !initialPurchase {
			quantity := math.Floor((cash * params.TargetAllocation) / currentPrice)
			if quantity > 0 {
				amount := quantity * currentPrice
				commission := amount * be.commissionRate
				totalCost := amount + commission

				if totalCost <= cash {
					trade := models.Trade{
						Date:       currentDate,
						Action:     "buy",
						Price:      currentPrice,
						Quantity:   quantity,
						Amount:     amount,
						Commission: commission,
					}
					trades = append(trades, trade)

					cash -= totalCost
					position.Quantity = quantity
					position.AvgPrice = currentPrice
					initialPurchase = true
				}
			}
		}

		// Check for rebalancing
		if params.RebalanceFrequency != "never" && params.RebalanceFrequency != "" {
			if be.shouldRebalance(currentDate, params.RebalanceFrequency, &lastRebalanceMonth) {
				// Simple rebalancing: maintain target allocation
				currentPortfolioValue := cash + (position.Quantity * currentPrice)
				targetPositionValue := currentPortfolioValue * params.TargetAllocation
				currentPositionValue := position.Quantity * currentPrice

				// If position value differs significantly from target, rebalance
				if math.Abs(currentPositionValue-targetPositionValue) > currentPortfolioValue*0.05 { // 5% threshold
					if currentPositionValue > targetPositionValue {
						// Sell excess
						excessValue := currentPositionValue - targetPositionValue
						excessQuantity := math.Floor(excessValue / currentPrice)
						if excessQuantity > 0 && excessQuantity <= position.Quantity {
							amount := excessQuantity * currentPrice
							commission := amount * be.commissionRate
							netAmount := amount - commission

							trade := models.Trade{
								Date:       currentDate,
								Action:     "sell",
								Price:      currentPrice,
								Quantity:   excessQuantity,
								Amount:     amount,
								Commission: commission,
							}
							trades = append(trades, trade)

							cash += netAmount
							position.Quantity -= excessQuantity
						}
					} else {
						// Buy more
						deficitValue := targetPositionValue - currentPositionValue
						additionalQuantity := math.Floor(deficitValue / currentPrice)
						if additionalQuantity > 0 {
							amount := additionalQuantity * currentPrice
							commission := amount * be.commissionRate
							totalCost := amount + commission

							if totalCost <= cash {
								trade := models.Trade{
									Date:       currentDate,
									Action:     "buy",
									Price:      currentPrice,
									Quantity:   additionalQuantity,
									Amount:     amount,
									Commission: commission,
								}
								trades = append(trades, trade)

								cash -= totalCost
								position.Quantity += additionalQuantity
								// Update average price
								position.AvgPrice = ((position.AvgPrice * (position.Quantity - additionalQuantity)) + (currentPrice * additionalQuantity)) / position.Quantity
							}
						}
					}
				}
			}
		}

		// Update position market value if holding
		if position.Quantity > 0 {
			position.MarketValue = position.Quantity * currentPrice
			position.UnrealizedPL = position.MarketValue - (position.Quantity * position.AvgPrice)
		}

		// Calculate daily portfolio value
		portfolioValue := cash + position.MarketValue
		dailyReturn := 0.0
		cumulativeReturn := (portfolioValue - initialPortfolioValue) / initialPortfolioValue

		// Calculate daily return correctly: only after first day
		if i > 0 {
			prevValue := dailyReturns[i-1].PortfolioValue
			if prevValue > 0 {
				dailyReturn = (portfolioValue - prevValue) / prevValue
			}
		}

		dailyReturns = append(dailyReturns, models.DailyReturn{
			Date:             currentDate,
			PortfolioValue:   portfolioValue,
			DailyReturn:      dailyReturn,
			CumulativeReturn: cumulativeReturn,
			Cash:             cash,
			Position:         position,
		})
	}

	// Calculate drawdown for each day
	be.calculateDrawdown(dailyReturns)

	return trades, dailyReturns, nil
}

// parseBuyAndHoldParams parses buy and hold strategy parameters
func (be *BacktestEngine) parseBuyAndHoldParams(parameters map[string]interface{}) (*models.BuyAndHoldParams, error) {
	params := &models.BuyAndHoldParams{
		RebalanceFrequency: "never",
		DividendReinvest:   false,
		TargetAllocation:   1.0, // 100% allocation by default
	}

	if rebalanceFreq, ok := parameters["rebalance_frequency"]; ok {
		if freq, ok := rebalanceFreq.(string); ok {
			params.RebalanceFrequency = freq
		}
	}

	if dividendReinvest, ok := parameters["dividend_reinvest"]; ok {
		if reinvest, ok := dividendReinvest.(bool); ok {
			params.DividendReinvest = reinvest
		}
	}

	if targetAllocation, ok := parameters["target_allocation"]; ok {
		if allocation, ok := targetAllocation.(float64); ok {
			if allocation > 0 && allocation <= 1.0 {
				params.TargetAllocation = allocation
			}
		}
	}

	return params, nil
}

// shouldRebalance determines if rebalancing should occur
func (be *BacktestEngine) shouldRebalance(currentDate time.Time, frequency string, lastRebalanceMonth *int) bool {
	currentMonth := int(currentDate.Month())
	currentYear := currentDate.Year()

	switch frequency {
	case "monthly":
		if *lastRebalanceMonth != currentMonth {
			*lastRebalanceMonth = currentMonth
			return true
		}
	case "quarterly":
		if currentMonth%3 == 1 && *lastRebalanceMonth != currentMonth { // January, April, July, October
			*lastRebalanceMonth = currentMonth
			return true
		}
	case "yearly":
		if currentMonth == 1 && *lastRebalanceMonth != currentYear*12+currentMonth {
			*lastRebalanceMonth = currentYear*12 + currentMonth
			return true
		}
	}

	return false
}
