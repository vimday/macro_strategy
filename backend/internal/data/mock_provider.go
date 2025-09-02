package data

import (
	"fmt"
	"macro_strategy/internal/models"
	"math"
	"math/rand"
	"time"
)

// MockDataProvider implements DataProvider interface for testing
// This simulates realistic market data without requiring external APIs
type MockDataProvider struct {
	seed int64
}

// NewMockDataProvider creates a new mock data provider
func NewMockDataProvider() *MockDataProvider {
	return &MockDataProvider{
		seed: time.Now().UnixNano(),
	}
}

// GetHistoricalData generates simulated historical data
func (m *MockDataProvider) GetHistoricalData(symbol string, startDate, endDate time.Time) ([]models.OHLCV, error) {
	if !m.IsValidSymbol(symbol) {
		return nil, fmt.Errorf("invalid symbol: %s", symbol)
	}

	// Use symbol as seed for consistent data generation
	rng := rand.New(rand.NewSource(m.generateSeed(symbol)))

	var data []models.OHLCV
	currentDate := startDate
	basePrice := m.getBasePriceForSymbol(symbol)

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		// Skip weekends for stock data
		if currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		ohlcv := m.generateDailyData(rng, currentDate, basePrice)
		data = append(data, ohlcv)

		// Update base price for next day (with some trend)
		basePrice = ohlcv.Close
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return data, nil
}

// GetLatestPrice returns a simulated latest price
func (m *MockDataProvider) GetLatestPrice(symbol string) (float64, error) {
	if !m.IsValidSymbol(symbol) {
		return 0, fmt.Errorf("invalid symbol: %s", symbol)
	}

	basePrice := m.getBasePriceForSymbol(symbol)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Add some random variation (-2% to +2%)
	variation := (rng.Float64() - 0.5) * 0.04
	return basePrice * (1 + variation), nil
}

// IsValidSymbol checks if a symbol is valid (for our predefined indexes)
func (m *MockDataProvider) IsValidSymbol(symbol string) bool {
	validSymbols := map[string]bool{
		"000300.SH": true, // CSI 300
		"000016.SH": true, // SSE 50
		"000905.SH": true, // CSI 500
		"000852.SH": true, // CSI 1000
		"000688.SH": true, // STAR 50
		"399006.SZ": true, // ChiNext
		"399330.SZ": true, // SZSE 100
	}
	return validSymbols[symbol]
}

// generateSeed creates a consistent seed based on symbol
func (m *MockDataProvider) generateSeed(symbol string) int64 {
	hash := int64(0)
	for _, c := range symbol {
		hash = hash*31 + int64(c)
	}
	return hash + m.seed
}

// getBasePriceForSymbol returns a realistic base price for different indexes
func (m *MockDataProvider) getBasePriceForSymbol(symbol string) float64 {
	basePrices := map[string]float64{
		"000300.SH": 4200.0, // CSI 300
		"000016.SH": 3100.0, // SSE 50
		"000905.SH": 6800.0, // CSI 500
		"000852.SH": 6500.0, // CSI 1000
		"000688.SH": 1200.0, // STAR 50
		"399006.SZ": 2400.0, // ChiNext
		"399330.SZ": 5200.0, // SZSE 100
	}

	if price, exists := basePrices[symbol]; exists {
		return price
	}
	return 3000.0 // Default price
}

// generateDailyData creates realistic OHLCV data for a single day
func (m *MockDataProvider) generateDailyData(rng *rand.Rand, date time.Time, basePrice float64) models.OHLCV {
	// Generate daily return with some volatility (normally distributed)
	dailyReturn := rng.NormFloat64() * 0.02 // 2% daily volatility

	// Calculate close price
	close := basePrice * (1 + dailyReturn)

	// Generate intraday volatility
	intradayVol := rng.Float64() * 0.01 // 1% intraday volatility

	// Generate open (with gap)
	gapReturn := (rng.Float64() - 0.5) * 0.005 // Small overnight gap
	open := basePrice * (1 + gapReturn)

	// Generate high and low based on close and volatility
	high := math.Max(open, close) * (1 + intradayVol)
	low := math.Min(open, close) * (1 - intradayVol)

	// Ensure OHLC relationships are valid
	high = math.Max(high, math.Max(open, close))
	low = math.Min(low, math.Min(open, close))

	// Generate volume (with some randomness)
	baseVolume := int64(100000000)          // 100M base volume
	volumeMultiplier := 0.5 + rng.Float64() // 0.5x to 1.5x variation
	volume := int64(float64(baseVolume) * volumeMultiplier)

	return models.OHLCV{
		Date:   date,
		Open:   math.Round(open*100) / 100, // Round to 2 decimal places
		High:   math.Round(high*100) / 100,
		Low:    math.Round(low*100) / 100,
		Close:  math.Round(close*100) / 100,
		Volume: volume,
	}
}
