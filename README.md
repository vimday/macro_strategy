# Macro Strategy | 巨策略

🚀 **A comprehensive multi-market trading strategy backtesting platform with real market data integration**

Built with modern TypeScript frontend and high-performance Go backend, providing professional-grade backtesting capabilities with AKShare real-time data integration, multi-market support, and interactive visualization.

## ✨ Key Features

- 🌍 **Multi-Market Support**: A-shares, US stocks, HK stocks, and cryptocurrencies
- 🏛️ **A-Share Focus**: Professional A-share index and individual stock backtesting
- 📊 **Real Data Integration**: AKShare, Yahoo Finance, and Binance integrations
- ⚡ **High Performance**: Go backend with efficient data processing and calculation
- 🎯 **Multiple Strategies**: Buy-and-hold, monthly rotation, and customizable strategies
- 📈 **Strategy Comparison**: Compare multiple strategies side-by-side
- 🖥️ **Modern UI**: Next.js 14 + React 18 + Ant Design 5 responsive interface
- 📱 **Interactive Charts**: ECharts-powered visualization with multiple chart types
- 🔧 **Easy Setup**: One-command startup with automated environment configuration

## 📋 Recent Fixes and Improvements ✅

### 🐛 **Critical Bug Fixes** (September 4, 2025)

- **Fixed 500 Error**: Resolved AxiosError "Request failed with status code 500" when running backtests
- **AKShare Integration**: Fixed missing `get_stock_zh_index_daily` function in Python client script
- **Path Configuration**: Corrected file paths in provider configuration for proper Python virtual environment usage
- **Accuracy Issues**: Fixed win rate calculation and P&L computation accuracy problems

These fixes ensure the platform now works correctly for all supported markets and strategies.

## 🏗️ Architecture

### Technology Stack

- **Frontend**: Next.js 14 + React 18 + TypeScript + Ant Design 5
- **Backend**: Go 1.19 + Gin Framework + Modular Architecture  
- **Data Sources**: 
  - AKShare (Python) for A-share market data
  - Yahoo Finance API for US/HK market data
  - Binance API for cryptocurrency data
- **Charts**: ECharts for interactive data visualization
- **State Management**: TanStack Query (React Query) for server state
- **Styling**: Tailwind CSS for utility-first styling

### System Architecture

```
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│    Web Interface    │    │   Backend Engine    │    │   Data Providers    │
│    (Next.js 14)     │◄──►│     (Go + Gin)      │◄──►│    (AKShare)        │
├─────────────────────┤    ├─────────────────────┤    ├─────────────────────┤
│ • Strategy Form     │    │ • Backtesting Core  │    │ • Real A-Share Data │
│ • Performance UI    │    │ • Metrics Engine    │    │ • CSI 300, CSI 1000│
│ • Interactive Charts│    │ • Trade Analytics   │    │ • STAR 50, SSE 50  │
│ • Results Dashboard │    │ • API Endpoints     │    │ • Real-time Feeds   │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘
```

## 🌍 Supported Markets and Assets

The platform now supports multiple global markets with real data integration:

| Market Type | Assets | Data Source | Description |
|-------------|--------|-------------|-------------|
| **A-Share Indexes** | CSI 300, SSE 50, CSI 500, CSI 1000, STAR 50, ChiNext, SZSE 100 | AKShare | Major Chinese market indexes |
| **A-Share Stocks** | Individual stocks (000858 Wuliangye, 000001 Ping An Bank, etc.) | AKShare | Chinese individual companies |
| **US Indexes** | S&P 500 ETF (SPY), NASDAQ-100 ETF (QQQ) | Yahoo Finance | Major US market ETFs |
| **US Stocks** | Apple (AAPL), Microsoft (MSFT), etc. | Yahoo Finance | US individual companies |
| **Cryptocurrencies** | Bitcoin (BTC), Ethereum (ETH) | Binance | Major digital currencies |
| **HK Indexes** | Hang Seng Index (HSI) | Yahoo Finance | Hong Kong market indexes |
| **HK Stocks** | Tencent (00700.HK) | Yahoo Finance | Hong Kong individual companies |

