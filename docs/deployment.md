# AI辅助项目开发平台 - 部署文档

## 1. 部署概述

### 1.1 系统架构
```
Internet
    ↓
Nginx (负载均衡/反向代理)
    ↓
┌─────────────────────────────────────┐
│          Docker Compose             │
├─────────────────────────────────────┤
│  Go Backend      │  React Frontend  │
│  (Port 8080)     │  (Nginx:80)     │
├─────────────────────────────────────┤
│  MySQL 8.0       │  Redis 7.0      │
│  (Port 3306)     │  (Port 6379)    │
└─────────────────────────────────────┘
```

### 1.2 部署环境类型
- **开发环境 (Development)**：开发人员本地环境
- **测试环境 (Testing)**：功能测试和集成测试
- **预生产环境 (Staging)**：模拟生产环境
- **生产环境 (Production)**：线上正式环境

### 1.3 技术栈
- **容器化**：Docker + Docker Compose
- **反向代理**：Nginx
- **应用服务**：Go 1.21+
- **前端**：React 18 + TypeScript
- **数据库**：MySQL 8.0
- **缓存**：Redis 7.0
- **监控**：Prometheus + Grafana

## 2. 环境要求

### 2.1 硬件要求

#### 开发环境
- **CPU**: 2核心以上
- **内存**: 8GB以上
- **存储**: 50GB可用空间
- **网络**: 稳定互联网连接

#### 测试环境
- **CPU**: 4核心以上
- **内存**: 16GB以上
- **存储**: 100GB可用空间
- **网络**: 稳定内网连接

#### 生产环境
- **CPU**: 8核心以上
- **内存**: 32GB以上
- **存储**: 500GB SSD (数据库另配)
- **网络**: 高带宽、低延迟
- **备用**：建议主备双机

### 2.2 软件要求

#### 基础软件
```bash
# 操作系统
Ubuntu 20.04 LTS / CentOS 8+ / Amazon Linux 2

# 容器环境
Docker 24.0+
Docker Compose 2.20+

# 系统工具
git
curl
wget
unzip
```

#### 可选组件
```bash
# 监控工具
htop
iotop
nethogs

# 网络工具
netstat
ss
tcpdump
```

## 3. 项目结构

### 3.1 目录结构
```
ai-dev-platform/
├── cmd/
│   └── server/                 # Go应用入口
│       └── main.go
├── internal/                   # Go内部包
│   ├── api/                   # API层
│   ├── service/               # 业务逻辑层
│   ├── repository/            # 数据访问层
│   ├── ai/                    # AI服务集成
│   ├── model/                 # 数据模型
│   ├── config/                # 配置管理
│   └── utils/                 # 工具函数
├── web/                       # 前端项目
│   ├── src/                   # React源码
│   ├── public/                # 静态资源
│   ├── dist/                  # 构建输出
│   ├── package.json
│   └── vite.config.ts
├── deployments/               # 部署文件
│   ├── docker/                # Docker文件
│   ├── kubernetes/            # K8s配置
│   └── nginx/                 # Nginx配置
├── scripts/                   # 部署脚本
│   ├── deploy.sh              # 部署脚本
│   ├── backup.sh              # 备份脚本
│   └── init_db.sql            # 数据库初始化
├── .env.example               # 环境变量模板
├── docker-compose.yml         # Docker编排文件
├── docker-compose.prod.yml    # 生产环境配置
├── Dockerfile                 # Go应用Docker文件
├── web.Dockerfile             # 前端Docker文件
└── README.md
```

## 4. Docker容器化配置

### 4.1 Go后端Dockerfile
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 生产镜像
FROM alpine:3.18

# 安装ca证书
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./main"]
```

### 4.2 前端Dockerfile
```dockerfile
# web.Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

# 复制package文件
COPY web/package*.json ./
RUN npm ci --only=production

# 复制源码并构建
COPY web/ .
RUN npm run build

