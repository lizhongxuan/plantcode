#!/bin/bash

# PlantUML 中文支持服务启动脚本
# 作者: AI开发平台
# 功能: 启动支持中文显示的PlantUML服务器

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行，请启动 Docker 服务"
        exit 1
    fi
    
    log_success "Docker 检查通过"
}

# 检查docker-compose是否安装
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "docker-compose 未安装，请先安装 docker-compose"
        exit 1
    fi
    
    log_success "docker-compose 检查通过"
}

# 停止现有的PlantUML服务
stop_existing_service() {
    log_info "检查并停止现有的 PlantUML 服务..."
    
    # 停止可能存在的容器
    if docker ps -q --filter "name=plantuml-server" | grep -q .; then
        log_warning "发现运行中的 plantuml-server 容器，正在停止..."
        docker stop plantuml-server || true
        docker rm plantuml-server || true
    fi
    
    if docker ps -q --filter "name=plantuml-server-chinese" | grep -q .; then
        log_warning "发现运行中的 plantuml-server-chinese 容器，正在停止..."
        docker stop plantuml-server-chinese || true
    fi
    
    # 检查端口8888是否被占用
    if lsof -Pi :8888 -sTCP:LISTEN -t >/dev/null 2>&1; then
        log_warning "端口 8888 被占用，请检查是否有其他服务在使用该端口"
        log_info "您可以使用以下命令查看占用端口的进程："
        echo "  lsof -i :8888"
        read -p "是否继续启动？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "用户取消启动"
            exit 0
        fi
    fi
}

# 构建并启动服务
start_service() {
    log_info "构建并启动 PlantUML 中文支持服务..."
    
    # 进入docker目录
    cd "$(dirname "$0")"
    
    # 使用docker-compose构建并启动
    if command -v docker-compose &> /dev/null; then
        docker-compose up -d --build
    else
        docker compose up -d --build
    fi
    
    log_success "PlantUML 服务启动成功！"
}

# 等待服务就绪
wait_for_service() {
    log_info "等待服务启动完成..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "http://localhost:8888/" > /dev/null 2>&1; then
            log_success "服务已就绪！"
            return 0
        fi
        
        log_info "等待服务启动... ($attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
    
    log_error "服务启动超时，请检查容器日志："
    echo "  docker logs plantuml-server-chinese"
    return 1
}

# 显示服务信息
show_service_info() {
    echo
    echo "========================================"
    echo "PlantUML 中文支持服务已启动"
    echo "========================================"
    echo "🌐 访问地址: http://localhost:8888"
    echo "📝 测试地址: http://localhost:8888/uml/SyfFKj2rKt3CoKnELR1Io4ZDoSa70000"
    echo
    echo "📋 常用命令:"
    echo "  查看容器状态: docker ps | grep plantuml"
    echo "  查看容器日志: docker logs plantuml-server-chinese"
    echo "  停止服务:     docker stop plantuml-server-chinese"
    echo "  重启服务:     docker restart plantuml-server-chinese"
    echo
    echo "🔧 测试中文支持:"
    echo "  可以访问上面的测试地址，或在代码中添加以下配置："
    echo "  !theme plain"
    echo "  skinparam defaultFontName \"WenQuanYi Zen Hei,WenQuanYi Micro Hei,sans-serif\""
    echo
}

# 主函数
main() {
    echo "========================================"
    echo "PlantUML 中文支持服务 启动脚本"
    echo "========================================"
    
    check_docker
    check_docker_compose
    stop_existing_service
    start_service
    
    if wait_for_service; then
        show_service_info
    else
        log_error "服务启动失败"
        exit 1
    fi
}

# 执行主函数
main "$@"