### Data Features
- ✅ **Real Market Data**: Direct integration with multiple data providers
- ✅ **Daily OHLCV**: Complete open, high, low, close, and volume data
- ✅ **Commission Handling**: Realistic transaction cost calculations
- ✅ **Date Range Filtering**: Flexible historical period selection
- ✅ **Multi-Currency Support**: CNY, USD, HKD, BTC, ETH, and more

## 📋 Multiple Trading Strategies

### Strategy Implementation

The platform features multiple fully implemented trading strategies:

```
📈 Buy-and-Hold Strategy:
   • Initial purchase with configurable allocation
   • Optional periodic rebalancing (monthly, quarterly, yearly)
   • Simple and effective long-term investing approach

📅 Monthly Rotation Strategy:
   • Buy signal: N days before month-end
   • Full cash allocation to selected asset
   • Sell signal: M days after month-start  
   • Complete position liquidation
   • Repeat cycle for entire backtest period

🔄 Future Strategies (Planned):
   • Grid Trading Strategy
   • Mean Reversion Strategy
   • Momentum Strategy
   • DCA (Dollar Cost Averaging)
   • Multi-Factor Strategies
```

### Strategy Parameters

**Buy-and-Hold Parameters**:
- `target_allocation`: Position size as percentage of portfolio (0.1-1.0)
- `rebalance_frequency`: How often to rebalance ("never", "monthly", "quarterly", "yearly")
- `dividend_reinvest`: Whether to reinvest dividends (future feature)

**Monthly Rotation Parameters**:
- `buy_days_before_month_end`: Entry timing (default: 1 day)
- `sell_days_after_month_start`: Exit timing (default: 1 day)

### Example Use Case: CSI 1000 Monthly Rotation

**Scenario**: Test monthly rotation on CSI 1000 index
- **Period**: 2024-01-01 to 2024-03-31
- **Entry**: 1 day before month-end
- **Exit**: 1 day after month-start
- **Capital**: 1,000,000 CNY

**Results** (✅ **Fixed accuracy issues**):
- ✅ **Win Rate**: 50.00% (1 win, 1 loss)
- ✅ **Total Return**: 0.4965%
- ✅ **Sharpe Ratio**: -0.665
- ✅ **Max Drawdown**: 0.0599%

## 🆚 Strategy Comparison and Analysis

### Multi-Strategy Backtesting

The platform now supports comparing multiple strategies side-by-side:

```
📊 Multi-Strategy Features:
   • Run multiple strategies on the same asset
   • Compare performance metrics across strategies
   • Generate correlation matrices between strategies
   • Rank strategies by various performance criteria
   • Visualize comparative performance charts
   • Include benchmark asset for reference
```

### ✅ **Enhanced Analytics Implemented**

| Category | Metrics | Status | Description |
|----------|---------|---------|-------------|
| **Returns** | Total Return, Annualized Return | ✅ Working | Absolute performance calculation |
| **Risk** | Max Drawdown, Volatility | ✅ Working | Downside risk and volatility measures |
| **Ratios** | Sharpe Ratio, Sortino Ratio, Calmar Ratio | ✅ Working | Risk-adjusted return ratios |
| **Trade** | Win Rate, Profit Factor, Trade Count | ✅ Fixed | Transaction-level statistics |
| **Advanced** | Max Drawdown Period, Recovery Period | ✅ Working | Temporal risk analysis |
| **Comparison** | Strategy Rankings, Correlations | ✅ New | Multi-strategy analysis |

### 🔍 **Critical Improvements**

**Recent Enhancements** (✅ **Completed**):
- **Multi-Market Support**: Added US, HK, and crypto market data providers
- **Individual Stocks**: Support for A-share and US individual stocks
- **Multiple Strategies**: Buy-and-hold strategy implementation
- **Strategy Comparison**: Side-by-side strategy analysis
- **Enhanced API**: New endpoints for markets, strategies, and multi-backtesting
- **Trade Pairing Fix**: Corrected Go range loop pointer problem
- **Win Rate Calculation**: Fixed from 0% to accurate percentage
- **P&L Calculation**: Fixed percentage return calculation

