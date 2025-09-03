package models

// PredefinedIndexes contains the list of popular A-share indexes with enhanced metadata
var PredefinedIndexes = []Index{
	{
		ID:          "csi300",
		Name:        "沪深300",
		Symbol:      "000300.SH",
		MarketType:  MarketTypeAShare,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyCNY,
		Description: "沪深300指数由沪深市场中市值大、流动性好的300只股票组成",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Shanghai",
			OpenTime:    "09:30",
			CloseTime:   "15:00",
			WeekendDays: []int{0, 6}, // Sunday and Saturday
		},
	},
	{
		ID:          "sse50",
		Name:        "上证50",
		Symbol:      "000016.SH",
		MarketType:  MarketTypeAShare,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyCNY,
		Description: "上证50指数挑选上海证券市场规模大、流动性好的最具代表性的50只股票",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Shanghai",
			OpenTime:    "09:30",
			CloseTime:   "15:00",
			WeekendDays: []int{0, 6},
		},
	},
	{
		ID:          "csi500",
		Name:        "中证500",
		Symbol:      "000905.SH",
		MarketType:  MarketTypeAShare,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyCNY,
		Description: "中证500指数样本空间内股票是扣除沪深300指数样本股及最近一年日均总市值排名前300名的股票",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Shanghai",
			OpenTime:    "09:30",
			CloseTime:   "15:00",
			WeekendDays: []int{0, 6},
		},
	},
	{
		ID:          "csi1000",
		Name:        "中证1000",
		Symbol:      "000852.SH",
		MarketType:  MarketTypeAShare,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyCNY,
		Description: "中证1000指数反映中国A股市场中一批中小市值公司的股票价格表现",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Shanghai",
			OpenTime:    "09:30",
			CloseTime:   "15:00",
			WeekendDays: []int{0, 6},
		},
	},
	{
		ID:          "star50",
		Name:        "科创50",
		Symbol:      "000688.SH",
		MarketType:  MarketTypeAShare,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyCNY,
		Description: "科创50指数由科创板中市值大、流动性好的50只证券组成",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Shanghai",
			OpenTime:    "09:30",
			CloseTime:   "15:00",
			WeekendDays: []int{0, 6},
		},
	},
	{
		ID:          "chinext",
		Name:        "创业板指",
		Symbol:      "399006.SZ",
		MarketType:  MarketTypeAShare,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyCNY,
		Description: "创业板指数反映创业板市场层次的运行情况和发展趋势",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Shanghai",
			OpenTime:    "09:30",
			CloseTime:   "15:00",
			WeekendDays: []int{0, 6},
		},
	},
	{
		ID:          "szse100",
		Name:        "深证100",
		Symbol:      "399330.SZ",
		MarketType:  MarketTypeAShare,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyCNY,
		Description: "深证100指数选取深圳市场A股流通市值最大、成交最活跃的100只成份股",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Shanghai",
			OpenTime:    "09:30",
			CloseTime:   "15:00",
			WeekendDays: []int{0, 6},
		},
	},
}

// CryptoAssets contains popular cryptocurrency assets for future expansion
var CryptoAssets = []Index{
	{
		ID:          "btc",
		Name:        "Bitcoin",
		Symbol:      "BTC/USDT",
		MarketType:  MarketTypeCrypto,
		AssetClass:  AssetClassCrypto,
		Currency:    CurrencyBTC,
		Description: "Bitcoin is the first decentralized digital currency",
		TradingHours: &TradingHours{
			Timezone:    "UTC",
			OpenTime:    "00:00",
			CloseTime:   "23:59",
			WeekendDays: []int{}, // 24/7 trading
		},
		Metadata: map[string]interface{}{
			"market_cap_rank": 1,
			"category":        "cryptocurrency",
		},
	},
	{
		ID:          "eth",
		Name:        "Ethereum",
		Symbol:      "ETH/USDT",
		MarketType:  MarketTypeCrypto,
		AssetClass:  AssetClassCrypto,
		Currency:    CurrencyETH,
		Description: "Ethereum is a decentralized platform for smart contracts",
		TradingHours: &TradingHours{
			Timezone:    "UTC",
			OpenTime:    "00:00",
			CloseTime:   "23:59",
			WeekendDays: []int{}, // 24/7 trading
		},
		Metadata: map[string]interface{}{
			"market_cap_rank": 2,
			"category":        "cryptocurrency",
		},
	},
}

// HKUSAssets contains Hong Kong and US market assets for future expansion
var HKUSAssets = []Index{
	{
		ID:          "spy",
		Name:        "SPDR S&P 500 ETF",
		Symbol:      "SPY",
		MarketType:  MarketTypeHKUS,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyUSD,
		Description: "SPDR S&P 500 ETF Trust tracks the S&P 500 index",
		TradingHours: &TradingHours{
			Timezone:    "America/New_York",
			OpenTime:    "09:30",
			CloseTime:   "16:00",
			WeekendDays: []int{0, 6}, // Sunday and Saturday
		},
		Metadata: map[string]interface{}{
			"exchange": "NYSE Arca",
			"category": "us_etf",
		},
	},
	{
		ID:          "hsi",
		Name:        "Hang Seng Index",
		Symbol:      "HSI",
		MarketType:  MarketTypeHKUS,
		AssetClass:  AssetClassIndex,
		Currency:    CurrencyHKD,
		Description: "Hang Seng Index tracks the largest companies in Hong Kong",
		TradingHours: &TradingHours{
			Timezone:    "Asia/Hong_Kong",
			OpenTime:    "09:30",
			CloseTime:   "16:00",
			WeekendDays: []int{0, 6}, // Sunday and Saturday
		},
		Metadata: map[string]interface{}{
			"exchange": "HKEX",
			"category": "hk_index",
		},
	},
}

// GetAllAssets returns all available assets across all market types
func GetAllAssets() []Index {
	var allAssets []Index
	allAssets = append(allAssets, PredefinedIndexes...)
	// Future expansion:
	// allAssets = append(allAssets, CryptoAssets...)
	// allAssets = append(allAssets, HKUSAssets...)
	return allAssets
}

// GetAssetsByClass returns assets filtered by asset class
func GetAssetsByClass(assetClass AssetClass) []Index {
	var result []Index
	allAssets := GetAllAssets()
	for _, asset := range allAssets {
		if asset.AssetClass == assetClass {
			result = append(result, asset)
		}
	}
	return result
}

// GetAssetsByCurrency returns assets filtered by currency
func GetAssetsByCurrency(currency Currency) []Index {
	var result []Index
	allAssets := GetAllAssets()
	for _, asset := range allAssets {
		if asset.Currency == currency {
			result = append(result, asset)
		}
	}
	return result
}

// GetIndexByID returns an index by its ID
func GetIndexByID(id string) *Index {
	for _, index := range PredefinedIndexes {
		if index.ID == id {
			return &index
		}
	}
	return nil
}

// GetIndexesByMarketType returns indexes filtered by market type
func GetIndexesByMarketType(marketType MarketType) []Index {
	var result []Index
	for _, index := range PredefinedIndexes {
		if index.MarketType == marketType {
			result = append(result, index)
		}
	}
	return result
}
