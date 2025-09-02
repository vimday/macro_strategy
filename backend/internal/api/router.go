package api

import (
	"macro_strategy/internal/backtesting"
	"macro_strategy/internal/data"
	"macro_strategy/internal/models"
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
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Initialize services
	dataManager := data.NewDataSourceManager()
	dataManager.RegisterProvider(models.MarketTypeAShare, data.NewMockDataProvider())

	backtestEngine := backtesting.NewBacktestEngine()
	backtestService := services.NewBacktestService(dataManager, backtestEngine)

	// Initialize handlers
	handlers := NewHandlers(backtestService)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", handlers.HealthCheck)

		// Index endpoints
		v1.GET("/indexes", handlers.GetIndexes)
		v1.GET("/indexes/market/:market_type", handlers.GetIndexesByMarketType)
		v1.GET("/indexes/data/:id", handlers.GetIndexData)

		// Backtest endpoints
		v1.POST("/backtest", handlers.RunBacktest)
		v1.GET("/backtest/:id", handlers.GetBacktestResult)
	}

	return router
}
