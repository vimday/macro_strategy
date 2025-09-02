# Macro Strategy

A comprehensive platform for testing and comparing macro trading strategies across multiple markets. Built with TypeScript frontend and Go backend, providing robust backtesting capabilities with configurable data sources and interactive visualization.

## Features

- **Multi-Market Coverage**: Support for A-shares (Chinese stocks), cryptocurrencies (Bitcoin-focused), and Hong Kong/US equities
- **Configurable Data Sources**: Flexible integration with local databases and external data providers
- **Extensible Architecture**: Modular design allowing easy extension to futures and other asset classes
- **Comprehensive Backtesting Engine**: Professional-grade strategy testing with detailed performance analytics
- **Interactive Visualization**: Web-based interface for strategy configuration, monitoring, and results analysis
- **Advanced Performance Metrics**: Maximum drawdown, Sharpe ratio, Sortino ratio, win rate, and comprehensive risk analytics

## Architecture

### Technology Stack

- **Frontend**: TypeScript with modern web frameworks
- **Backend**: Go (Golang) for high-performance data processing
- **Data Layer**: Configurable adapters for multiple data sources
- **Visualization**: Interactive charting and analytics dashboard

### System Overview

```
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│    Web Interface    │    │   Backend Engine    │    │   Data Providers    │
│    (TypeScript)     │◄──►│     (Golang)        │◄──►│   (Configurable)    │
├─────────────────────┤    ├─────────────────────┤    ├─────────────────────┤
│ • Strategy Builder  │    │ • Backtesting Core  │    │ • Local Database    │
│ • Parameter Config  │    │ • Risk Analytics    │    │ • Market Data APIs  │
│ • Results Dashboard │    │ • Performance Calc  │    │ • Alternative Data  │
│ • Portfolio View    │    │ • Data Pipeline     │    │ • Real-time Feeds   │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘
```

## Use Cases

### Example: CSI 500 Monthly Rotation Strategy

**Strategy Logic**:
- Entry: Purchase CSI 500 index at month-end
- Exit: Liquidate position at beginning of following month
- Capital: Full allocation with rebalancing

**Performance Metrics**:
- Total Return and Annualized Return
- Maximum Drawdown and Recovery Time
- Sharpe Ratio and Sortino Ratio
- Win/Loss Ratio and Average Trade Duration
- Volatility and Beta Analysis

## Data Sources

### Supported Markets

| Market | Assets | Data Providers | Frequency |
|--------|--------|----------------|----------|
| A-Shares | Stocks, Indices, ETFs | Tushare, Wind, Local DB | Daily, Intraday |
| Crypto | Bitcoin, Major Altcoins | Binance, CoinGecko, APIs | Tick, Minute, Daily |
| HK/US | Equities, ETFs, ADRs | Yahoo Finance, Alpha Vantage | Daily, Real-time |

### Configuration Options

- **Local Database**: Historical data storage and management
- **API Integration**: Real-time and historical data feeds
- **Custom Adapters**: Extensible framework for new data sources
- **Data Quality**: Built-in validation and cleaning pipelines

## Performance Analytics

### Core Metrics

| Category | Metrics | Description |
|----------|---------|-------------|
| **Returns** | Total Return, CAGR, Rolling Returns | Absolute and risk-adjusted performance |
| **Risk** | Max Drawdown, VaR, CVaR, Volatility | Downside risk and volatility measures |
| **Ratios** | Sharpe, Sortino, Calmar, Information | Risk-adjusted return ratios |
| **Trade** | Win Rate, Profit Factor, Expectancy | Transaction-level statistics |

### Advanced Analytics

- **Factor Attribution**: Performance decomposition by market factors
- **Regime Analysis**: Strategy performance across different market conditions
- **Monte Carlo Simulation**: Probabilistic scenario analysis
- **Stress Testing**: Performance under extreme market conditions

## Installation

### Prerequisites

- Go 1.19 or higher
- Node.js 16 or higher
- Modern web browser with JavaScript enabled

