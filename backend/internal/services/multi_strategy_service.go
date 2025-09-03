package services

import (
	"fmt"
	"macro_strategy/internal/backtesting"
	"macro_strategy/internal/data"
	"macro_strategy/internal/models"
	"math"
	"sort"
	"time"
)

// MultiStrategyService handles multi-strategy comparison backtesting
type MultiStrategyService struct {
	backtestEngine *backtesting.BacktestEngine
	dataManager   *data.DataSourceManager
}

// NewMultiStrategyService creates a new multi-strategy service
func NewMultiStrategyService(backtestEngine *backtesting.BacktestEngine, dataManager *data.DataSourceManager) *MultiStrategyService {
	return &MultiStrategyService{
		backtestEngine: backtestEngine,
		dataManager:   dataManager,
	}
}

// RunMultiStrategyBacktest runs backtests for multiple strategies and compares results
func (mss *MultiStrategyService) RunMultiStrategyBacktest(request models.MultiStrategyBacktestRequest) (*models.MultiStrategyBacktestResult, error) {
	startTime := time.Now()

	// Validate request
	if err := mss.validateMultiStrategyRequest(request); err != nil {
		return nil, fmt.Errorf("invalid multi-strategy request: %w", err)
	}

	// Get asset information
	asset := models.GetIndexByID(request.AssetID)
	if asset == nil {
		return nil, fmt.Errorf("asset not found: %s", request.AssetID)
	}

	// Fetch market data
	marketData, err := mss.dataManager.GetMarketData(asset, request.StartDate, request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market data: %w", err)
	}

	// Run backtests for each strategy
	var results []models.BacktestResult
	for i, strategy := range request.Strategies {
		// Create individual backtest request
		backtestRequest := models.BacktestRequest{
			AssetID:     request.AssetID,
			IndexID:     request.AssetID, // For backward compatibility
			Strategy:    strategy,
			StartDate:   request.StartDate,
			EndDate:     request.EndDate,
			InitialCash: request.InitialCash,
			Benchmark:   request.Benchmark,
			DataSource:  request.DataSource,
			Metadata: map[string]interface{}{
				"strategy_index": i,
				"strategy_name":  fmt.Sprintf("%s_%d", strategy.Type, i+1),
			},
		}

		// Run backtest
		result, err := mss.backtestEngine.RunBacktest(backtestRequest, marketData)
		if err != nil {
			return nil, fmt.Errorf("failed to run backtest for strategy %d (%s): %w", i+1, strategy.Type, err)
		}

		results = append(results, *result)
	}

	// Run benchmark backtest if specified
	var benchmarkResult *models.BacktestResult
	if request.Benchmark != "" && request.Benchmark != request.AssetID {
		benchmarkAsset := models.GetIndexByID(request.Benchmark)
		if benchmarkAsset != nil {
			benchmarkData, err := mss.dataManager.GetMarketData(benchmarkAsset, request.StartDate, request.EndDate)
			if err == nil {
				// Create a simple buy-and-hold strategy for benchmark
				benchmarkStrategy := models.StrategyConfig{
					Type: models.StrategyTypeBuyAndHold,
					Parameters: map[string]interface{}{
						"target_allocation": 1.0,
						"rebalance_frequency": "never",
					},
					Description: fmt.Sprintf("Benchmark: %s", request.Benchmark),
				}

				benchmarkRequest := models.BacktestRequest{
					AssetID:     request.Benchmark,
					IndexID:     request.Benchmark,
					Strategy:    benchmarkStrategy,
					StartDate:   request.StartDate,
					EndDate:     request.EndDate,
					InitialCash: request.InitialCash,
					Metadata: map[string]interface{}{
						"is_benchmark": true,
					},
				}

				benchmarkRes, err := mss.backtestEngine.RunBacktest(benchmarkRequest, benchmarkData)
				if err == nil {
					benchmarkResult = benchmarkRes
				}
			}
		}
	}

	// Perform strategy comparison analysis
	comparison := mss.performStrategyComparison(results, request.ComparisonOpt)

	// Create multi-strategy result
	result := &models.MultiStrategyBacktestResult{
		ID:              generateMultiStrategyID(),
		Request:         request,
		Results:         results,
		Comparison:      comparison,
		BenchmarkResult: benchmarkResult,
		CreatedAt:       startTime,
		Duration:        time.Since(startTime),
	}

	return result, nil
}

