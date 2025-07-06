#!/bin/bash

# AI开发平台构建脚本

echo "🔨 构建AI开发平台..."

# 设置构建参数
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION=${VERSION:-"1.0.0"}

# 创建构建目录
mkdir -p bin

# 设置Go构建标志
export CGO_ENABLED=0
export GOOS=${GOOS:-linux}
export GOARCH=${GOARCH:-amd64}

echo "📍 构建目标: $GOOS/$GOARCH"
echo "📦 版本: $VERSION"
echo "🕒 构建时间: $BUILD_TIME"
echo "📋 Git提交: $GIT_COMMIT"

# 下载依赖
echo "📦 下载依赖..."
go mod tidy

# 构建二进制文件
echo "🔨 编译服务器..."
go build -ldflags "-s -w" -o bin/ai-dev-platform ./cmd/server

if [ $? -eq 0 ]; then
    echo "✅ 构建成功"
    
    # 显示构建信息
    echo ""
    echo "📁 构建输出:"
    ls -lh bin/
    echo ""
    
    echo "🚀 部署说明:"
    echo "1. 将 bin/ai-dev-platform 复制到目标服务器"
    echo "2. 设置环境变量或创建 .env 文件"
    echo "3. 运行: ./ai-dev-platform"
    
else
    echo "❌ 构建失败"
    exit 1
fi 