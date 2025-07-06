#!/bin/bash

# AI开发平台启动脚本

echo "🚀 启动AI开发平台..."

# 检查环境配置文件
if [ ! -f .env ]; then
    echo "❌ 未找到 .env 配置文件"
    echo "请先运行: ./scripts/setup.sh"
    exit 1
fi

# 加载环境变量
set -a
source .env
set +a

echo "✅ 已加载环境配置"
echo "📍 环境: ${ENV:-development}"
echo "🌐 端口: ${PORT:-8080}"

# 检查必需的环境变量
if [ -z "$DB_PASSWORD" ] || [ "$DB_PASSWORD" = "your_password" ]; then
    echo "⚠️ 警告: 请在 .env 文件中设置正确的数据库密码"
fi

if [ -z "$JWT_SECRET" ] || [ "$JWT_SECRET" = "your-super-secret-jwt-key-change-this-in-production" ]; then
    echo "⚠️ 警告: 请在 .env 文件中设置安全的JWT密钥"
fi

# 创建日志目录
mkdir -p logs

# 检查Go模块
if [ ! -f go.mod ]; then
    echo "❌ 未找到 go.mod 文件"
    exit 1
fi

# 下载依赖
echo "📦 检查依赖..."
go mod tidy

# 构建并运行
echo "🔨 构建项目..."
if go build -o bin/server ./cmd/server; then
    echo "✅ 构建成功"
    echo ""
    echo "🎯 启动服务器..."
    echo "📱 访问: http://localhost:${PORT:-8080}"
    echo "🏥 健康检查: http://localhost:${PORT:-8080}/health"
    echo ""
    echo "按 Ctrl+C 停止服务器"
    echo "=========================="
    
    # 启动服务器
    ./bin/server
else
    echo "❌ 构建失败"
    exit 1
fi 