### 📊 **Available Chart Types**

1. **Portfolio Value Chart**: Track portfolio growth over time
2. **Daily Returns Chart**: Visualize daily performance fluctuations  
3. **Drawdown Chart**: Monitor risk exposure periods
4. **Trade Distribution**: Monthly trade frequency analysis

## 📦 Enhanced API Reference

### **New and Updated Endpoints**

```bash
# Health Check
GET /api/v1/health

# Asset Management (enhanced from indexes)  
GET /api/v1/assets                    # Get all available assets
GET /api/v1/assets/market/:type       # Get assets by market type
GET /api/v1/assets/data/:id           # Get market data for asset
GET /api/v1/markets                   # Get all supported markets

# Strategy Management
GET /api/v1/strategies                # Get all supported strategies

# Single Strategy Backtesting
POST /api/v1/backtest                 # Run single strategy backtest
GET /api/v1/backtest/:id              # Get backtest results

# Multi-Strategy Comparison (NEW)
POST /api/v1/backtest/multi           # Run multi-strategy comparison
GET /api/v1/backtest/multi/:id        # Get multi-strategy results

# Backward Compatibility
GET /api/v1/indexes                   # Get all indexes (legacy)
GET /api/v1/indexes/market/:type      # Get indexes by market type (legacy)
GET /api/v1/indexes/data/:id          # Get market data for index (legacy)
```

### **Multi-Strategy Request Example**

```json
{
  "asset_id": "csi1000",
  "strategies": [
    {
      "type": "buy_and_hold",
      "parameters": {
        "target_allocation": 1.0,
        "rebalance_frequency": "never"
      }
    },
    {
      "type": "monthly_rotation",
      "parameters": {
        "buy_days_before_month_end": 1,
        "sell_days_after_month_start": 1
      }
    }
  ],
  "start_date": "2024-01-01",
  "end_date": "2024-03-31",
  "initial_cash": 1000000,
  "benchmark": "csi300",
  "comparison_opt": {
    "show_benchmark": true,
    "normalize_returns": false,
    "show_drawdown": true,
    "metrics": ["total_return", "sharpe_ratio", "max_drawdown", "win_rate"]
  }
}
```

## 🚀 Quick Start

### 📚 **Prerequisites**

- **Go 1.19+** for backend development
- **Node.js 16+** for frontend development  
- **Python 3.8+** for AKShare data integration
- **Git** for version control

### 🔧 **One-Command Setup**

```bash
# Clone the repository
git clone https://github.com/vimday/macro_strategy.git
cd macro_strategy

# Setup AKShare environment (first time only)
./setup_akshare.sh

# Start both backend and frontend
./start_dev.sh
```

### 🏭 **Manual Setup**

If you prefer manual control:

```bash
# Backend setup
cd backend
go mod tidy
go run cmd/main.go                    # Starts on :8080

# Frontend setup (new terminal)
cd frontend  
npm install
npm run dev                           # Starts on :3000

# AKShare setup (if needed)
source venv/bin/activate
python3 backend/scripts/akshare_client.py get_stock_zh_a_hist sh000300 20240101 20240105
```

### 🌐 **Access URLs**

- **Frontend Interface**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/api/v1/health
- **API Documentation**: Available via endpoint testing

### ✅ **Verify Installation**

```bash
# Test backend health
curl http://localhost:8080/api/v1/health

# Test data retrieval
curl "http://localhost:8080/api/v1/indexes"

# Run quick backtest (optional)
python3 test_win_rate.py
```

## 🗺️ Project Structure