// validateMultiStrategyRequest validates the multi-strategy backtest request
func (mss *MultiStrategyService) validateMultiStrategyRequest(request models.MultiStrategyBacktestRequest) error {
	if request.AssetID == "" {
		return fmt.Errorf("asset_id is required")
	}
	if len(request.Strategies) == 0 {
		return fmt.Errorf("at least one strategy is required")
	}
	if len(request.Strategies) > 10 {
		return fmt.Errorf("maximum 10 strategies allowed for comparison")
	}
	if request.InitialCash <= 0 {
		return fmt.Errorf("initial_cash must be positive")
	}
	if request.StartDate.After(request.EndDate) {
		return fmt.Errorf("start_date must be before end_date")
	}

	// Validate each strategy
	for i, strategy := range request.Strategies {
		if strategy.Type == "" {
			return fmt.Errorf("strategy %d: type is required", i+1)
		}
	}

	return nil
}

// performStrategyComparison analyzes and compares strategy results
func (mss *MultiStrategyService) performStrategyComparison(results []models.BacktestResult, options *models.ComparisonOptions) *models.StrategyComparison {
	if options == nil {
		options = &models.ComparisonOptions{
			ShowBenchmark:      true,
			NormalizeReturns:   false,
			ShowDrawdown:       true,
			ShowRollingMetrics: false,
			RollingWindow:      30,
			Metrics:            []string{"total_return", "sharpe_ratio", "max_drawdown", "win_rate"},
		}
	}

	comparison := &models.StrategyComparison{
		MetricsComparison: make(map[string][]float64),
		Rankings:          make(map[string][]int),
		CorrelationMatrix: make([][]float64, len(results)),
	}

	// Initialize correlation matrix
	for i := range comparison.CorrelationMatrix {
		comparison.CorrelationMatrix[i] = make([]float64, len(results))
	}

	// Extract metrics for comparison
	metrics := []string{"total_return", "annualized_return", "sharpe_ratio", "sortino_ratio", 
		"max_drawdown", "volatility", "win_rate", "profit_factor", "calmar_ratio"}

	for _, metric := range metrics {
		var values []float64
		for _, result := range results {
			value := mss.extractMetricValue(result.PerformanceMetrics, metric)
			values = append(values, value)
		}
		comparison.MetricsComparison[metric] = values

		// Calculate rankings (higher is better for most metrics, except max_drawdown)
		rankings := mss.calculateRankings(values, metric == "max_drawdown" || metric == "volatility")
		comparison.Rankings[metric] = rankings
	}

	// Calculate correlation matrix based on daily returns
	for i := 0; i < len(results); i++ {
		for j := 0; j < len(results); j++ {
			if i == j {
				comparison.CorrelationMatrix[i][j] = 1.0
			} else if i < j {
				correlation := mss.calculateCorrelation(results[i].DailyReturns, results[j].DailyReturns)
				comparison.CorrelationMatrix[i][j] = correlation
				comparison.CorrelationMatrix[j][i] = correlation
			}
		}
	}

	// Determine best and worst strategies
	comparison.BestStrategy, comparison.WorstStrategy = mss.determineBestWorstStrategies(results)

	// Generate summary
	comparison.Summary = mss.generateComparisonSummary(results, comparison)

	return comparison
}

// extractMetricValue extracts a specific metric value from performance metrics
func (mss *MultiStrategyService) extractMetricValue(metrics models.PerformanceMetrics, metricName string) float64 {
	switch metricName {
	case "total_return":
		return metrics.TotalReturn
	case "annualized_return":
		return metrics.AnnualizedReturn
	case "sharpe_ratio":
		return metrics.SharpeRatio
	case "sortino_ratio":
		return metrics.SortinoRatio
	case "max_drawdown":
		return metrics.MaxDrawdown
	case "volatility":
		return metrics.Volatility
	case "win_rate":
		return metrics.WinRate
	case "profit_factor":
		return metrics.ProfitFactor
	case "calmar_ratio":
		return metrics.CalmarRatio
	default:
		return 0.0
	}
}

// calculateRankings calculates rankings for a set of values
func (mss *MultiStrategyService) calculateRankings(values []float64, lowerIsBetter bool) []int {
	type valueIndex struct {
		value float64
		index int
	}

	var valueIndices []valueIndex
	for i, v := range values {
		valueIndices = append(valueIndices, valueIndex{value: v, index: i})
	}

	// Sort values
	if lowerIsBetter {
		sort.Slice(valueIndices, func(i, j int) bool {
			return valueIndices[i].value < valueIndices[j].value
		})
	} else {
		sort.Slice(valueIndices, func(i, j int) bool {
			return valueIndices[i].value > valueIndices[j].value
		})
	}

	// Assign rankings
	rankings := make([]int, len(values))
	for rank, vi := range valueIndices {
		rankings[vi.index] = rank + 1
	}

	return rankings
}

