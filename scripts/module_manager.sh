#!/bin/bash

# 模块化管理工具
# 用于管理项目模块、依赖关系和测试

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取模块路径
get_module_path() {
    local module=$1
    case $module in
        "api") echo "internal/api" ;;
        "service") echo "internal/service" ;;
        "repository") echo "internal/repository" ;;
        "model") echo "internal/model" ;;
        "config") echo "internal/config" ;;
        "utils") echo "internal/utils" ;;
        "ai") echo "internal/ai" ;;
        *) echo "" ;;
    esac
}

# 获取模块依赖
get_module_deps() {
    local module=$1
    case $module in
        "api") echo "service config" ;;
        "service") echo "repository model config utils ai" ;;
        "repository") echo "model config" ;;
        "model") echo "" ;;
        "config") echo "" ;;
        "utils") echo "" ;;
        "ai") echo "model config" ;;
        *) echo "" ;;
    esac
}

# 获取所有模块
get_all_modules() {
    echo "api service repository model config utils ai"
}

# 打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 检查模块依赖
check_module_deps() {
    local module=$1
    local deps=$(get_module_deps $module)
    
    print_message $BLUE "检查模块 $module 的依赖: $deps"
    
    for dep in $deps; do
        local dep_path=$(get_module_path $dep)
        if [[ -z "$dep_path" ]]; then
            print_message $RED "错误: 未知依赖模块 $dep"
            return 1
        fi
        
        if [[ ! -d "$dep_path" ]]; then
            print_message $RED "错误: 依赖模块 $dep 不存在: $dep_path"
            return 1
        fi
    done
    
    print_message $GREEN "模块 $module 的依赖检查通过"
}

# 分析模块复杂度
analyze_module_complexity() {
    local module=$1
    local module_path=$(get_module_path $module)
    
    if [[ -z "$module_path" ]]; then
        print_message $RED "错误: 未知模块 $module"
        return 1
    fi
    
    if [[ ! -d "$module_path" ]]; then
        print_message $RED "错误: 模块路径不存在: $module_path"
        return 1
    fi
    
    print_message $BLUE "分析模块 $module 的复杂度..."
    
    # 统计代码行数
    local total_lines=$(find $module_path -name "*.go" -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}' || echo "0")
    
    # 统计文件数量
    local file_count=$(find $module_path -name "*.go" 2>/dev/null | wc -l)
    
    # 统计函数数量
    local function_count=$(grep -r "^func " $module_path --include="*.go" 2>/dev/null | wc -l)
    
    # 统计结构体数量
    local struct_count=$(grep -r "^type.*struct" $module_path --include="*.go" 2>/dev/null | wc -l)
    
    # 统计接口数量
    local interface_count=$(grep -r "^type.*interface" $module_path --include="*.go" 2>/dev/null | wc -l)
    
    echo "模块复杂度报告:"
    echo "  代码行数: $total_lines"
    echo "  文件数量: $file_count"
    echo "  函数数量: $function_count"
    echo "  结构体数量: $struct_count"
    echo "  接口数量: $interface_count"
    
    # 计算复杂度分数（简单算法）
    local complexity_score=$((total_lines / 100 + file_count + function_count / 10))
    
    if [[ $complexity_score -lt 50 ]]; then
        print_message $GREEN "复杂度: 低 (分数: $complexity_score)"
    elif [[ $complexity_score -lt 100 ]]; then
        print_message $YELLOW "复杂度: 中 (分数: $complexity_score)"
    else
        print_message $RED "复杂度: 高 (分数: $complexity_score) - 建议重构"
    fi
}

# 运行模块测试
run_module_tests() {
    local module=$1
    
    print_message $BLUE "运行模块 $module 的测试..."
    
    # 检查是否存在测试文件
    local test_file="tests/${module}_test.go"
    if [[ ! -f "$test_file" ]]; then
        print_message $YELLOW "警告: 模块 $module 没有测试文件 $test_file"
        return 1
    fi
    
    # 运行测试
    go test -v -run "Test.*${module^}" ./tests/
    
    if [[ $? -eq 0 ]]; then
        print_message $GREEN "模块 $module 测试通过"
    else
        print_message $RED "模块 $module 测试失败"
        return 1
    fi
}

