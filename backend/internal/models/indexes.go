package models

// PredefinedIndexes contains the list of popular A-share indexes
var PredefinedIndexes = []Index{
	{
		ID:          "csi300",
		Name:        "沪深300",
		Symbol:      "000300.SH",
		MarketType:  MarketTypeAShare,
		Description: "沪深300指数由沪深市场中市值大、流动性好的300只股票组成",
	},
	{
		ID:          "sse50",
		Name:        "上证50",
		Symbol:      "000016.SH",
		MarketType:  MarketTypeAShare,
		Description: "上证50指数挑选上海证券市场规模大、流动性好的最具代表性的50只股票",
	},
	{
		ID:          "csi500",
		Name:        "中证500",
		Symbol:      "000905.SH",
		MarketType:  MarketTypeAShare,
		Description: "中证500指数样本空间内股票是扣除沪深300指数样本股及最近一年日均总市值排名前300名的股票",
	},
	{
		ID:          "csi1000",
		Name:        "中证1000",
		Symbol:      "000852.SH",
		MarketType:  MarketTypeAShare,
		Description: "中证1000指数反映中国A股市场中一批中小市值公司的股票价格表现",
	},
	{
		ID:          "star50",
		Name:        "科创50",
		Symbol:      "000688.SH",
		MarketType:  MarketTypeAShare,
		Description: "科创50指数由科创板中市值大、流动性好的50只证券组成",
	},
	{
		ID:          "chinext",
		Name:        "创业板指",
		Symbol:      "399006.SZ",
		MarketType:  MarketTypeAShare,
		Description: "创业板指数反映创业板市场层次的运行情况和发展趋势",
	},
	{
		ID:          "szse100",
		Name:        "深证100",
		Symbol:      "399330.SZ",
		MarketType:  MarketTypeAShare,
		Description: "深证100指数选取深圳市场A股流通市值最大、成交最活跃的100只成份股",
	},
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
