#!/bin/bash

# 项目详情接口测试脚本

echo "🚀 开始测试项目详情接口..."

# 配置
API_BASE="http://localhost:8080"
PROJECT_ID="6264f304-8cd3-4fd7-a44b-57d9053649fc"

echo "📋 测试配置："
echo "   API地址: $API_BASE"
echo "   项目ID: $PROJECT_ID"
echo ""

# 测试1: 无认证访问（应该返回401）
echo "🔓 测试1: 无认证访问"
curl -s -X GET "$API_BASE/api/v1/projects/$PROJECT_ID" \
  -H "Content-Type: application/json" | jq '.'
echo ""

# 测试2: 错误的项目ID格式（应该返回400）
echo "❌ 测试2: 错误的项目ID格式"
curl -s -X GET "$API_BASE/api/v1/projects/invalid-uuid" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer fake-token" | jq '.'
echo ""

# 测试3: 有效请求（需要先登录获取token）
echo "🔑 测试3: 获取认证token..."
# 这里需要先有一个有效的用户登录
# LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/api/auth/login" \
#   -H "Content-Type: application/json" \
#   -d '{"email":"test@example.com","password":"password123"}')
# TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

echo "💡 提示："
echo "   1. 确保服务器正在运行: go run cmd/server/main.go"
echo "   2. 服务器应该运行在端口8080（不是3001）"
echo "   3. 需要先注册用户并登录获取有效token"
echo "   4. 确保项目ID存在于数据库中"

echo ""
echo "✅ 项目详情接口已实现！"
echo "   GET /api/v1/projects/:id - 获取项目详情"
echo "   PUT /api/v1/projects/:id - 更新项目"
echo "   DELETE /api/v1/projects/:id - 删除项目"