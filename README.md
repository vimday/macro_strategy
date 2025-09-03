# Macro Strategy | 巨策略

🚀 **A comprehensive A-share trading strategy backtesting platform with real market data integration**

Built with modern TypeScript frontend and high-performance Go backend, providing professional-grade backtesting capabilities with AKShare real-time data integration and interactive visualization.

## ✨ Key Features

- 🏛️ **A-Share Focus**: Professional A-share index backtesting with real market data
- 📊 **Real Data Integration**: AKShare integration for authentic Chinese market data  
- ⚡ **High Performance**: Go backend with efficient data processing and calculation
- 🎯 **Monthly Rotation Strategy**: Built-in month-end rotation strategy implementation
- 📈 **Professional Metrics**: Comprehensive performance analytics (Sharpe ratio, max drawdown, win rate, etc.)
- 🖥️ **Modern UI**: Next.js 14 + React 18 + Ant Design responsive interface
- 📱 **Interactive Charts**: ECharts-powered visualization with multiple chart types
- 🔧 **Easy Setup**: One-command startup with automated environment configuration

## 🏗️ Architecture

### Technology Stack

- **Frontend**: Next.js 14 + React 18 + TypeScript + Ant Design 5
- **Backend**: Go 1.19 + Gin Framework + Modular Architecture  
- **Data Source**: AKShare (Python) for real A-share market data
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

## 📈 Supported A-Share Indexes

The platform currently supports the following major Chinese market indexes with real AKShare data:

| Index | Name | Symbol | Description |
|-------|------|--------|--------------|
| **CSI 300** | 汪深300 | 000300.SH | 中国最具代表性的300只大盘股 |
| **SSE 50** | 上证50 | 000016.SH | 上海证券市场最具代表性的50只股票 |
| **CSI 500** | 中证500 | 000905.SH | 中小市值代表性指数 |
| **CSI 1000** | 中证1000 | 000852.SH | 中小市值股票价格表现 |
| **STAR 50** | 科创50 | 000688.SH | 科创板最具代表性的50只证券 |
| **ChiNext** | 创业板指 | 399006.SZ | 创业板市场运行情况 |
| **SZSE 100** | 深证100 | 399330.SZ | 深圳市场最活跃100只成份股 |

### Data Features
- ✅ **Real Market Data**: Direct integration with AKShare for authentic historical data
- ✅ **Daily OHLCV**: Complete open, high, low, close, and volume data
- ✅ **Commission Handling**: Realistic transaction cost calculations
- ✅ **Date Range Filtering**: Flexible historical period selection

## 📋 Monthly Rotation Strategy

### Strategy Implementation

The platform features a fully implemented **Monthly Rotation Strategy** with the following logic:

```
📅 Month-End Entry:
   • Buy signal: N days before month-end
   • Full cash allocation to selected index
   • Market order execution with commission

📅 Month-Start Exit:
   • Sell signal: M days after month-start  
   • Complete position liquidation
   • Return to cash until next signal

🔄 Repeat cycle for entire backtest period
```

**Configurable Parameters**:
- `buy_days_before_month_end`: Entry timing (default: 1 day)
- `sell_days_after_month_start`: Exit timing (default: 1 day)
- `initial_cash`: Starting capital (default: 1,000,000 CNY)
- `commission_rate`: Transaction costs (configurable per trade)

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

## 📈 Performance Analytics

The platform provides comprehensive performance metrics with accurate calculations:

### ✅ **Core Metrics Implemented**

| Category | Metrics | Status | Description |
|----------|---------|---------|-------------|
| **Returns** | Total Return, Annualized Return | ✅ Working | Absolute performance calculation |
| **Risk** | Max Drawdown, Volatility | ✅ Working | Downside risk and volatility measures |
| **Ratios** | Sharpe Ratio, Sortino Ratio, Calmar Ratio | ✅ Working | Risk-adjusted return ratios |
| **Trade** | Win Rate, Profit Factor, Trade Count | ✅ **Fixed** | Transaction-level statistics |
| **Advanced** | Max Drawdown Period, Recovery Period | ✅ Working | Temporal risk analysis |

### 🔍 **Critical Bug Fixes**

**Recent Fixes** (✅ **Completed**):
- **Trade Pairing Issue**: Fixed Go range loop pointer problem causing incorrect trade matching
- **Win Rate Calculation**: Corrected from 0% to accurate percentage (e.g., 50%)
- **P&L Calculation**: Fixed percentage return calculation with proper commission handling
- **Metrics Display**: All metrics now show correct values in the frontend

### 📊 **Available Chart Types**

1. **Portfolio Value Chart**: Track portfolio growth over time
2. **Daily Returns Chart**: Visualize daily performance fluctuations  
3. **Drawdown Chart**: Monitor risk exposure periods
4. **Trade Distribution**: Monthly trade frequency analysis

## 📦 API Reference

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
  - Go + Gin REST API server
  - Modular architecture with clean separation
  - Comprehensive error handling and validation
  - CORS configuration for frontend integration

- **✅ Data Integration**
  - AKShare real market data provider
  - Mock data provider for testing
  - Unified data model for multiple markets
  - Date range filtering and validation

- **✅ Backtesting Engine**
  - Monthly rotation strategy implementation
  - Trade execution with commission calculation
  - Daily portfolio value tracking
  - Position management and cash handling

- **✅ Performance Metrics** (✅ **Bug fixes completed**)
  - Accurate win rate calculation (fixed pointer issues)
  - Sharpe ratio, Sortino ratio, Calmar ratio
  - Maximum drawdown and recovery periods
  - Trade-level statistics and P&L analysis

- **✅ Frontend Interface**
  - Next.js 14 with React 18 and TypeScript
  - Ant Design components for professional UI
  - TanStack Query for efficient server state
  - Interactive ECharts for data visualization
  - Responsive design with Tailwind CSS

- **✅ Development Tools**
  - Automated development environment setup
  - Testing utilities for verification
  - Debug scripts for troubleshooting

### 🔄 **Future Roadmap**

- **🔄 Additional Strategies**
  - Buy and hold strategy
  - Grid trading strategy
  - Mean reversion strategy
  - Momentum-based strategies

- **🔄 Extended Market Support**
  - Cryptocurrency data integration (Binance API)
  - Hong Kong and US equity markets
  - Commodity futures and bond ETFs

- **🔄 Advanced Features**
  - Strategy portfolio backtesting
  - Risk management modules
  - Real-time signal alerts
  - Strategy optimization tools
  - Multi-timeframe analysis

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
