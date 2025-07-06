#!/bin/bash

# 前端启动脚本
echo "🎨 启动前端开发服务器..."

# 检查是否在项目根目录
if [ ! -f "package.json" ] && [ ! -d "web" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

# 检查web目录是否存在
if [ ! -d "web" ]; then
    echo "❌ web目录不存在"
    exit 1
fi

# 进入前端目录
cd web

# 检查package.json是否存在
if [ ! -f "package.json" ]; then
    echo "❌ package.json不存在，请先运行 ./scripts/setup-frontend.sh"
    exit 1
fi

# 检查node_modules是否存在
if [ ! -d "node_modules" ]; then
    echo "📦 node_modules不存在，正在安装依赖..."
    npm install
fi

# 启动开发服务器
echo "🚀 启动前端开发服务器..."
echo "前端服务将在 http://localhost:3000 启动"
echo "按 Ctrl+C 停止服务"
echo ""

npm run dev 