# 生产镜像
FROM nginx:1.24-alpine

# 复制构建文件
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY deployments/nginx/default.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 4.3 Docker Compose配置

#### 开发环境 (docker-compose.yml)
```yaml
version: '3.8'

services:
  # 后端服务
  backend:
    build: 
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - GO_ENV=development
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_NAME=ai_dev_platform
      - DB_USER=root
      - DB_PASSWORD=rootpassword
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=dev-jwt-secret
      - CORS_ORIGINS=http://localhost:3000
    depends_on:
      - mysql
      - redis
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped

  # 前端服务
  frontend:
    build:
      context: .
      dockerfile: web.Dockerfile
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped

  # MySQL数据库
  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=ai_dev_platform
      - MYSQL_CHARACTER_SET_SERVER=utf8mb4
      - MYSQL_COLLATION_SERVER=utf8mb4_unicode_ci
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped

  # Redis缓存
  redis:
    image: redis:7.0-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

## 5. 环境变量配置

### 5.1 环境变量模板 (.env.example)
```bash
# 应用环境
GO_ENV=production
PORT=8080

# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_NAME=ai_dev_platform
DB_USER=ai_user
DB_PASSWORD=your_db_password
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE=10

# Redis配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0
REDIS_MAX_CONNECTIONS=50

# JWT配置
JWT_SECRET=your_super_secret_jwt_key
JWT_EXPIRES_IN=3600

# AI服务配置
AI_PROVIDER=openai
OPENAI_API_KEY=your_openai_api_key
CLAUDE_API_KEY=your_claude_api_key
AI_TIMEOUT=30s
AI_MAX_RETRIES=3

# CORS配置
CORS_ORIGINS=https://ai-dev-platform.com,https://www.ai-dev-platform.com

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json

# 监控配置
METRICS_ENABLED=true
METRICS_PORT=9091

# 邮件配置 (可选)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_email_password

# 对象存储配置 (可选)
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=us-west-2
S3_BUCKET=ai-dev-platform-files
```

## 6. 部署脚本

### 6.1 自动部署脚本 (scripts/deploy.sh)
```bash
#!/bin/bash

# AI辅助项目开发平台部署脚本
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置参数
PROJECT_NAME="ai-dev-platform"
DEPLOY_ENV=${1:-"production"}
BACKUP_DIR="/backup/$(date +%Y%m%d_%H%M%S)"

echo -e "${GREEN}开始部署 AI辅助项目开发平台${NC}"
echo -e "${YELLOW}部署环境: $DEPLOY_ENV${NC}"

# 函数：打印信息
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 函数：检查系统要求
check_requirements() {
    log_info "检查系统要求..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    log_info "系统要求检查通过"
}

# 函数：拉取最新代码
pull_code() {
    log_info "拉取最新代码..."
    
    # 检查Git仓库
    if [ ! -d ".git" ]; then
        log_error "当前目录不是Git仓库"
        exit 1
    fi
    
    # 拉取最新代码
    git fetch origin
    git checkout main
    git pull origin main
    
    log_info "代码更新完成"
}

# 函数：构建应用
build_application() {
    log_info "构建应用..."
    
    # 构建Docker镜像
    case $DEPLOY_ENV in
        "development")
            docker-compose build
            ;;
        "production")
            docker-compose -f docker-compose.yml build
            ;;
        *)
            log_error "未知的部署环境: $DEPLOY_ENV"
            exit 1
            ;;
    esac
    
    log_info "应用构建完成"
}

# 函数：启动服务
start_services() {
    log_info "启动服务..."
    
    case $DEPLOY_ENV in
        "development")
            docker-compose up -d
            ;;
        "production")
            docker-compose up -d
            ;;
    esac
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 30
    
    log_info "服务启动完成"
}

