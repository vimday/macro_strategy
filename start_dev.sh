#!/bin/bash

# Macro Strategy 开发环境启动脚本

echo "🚀 启动 Macro Strategy 开发环境..."

# 检查依赖
echo "📋 检查依赖..."

# 检查 Go
if ! command -v go &> /dev/null; then
    echo "❌ 未找到 Go，请先安装 Go 1.19+"
    exit 1
fi

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ 未找到 Node.js，请先安装 Node.js 16+"
    exit 1
fi

# 检查 Python（可选，用于 AKShare）
if command -v python3 &> /dev/null; then
    echo "✅ Python3 已安装，可以使用 AKShare 数据源"
else
    echo "⚠️  未找到 Python3，将使用 Mock 数据源"
fi

echo ""

# 启动后端
echo "🔧 启动后端服务..."
cd backend

# 安装依赖
echo "📦 安装 Go 依赖..."
go mod tidy

# 启动后端（后台运行）
echo "▶️  启动后端服务（端口：8080）..."
go run cmd/main.go &
BACKEND_PID=$!

# 等待后端启动
sleep 3

# 检查后端是否启动成功
if curl -f -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ 后端服务启动成功"
else
    echo "❌ 后端服务启动失败"
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi

echo ""

# 启动前端
echo "🎨 启动前端服务..."
cd ../frontend

# 安装依赖
echo "📦 安装 Node.js 依赖..."
npm install

# 启动前端
echo "▶️  启动前端服务（端口：3000）..."
npm run dev &
FRONTEND_PID=$!

echo ""
echo "🎉 启动完成！"
echo ""
echo "📊 访问地址："
echo "   前端：http://localhost:3000"
echo "   后端：http://localhost:8080"
echo "   健康检查：http://localhost:8080/api/v1/health"
echo ""
echo "🛑 停止服务："
echo "   按 Ctrl+C 停止此脚本，或运行："
echo "   kill $BACKEND_PID $FRONTEND_PID"
echo ""

# 等待用户中断
trap 'echo ""; echo "🛑 正在停止服务..."; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0' INT

# 保持脚本运行
wait