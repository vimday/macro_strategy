package data

import (
	"encoding/json"
	"fmt"
	"io"
	"macro_strategy/internal/models"
	"net/http"
	"strconv"
	"time"
)

// BinanceProvider implements DataProvider for Binance API
type BinanceProvider struct {
	baseURL    string
	httpClient *http.Client
	rateLimit  time.Duration
	lastCall   time.Time
}

// NewBinanceProvider creates a new Binance data provider
func NewBinanceProvider() *BinanceProvider {
	return &BinanceProvider{
		baseURL:    "https://api.binance.com/api/v3",
		httpClient: &http.Client{Timeout: 30 * time.Second},
		rateLimit:  100 * time.Millisecond, // 10 requests per second
	}
}

// BinanceKlineResponse represents the response from Binance klines API
type BinanceKlineResponse [][]interface{}

// BinanceTickerResponse represents the response from Binance ticker API
type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// BinanceExchangeInfoResponse represents exchange info response
type BinanceExchangeInfoResponse struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	Symbols    []struct {
		Symbol     string `json:"symbol"`
		Status     string `json:"status"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
	} `json:"symbols"`
}

// GetHistoricalData fetches historical data from Binance
func (bp *BinanceProvider) GetHistoricalData(symbol string, startDate, endDate time.Time) ([]models.OHLCV, error) {
	// Rate limiting
	if time.Since(bp.lastCall) < bp.rateLimit {
		time.Sleep(bp.rateLimit - time.Since(bp.lastCall))
	}
	bp.lastCall = time.Now()

	// Convert symbol format (e.g., "BTC/USDT" -> "BTCUSDT")
	binanceSymbol := convertToBinanceSymbol(symbol)

	// Convert dates to milliseconds
	startTime := startDate.UnixMilli()
	endTime := endDate.UnixMilli()

	// Build API URL - using daily klines
	url := fmt.Sprintf("%s/klines?symbol=%s&interval=1d&startTime=%d&endTime=%d&limit=1000",
		bp.baseURL, binanceSymbol, startTime, endTime)

	// Make HTTP request
	resp, err := bp.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from Binance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Binance API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var klines BinanceKlineResponse
	if err := json.Unmarshal(body, &klines); err != nil {
		return nil, fmt.Errorf("failed to parse Binance response: %w", err)
	}

	if len(klines) == 0 {
		return nil, fmt.Errorf("no data found for symbol %s", symbol)
	}

	// Convert to OHLCV format
	var ohlcvData []models.OHLCV
	var prevClose float64

	for _, kline := range klines {
		if len(kline) < 11 {
			continue
		}

		// Parse timestamp (index 0)
		timestampMs, ok := kline[0].(float64)
		if !ok {
			continue
		}
		timestamp := time.UnixMilli(int64(timestampMs))

		// Parse OHLCV data
		open, err := parseFloat(kline[1])
		if err != nil {
			continue
		}
		high, err := parseFloat(kline[2])
		if err != nil {
			continue
		}
		low, err := parseFloat(kline[3])
		if err != nil {
			continue
		}
		close, err := parseFloat(kline[4])
		if err != nil {
			continue
		}
		volume, err := parseFloat(kline[5])
		if err != nil {
			continue
		}

		// Calculate percentage change
		var pctChg float64
		if prevClose > 0 {
			pctChg = ((close - prevClose) / prevClose) * 100
		}

		ohlcv := models.OHLCV{
			Date:   timestamp,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: int64(volume),
			PctChg: pctChg,
		}

		ohlcvData = append(ohlcvData, ohlcv)
		prevClose = close
	}

	if len(ohlcvData) == 0 {
		return nil, fmt.Errorf("no valid OHLCV data extracted for symbol %s", symbol)
	}

	return ohlcvData, nil
}

// GetLatestPrice fetches the latest price for a symbol
func (bp *BinanceProvider) GetLatestPrice(symbol string) (float64, error) {
	// Rate limiting
	if time.Since(bp.lastCall) < bp.rateLimit {
		time.Sleep(bp.rateLimit - time.Since(bp.lastCall))
	}
	bp.lastCall = time.Now()

	// Convert symbol format
	binanceSymbol := convertToBinanceSymbol(symbol)

	// Build API URL
	url := fmt.Sprintf("%s/ticker/price?symbol=%s", bp.baseURL, binanceSymbol)

	// Make HTTP request
	resp, err := bp.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch latest price from Binance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Binance API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var ticker BinanceTickerResponse
	if err := json.Unmarshal(body, &ticker); err != nil {
		return 0, fmt.Errorf("failed to parse Binance response: %w", err)
	}

	// Convert price string to float64
	price, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}

// IsValidSymbol checks if a symbol is valid for Binance
func (bp *BinanceProvider) IsValidSymbol(symbol string) bool {
	// Try to get latest price as a validation check
	_, err := bp.GetLatestPrice(symbol)
	return err == nil
}

// GetExchangeInfo returns exchange information for Binance
func (bp *BinanceProvider) GetExchangeInfo(symbol string) (map[string]interface{}, error) {
	// Rate limiting
	if time.Since(bp.lastCall) < bp.rateLimit {
		time.Sleep(bp.rateLimit - time.Since(bp.lastCall))
	}
	bp.lastCall = time.Now()

	// Build API URL
	url := fmt.Sprintf("%s/exchangeInfo", bp.baseURL)

	// Make HTTP request
	resp, err := bp.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange info from Binance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Binance API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var exchangeInfo BinanceExchangeInfoResponse
	if err := json.Unmarshal(body, &exchangeInfo); err != nil {
		return nil, fmt.Errorf("failed to parse Binance response: %w", err)
	}

	// Convert symbol format
	binanceSymbol := convertToBinanceSymbol(symbol)

	// Find symbol info
	for _, symbolInfo := range exchangeInfo.Symbols {
		if symbolInfo.Symbol == binanceSymbol {
			return map[string]interface{}{
				"symbol":      symbolInfo.Symbol,
				"status":      symbolInfo.Status,
				"base_asset":  symbolInfo.BaseAsset,
				"quote_asset": symbolInfo.QuoteAsset,
				"timezone":    exchangeInfo.Timezone,
				"server_time": time.UnixMilli(exchangeInfo.ServerTime),
				"exchange":    "Binance",
			}, nil
		}
	}

	return map[string]interface{}{
		"timezone":    exchangeInfo.Timezone,
		"server_time": time.UnixMilli(exchangeInfo.ServerTime),
		"exchange":    "Binance",
	}, nil
}

// Helper functions

// convertToBinanceSymbol converts symbol format for Binance API
func convertToBinanceSymbol(symbol string) string {
	// Convert "BTC/USDT" to "BTCUSDT"
	if len(symbol) > 4 && symbol[len(symbol)-5:len(symbol)-4] == "/" {
		return symbol[:len(symbol)-5] + symbol[len(symbol)-4:]
	}
	// If already in correct format, return as is
	return symbol
}

// parseFloat safely parses a float from interface{}
func parseFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("cannot parse %v as float64", val)
	}
}