### Quick Start

```bash
# Clone repository
git clone https://github.com/your-org/macro_strategy.git
cd macro_strategy

# Backend setup
cd backend
go mod download
go build -o main .

# Frontend setup
cd ../frontend
npm install
npm run build

# Start services
./backend/main &
npm run serve
```

### Configuration

1. **Data Sources**: Configure `config/data_sources.yaml`
2. **Backend Settings**: Modify `config/server.yaml`
3. **Frontend Config**: Update `frontend/src/config.ts`

## Project Structure

```
macro_strategy/
├── backend/
│   ├── cmd/                    # Application entry points
│   ├── internal/
│   │   ├── api/               # REST API handlers
│   │   ├── backtesting/       # Core backtesting engine
│   │   ├── data/              # Data access layer
│   │   ├── models/            # Business domain models
│   │   └── services/          # Business logic services
│   ├── pkg/                   # Shared packages
│   └── config/                # Configuration files
├── frontend/
│   ├── src/
│   │   ├── components/        # React components
│   │   ├── pages/             # Application pages
│   │   ├── services/          # API service layer
│   │   ├── utils/             # Utility functions
│   │   └── types/             # TypeScript definitions
│   └── public/                # Static assets
├── docs/                      # Documentation
└── scripts/                   # Build and deployment scripts
```

## API Reference

### Core Endpoints

- `POST /api/v1/strategies` - Create new strategy
- `GET /api/v1/strategies/{id}/backtest` - Run backtest
- `GET /api/v1/strategies/{id}/results` - Retrieve results
- `GET /api/v1/data/markets` - Available markets
- `POST /api/v1/data/configure` - Configure data sources

### WebSocket Events

- `backtest.progress` - Real-time backtest progress
- `data.update` - Live data updates
- `strategy.notification` - Strategy alerts

## Contributing

We welcome contributions from the quantitative finance community. Please read our contributing guidelines and submit pull requests for review.

### Development Workflow

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.

## Risk Disclaimer

⚠️ **Important**: This software is for research and educational purposes only. Past performance does not guarantee future results. Users are responsible for validating all strategies and data before live trading.

---

# 巨策略

一个用于测试和比较宏观交易策略的综合平台，支持多市场策略回测。采用TypeScript前端和Go后端构建，提供强大的回测能力、可配置数据源和交互式可视化功能。

## 功能特性

- **多市场覆盖**：支持A股、数字货币（以比特币为主）和港美股市场
- **可配置数据源**：灵活集成本地数据库和外部数据提供商
- **可扩展架构**：模块化设计，便于扩展期货等其他资产类别
- **专业级回测引擎**：提供详细性能分析的专业策略测试
- **交互式可视化**：基于Web的界面，用于策略配置、监控和结果分析
- **高级性能指标**：最大回撤、夏普比率、索提诺比率、胜率和综合风险分析

## 系统架构

### 技术栈

- **前端**：TypeScript配合现代Web框架
- **后端**：Go语言实现高性能数据处理
- **数据层**：支持多数据源的可配置适配器
- **可视化**：交互式图表和分析仪表板

### 系统概览

```
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│      Web界面        │    │     后端引擎        │    │     数据提供商      │
│    (TypeScript)     │◄──►│     (Golang)        │◄──►│     (可配置)        │
├─────────────────────┤    ├─────────────────────┤    ├─────────────────────┤
│ • 策略构建器        │    │ • 回测核心          │    │ • 本地数据库        │
│ • 参数配置          │    │ • 风险分析          │    │ • 市场数据API       │
│ • 结果仪表板        │    │ • 性能计算          │    │ • 替代数据源        │
│ • 组合视图          │    │ • 数据管道          │    │ • 实时数据流        │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘
```

## 应用场景

### 示例：中证500月度轮动策略

**策略逻辑**：
- 入场：月末买入中证500指数
- 出场：次月初卖出全部仓位
- 资金：满仓配置并定期再平衡

