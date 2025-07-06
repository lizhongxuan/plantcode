#!/bin/bash

# AI开发平台环境设置脚本

echo "🚀 设置AI开发平台开发环境..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请先安装Go 1.21或更高版本"
    exit 1
fi

echo "✅ Go版本: $(go version)"

# 检查MySQL是否运行
if ! command -v mysql &> /dev/null; then
    echo "⚠️ MySQL未安装或未在PATH中，请确保MySQL可用"
fi

# 检查Redis是否运行
if ! command -v redis-cli &> /dev/null; then
    echo "⚠️ Redis未安装或未在PATH中，请确保Redis可用"
fi

# 创建环境配置文件
if [ ! -f .env ]; then
    echo "📝 创建环境配置文件..."
    cat > .env << 'EOF'
# 服务器配置
PORT=8080
ENV=development

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ai_dev_platform
DB_MAX_CONNECTIONS=50
DB_MAX_IDLE_CONN=10
DB_CONN_MAX_LIFETIME=3600

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5

# JWT配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=86400

# AI服务配置
AI_API_KEY=your-ai-api-key
AI_BASE_URL=https://api.openai.com/v1
AI_MODEL=gpt-3.5-turbo
AI_MAX_TOKENS=2048
AI_TEMPERATURE=0.7

# CORS配置
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_HEADERS=Content-Type,Authorization,X-Requested-With
CORS_CREDENTIALS=true
EOF
    echo "✅ 环境配置文件已创建: .env"
    echo "⚠️ 请编辑 .env 文件设置正确的数据库密码和API密钥"
else
    echo "✅ 环境配置文件已存在"
fi

# 下载Go依赖
echo "📦 下载Go依赖..."
go mod tidy

# 创建必要的目录
echo "📁 创建项目目录..."
mkdir -p {web/src,web/public,logs,data}

# 设置权限
chmod +x scripts/*.sh

echo ""
echo "🎉 环境设置完成！"
echo ""
echo "📋 接下来的步骤："
echo "1. 编辑 .env 文件设置数据库和API配置"
echo "2. 启动MySQL和Redis服务"
echo "3. 运行: ./scripts/run.sh"
echo ""
echo "🔗 有用的命令："
echo "  启动服务: ./scripts/run.sh"
echo "  构建项目: ./scripts/build.sh"
echo "  运行测试: go test ./..."
echo "  查看日志: tail -f logs/app.log" 