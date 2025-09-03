package data

import (
	"fmt"
	"macro_strategy/internal/models"
	"time"
)

// DataProvider interface defines methods for fetching market data
type DataProvider interface {
	GetHistoricalData(symbol string, startDate, endDate time.Time) ([]models.OHLCV, error)
	GetLatestPrice(symbol string) (float64, error)
	IsValidSymbol(symbol string) bool
}

// DataSourceManager manages different data providers
type DataSourceManager struct {
	providers map[models.MarketType]DataProvider
}

// NewDataSourceManager creates a new data source manager with all providers registered
func NewDataSourceManager() *DataSourceManager {
	dsm := &DataSourceManager{
		providers: make(map[models.MarketType]DataProvider),
	}

	// Register A-share providers
	// Use the virtual environment Python path (relative to backend directory)
	akshareProvider := NewAKShareProvider("../akshare_env/bin/python3", "../backend/scripts/akshare_client.py")
	dsm.RegisterProvider(models.MarketTypeAShareIndex, akshareProvider)
	dsm.RegisterProvider(models.MarketTypeAShareStock, akshareProvider)

	// Register US market provider
	yahooProvider := NewYahooProvider()
	dsm.RegisterProvider(models.MarketTypeUSIndex, yahooProvider)
	dsm.RegisterProvider(models.MarketTypeUSStock, yahooProvider)

	// Register crypto provider
	binanceProvider := NewBinanceProvider()
	dsm.RegisterProvider(models.MarketTypeCrypto, binanceProvider)

	// Register HK market provider (also using Yahoo)
	dsm.RegisterProvider(models.MarketTypeHKIndex, yahooProvider)
	dsm.RegisterProvider(models.MarketTypeHKStock, yahooProvider)

	// Register ETF provider (Yahoo can handle most ETFs)
	dsm.RegisterProvider(models.MarketTypeETF, yahooProvider)

	return dsm
}

// RegisterProvider registers a data provider for a specific market type
func (dsm *DataSourceManager) RegisterProvider(marketType models.MarketType, provider DataProvider) {
	dsm.providers[marketType] = provider
}

// GetProvider returns the data provider for a specific market type
func (dsm *DataSourceManager) GetProvider(marketType models.MarketType) (DataProvider, error) {
	provider, exists := dsm.providers[marketType]
	if !exists {
		return nil, fmt.Errorf("no data provider registered for market type: %s", marketType)
	}
	return provider, nil
}

// GetMarketData fetches historical data for an asset with enhanced support
func (dsm *DataSourceManager) GetMarketData(index *models.Index, startDate, endDate time.Time) (*models.MarketData, error) {
	provider, err := dsm.GetProvider(index.MarketType)
	if err != nil {
		return nil, err
	}

	data, err := provider.GetHistoricalData(index.Symbol, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data for %s: %w", index.Symbol, err)
	}

	return &models.MarketData{
		AssetID:    index.ID,
		Symbol:     index.Symbol,
		MarketType: index.MarketType,
		AssetClass: index.AssetClass,
		Currency:   index.Currency,
		Data:       data,
		LastUpdate: time.Now(),
		Metadata:   make(map[string]interface{}),
	}, nil
}
