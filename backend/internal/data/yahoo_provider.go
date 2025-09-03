package data

import (
	"encoding/json"
	"fmt"
	"io"
	"macro_strategy/internal/models"
	"net/http"
	"time"
)

// YahooProvider implements DataProvider for Yahoo Finance API
type YahooProvider struct {
	baseURL    string
	httpClient *http.Client
	rateLimit  time.Duration
	lastCall   time.Time
}

// NewYahooProvider creates a new Yahoo Finance data provider
func NewYahooProvider() *YahooProvider {
	return &YahooProvider{
		baseURL:    "https://query1.finance.yahoo.com/v8/finance/chart",
		httpClient: &http.Client{Timeout: 30 * time.Second},
		rateLimit:  100 * time.Millisecond, // 10 requests per second
	}
}

// YahooResponse represents the response structure from Yahoo Finance API
type YahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int64   `json:"firstTradeDate"`
				RegularMarketTime    int64   `json:"regularMarketTime"`
				Gmtoffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
				Adjclose []struct {
					Adjclose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

// GetHistoricalData fetches historical data from Yahoo Finance
func (yp *YahooProvider) GetHistoricalData(symbol string, startDate, endDate time.Time) ([]models.OHLCV, error) {
	// Rate limiting
	if time.Since(yp.lastCall) < yp.rateLimit {
		time.Sleep(yp.rateLimit - time.Since(yp.lastCall))
	}
	yp.lastCall = time.Now()

	// Convert dates to Unix timestamps
	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()

	// Build API URL
	url := fmt.Sprintf("%s/%s?period1=%d&period2=%d&interval=1d&includePrePost=false",
		yp.baseURL, symbol, startTimestamp, endTimestamp)

	// Make HTTP request
	resp, err := yp.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from Yahoo Finance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Yahoo Finance API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var yahooResp YahooResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return nil, fmt.Errorf("failed to parse Yahoo Finance response: %w", err)
	}

	// Check for API errors
	if yahooResp.Chart.Error != nil {
		return nil, fmt.Errorf("Yahoo Finance API error: %v", yahooResp.Chart.Error)
	}

	// Check if we have results
	if len(yahooResp.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data found for symbol %s", symbol)
	}

	result := yahooResp.Chart.Result[0]
	if len(result.Timestamp) == 0 {
		return nil, fmt.Errorf("no historical data available for symbol %s", symbol)
	}

	// Extract OHLCV data
	var ohlcvData []models.OHLCV
	quotes := result.Indicators.Quote[0]

	for i, timestamp := range result.Timestamp {
		// Skip if any required data is missing
		if i >= len(quotes.Open) || i >= len(quotes.High) ||
			i >= len(quotes.Low) || i >= len(quotes.Close) || i >= len(quotes.Volume) {
			continue
		}

		// Skip if any value is nil/invalid
		if quotes.Open[i] == 0 || quotes.High[i] == 0 ||
			quotes.Low[i] == 0 || quotes.Close[i] == 0 {
			continue
		}

		ohlcv := models.OHLCV{
			Date:   time.Unix(timestamp, 0),
			Open:   quotes.Open[i],
			High:   quotes.High[i],
			Low:    quotes.Low[i],
			Close:  quotes.Close[i],
			Volume: quotes.Volume[i],
		}

		// Calculate percentage change if possible
		if i > 0 && quotes.Close[i-1] != 0 {
			ohlcv.PctChg = ((quotes.Close[i] - quotes.Close[i-1]) / quotes.Close[i-1]) * 100
		}

		ohlcvData = append(ohlcvData, ohlcv)
	}

	if len(ohlcvData) == 0 {
		return nil, fmt.Errorf("no valid OHLCV data extracted for symbol %s", symbol)
	}

	return ohlcvData, nil
}

// GetLatestPrice fetches the latest price for a symbol
func (yp *YahooProvider) GetLatestPrice(symbol string) (float64, error) {
	// Rate limiting
	if time.Since(yp.lastCall) < yp.rateLimit {
		time.Sleep(yp.rateLimit - time.Since(yp.lastCall))
	}
	yp.lastCall = time.Now()

	// Get data for the last day
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -1)

	url := fmt.Sprintf("%s/%s?period1=%d&period2=%d&interval=1d&includePrePost=false",
		yp.baseURL, symbol, startDate.Unix(), endDate.Unix())

	resp, err := yp.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch latest price from Yahoo Finance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Yahoo Finance API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var yahooResp YahooResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return 0, fmt.Errorf("failed to parse Yahoo Finance response: %w", err)
	}

	if yahooResp.Chart.Error != nil {
		return 0, fmt.Errorf("Yahoo Finance API error: %v", yahooResp.Chart.Error)
	}

	if len(yahooResp.Chart.Result) == 0 || len(yahooResp.Chart.Result[0].Indicators.Quote) == 0 {
		return 0, fmt.Errorf("no price data found for symbol %s", symbol)
	}

	result := yahooResp.Chart.Result[0]

	// Try to get from meta first (most recent price)
	if result.Meta.RegularMarketPrice > 0 {
		return result.Meta.RegularMarketPrice, nil
	}

	// Otherwise get the latest close price
	quotes := result.Indicators.Quote[0]
	if len(quotes.Close) > 0 {
		for i := len(quotes.Close) - 1; i >= 0; i-- {
			if quotes.Close[i] > 0 {
				return quotes.Close[i], nil
			}
		}
	}

	return 0, fmt.Errorf("no valid price data found for symbol %s", symbol)
}

// IsValidSymbol checks if a symbol is valid for Yahoo Finance
func (yp *YahooProvider) IsValidSymbol(symbol string) bool {
	// Try to get latest price as a validation check
	_, err := yp.GetLatestPrice(symbol)
	return err == nil
}

// GetExchangeInfo returns exchange information for a symbol
func (yp *YahooProvider) GetExchangeInfo(symbol string) (map[string]interface{}, error) {
	// Rate limiting
	if time.Since(yp.lastCall) < yp.rateLimit {
		time.Sleep(yp.rateLimit - time.Since(yp.lastCall))
	}
	yp.lastCall = time.Now()

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -1)

	url := fmt.Sprintf("%s/%s?period1=%d&period2=%d&interval=1d&includePrePost=false",
		yp.baseURL, symbol, startDate.Unix(), endDate.Unix())

	resp, err := yp.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange info from Yahoo Finance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Yahoo Finance API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var yahooResp YahooResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return nil, fmt.Errorf("failed to parse Yahoo Finance response: %w", err)
	}

	if yahooResp.Chart.Error != nil {
		return nil, fmt.Errorf("Yahoo Finance API error: %v", yahooResp.Chart.Error)
	}

	if len(yahooResp.Chart.Result) == 0 {
		return nil, fmt.Errorf("no exchange info found for symbol %s", symbol)
	}

	meta := yahooResp.Chart.Result[0].Meta
	return map[string]interface{}{
		"currency":               meta.Currency,
		"exchange_name":          meta.ExchangeName,
		"instrument_type":        meta.InstrumentType,
		"timezone":               meta.Timezone,
		"exchange_timezone_name": meta.ExchangeTimezoneName,
		"first_trade_date":       time.Unix(meta.FirstTradeDate, 0),
		"regular_market_time":    time.Unix(meta.RegularMarketTime, 0),
		"regular_market_price":   meta.RegularMarketPrice,
	}, nil
}