```
macro_strategy/
├── backend/                     # Go backend application
│   ├── cmd/                     # Application entry points
│   │   └── main.go             # Main server executable
│   ├── internal/
│   │   ├── api/                # REST API handlers & routing
│   │   ├── backtesting/        # Core backtesting engine
│   │   ├── data/               # Data providers (AKShare, Mock)
│   │   ├── models/             # Business domain models & types
│   │   └── services/           # Business logic services
│   ├── scripts/                # Python AKShare integration
│   └── go.mod                  # Go dependencies
├── frontend/                    # Next.js frontend application
│   ├── src/
│   │   ├── app/                # Next.js 14 app directory
│   │   ├── components/         # React components
│   │   ├── hooks/              # Custom React hooks
│   │   ├── lib/                # API services & utilities
│   │   └── types/              # TypeScript type definitions
│   └── package.json            # Node.js dependencies
├── venv/                       # Python virtual environment
├── setup_akshare.sh            # AKShare environment setup
├── start_dev.sh                # Development startup script
├── test_win_rate.py            # Win rate verification script
└── debug_backtest.py           # Backtest debugging utility
```

## 🛠️ Current Implementation Status

### ✅ **Completed Features**

- **✅ Backend Infrastructure**
  - Go + Gin REST API server with modular architecture
  - Comprehensive error handling and validation
  - CORS configuration for frontend integration
  - Multi-market data provider system

- **✅ Data Integration**
  - AKShare real market data provider for A-shares
  - Yahoo Finance API provider for US/HK markets
  - Binance API provider for cryptocurrency data
  - Mock data provider for testing
  - Unified data model for multiple markets
  - Date range filtering and validation

- **✅ Backtesting Engine**
  - Multiple strategy implementations (Buy-and-Hold, Monthly Rotation)
  - Trade execution with commission calculation
  - Daily portfolio value tracking
  - Position management and cash handling
  - Multi-strategy comparison capabilities

- **✅ Performance Metrics** 
  - Accurate win rate calculation (fixed pointer issues)
  - Sharpe ratio, Sortino ratio, Calmar ratio
  - Maximum drawdown and recovery periods
  - Trade-level statistics and P&L analysis
  - Strategy comparison and ranking

- **✅ Frontend Interface**
  - Next.js 14 with React 18 and TypeScript
  - Ant Design components for professional UI
  - TanStack Query for efficient server state
  - Interactive ECharts for data visualization
  - Responsive design with Tailwind CSS
  - Multi-strategy configuration UI

- **✅ Development Tools**
  - Automated development environment setup
  - Testing utilities for verification
  - Debug scripts for troubleshooting
  - Comprehensive API documentation

### 🔄 **Future Roadmap**

- **🔄 Additional Strategies**
  - Grid trading strategy
  - Mean reversion strategy
  - Momentum-based strategies
  - DCA (Dollar Cost Averaging)
  - Multi-factor strategies
  - Machine learning strategies

- **🔄 Extended Market Support**
  - Commodity futures and bond ETFs
  - More individual stocks and ETFs
  - Additional cryptocurrency exchanges
  - Global market expansion (EU, Japan, etc.)

- **🔄 Advanced Features**
  - Strategy portfolio backtesting
  - Risk management modules
  - Real-time signal alerts
  - Strategy optimization tools
  - Multi-timeframe analysis
  - Dividend and corporate action handling
  - Advanced visualization dashboards

## 🔌 **API Reference**

### **Implemented Endpoints**

```bash
# Health Check
GET /api/v1/health

# Index Management  
GET /api/v1/indexes                    # Get all available indexes
GET /api/v1/indexes/market/:type       # Get indexes by market type
GET /api/v1/indexes/data/:id           # Get market data for index

# Backtesting
POST /api/v1/backtest                  # Run backtest
GET /api/v1/backtest/:id               # Get backtest results
```

### **Request Example**

```json
{
  "index_id": "csi1000",
  "strategy": {
    "type": "monthly_rotation",
    "parameters": {
      "buy_days_before_month_end": 1,
      "sell_days_after_month_start": 1
    }
  },
  "start_date": "2024-01-01",
  "end_date": "2024-03-31",
  "initial_cash": 1000000
}
```

## 🤝 **Contributing**

We welcome contributions from the quantitative finance community. Please read our contributing guidelines and submit pull requests for review.