# 函数：健康检查
health_check() {
    log_info "进行健康检查..."
    
    # 检查后端服务
    max_attempts=30
    attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f http://localhost:8080/health &> /dev/null; then
            log_info "后端服务健康检查通过"
            break
        fi
        
        if [ $attempt -eq $max_attempts ]; then
            log_error "后端服务健康检查失败"
            docker-compose logs backend
            exit 1
        fi
        
        log_warn "等待后端服务启动... ($attempt/$max_attempts)"
        sleep 5
        ((attempt++))
    done
    
    # 检查前端服务
    if curl -f http://localhost:3000 &> /dev/null; then
        log_info "前端服务健康检查通过"
    else
        log_error "前端服务健康检查失败"
        docker-compose logs frontend
        exit 1
    fi
    
    log_info "所有服务健康检查通过"
}

# 主部署流程
main() {
    echo "=========================================="
    echo "AI辅助项目开发平台部署脚本"
    echo "部署环境: $DEPLOY_ENV"
    echo "=========================================="
    
    # 检查系统要求
    check_requirements
    
    # 拉取最新代码
    pull_code
    
    # 构建应用
    build_application
    
    # 停止旧服务
    log_info "停止旧服务..."
    docker-compose down || true
    
    # 启动新服务
    start_services
    
    # 健康检查
    health_check
    
    echo "=========================================="
    echo -e "${GREEN}部署完成！${NC}"
    echo "前端访问地址: http://localhost:3000"
    echo "后端API地址: http://localhost:8080"
    echo "=========================================="
}

# 执行主函数
main "$@"
```

### 6.2 备份脚本 (scripts/backup.sh)
```bash
#!/bin/bash

# 数据备份脚本
set -e

# 配置
BACKUP_DIR="/backup/$(date +%Y%m%d_%H%M%S)"
MYSQL_CONTAINER="ai-dev-platform-mysql-1"
REDIS_CONTAINER="ai-dev-platform-redis-1"
RETENTION_DAYS=30

# 创建备份目录
mkdir -p $BACKUP_DIR

echo "开始备份数据到: $BACKUP_DIR"

# 备份MySQL数据库
echo "备份MySQL数据库..."
docker exec $MYSQL_CONTAINER mysqldump -u root -p${MYSQL_ROOT_PASSWORD} ai_dev_platform > $BACKUP_DIR/mysql_backup.sql

# 备份Redis数据
echo "备份Redis数据..."
docker exec $REDIS_CONTAINER redis-cli BGSAVE
docker cp $REDIS_CONTAINER:/data/dump.rdb $BACKUP_DIR/redis_backup.rdb

# 备份应用文件
echo "备份应用文件..."
tar -czf $BACKUP_DIR/app_files.tar.gz ./uploads ./logs

# 备份配置文件
echo "备份配置文件..."
cp .env $BACKUP_DIR/env_backup

# 压缩备份文件
echo "压缩备份文件..."
cd /backup
tar -czf $(basename $BACKUP_DIR).tar.gz $(basename $BACKUP_DIR)
rm -rf $BACKUP_DIR

# 清理过期备份
echo "清理过期备份..."
find /backup -name "*.tar.gz" -type f -mtime +$RETENTION_DAYS -delete

echo "备份完成!"
```

## 7. 运维管理

### 7.1 日常维护

#### 检查服务状态
```bash
# 查看容器状态
docker-compose ps

# 查看容器日志
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mysql
docker-compose logs -f redis
```

#### 系统监控
```bash
# 查看系统资源
htop
df -h
free -h

# 查看网络连接
netstat -tlnp
ss -tlnp
```

### 7.2 故障处理

#### 常见问题排查
```bash
# 重启服务
docker-compose restart backend
docker-compose restart frontend

# 查看错误日志
docker-compose logs backend | grep ERROR

# 完全重建
docker-compose down
docker-compose up -d
```

#### 数据恢复
```bash
# 恢复MySQL数据库
docker exec -i mysql_container mysql -u root -p${PASSWORD} ai_dev_platform < backup.sql

