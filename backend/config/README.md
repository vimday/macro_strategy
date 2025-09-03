# Configuration Guide

## 数据源配置

### 使用 AKShare 获取 A 股数据

1. **安装依赖**
   ```bash
   # 确保安装了 Python 3.7+
   pip install akshare pandas
   ```

2. **配置数据提供者**
   在 `internal/api/router.go` 中，将 Mock 数据提供者替换为 AKShare：
   
   ```go
   // 注释掉 Mock 提供者
   // dataManager.RegisterProvider(models.MarketTypeAShare, data.NewMockDataProvider())
   
   // 启用 AKShare 提供者
   akshareProvider := data.NewAKShareProvider("python3", "./scripts/akshare_client.py")
   dataManager.RegisterProvider(models.MarketTypeAShare, akshareProvider)
   ```

3. **运行说明**
   - Python 脚本位于 `scripts/akshare_client.py`
   - 支持的指数代码：sh000300（沪深300）、sh000905（中证500）等
   - 数据格式自动转换为标准 OHLCV 格式

## 策略扩展

### 添加新策略类型

1. **后端扩展**
   - 在 `models/types.go` 中添加新的 `StrategyType`
   - 在 `backtesting/engine.go` 中实现策略逻辑
   - 在 `services/backtest_service.go` 中添加验证

2. **前端扩展**
   - 在 `types/index.ts` 中添加新的策略类型
   - 在 `components/StrategyForm.tsx` 中添加表单字段
   - 更新相关组件

## 部署配置

### 开发环境
```bash
# 后端
cd backend
go mod tidy
go run cmd/main.go

# 前端
cd frontend
npm install
npm run dev
```

### 生产环境
```bash
# 后端构建
cd backend
go build -o macro_strategy cmd/main.go

# 前端构建
cd frontend
npm run build
npm start
```

## 注意事项

1. **数据限制**：AKShare 可能有请求频率限制，建议添加缓存机制
2. **错误处理**：网络异常时会自动降级到 Mock 数据
3. **扩展性**：当前架构支持轻松添加新的市场类型和数据源
4. **安全性**：生产环境请添加 API 认证和限流机制