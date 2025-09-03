package api

import (
	"macro_strategy/internal/models"
	"macro_strategy/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handlers contains all API handlers
type Handlers struct {
	backtestService *services.BacktestService
}

// NewHandlers creates a new handlers instance
func NewHandlers(backtestService *services.BacktestService) *Handlers {
	return &Handlers{
		backtestService: backtestService,
	}
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "macro-strategy-backend",
	})
}

// GetIndexes handles requests to get all available indexes
func (h *Handlers) GetIndexes(c *gin.Context) {
	indexes := models.PredefinedIndexes

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    indexes,
	})
}

// GetIndexesByMarketType handles requests to get indexes by market type
func (h *Handlers) GetIndexesByMarketType(c *gin.Context) {
	marketTypeStr := c.Param("market_type")
	marketType := models.MarketType(marketTypeStr)

	// Validate market type (legacy check, now we support all market types)
	// Just validate that it's a valid market type by checking if we get any results
	assets := models.GetIndexesByMarketType(marketType)
	if len(assets) == 0 && marketType != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No assets found for market type: " + string(marketType),
		})
		return
	}

	indexes := models.GetIndexesByMarketType(marketType)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    indexes,
	})
}

// GetIndexData handles requests to get historical data for an index
func (h *Handlers) GetIndexData(c *gin.Context) {
	indexID := c.Param("id")

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid start_date format, use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid end_date format, use YYYY-MM-DD",
		})
		return
	}

	// Get market data
	marketData, err := h.backtestService.GetMarketData(indexID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    marketData,
	})
}

// BacktestRequestJSON represents the JSON structure for backtest requests
type BacktestRequestJSON struct {
	IndexID     string             `json:"index_id" binding:"required"`
	Strategy    StrategyConfigJSON `json:"strategy" binding:"required"`
	StartDate   string             `json:"start_date" binding:"required"`
	EndDate     string             `json:"end_date" binding:"required"`
	InitialCash float64            `json:"initial_cash" binding:"required"`
}

// StrategyConfigJSON represents the JSON structure for strategy configuration
type StrategyConfigJSON struct {
	Type        string                 `json:"type" binding:"required"`
	Parameters  map[string]interface{} `json:"parameters"`
	Description string                 `json:"description"`
}

// RunBacktest handles backtest execution requests
func (h *Handlers) RunBacktest(c *gin.Context) {
	var requestJSON BacktestRequestJSON

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", requestJSON.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid start_date format, use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", requestJSON.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid end_date format, use YYYY-MM-DD",
		})
		return
	}

	// Convert to internal request format
	request := models.BacktestRequest{
		IndexID: requestJSON.IndexID,
		Strategy: models.StrategyConfig{
			Type:        models.StrategyType(requestJSON.Strategy.Type),
			Parameters:  requestJSON.Strategy.Parameters,
			Description: requestJSON.Strategy.Description,
		},
		StartDate:   startDate,
		EndDate:     endDate,
		InitialCash: requestJSON.InitialCash,
	}

	// Run backtest
	result, err := h.backtestService.RunBacktest(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetBacktestResult handles requests to get backtest results by ID
func (h *Handlers) GetBacktestResult(c *gin.Context) {
	backtestID := c.Param("id")

	result, err := h.backtestService.GetBacktestResult(backtestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Backtest result not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetAssets handles requests to get all available assets (updated from GetIndexes)
func (h *Handlers) GetAssets(c *gin.Context) {
	assets := models.GetAllAssets()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    assets,
	})
}

// GetAssetsByMarketType handles requests to get assets by market type (updated from GetIndexesByMarketType)
func (h *Handlers) GetAssetsByMarketType(c *gin.Context) {
	marketTypeStr := c.Param("market_type")
	marketType := models.MarketType(marketTypeStr)

	// Get assets for the market type
	assets := models.GetIndexesByMarketType(marketType)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    assets,
	})
}

// GetAssetData handles requests to get historical data for an asset (alias for GetIndexData)
func (h *Handlers) GetAssetData(c *gin.Context) {
	// Reuse the existing logic
	h.GetIndexData(c)
}

