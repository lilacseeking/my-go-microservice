#!/bin/bash
# API 测试脚本 - 使用 curl 命令测试 HTTP 接口
# 用法: ./scripts/test-api.sh [base_url]
# 默认 base_url: http://localhost:8080

BASE_URL="${1:-http://localhost:8080}"

echo "========================================="
echo "API 测试脚本"
echo "Base URL: $BASE_URL"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数器
TOTAL=0
PASSED=0
FAILED=0

# 测试结果函数
print_result() {
    TOTAL=$((TOTAL + 1))
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ 测试通过${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗ 测试失败${NC}"
        FAILED=$((FAILED + 1))
    fi
    echo ""
}

# ==========================================
# 测试 1: 健康检查
# ==========================================
HEALTH_URL="${BASE_URL//8080/8081}"
echo -e "${YELLOW}[测试 1] 健康检查${NC}"
echo "GET $HEALTH_URL/healthz"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$HEALTH_URL/healthz")
RESPONSE=$(curl -s "$HEALTH_URL/healthz")
echo "HTTP Status: $HTTP_CODE"
echo "Response: $RESPONSE"
if [ "$HTTP_CODE" = "200" ]; then
    print_result 0
else
    print_result 1
fi

# ==========================================
# 测试 2: 创建用户 - 成功场景
# ==========================================
echo -e "${YELLOW}[测试 2] 创建用户 - 成功场景${NC}"
echo "POST $BASE_URL/v1/users"
CREATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com"
  }')
HTTP_CODE=$(echo "$CREATE_RESPONSE" | tail -n1)
BODY=$(echo "$CREATE_RESPONSE" | sed '$d')
echo "HTTP Status: $HTTP_CODE"
echo "Response: $BODY"
if [ "$HTTP_CODE" = "201" ]; then
    print_result 0
    # 提取用户ID供后续测试使用
    USER_ID=$(echo "$BODY" | grep -o '"user_id":"[^"]*"' | cut -d'"' -f4)
else
    print_result 1
    USER_ID=""
fi

# ==========================================
# 测试 3: 创建用户 - 参数验证失败
# ==========================================
echo -e "${YELLOW}[测试 3] 创建用户 - 缺少必填字段${NC}"
echo "POST $BASE_URL/v1/users"
INVALID_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "email": "invalid-email"
  }')
HTTP_CODE=$(echo "$INVALID_RESPONSE" | tail -n1)
BODY=$(echo "$INVALID_RESPONSE" | sed '$d')
echo "HTTP Status: $HTTP_CODE"
echo "Response: $BODY"
if [ "$HTTP_CODE" = "400" ]; then
    print_result 0
else
    print_result 1
fi

# ==========================================
# 测试 4: 获取用户 - 使用测试数据中的ID
# ==========================================
echo -e "${YELLOW}[测试 4] 获取用户 - test-user-001${NC}"
echo "GET $BASE_URL/v1/users/test-user-001"
GET_RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/users/test-user-001")
HTTP_CODE=$(echo "$GET_RESPONSE" | tail -n1)
BODY=$(echo "$GET_RESPONSE" | sed '$d')
echo "HTTP Status: $HTTP_CODE"
echo "Response: $BODY"
if [ "$HTTP_CODE" = "200" ]; then
    print_result 0
else
    print_result 1
fi

# ==========================================
# 测试 5: 获取用户 - 使用刚创建的用户ID
# ==========================================
if [ -n "$USER_ID" ]; then
    echo -e "${YELLOW}[测试 5] 获取用户 - 新创建的用户${NC}"
    echo "GET $BASE_URL/v1/users/$USER_ID"
    GET_NEW_RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/users/$USER_ID")
    HTTP_CODE=$(echo "$GET_NEW_RESPONSE" | tail -n1)
    BODY=$(echo "$GET_NEW_RESPONSE" | sed '$d')
    echo "HTTP Status: $HTTP_CODE"
    echo "Response: $BODY"
    if [ "$HTTP_CODE" = "200" ]; then
        print_result 0
    else
        print_result 1
    fi
fi

# ==========================================
# 测试 6: 获取用户 - 不存在的用户
# ==========================================
echo -e "${YELLOW}[测试 6] 获取用户 - 不存在的用户${NC}"
echo "GET $BASE_URL/v1/users/non-existent-id"
NOT_FOUND_RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/users/non-existent-id")
HTTP_CODE=$(echo "$NOT_FOUND_RESPONSE" | tail -n1)
BODY=$(echo "$NOT_FOUND_RESPONSE" | sed '$d')
echo "HTTP Status: $HTTP_CODE"
echo "Response: $BODY"
if [ "$HTTP_CODE" = "404" ]; then
    print_result 0
else
    print_result 1
fi

# ==========================================
# 测试 7: 创建第二个用户
# ==========================================
echo -e "${YELLOW}[测试 7] 创建第二个用户${NC}"
echo "POST $BASE_URL/v1/users"
CREATE_RESPONSE2=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "李四",
    "email": "lisi@example.com"
  }')
HTTP_CODE=$(echo "$CREATE_RESPONSE2" | tail -n1)
BODY=$(echo "$CREATE_RESPONSE2" | sed '$d')
echo "HTTP Status: $HTTP_CODE"
echo "Response: $BODY"
if [ "$HTTP_CODE" = "201" ]; then
    print_result 0
else
    print_result 1
fi

# ==========================================
# 测试汇总
# ==========================================
echo ""
echo "========================================="
echo "测试汇总"
echo "========================================="
echo -e "总测试数: ${TOTAL}"
echo -e "通过: ${GREEN}${PASSED}${NC}"
echo -e "失败: ${RED}${FAILED}${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过! ✓${NC}"
    exit 0
else
    echo -e "${RED}部分测试失败 ✗${NC}"
    exit 1
fi