**性能指标**：
- 总收益率和年化收益率
- 最大回撤和恢复时间
- 夏普比率和索提诺比率
- 胜负比和平均持仓时间
- 波动率和贝塔分析

## 数据源配置

### 支持市场

| 市场 | 资产类别 | 数据提供商 | 更新频率 |
|------|----------|------------|----------|
| A股 | 股票、指数、ETF | Tushare、Wind、本地数据库 | 日频、分钟 |
| 数字货币 | 比特币、主流代币 | Binance、CoinGecko、API | Tick、分钟、日频 |
| 港美股 | 股票、ETF、ADR | Yahoo Finance、Alpha Vantage | 日频、实时 |

### 配置选项

- **本地数据库**：历史数据存储和管理
- **API集成**：实时和历史数据源
- **自定义适配器**：支持新数据源的可扩展框架
- **数据质量**：内置验证和清洗流程

## 性能分析

### 核心指标

| 类别 | 指标 | 描述 |
|------|------|------|
| **收益** | 总收益、年化收益、滚动收益 | 绝对收益和风险调整收益 |
| **风险** | 最大回撤、VaR、CVaR、波动率 | 下行风险和波动性测量 |
| **比率** | 夏普、索提诺、卡玛、信息比率 | 风险调整收益比率 |
| **交易** | 胜率、盈亏比、期望值 | 交易级别统计 |

### 高级分析

- **因子归因**：按市场因子进行绩效分解
- **市场状态分析**：不同市场环境下的策略表现
- **蒙特卡罗模拟**：概率情景分析
- **压力测试**：极端市场条件下的表现

## 安装指南

### 系统要求

- Go 1.19或更高版本
- Node.js 16或更高版本
- 支持JavaScript的现代浏览器

### 快速开始

```bash
# 克隆仓库
git clone https://github.com/your-org/macro_strategy.git
cd macro_strategy

# 后端设置
cd backend
go mod download
go build -o main .

# 前端设置
cd ../frontend
npm install
npm run build

# 启动服务
./backend/main &
npm run serve
```

### 配置说明

1. **数据源配置**：编辑 `config/data_sources.yaml`
2. **后端设置**：修改 `config/server.yaml`
3. **前端配置**：更新 `frontend/src/config.ts`

## 项目结构

```
macro_strategy/
├── backend/
│   ├── cmd/                    # 应用程序入口
│   ├── internal/
│   │   ├── api/               # REST API处理器
│   │   ├── backtesting/       # 核心回测引擎
│   │   ├── data/              # 数据访问层
│   │   ├── models/            # 业务领域模型
│   │   └── services/          # 业务逻辑服务
│   ├── pkg/                   # 共享包
│   └── config/                # 配置文件
├── frontend/
│   ├── src/
│   │   ├── components/        # React组件
│   │   ├── pages/             # 应用页面
│   │   ├── services/          # API服务层
│   │   ├── utils/             # 工具函数
│   │   └── types/             # TypeScript类型定义
│   └── public/                # 静态资源
├── docs/                      # 文档
└── scripts/                   # 构建和部署脚本
```

## API接口

### 核心端点

- `POST /api/v1/strategies` - 创建新策略
- `GET /api/v1/strategies/{id}/backtest` - 运行回测
- `GET /api/v1/strategies/{id}/results` - 获取结果
- `GET /api/v1/data/markets` - 可用市场
- `POST /api/v1/data/configure` - 配置数据源

### WebSocket事件

- `backtest.progress` - 实时回测进度
- `data.update` - 实时数据更新
- `strategy.notification` - 策略提醒

## 开发贡献

我们欢迎量化金融社区的贡献。请阅读我们的贡献指南并提交拉取请求进行审核。

### 开发流程

1. Fork仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建拉取请求

## 开源协议

本项目采用MIT协议。详见 [LICENSE](LICENSE) 文件。

## 风险提示

⚠️ **重要提示**：本软件仅用于研究和教育目的。历史业绩不代表未来收益。用户在实盘交易前有责任验证所有策略和数据。