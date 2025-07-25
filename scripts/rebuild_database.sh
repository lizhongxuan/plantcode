#!/bin/bash

# 数据库重建脚本
# 解决外键约束冲突问题

set -e  # 遇到错误立即退出

# 数据库配置
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD="lzx234258"
DB_NAME="aicode"

echo "🔄 开始数据库重建过程..."

# 检查MySQL连接
echo "📡 检查数据库连接..."
if ! mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ 无法连接到MySQL数据库，请检查："
    echo "   - MySQL服务是否启动"
    echo "   - 用户名密码是否正确: $DB_USER/$DB_PASSWORD"
    echo "   - 主机端口是否正确: $DB_HOST:$DB_PORT"
    exit 1
fi
echo "✅ 数据库连接成功"

# 1. 删除所有表（解决外键约束）
echo "🗑️  删除所有现有表..."
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < scripts/rebuild_database.sql
if [ $? -eq 0 ]; then
    echo "✅ 现有表删除完成"
else
    echo "❌ 删除表失败"
    exit 1
fi

# 2. 启动Go服务器进行自动迁移
echo "🚀 启动Go服务器进行表结构自动创建..."
echo "   提示：服务器将自动创建所需的表结构"

# 设置超时，避免服务器一直运行
timeout 30s go run cmd/server/main.go &
SERVER_PID=$!

# 等待服务器启动和表创建
echo "⏳ 等待表结构创建（最多30秒）..."
sleep 10

# 检查服务器进程
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "🛑 停止服务器..."
    kill $SERVER_PID
    wait $SERVER_PID 2>/dev/null || true
fi

# 3. 验证表是否创建成功
echo "🔍 验证表结构..."
TABLES=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES;" -s)

if [ -n "$TABLES" ]; then
    echo "✅ 数据库重建成功！创建的表："
    echo "$TABLES" | sed 's/^/   - /'
    echo ""
    echo "🎉 现在可以正常启动服务器了！"
    echo "   运行: go run cmd/server/main.go"
else
    echo "❌ 表创建失败，请检查Go服务器日志"
    exit 1
fi