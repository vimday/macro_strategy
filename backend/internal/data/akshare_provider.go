package data

import (
	"encoding/json"
	"fmt"
	"macro_strategy/internal/models"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// AKShareProvider implements DataProvider interface using AKShare Python library
// This provider calls Python scripts to fetch A-share data via AKShare
type AKShareProvider struct {
	pythonPath   string
	scriptPath   string
	cacheEnabled bool
	cacheDir     string
}

// NewAKShareProvider creates a new AKShare data provider
func NewAKShareProvider(pythonPath, scriptPath string) *AKShareProvider {
	return &AKShareProvider{
		pythonPath:   pythonPath,
		scriptPath:   scriptPath,
		cacheEnabled: true,
		cacheDir:     "./data_cache",
	}
}

// GetHistoricalData fetches historical data using AKShare
func (a *AKShareProvider) GetHistoricalData(symbol string, startDate, endDate time.Time) ([]models.OHLCV, error) {
	// Convert symbol format for AKShare (e.g., "000300.SH" -> "sh000300")
	akSymbol := a.convertSymbolForAKShare(symbol)

	// Format dates for AKShare
	startDateStr := startDate.Format("20060102")
	endDateStr := endDate.Format("20060102")

	// Determine the appropriate AKShare command based on symbol type
	command := a.getAKShareCommand(symbol)

	// Call Python script with AKShare
	cmd := exec.Command(a.pythonPath, a.scriptPath, command,
		akSymbol, startDateStr, endDateStr)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute AKShare command: %w", err)
	}

	// Parse JSON output from Python script
	var rawData []map[string]interface{}
	if err := json.Unmarshal(output, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse AKShare output: %w", err)
	}

	// Convert to OHLCV format
	var data []models.OHLCV
	for _, row := range rawData {
		ohlcv, err := a.parseRowToOHLCV(row)
		if err != nil {
			continue // Skip invalid rows
		}
		data = append(data, ohlcv)
	}

	// Sort by date
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date.Before(data[j].Date)
	})

	return data, nil
}

// GetLatestPrice fetches the latest price using AKShare
func (a *AKShareProvider) GetLatestPrice(symbol string) (float64, error) {
	// Get recent data (last 5 days) and return the most recent close price
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -5)

	data, err := a.GetHistoricalData(symbol, startDate, endDate)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no recent data found for symbol: %s", symbol)
	}

	// Return the most recent close price
	return data[len(data)-1].Close, nil
}

// IsValidSymbol checks if a symbol is valid for A-share market
func (a *AKShareProvider) IsValidSymbol(symbol string) bool {
	// A-share symbols should match pattern: 6 digits + ".SH" or ".SZ"
	if len(symbol) != 9 {
		return false
	}

	parts := strings.Split(symbol, ".")
	if len(parts) != 2 {
		return false
	}

	code := parts[0]
	exchange := parts[1]

	// Check if code is 6 digits
	if len(code) != 6 {
		return false
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			return false
		}
	}

	// Check if exchange is valid
	return exchange == "SH" || exchange == "SZ"
}

// getAKShareCommand determines the appropriate AKShare command based on symbol type
func (a *AKShareProvider) getAKShareCommand(symbol string) string {
	if !a.IsValidSymbol(symbol) {
		return "get_stock_zh_a_hist" // default
	}

	parts := strings.Split(symbol, ".")
	code := parts[0]

	// Determine symbol type based on code prefix
	// Index codes: 000300, 000016, 000905, etc.
	// Stock codes: 600xxx (SH), 000xxx (SZ), 002xxx (SZ), 300xxx (SZ)
	if a.isIndexSymbol(code) {
		return "get_stock_zh_index_daily" // For indexes
	}

	return "get_stock_zh_a_hist" // For individual stocks
}

// isIndexSymbol checks if a code represents an index
func (a *AKShareProvider) isIndexSymbol(code string) bool {
	// Common A-share index patterns
	indexPatterns := []string{
		"000300", // CSI 300
		"000016", // SSE 50
		"000905", // CSI 500
		"000852", // CSI 1000
		"000688", // STAR 50
		"399006", // ChiNext
		"399330", // SZSE 100
		"000001", // SSE Composite (if .SH)
		"399001", // SZSE Composite (if .SZ)
	}

	for _, pattern := range indexPatterns {
		if code == pattern {
			return true
		}
	}

	// General index pattern check
	// Most A-share indexes start with 000, 399
	if strings.HasPrefix(code, "000") || strings.HasPrefix(code, "399") {
		return true
	}

	return false
}

// convertSymbolForAKShare converts symbol format from "000300.SH" to AKShare format
func (a *AKShareProvider) convertSymbolForAKShare(symbol string) string {
	parts := strings.Split(symbol, ".")
	if len(parts) != 2 {
		return symbol
	}

	code := parts[0]
	exchange := strings.ToLower(parts[1])

	// Convert to AKShare format: exchange + code
	return exchange + code
}

// parseRowToOHLCV converts a raw data row to OHLCV format
func (a *AKShareProvider) parseRowToOHLCV(row map[string]interface{}) (models.OHLCV, error) {
	var ohlcv models.OHLCV

	// Parse date
	dateStr, ok := row["日期"].(string)
	if !ok {
		return ohlcv, fmt.Errorf("invalid date field")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ohlcv, fmt.Errorf("failed to parse date: %w", err)
	}
	ohlcv.Date = date

	// Parse OHLC values
	if open, err := a.parseFloat(row["开盘"]); err == nil {
		ohlcv.Open = open
	}
	if high, err := a.parseFloat(row["最高"]); err == nil {
		ohlcv.High = high
	}
	if low, err := a.parseFloat(row["最低"]); err == nil {
		ohlcv.Low = low
	}
	if close, err := a.parseFloat(row["收盘"]); err == nil {
		ohlcv.Close = close
	}

	// Parse volume
	if volume, err := a.parseInt64(row["成交量"]); err == nil {
		ohlcv.Volume = volume
	}

	// Parse additional fields if available
	if amount, err := a.parseFloat(row["成交额"]); err == nil {
		ohlcv.Amount = amount
	}
	if turnover, err := a.parseFloat(row["换手率"]); err == nil {
		ohlcv.Turnover = turnover
	}
	if pctChg, err := a.parseFloat(row["涨跌幅"]); err == nil {
		ohlcv.PctChg = pctChg / 100 // Convert percentage to decimal
	}

	return ohlcv, nil
}

// parseFloat safely parses interface{} to float64
func (a *AKShareProvider) parseFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// parseInt64 safely parses interface{} to int64
func (a *AKShareProvider) parseInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// SetCacheEnabled enables or disables caching
func (a *AKShareProvider) SetCacheEnabled(enabled bool) {
	a.cacheEnabled = enabled
}

// SetCacheDir sets the cache directory
func (a *AKShareProvider) SetCacheDir(dir string) {
	a.cacheDir = dir
}
