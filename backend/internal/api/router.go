package api

import (
	"macro_strategy/internal/backtesting"
	"macro_strategy/internal/data"
	"macro_strategy/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter() *gin.Engine {
	// Create Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Initialize services
	dataManager := data.NewDataSourceManager() // This now auto-registers all providers

	backtestEngine := backtesting.NewBacktestEngine()
	backtestService := services.NewBacktestService(dataManager, backtestEngine)

	// Initialize handlers
	handlers := NewHandlers(backtestService)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", handlers.HealthCheck)

		// Asset and market endpoints
		v1.GET("/assets", handlers.GetAssets) // Updated: renamed from indexes
		v1.GET("/assets/market/:market_type", handlers.GetAssetsByMarketType)
		v1.GET("/assets/data/:id", handlers.GetAssetData)
		v1.GET("/markets", handlers.GetSupportedMarkets) // New: get all supported markets

		// Strategy endpoints
		v1.GET("/strategies", handlers.GetSupportedStrategies) // New: get all supported strategies

		// Single strategy backtest endpoints
		v1.POST("/backtest", handlers.RunBacktest)
		v1.GET("/backtest/:id", handlers.GetBacktestResult)

		// Multi-strategy comparison endpoints
		v1.POST("/backtest/multi", handlers.RunMultiStrategyBacktest)  // New: multi-strategy comparison
		v1.GET("/backtest/multi/:id", handlers.GetMultiStrategyResult) // New: get multi-strategy results

		// Backward compatibility - keep old index endpoints
		v1.GET("/indexes", handlers.GetIndexes)
		v1.GET("/indexes/market/:market_type", handlers.GetIndexesByMarketType)
		v1.GET("/indexes/data/:id", handlers.GetIndexData)
	}

	return router
}
