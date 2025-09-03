package services

import (
	"fmt"
	"macro_strategy/internal/backtesting"
	"macro_strategy/internal/data"
	"macro_strategy/internal/models"
	"sync"
	"time"
)

// BacktestService handles backtest-related business logic
type BacktestService struct {
	dataManager          *data.DataSourceManager
	backtestEngine       *backtesting.BacktestEngine
	multiStrategyService *MultiStrategyService
	resultCache          map[string]*models.BacktestResult
	multiResultCache     map[string]*models.MultiStrategyBacktestResult
	cacheMutex           sync.RWMutex
}

// NewBacktestService creates a new backtest service
func NewBacktestService(dataManager *data.DataSourceManager, backtestEngine *backtesting.BacktestEngine) *BacktestService {
	bs := &BacktestService{
		dataManager:      dataManager,
		backtestEngine:   backtestEngine,
		resultCache:      make(map[string]*models.BacktestResult),
		multiResultCache: make(map[string]*models.MultiStrategyBacktestResult),
	}

	// Initialize multi-strategy service
	bs.multiStrategyService = NewMultiStrategyService(backtestEngine, dataManager)

	return bs
}

// GetMarketData retrieves market data for a given index and date range
func (bs *BacktestService) GetMarketData(indexID string, startDate, endDate time.Time) (*models.MarketData, error) {
	// Find the index
	index := models.GetIndexByID(indexID)
	if index == nil {
		return nil, fmt.Errorf("index not found: %s", indexID)
	}

	// Get market data
	marketData, err := bs.dataManager.GetMarketData(index, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	return marketData, nil
}

// RunBacktest executes a backtest and returns the results
func (bs *BacktestService) RunBacktest(request models.BacktestRequest) (*models.BacktestResult, error) {
	// Support both AssetID and IndexID for backward compatibility
	assetID := request.AssetID
	if assetID == "" {
		assetID = request.IndexID
	}

	// Validate the asset exists
	index := models.GetIndexByID(assetID)
	if index == nil {
		return nil, fmt.Errorf("asset not found: %s", assetID)
	}

	// Get market data for the backtest period
	marketData, err := bs.dataManager.GetMarketData(index, request.StartDate, request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	// Run the backtest
	result, err := bs.backtestEngine.RunBacktest(request, marketData)
	if err != nil {
		return nil, fmt.Errorf("backtest execution failed: %w", err)
	}

	// Cache the result
	bs.cacheMutex.Lock()
	bs.resultCache[result.ID] = result
	bs.cacheMutex.Unlock()

	return result, nil
}

// GetBacktestResult retrieves a backtest result by ID
func (bs *BacktestService) GetBacktestResult(backtestID string) (*models.BacktestResult, error) {
	bs.cacheMutex.RLock()
	result, exists := bs.resultCache[backtestID]
	bs.cacheMutex.RUnlock()

	if !exists {
		return nil, nil // Not found, but not an error
	}

	return result, nil
}

// ListBacktestResults returns all cached backtest results
func (bs *BacktestService) ListBacktestResults() []*models.BacktestResult {
	bs.cacheMutex.RLock()
	defer bs.cacheMutex.RUnlock()

	results := make([]*models.BacktestResult, 0, len(bs.resultCache))
	for _, result := range bs.resultCache {
		results = append(results, result)
	}

	return results
}

// ClearCache clears all cached backtest results
func (bs *BacktestService) ClearCache() {
	bs.cacheMutex.Lock()
	bs.resultCache = make(map[string]*models.BacktestResult)
	bs.cacheMutex.Unlock()
}

// ValidateBacktestRequest validates a backtest request
func (bs *BacktestService) ValidateBacktestRequest(request models.BacktestRequest) error {
	// Support both AssetID and IndexID for backward compatibility
	assetID := request.AssetID
	if assetID == "" {
		assetID = request.IndexID
	}

	// Check if asset exists
	index := models.GetIndexByID(assetID)
	if index == nil {
		return fmt.Errorf("invalid asset ID: %s", assetID)
	}

	// Validate dates
	if request.StartDate.After(request.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}

	// Check if date range is reasonable (not too far in the past or future)
	now := time.Now()
	if request.EndDate.After(now) {
		return fmt.Errorf("end date cannot be in the future")
	}

	// Minimum backtest period (at least 30 days)
	if request.EndDate.Sub(request.StartDate) < 30*24*time.Hour {
		return fmt.Errorf("backtest period must be at least 30 days")
	}

	// Validate initial cash
	if request.InitialCash <= 0 {
		return fmt.Errorf("initial cash must be positive")
	}

	// Validate strategy
	if request.Strategy.Type == "" {
		return fmt.Errorf("strategy type is required")
	}

	// Strategy-specific validations
	switch request.Strategy.Type {
	case models.StrategyTypeMonthlyRotation:
		if err := bs.validateMonthlyRotationStrategy(request.Strategy); err != nil {
			return fmt.Errorf("invalid monthly rotation strategy: %w", err)
		}
	default:
		return fmt.Errorf("unsupported strategy type: %s", request.Strategy.Type)
	}

	return nil
}

// validateMonthlyRotationStrategy validates monthly rotation strategy parameters
func (bs *BacktestService) validateMonthlyRotationStrategy(strategy models.StrategyConfig) error {
	if strategy.Parameters == nil {
		return fmt.Errorf("strategy parameters are required")
	}

	// Check buy days parameter
	buyDaysRaw, ok := strategy.Parameters["buy_days_before_month_end"]
	if !ok {
		return fmt.Errorf("buy_days_before_month_end parameter is required")
	}

	buyDays, ok := buyDaysRaw.(float64)
	if !ok {
		return fmt.Errorf("buy_days_before_month_end must be a number")
	}

	if buyDays < 1 || buyDays > 20 {
		return fmt.Errorf("buy_days_before_month_end must be between 1 and 20")
	}

	// Check sell days parameter
	sellDaysRaw, ok := strategy.Parameters["sell_days_after_month_start"]
	if !ok {
		return fmt.Errorf("sell_days_after_month_start parameter is required")
	}

	sellDays, ok := sellDaysRaw.(float64)
	if !ok {
		return fmt.Errorf("sell_days_after_month_start must be a number")
	}

	if sellDays < 1 || sellDays > 20 {
		return fmt.Errorf("sell_days_after_month_start must be between 1 and 20")
	}

	return nil
}

// RunMultiStrategyBacktest executes multiple strategies and compares results
func (bs *BacktestService) RunMultiStrategyBacktest(request models.MultiStrategyBacktestRequest) (*models.MultiStrategyBacktestResult, error) {
	// Run the multi-strategy backtest
	result, err := bs.multiStrategyService.RunMultiStrategyBacktest(request)
	if err != nil {
		return nil, fmt.Errorf("multi-strategy backtest execution failed: %w", err)
	}

	// Cache the result
	bs.cacheMutex.Lock()
	bs.multiResultCache[result.ID] = result
	bs.cacheMutex.Unlock()

	return result, nil
}

// GetMultiStrategyResult retrieves a multi-strategy backtest result by ID
func (bs *BacktestService) GetMultiStrategyResult(backtestID string) (*models.MultiStrategyBacktestResult, error) {
	bs.cacheMutex.RLock()
	result, exists := bs.multiResultCache[backtestID]
	bs.cacheMutex.RUnlock()

	if !exists {
		return nil, nil // Not found, but not an error
	}

	return result, nil
}

// GetSupportedStrategies returns a list of supported strategy types with their parameters
func (bs *BacktestService) GetSupportedStrategies() map[string]interface{} {
	return map[string]interface{}{
		"monthly_rotation": map[string]interface{}{
			"name":        "Monthly Rotation Strategy",
			"description": "Buy before month-end, sell after month-start",
			"parameters": map[string]interface{}{
				"buy_days_before_month_end": map[string]interface{}{
					"type":        "integer",
					"default":     1,
					"range":       []int{1, 20},
					"description": "Number of days before month-end to buy",
				},
				"sell_days_after_month_start": map[string]interface{}{
					"type":        "integer",
					"default":     1,
					"range":       []int{1, 20},
					"description": "Number of days after month-start to sell",
				},
			},
		},
		"buy_and_hold": map[string]interface{}{
			"name":        "Buy and Hold Strategy",
			"description": "Buy and hold with optional rebalancing",
			"parameters": map[string]interface{}{
				"target_allocation": map[string]interface{}{
					"type":        "float",
					"default":     1.0,
					"range":       []float64{0.1, 1.0},
					"description": "Target allocation percentage (0.1 = 10%, 1.0 = 100%)",
				},
				"rebalance_frequency": map[string]interface{}{
					"type":        "string",
					"default":     "never",
					"options":     []string{"never", "monthly", "quarterly", "yearly"},
					"description": "How often to rebalance the portfolio",
				},
				"dividend_reinvest": map[string]interface{}{
					"type":        "boolean",
					"default":     false,
					"description": "Whether to reinvest dividends (future feature)",
				},
			},
		},
	}
}

// GetSupportedMarkets returns a list of supported markets with their assets
func (bs *BacktestService) GetSupportedMarkets() map[string]interface{} {
	return map[string]interface{}{
		"a_share_index": map[string]interface{}{
			"name":        "A-Share Indexes",
			"description": "Chinese stock market indexes",
			"assets":      models.GetIndexesByMarketType(models.MarketTypeAShareIndex),
			"data_source": "AKShare",
		},
		"a_share_stock": map[string]interface{}{
			"name":        "A-Share Individual Stocks",
			"description": "Chinese individual stocks",
			"assets":      models.GetIndexesByMarketType(models.MarketTypeAShareStock),
			"data_source": "AKShare",
		},
		"us_index": map[string]interface{}{
			"name":        "US Market Indexes",
			"description": "US stock market indexes and ETFs",
			"assets":      models.GetIndexesByMarketType(models.MarketTypeUSIndex),
			"data_source": "Yahoo Finance",
		},
		"us_stock": map[string]interface{}{
			"name":        "US Individual Stocks",
			"description": "US individual stocks",
			"assets":      models.GetIndexesByMarketType(models.MarketTypeUSStock),
			"data_source": "Yahoo Finance",
		},
		"crypto": map[string]interface{}{
			"name":        "Cryptocurrency",
			"description": "Digital currencies and tokens",
			"assets":      models.GetIndexesByMarketType(models.MarketTypeCrypto),
			"data_source": "Binance",
		},
		"hk_index": map[string]interface{}{
			"name":        "Hong Kong Indexes",
			"description": "Hong Kong stock market indexes",
			"assets":      models.GetIndexesByMarketType(models.MarketTypeHKIndex),
			"data_source": "Yahoo Finance",
		},
		"hk_stock": map[string]interface{}{
			"name":        "Hong Kong Stocks",
			"description": "Hong Kong individual stocks",
			"assets":      models.GetIndexesByMarketType(models.MarketTypeHKStock),
			"data_source": "Yahoo Finance",
		},
	}
}
