#!/bin/bash

# 开发环境启动脚本
echo "🚀 启动AI开发平台开发环境..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查.env文件
if [ ! -f ".env" ]; then
    echo -e "${RED}❌ 未找到.env文件${NC}"
    exit 1
fi

# 加载环境变量
source .env
echo -e "${GREEN}✅ 已加载环境配置${NC}"
echo -e "${GREEN}📍 环境: $GO_ENV${NC}"
echo -e "${GREEN}🌐 端口: $PORT${NC}"

# 检查数据库配置
if [ "$DB_PASSWORD" = "password" ]; then
    echo -e "${YELLOW}⚠️ 警告: 请在 .env 文件中设置正确的数据库密码${NC}"
fi

# 检查JWT密钥
if [ "$JWT_SECRET" = "ai-dev-platform-secret" ]; then
    echo -e "${YELLOW}⚠️ 警告: 请在 .env 文件中设置安全的JWT密钥${NC}"
fi

# 检查MySQL连接
echo "🔍 检查数据库连接..."
if command -v mysql &> /dev/null; then
    if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" &> /dev/null; then
        echo -e "${GREEN}✅ MySQL连接正常${NC}"
    else
        echo -e "${YELLOW}⚠️ MySQL连接失败，应用将在有限功能模式下运行${NC}"
        echo -e "${YELLOW}💡 提示: 请检查MySQL服务状态和配置${NC}"
    fi
else
    echo -e "${YELLOW}⚠️ 未找到MySQL客户端${NC}"
fi

# 检查Redis连接
echo "🔍 检查Redis连接..."
if command -v redis-cli &> /dev/null; then
    if redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping &> /dev/null; then
        echo -e "${GREEN}✅ Redis连接正常${NC}"
    else
        echo -e "${YELLOW}⚠️ Redis连接失败，缓存功能将不可用${NC}"
        echo -e "${YELLOW}💡 提示: 请安装并启动Redis服务${NC}"
    fi
else
    echo -e "${YELLOW}⚠️ 未找到Redis客户端${NC}"
fi

# 检查依赖
echo "📦 检查依赖..."
if ! go mod verify &> /dev/null; then
    echo -e "${YELLOW}⚠️ 依赖验证失败，正在重新下载...${NC}"
    go mod tidy
fi

# 构建项目
echo "🔨 构建项目..."
if go build -o bin/server cmd/server/main.go; then
    echo -e "${GREEN}✅ 构建成功${NC}"
else
    echo -e "${RED}❌ 构建失败${NC}"
    exit 1
fi

# 启动服务器
echo "🎯 启动服务器..."
echo -e "${GREEN}📱 访问: http://localhost:$PORT${NC}"
echo -e "${GREEN}🏥 健康检查: http://localhost:$PORT/health${NC}"
echo -e "${GREEN}📄 API文档: http://localhost:$PORT/api/docs${NC}"
echo ""
echo "按 Ctrl+C 停止服务器"
echo "=========================="

# 启动应用
exec ./bin/server 