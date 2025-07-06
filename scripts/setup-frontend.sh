#!/bin/bash

# 前端设置脚本
echo "🎨 设置前端环境..."

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js"
    exit 1
fi

# 检查npm是否安装
if ! command -v npm &> /dev/null; then
    echo "❌ npm 未安装，请先安装 npm"
    exit 1
fi

echo "✅ Node.js 版本: $(node --version)"
echo "✅ npm 版本: $(npm --version)"

# 进入前端目录
cd web

# 安装依赖
echo "📦 安装前端依赖..."
npm install

# 检查是否安装成功
if [ $? -eq 0 ]; then
    echo "✅ 前端依赖安装成功！"
else
    echo "❌ 前端依赖安装失败"
    exit 1
fi

echo "🎉 前端环境设置完成！"
echo ""
echo "使用以下命令启动前端开发服务器:"
echo "  cd web && npm run dev"
echo ""
echo "或者使用快捷脚本:"
echo "  ./scripts/run-frontend.sh" 