### **Development Workflow**

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m '🚀 feat: Add amazing feature'`)
4. Push branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### **Commit Message Format**

We follow the conventional commits specification with emojis:

- `🚀 feat:` - New features
- `🐛 fix:` - Bug fixes  
- `📈 perf:` - Performance improvements
- `🔥 refactor:` - Code refactoring
- `📝 docs:` - Documentation updates
- `✅ test:` - Testing improvements

## 📜 **License**

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.

## ⚠️ **Risk Disclaimer**

**Important**: This software is for research and educational purposes only. Past performance does not guarantee future results. Users are responsible for validating all strategies and data before live trading.

---

# 🗡️ 巨策略 | Macro Strategy

**中国 A 股交易策略回测专业平台**

采用现代化 TypeScript 前端和高性能 Go 后端构建，提供专业级回测能力、AKShare 真实数据集成和交互式可视化。

## ✨ 核心特性

- 🏛️ **A 股专注**：专业 A 股指数回测，真实市场数据
- 📊 **真实数据**：AKShare 集成，提供正宗中国市场数据
- ⚡ **高性能**：Go 后端，高效数据处理和计算
- 🎯 **月末轮动策略**：内置月末轮动策略实现
- 📈 **专业指标**：全面性能分析（夏普比率、最大回撤、胜率等）
- 🖥️ **现代界面**：Next.js 14 + React 18 + Ant Design 响应式界面
- 📱 **交互图表**：ECharts 驱动的可视化，多种图表类型
- 🔧 **简单设置**：一键启动，自动环境配置

## 📈 支持的 A 股指数

平台目前支持以下主要中国市场指数，均使用 AKShare 真实数据：

| 指数 | 名称 | 代码 | 描述 |
|------|------|--------|--------------|
| **汪深300** | CSI 300 | 000300.SH | 中国最具代表性的300只大盘股 |
| **上证50** | SSE 50 | 000016.SH | 上海证券市场最具代表性的50只股票 |
| **中证500** | CSI 500 | 000905.SH | 中小市值代表性指数 |
| **中证1000** | CSI 1000 | 000852.SH | 中小市值股票价格表现 |
| **科创50** | STAR 50 | 000688.SH | 科创板最具代表性的50只证券 |
| **创业板指** | ChiNext | 399006.SZ | 创业板市场运行情况 |
| **深证100** | SZSE 100 | 399330.SZ | 深圳市场最活跃100只成份股 |

### 数据特性
- ✅ **真实市场数据**：直接集成 AKShare 获取正宗历史数据
- ✅ **完整 OHLCV**：完整的开盘、最高、最低、收盘、成交量数据
- ✅ **手续费处理**：真实交易成本计算
- ✅ **日期范围筛选**：灵活历史周期选择

## 📋 月末轮动策略

### 策略实现

平台实现了完整的 **月末轮动策略**，逻辑如下：

```
📅 月末入场：
   • 买入信号：月末前 N 天
   • 全仓配置到选定指数
   • 市价单执行，含手续费

📅 月初出场：
   • 卖出信号：月初后 M 天
   • 完全清仓
   • 返回现金等待下次信号

🔄 在整个回测周期内重复循环
```

**可配置参数**：
- `buy_days_before_month_end`：入场时机（默认：1 天）
- `sell_days_after_month_start`：出场时机（默认：1 天）
- `initial_cash`：起始资金（默认：1,000,000 人民币）
- `commission_rate`：交易成本（每笔交易可配置）

### 用例：中证1000月末轮动

**场景**：测试中证1000指数的月末轮动
- **周期**：2024-01-01 至 2024-03-31
- **入场**：月末前 1 天
- **出场**：月初后 1 天
- **资金**：1,000,000 人民币

**结果** (✅ **已修复准确性问题**):
- ✅ **胜率**：50.00%（1 胜 1 负）
- ✅ **总收益率**：0.4965%
- ✅ **夏普比率**：-0.665
- ✅ **最大回撤**：0.0599%

```
