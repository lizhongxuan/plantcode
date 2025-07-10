#!/bin/bash

# AI开发平台认证流程集成测试脚本
set -e

echo "🚀 开始认证流程集成测试..."

# 配置
BASE_URL="http://localhost:8080"
TEST_USER="testflow$(date +%s)@example.com"
TEST_PASSWORD="password123"
TEST_USERNAME="testflow$(date +%s)"
TEST_FULLNAME="Test Flow User"

echo "📧 测试邮箱: $TEST_USER"

# 检查服务器是否运行
echo "🔍 检查服务器状态..."
if ! curl -s $BASE_URL/health > /dev/null; then
    echo "❌ 服务器未运行，请先启动服务器"
    exit 1
fi
echo "✅ 服务器运行正常"

# 1. 测试用户注册
echo "📝 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/auth/register \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USERNAME\",
        \"email\": \"$TEST_USER\",
        \"password\": \"$TEST_PASSWORD\",
        \"full_name\": \"$TEST_FULLNAME\"
    }")

echo "注册响应: $REGISTER_RESPONSE"

# 检查注册是否成功
if echo "$REGISTER_RESPONSE" | grep -q '"success":true'; then
    echo "✅ 用户注册成功"
else
    echo "❌ 用户注册失败"
    exit 1
fi

# 2. 测试用户登录
echo "🔐 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/auth/login \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$TEST_USER\",
        \"password\": \"$TEST_PASSWORD\"
    }")

echo "登录响应: $LOGIN_RESPONSE"

# 检查登录是否成功并提取token
if echo "$LOGIN_RESPONSE" | grep -q '"success":true'; then
    echo "✅ 用户登录成功"
    
    # 提取token（简单的文本处理）
    TOKEN=$(echo "$LOGIN_RESPONSE" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
    if [ -z "$TOKEN" ]; then
        echo "❌ 无法提取token"
        exit 1
    fi
    echo "📱 Token提取成功: ${TOKEN:0:50}..."
else
    echo "❌ 用户登录失败"
    exit 1
fi

# 3. 测试token验证
echo "🔍 测试token验证..."
VALIDATE_RESPONSE=$(curl -s -X GET $BASE_URL/api/auth/validate \
    -H "Authorization: Bearer $TOKEN")

echo "验证响应: $VALIDATE_RESPONSE"

if echo "$VALIDATE_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Token验证成功"
else
    echo "❌ Token验证失败"
    exit 1
fi

# 4. 测试受保护的API端点
echo "🛡️  测试受保护的API端点..."
PROFILE_RESPONSE=$(curl -s -X GET $BASE_URL/api/user/profile \
    -H "Authorization: Bearer $TOKEN")

echo "用户信息响应: $PROFILE_RESPONSE"

if echo "$PROFILE_RESPONSE" | grep -q '"success":true'; then
    echo "✅ 受保护端点访问成功"
else
    echo "❌ 受保护端点访问失败"
    exit 1
fi

# 5. 测试无效token
echo "🚫 测试无效token..."
INVALID_RESPONSE=$(curl -s -X GET $BASE_URL/api/auth/validate \
    -H "Authorization: Bearer invalid.token.here")

echo "无效token响应: $INVALID_RESPONSE"

if echo "$INVALID_RESPONSE" | grep -q '"success":false'; then
    echo "✅ 无效token被正确拒绝"
else
    echo "❌ 无效token验证行为异常"
    exit 1
fi

# 6. 测试无token访问受保护端点
echo "🔒 测试无token访问受保护端点..."
NO_TOKEN_RESPONSE=$(curl -s -X GET $BASE_URL/api/user/profile)

echo "无token响应: $NO_TOKEN_RESPONSE"

if echo "$NO_TOKEN_RESPONSE" | grep -q '"success":false'; then
    echo "✅ 无token访问被正确拒绝"
else
    echo "❌ 无token访问行为异常"
    exit 1
fi

echo ""
echo "🎉 所有认证流程测试通过！"
echo "✅ 用户注册 ✅ 用户登录 ✅ Token验证"
echo "✅ 受保护端点访问 ✅ 安全验证"
echo ""
echo "🔧 测试的功能："
echo "  - 用户注册功能"
echo "  - 用户登录功能"
echo "  - JWT Token生成和验证"
echo "  - 受保护API端点访问控制"
echo "  - 无效token和无token的安全处理"
echo ""
echo "📊 认证系统工作正常，可以正常使用登录功能！" 