# 恢复Redis数据
docker cp backup.rdb redis_container:/data/dump.rdb
docker-compose restart redis
```

### 7.3 性能优化

#### 数据库优化
```sql
-- 检查慢查询
SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;

-- 优化表
OPTIMIZE TABLE users, projects, chat_messages;

-- 分析表
ANALYZE TABLE users, projects, chat_messages;
```

#### 缓存优化
```bash
# Redis内存使用分析
docker exec redis redis-cli info memory

# 清理过期键
docker exec redis redis-cli --scan --pattern "expired:*" | xargs docker exec redis redis-cli del
```

## 8. 安全配置

### 8.1 防火墙设置
```bash
# Ubuntu UFW配置
ufw allow 22/tcp     # SSH
ufw allow 80/tcp     # HTTP
ufw allow 443/tcp    # HTTPS
ufw --force enable
```

### 8.2 SSL证书配置
```bash
# 使用Let's Encrypt
certbot --nginx -d ai-dev-platform.com

# 自动续期
echo "0 12 * * * /usr/bin/certbot renew --quiet" | crontab -
```

### 8.3 访问控制
```nginx
# Nginx访问限制
location /admin {
    allow 192.168.1.0/24;
    deny all;
}

# 限制请求频率
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
location /api {
    limit_req zone=api burst=20;
}
```

## 9. 监控告警

### 9.1 基础监控
- 服务可用性监控
- 系统资源监控  
- 应用性能监控
- 错误日志监控

### 9.2 告警设置
```bash
# 磁盘空间告警
if [ $(df / | tail -1 | awk '{print $5}' | sed 's/%//') -gt 80 ]; then
    echo "磁盘空间不足" | mail -s "服务器告警" admin@example.com
fi

# 服务状态检查
if ! curl -f http://localhost:8080/health; then
    echo "后端服务异常" | mail -s "服务告警" admin@example.com
fi
```

## 10. 部署流程

### 10.1 首次部署

#### 准备阶段
```bash
# 1. 安装Docker和Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 2. 克隆项目代码
git clone https://github.com/your-org/ai-dev-platform.git
cd ai-dev-platform

# 3. 配置环境变量
cp .env.example .env
vim .env  # 修改配置

# 4. 创建必要目录
mkdir -p logs uploads backup
```

#### 执行部署
```bash
# 1. 赋予执行权限
chmod +x scripts/deploy.sh scripts/backup.sh

# 2. 执行部署
./scripts/deploy.sh production

# 3. 验证部署
curl http://localhost:8080/health
curl http://localhost:3000
```

### 10.2 更新部署流程

#### 滚动更新
```bash
# 1. 备份当前数据
./scripts/backup.sh

# 2. 拉取最新代码
git pull origin main

# 3. 重新构建并部署
./scripts/deploy.sh production

# 4. 验证更新
curl http://localhost:8080/health
```

## 11. 故障恢复

### 11.1 数据恢复
```bash
# 恢复数据库
docker exec -i mysql-container mysql -u root -p ai_dev_platform < backup/database.sql

# 恢复Redis
docker cp backup/dump.rdb redis-container:/data/
docker-compose restart redis

# 恢复文件
tar -xzf backup/app_files.tar.gz
```

### 11.2 服务恢复
```bash
# 完全重建
docker-compose down
docker system prune -f
docker-compose up -d

# 检查状态
docker-compose ps
```

## 12. 扩展部署

### 12.1 负载均衡
```nginx
upstream backend {
    least_conn;
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}
```

### 12.2 水平扩展
```yaml
# docker-compose.scale.yml
services:
  backend:
    deploy:
      replicas: 3
      
  nginx:
    depends_on:
      - backend
    ports:
      - "80:80"
```

这个部署文档提供了从开发到生产的完整部署指南，涵盖了容器化、安全、监控、运维等关键方面，为项目的成功部署和稳定运行提供了全面指导。
