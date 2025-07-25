#!/bin/bash

# 获取数据库配置
DB_HOST=$(grep DB_HOST .env | cut -d '=' -f2)
DB_PORT=$(grep DB_PORT .env | cut -d '=' -f2)
DB_USER=$(grep DB_USER .env | cut -d '=' -f2)
DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d '=' -f2)
DB_NAME=$(grep DB_NAME .env | cut -d '=' -f2)

# 如果配置为空，使用默认值
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-root}
DB_NAME=${DB_NAME:-aicode}

echo "正在清理数据库 $DB_NAME..."

# 执行SQL脚本
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < scripts/clean_database.sql

if [ $? -eq 0 ]; then
  echo "数据库清理成功！"
else
  echo "数据库清理失败，请检查错误信息。"
  exit 1
fi

echo "正在重新创建必要的表结构..."
# 重启服务器以重新创建表结构
go run cmd/server/main.go &
SERVER_PID=$!

# 等待服务器启动
sleep 5

# 停止服务器
kill $SERVER_PID

echo "数据库初始化完成！"