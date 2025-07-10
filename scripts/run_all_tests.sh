#!/bin/bash

# 运行所有单元测试脚本
# 使用方法: ./scripts/run_all_tests.sh [options]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认设置
VERBOSE=false
COVERAGE=false
RACE=false
INTEGRATION=false
SPECIFIC_MODULE=""

# 显示帮助信息
show_help() {
    echo "运行所有单元测试脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -v, --verbose      显示详细输出"
    echo "  -c, --coverage     生成覆盖率报告"
    echo "  -r, --race         启用竞态检测"
    echo "  -i, --integration  运行集成测试"
    echo "  -m, --module NAME  只运行指定模块的测试"
    echo "  -h, --help         显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -c -v          # 运行所有测试并生成覆盖率报告"
    echo "  $0 -m config      # 只运行config模块的测试"
    echo "  $0 -r             # 启用竞态检测运行测试"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -r|--race)
            RACE=true
            shift
            ;;
        -i|--integration)
            INTEGRATION=true
            shift
            ;;
        -m|--module)
            SPECIFIC_MODULE="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 进入项目根目录
cd "$(dirname "$0")/.."

echo -e "${BLUE}=== AI开发平台单元测试 ===${NC}"
echo "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到Go命令${NC}"
    exit 1
fi

echo -e "${BLUE}Go版本:${NC}"
go version
echo ""

# 确保依赖已安装
echo -e "${YELLOW}检查依赖...${NC}"
go mod tidy
go mod download
echo ""

# 构建测试标志
TEST_FLAGS=""
if [ "$VERBOSE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -v"
fi

if [ "$RACE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -race"
fi

if [ "$COVERAGE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -cover -coverprofile=coverage.out"
fi

# 定义测试模块和路径
TEST_MODULES=(
    "config:./internal/config"
    "utils:./internal/utils"
    "models:./internal/model"
    "repository:./internal/repository"
    "services:./internal/service"
    "api:./internal/api"
    "ai:./internal/ai"
    "tests:./tests"
)

# 获取模块路径
get_module_path() {
    local module_name=$1
    for item in "${TEST_MODULES[@]}"; do
        if [[ $item == "$module_name:"* ]]; then
            echo "${item#*:}"
            return 0
        fi
    done
    echo ""
}

# 运行特定模块测试
run_module_test() {
    local module_name=$1
    local module_path=$2
    
    echo -e "${YELLOW}运行 ${module_name} 模块测试...${NC}"
    
    if [ -d "$module_path" ]; then
        # 检查是否有测试文件
        if ls ${module_path}/*_test.go 1> /dev/null 2>&1; then
            if go test $TEST_FLAGS $module_path; then
                echo -e "${GREEN}✓ ${module_name} 测试通过${NC}"
                return 0
            else
                echo -e "${RED}✗ ${module_name} 测试失败${NC}"
                return 1
            fi
        else
            echo -e "${YELLOW}! ${module_name} 模块没有测试文件${NC}"
            return 0
        fi
    else
        echo -e "${YELLOW}! ${module_name} 模块目录不存在: $module_path${NC}"
        return 0
    fi
}

# 统计变量
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# 如果指定了特定模块
if [ -n "$SPECIFIC_MODULE" ]; then
    module_path=$(get_module_path "$SPECIFIC_MODULE")
    if [ -n "$module_path" ]; then
        run_module_test "$SPECIFIC_MODULE" "$module_path"
        exit $?
    else
        echo -e "${RED}错误: 未知模块 '$SPECIFIC_MODULE'${NC}"
        echo "可用模块:"
        for item in "${TEST_MODULES[@]}"; do
            echo "  ${item%:*}"
        done
        exit 1
    fi
fi

echo -e "${BLUE}运行所有模块测试...${NC}"
echo ""

# 运行所有模块测试
for item in "${TEST_MODULES[@]}"; do
    module_name="${item%:*}"
    module_path="${item#*:}"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if run_module_test "$module_name" "$module_path"; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""
done

# 运行集成测试
if [ "$INTEGRATION" = true ]; then
    echo -e "${YELLOW}运行集成测试...${NC}"
    if [ -d "./tests/integration" ]; then
        if go test $TEST_FLAGS ./tests/integration/...; then
            echo -e "${GREEN}✓ 集成测试通过${NC}"
        else
            echo -e "${RED}✗ 集成测试失败${NC}"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    else
        echo -e "${YELLOW}! 未找到集成测试目录${NC}"
    fi
    echo ""
fi

# 生成覆盖率报告
if [ "$COVERAGE" = true ] && [ -f "coverage.out" ]; then
    echo -e "${YELLOW}生成覆盖率报告...${NC}"
    
    # 计算总覆盖率
    COVERAGE_PERCENT=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo -e "${BLUE}总覆盖率: ${COVERAGE_PERCENT}${NC}"
    
    # 生成HTML报告
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ HTML覆盖率报告已生成: coverage.html${NC}"
    
    # 显示覆盖率详情
    if [ "$VERBOSE" = true ]; then
        echo ""
        echo -e "${BLUE}覆盖率详情:${NC}"
        go tool cover -func=coverage.out
    fi
    echo ""
fi

# 显示测试总结
echo -e "${BLUE}=== 测试总结 ===${NC}"
echo "结束时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
echo -e "总模块数: ${TOTAL_TESTS}"
echo -e "${GREEN}通过: ${PASSED_TESTS}${NC}"
echo -e "${RED}失败: ${FAILED_TESTS}${NC}"
echo -e "${YELLOW}跳过: ${SKIPPED_TESTS}${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试都通过了！${NC}"
    exit 0
else
    echo -e "${RED}❌ 有 $FAILED_TESTS 个模块测试失败${NC}"
    exit 1
fi 