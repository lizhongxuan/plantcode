#!/bin/bash

# 一键启动脚本
echo "🚀 启动AI开发平台..."

# 检查是否在项目根目录
if [ ! -f "go.mod" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

# 创建日志目录
mkdir -p logs

# 启动后端服务
echo "🔧 启动后端服务..."
./scripts/run.sh &
BACKEND_PID=$!
echo "后端服务启动，PID: $BACKEND_PID"

# 等待后端服务启动
sleep 5

# 启动前端服务
echo "🎨 启动前端服务..."
./scripts/run-frontend.sh &
FRONTEND_PID=$!
echo "前端服务启动，PID: $FRONTEND_PID"

echo ""
echo "🎉 服务启动完成！"
echo "后端服务: http://localhost:8080"
echo "前端服务: http://localhost:3000"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 等待用户中断
trap 'echo "正在停止服务..."; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0' SIGINT SIGTERM

wait 