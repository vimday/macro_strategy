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

// NewDataSourceManager creates a new data source manager
func NewDataSourceManager() *DataSourceManager {
	return &DataSourceManager{
		providers: make(map[models.MarketType]DataProvider),
	}
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
