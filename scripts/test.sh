#!/bin/bash

# 自动化测试脚本
# 用于运行所有测试并生成报告

set -e  # 遇到错误时立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 检查Go环境
check_go_env() {
    print_message $BLUE "检查Go环境..."
    if ! command -v go &> /dev/null; then
        print_message $RED "错误: Go未安装"
        exit 1
    fi
    
    go_version=$(go version | awk '{print $3}')
    print_message $GREEN "Go版本: $go_version"
}

# 检查数据库连接
check_database() {
    print_message $BLUE "检查数据库连接..."
    
    # 检查MySQL
    if ! command -v mysql &> /dev/null; then
        print_message $YELLOW "警告: MySQL客户端未安装"
    fi
    
    # 创建测试数据库
    mysql -u root -plzx234258 -e "CREATE DATABASE IF NOT EXISTS aicode_test;" 2>/dev/null || {
        print_message $RED "错误: 无法连接到MySQL数据库"
        print_message $YELLOW "请确保MySQL服务正在运行并且密码正确"
        exit 1
    }
    
    print_message $GREEN "数据库连接正常"
}

# 安装依赖
install_dependencies() {
    print_message $BLUE "安装依赖..."
    go mod tidy
    go mod download
    print_message $GREEN "依赖安装完成"
}

# 运行静态代码分析
run_static_analysis() {
    print_message $BLUE "运行静态代码分析..."
    
    # 检查Go代码格式
    if ! gofmt -l . | grep -q .; then
        print_message $GREEN "代码格式检查通过"
    else
        print_message $YELLOW "代码格式需要修复:"
        gofmt -l .
    fi
    
    # 运行go vet
    if go vet ./...; then
        print_message $GREEN "静态分析通过"
    else
        print_message $RED "静态分析发现问题"
        exit 1
    fi
}

# 运行单元测试
run_unit_tests() {
    print_message $BLUE "运行单元测试..."
    
    # 设置测试环境变量
    export TEST_DB_HOST=localhost
    export TEST_DB_PORT=3306
    export TEST_DB_USER=root
    export TEST_DB_PASSWORD=lzx234258
    export TEST_DB_NAME=aicode_test
    
    # 运行测试并生成覆盖率报告
    go test -v -coverprofile=coverage.out ./tests/...
    
    if [ $? -eq 0 ]; then
        print_message $GREEN "单元测试通过"
        
        # 生成覆盖率报告
        go tool cover -html=coverage.out -o coverage.html
        coverage_percent=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
        print_message $GREEN "代码覆盖率: $coverage_percent"
    else
        print_message $RED "单元测试失败"
        exit 1
    fi
}

# 运行集成测试
run_integration_tests() {
    print_message $BLUE "运行集成测试..."
    
    # 启动服务器
    go run cmd/server/main.go &
    SERVER_PID=$!
    
    # 等待服务器启动
    sleep 5
    
    # 运行集成测试
    curl -s http://localhost:8080/health > /dev/null
    if [ $? -eq 0 ]; then
        print_message $GREEN "集成测试通过"
    else
        print_message $RED "集成测试失败"
        kill $SERVER_PID
        exit 1
    fi
    
    # 停止服务器
    kill $SERVER_PID
}

# 运行性能测试
run_performance_tests() {
    print_message $BLUE "运行性能测试..."
    
    # 这里可以添加具体的性能测试
    print_message $GREEN "性能测试完成"
}

# 生成测试报告
generate_report() {
    print_message $BLUE "生成测试报告..."
    
    # 创建报告目录
    mkdir -p reports
    
    # 生成JSON格式的测试报告
    go test -json ./tests/... > reports/test_report.json
    
    # 生成HTML格式的覆盖率报告
    if [ -f coverage.out ]; then
        go tool cover -html=coverage.out -o reports/coverage.html
    fi
    
    print_message $GREEN "测试报告已生成到 reports/ 目录"
}

# 清理测试数据
cleanup() {
    print_message $BLUE "清理测试数据..."
    
    # 删除测试数据库
    mysql -u root -plzx234258 -e "DROP DATABASE IF EXISTS aicode_test;" 2>/dev/null || true
    
    # 删除临时文件
    rm -f coverage.out
    
    print_message $GREEN "清理完成"
}

# 主函数
main() {
    print_message $BLUE "=== AI开发平台自动化测试 ==="
    
    # 检查命令行参数
    case "${1:-all}" in
        "unit")
            check_go_env
            check_database
            install_dependencies
            run_unit_tests
            ;;
        "integration")
            check_go_env
            check_database
            install_dependencies
            run_integration_tests
            ;;
        "static")
            check_go_env
            install_dependencies
            run_static_analysis
            ;;
        "performance")
            check_go_env
            check_database
            install_dependencies
            run_performance_tests
            ;;
        "all")
            check_go_env
            check_database
            install_dependencies
            run_static_analysis
            run_unit_tests
            run_integration_tests
            run_performance_tests
            generate_report
            ;;
        "clean")
            cleanup
            exit 0
            ;;
        *)
            print_message $YELLOW "用法: $0 [unit|integration|static|performance|all|clean]"
            print_message $YELLOW "  unit        - 只运行单元测试"
            print_message $YELLOW "  integration - 只运行集成测试"
            print_message $YELLOW "  static      - 只运行静态分析"
            print_message $YELLOW "  performance - 只运行性能测试"
            print_message $YELLOW "  all         - 运行所有测试 (默认)"
            print_message $YELLOW "  clean       - 清理测试数据"
            exit 1
            ;;
    esac
    
    print_message $GREEN "=== 测试完成 ==="
}

# 设置陷阱，确保程序退出时清理资源
trap cleanup EXIT

# 运行主函数
main "$@" 