# 生成模块依赖图
generate_dependency_graph() {
    print_message $BLUE "生成模块依赖图..."
    
    # 创建 .dot 文件
    cat > module_dependencies.dot << 'EOF'
digraph module_dependencies {
    node [shape=box, style=filled, fillcolor=lightblue];
    
EOF
    
    # 添加节点和边
    for module in $(get_all_modules); do
        echo "    $module;" >> module_dependencies.dot
        
        local deps=$(get_module_deps $module)
        for dep in $deps; do
            echo "    $dep -> $module;" >> module_dependencies.dot
        done
    done
    
    echo "}" >> module_dependencies.dot
    
    # 生成图片（需要安装graphviz）
    if command -v dot &> /dev/null; then
        dot -Tpng module_dependencies.dot -o module_dependencies.png
        print_message $GREEN "依赖图已生成: module_dependencies.png"
    else
        print_message $YELLOW "警告: 未安装graphviz，无法生成图片"
        print_message $YELLOW "请运行: brew install graphviz (macOS) 或 apt-get install graphviz (Ubuntu)"
    fi
    
    print_message $GREEN "依赖图定义已生成: module_dependencies.dot"
}

# 创建新模块
create_module() {
    local module_name=$1
    local module_path="internal/$module_name"
    
    if [[ -d "$module_path" ]]; then
        print_message $YELLOW "警告: 模块 $module_name 已存在"
        return 1
    fi
    
    print_message $BLUE "创建新模块: $module_name"
    
    # 创建目录
    mkdir -p "$module_path"
    
    # 创建基础文件
    cat > "$module_path/${module_name}.go" << EOF
package $module_name

// ${module_name^} 模块的基本结构和接口定义
type ${module_name^}Service interface {
    // TODO: 定义接口方法
}

type ${module_name}Service struct {
    // TODO: 定义结构体字段
}

func New${module_name^}Service() ${module_name^}Service {
    return &${module_name}Service{}
}
EOF
    
    # 创建测试文件
    cat > "tests/${module_name}_test.go" << EOF
package tests

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type ${module_name^}TestSuite struct {
    TestSuite
}

func (s *${module_name^}TestSuite) SetupSuite() {
    s.TestSuite.SetupSuite()
}

func (s *${module_name^}TestSuite) TestExample() {
    // TODO: 添加测试用例
    s.True(true)
}

func Test${module_name^}(t *testing.T) {
    suite.Run(t, new(${module_name^}TestSuite))
}
EOF
    
    print_message $GREEN "模块 $module_name 创建完成"
}

# 运行回归测试
run_regression_tests() {
    print_message $BLUE "运行回归测试..."
    
    # 检查所有模块的依赖
    local all_passed=true
    for module in $(get_all_modules); do
        if ! check_module_deps "$module"; then
            all_passed=false
        fi
    done
    
    if [[ "$all_passed" == false ]]; then
        print_message $RED "依赖检查失败，停止回归测试"
        return 1
    fi
    
    # 运行所有模块的测试
    for module in $(get_all_modules); do
        if ! run_module_tests "$module"; then
            print_message $RED "模块 $module 测试失败"
            all_passed=false
        fi
    done
    
    if [[ "$all_passed" == true ]]; then
        print_message $GREEN "所有回归测试通过"
    else
        print_message $RED "部分回归测试失败"
        return 1
    fi
}

# 显示帮助信息
show_help() {
    echo "模块化管理工具"
    echo ""
    echo "用法: $0 [command] [options]"
    echo ""
    echo "命令:"
    echo "  list                    - 列出所有模块"
    echo "  deps <module>           - 检查模块依赖"
    echo "  analyze <module>        - 分析模块复杂度"
    echo "  test <module>           - 运行模块测试"
    echo "  test-all                - 运行所有模块测试"
    echo "  graph                   - 生成依赖图"
    echo "  create <module>         - 创建新模块"
    echo "  regression              - 运行回归测试"
    echo "  help                    - 显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 list"
    echo "  $0 deps api"
    echo "  $0 analyze service"
    echo "  $0 test user"
    echo "  $0 create payment"
}

# 主函数
main() {
    case "${1:-help}" in
        "list")
            print_message $BLUE "项目模块列表:"
            for module in $(get_all_modules); do
                local module_path=$(get_module_path $module)
                echo "  $module -> $module_path"
            done
            ;;
        "deps")
            if [[ -z "$2" ]]; then
                print_message $RED "错误: 请指定模块名"
                exit 1
            fi
            local module_path=$(get_module_path $2)
            if [[ -z "$module_path" ]]; then
                print_message $RED "错误: 未知模块 $2"
                exit 1
            fi
            check_module_deps "$2"
            ;;
        "analyze")
            if [[ -z "$2" ]]; then
                print_message $RED "错误: 请指定模块名"
                exit 1
            fi
            local module_path=$(get_module_path $2)
            if [[ -z "$module_path" ]]; then
                print_message $RED "错误: 未知模块 $2"
                exit 1
            fi
            analyze_module_complexity "$2"
            ;;
        "test")
            if [[ -z "$2" ]]; then
                print_message $RED "错误: 请指定模块名"
                exit 1
            fi
            run_module_tests "$2"
            ;;
        "test-all")
            for module in $(get_all_modules); do
                run_module_tests "$module" || true
            done
            ;;
        "graph")
            generate_dependency_graph
            ;;
        "create")
            if [[ -z "$2" ]]; then
                print_message $RED "错误: 请指定模块名"
                exit 1
            fi
            create_module "$2"
            ;;
        "regression")
            run_regression_tests
            ;;
        "help")
            show_help
            ;;
        *)
            print_message $RED "错误: 未知命令 $1"
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@" 