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

	// Validate market type
	if marketType != models.MarketTypeAShare &&
		marketType != models.MarketTypeCrypto &&
		marketType != models.MarketTypeHKUS {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid market type",
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