// calculateCorrelation calculates correlation between two daily return series
func (mss *MultiStrategyService) calculateCorrelation(returns1, returns2 []models.DailyReturn) float64 {
	if len(returns1) != len(returns2) || len(returns1) < 2 {
		return 0.0
	}

	// Extract daily return values
	var x, y []float64
	for i := 0; i < len(returns1); i++ {
		x = append(x, returns1[i].DailyReturn)
		y = append(y, returns2[i].DailyReturn)
	}

	// Calculate correlation coefficient
	return mss.pearsonCorrelation(x, y)
}

// pearsonCorrelation calculates Pearson correlation coefficient
func (mss *MultiStrategyService) pearsonCorrelation(x, y []float64) float64 {
	n := float64(len(x))
	if n < 2 {
		return 0.0
	}

	// Calculate means
	var sumX, sumY float64
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
	}
	meanX := sumX / n
	meanY := sumY / n

	// Calculate correlation
	var numerator, denomX, denomY float64
	for i := 0; i < len(x); i++ {
		deltaX := x[i] - meanX
		deltaY := y[i] - meanY
		numerator += deltaX * deltaY
		denomX += deltaX * deltaX
		denomY += deltaY * deltaY
	}

	if denomX == 0 || denomY == 0 {
		return 0.0
	}

	return numerator / math.Sqrt(denomX*denomY)
}

// determineBestWorstStrategies determines the best and worst performing strategies
func (mss *MultiStrategyService) determineBestWorstStrategies(results []models.BacktestResult) (string, string) {
	if len(results) == 0 {
		return "", ""
	}

	bestIndex := 0
	worstIndex := 0
	bestScore := mss.calculateOverallScore(results[0].PerformanceMetrics)
	worstScore := bestScore

	for i := 1; i < len(results); i++ {
		score := mss.calculateOverallScore(results[i].PerformanceMetrics)
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
		if score < worstScore {
			worstScore = score
			worstIndex = i
		}
	}

	bestName := fmt.Sprintf("%s_%d", results[bestIndex].Request.Strategy.Type, bestIndex+1)
	worstName := fmt.Sprintf("%s_%d", results[worstIndex].Request.Strategy.Type, worstIndex+1)

	return bestName, worstName
}

// calculateOverallScore calculates an overall score for strategy performance
func (mss *MultiStrategyService) calculateOverallScore(metrics models.PerformanceMetrics) float64 {
	// Weighted scoring: return 40%, risk-adjusted return 40%, risk 20%
	returnScore := metrics.TotalReturn * 0.4
	riskAdjustedScore := metrics.SharpeRatio * 0.4
	riskScore := (1.0 - metrics.MaxDrawdown) * 0.2 // Lower drawdown is better

	return returnScore + riskAdjustedScore + riskScore
}

// generateComparisonSummary generates a text summary of the strategy comparison
func (mss *MultiStrategyService) generateComparisonSummary(results []models.BacktestResult, comparison *models.StrategyComparison) string {
	if len(results) == 0 {
		return "No strategies to compare."
	}

	summary := fmt.Sprintf("Strategy Comparison Summary (%d strategies):\n\n", len(results))

	// Best and worst performers
	summary += fmt.Sprintf("🏆 Best Strategy: %s\n", comparison.BestStrategy)
	summary += fmt.Sprintf("📉 Worst Strategy: %s\n\n", comparison.WorstStrategy)

	// Performance overview
	returns := comparison.MetricsComparison["total_return"]
	sharpeRatios := comparison.MetricsComparison["sharpe_ratio"]
	maxDrawdowns := comparison.MetricsComparison["max_drawdown"]

	if len(returns) > 0 {
		avgReturn := mss.average(returns)
		maxReturn := mss.max(returns)
		minReturn := mss.min(returns)

		summary += fmt.Sprintf("📈 Returns: Avg %.2f%%, Range %.2f%% to %.2f%%\n", 
			avgReturn*100, minReturn*100, maxReturn*100)
	}

	if len(sharpeRatios) > 0 {
		avgSharpe := mss.average(sharpeRatios)
		summary += fmt.Sprintf("⚖️ Avg Sharpe Ratio: %.3f\n", avgSharpe)
	}

	if len(maxDrawdowns) > 0 {
		avgDrawdown := mss.average(maxDrawdowns)
		summary += fmt.Sprintf("🔻 Avg Max Drawdown: %.2f%%\n", avgDrawdown*100)
	}

	return summary
}

// Helper functions
func (mss *MultiStrategyService) average(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (mss *MultiStrategyService) max(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func (mss *MultiStrategyService) min(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// generateMultiStrategyID generates a unique ID for multi-strategy backtest
func generateMultiStrategyID() string {
	return fmt.Sprintf("ms_%d", time.Now().UnixNano())
}