// GetSupportedMarkets handles requests to get all supported markets
func (h *Handlers) GetSupportedMarkets(c *gin.Context) {
	markets := h.backtestService.GetSupportedMarkets()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    markets,
	})
}

// GetSupportedStrategies handles requests to get all supported strategies
func (h *Handlers) GetSupportedStrategies(c *gin.Context) {
	strategies := h.backtestService.GetSupportedStrategies()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    strategies,
	})
}

// MultiStrategyBacktestRequestJSON represents the JSON structure for multi-strategy backtest requests
type MultiStrategyBacktestRequestJSON struct {
	AssetID       string                 `json:"asset_id" binding:"required"`
	Strategies    []StrategyConfigJSON   `json:"strategies" binding:"required"`
	StartDate     string                 `json:"start_date" binding:"required"`
	EndDate       string                 `json:"end_date" binding:"required"`
	InitialCash   float64                `json:"initial_cash" binding:"required"`
	Benchmark     string                 `json:"benchmark,omitempty"`
	DataSource    string                 `json:"data_source,omitempty"`
	ComparisonOpt *ComparisonOptionsJSON `json:"comparison_opt,omitempty"`
}

// ComparisonOptionsJSON represents the JSON structure for comparison options
type ComparisonOptionsJSON struct {
	ShowBenchmark      bool     `json:"show_benchmark"`
	NormalizeReturns   bool     `json:"normalize_returns"`
	ShowDrawdown       bool     `json:"show_drawdown"`
	ShowRollingMetrics bool     `json:"show_rolling_metrics"`
	RollingWindow      int      `json:"rolling_window"`
	Metrics            []string `json:"metrics"`
}

// RunMultiStrategyBacktest handles multi-strategy backtest execution requests
func (h *Handlers) RunMultiStrategyBacktest(c *gin.Context) {
	var requestJSON MultiStrategyBacktestRequestJSON

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", requestJSON.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid start_date format, use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", requestJSON.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid end_date format, use YYYY-MM-DD",
		})
		return
	}

	// Convert strategies
	var strategies []models.StrategyConfig
	for _, strategyJSON := range requestJSON.Strategies {
		strategy := models.StrategyConfig{
			Type:        models.StrategyType(strategyJSON.Type),
			Parameters:  strategyJSON.Parameters,
			Description: strategyJSON.Description,
		}
		strategies = append(strategies, strategy)
	}

	// Convert comparison options
	var comparisonOpt *models.ComparisonOptions
	if requestJSON.ComparisonOpt != nil {
		comparisonOpt = &models.ComparisonOptions{
			ShowBenchmark:      requestJSON.ComparisonOpt.ShowBenchmark,
			NormalizeReturns:   requestJSON.ComparisonOpt.NormalizeReturns,
			ShowDrawdown:       requestJSON.ComparisonOpt.ShowDrawdown,
			ShowRollingMetrics: requestJSON.ComparisonOpt.ShowRollingMetrics,
			RollingWindow:      requestJSON.ComparisonOpt.RollingWindow,
			Metrics:            requestJSON.ComparisonOpt.Metrics,
		}
	}

	// Convert to internal request format
	request := models.MultiStrategyBacktestRequest{
		AssetID:       requestJSON.AssetID,
		Strategies:    strategies,
		StartDate:     startDate,
		EndDate:       endDate,
		InitialCash:   requestJSON.InitialCash,
		Benchmark:     requestJSON.Benchmark,
		DataSource:    requestJSON.DataSource,
		ComparisonOpt: comparisonOpt,
	}

	// Run multi-strategy backtest
	result, err := h.backtestService.RunMultiStrategyBacktest(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetMultiStrategyResult handles requests to get multi-strategy backtest results by ID
func (h *Handlers) GetMultiStrategyResult(c *gin.Context) {
	backtestID := c.Param("id")

	result, err := h.backtestService.GetMultiStrategyResult(backtestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Multi-strategy backtest result not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
