# AI辅助项目开发平台 - 技术规范文档

## 1. 技术栈选型

### 1.1 后端技术栈

#### 核心框架
- **编程语言**：Go 1.21+
- **Web框架**：标准库 `net/http`
- **路由管理**：标准库 `http.ServeMux` + 自定义路由器
- **JSON处理**：标准库 `encoding/json`
- **HTTP客户端**：标准库 `net/http`

#### 数据存储
- **关系型数据库**：MySQL 8.0+
- **数据库驱动**：`github.com/go-sql-driver/mysql`
- **缓存系统**：Redis 7.0+
- **Redis客户端**：`github.com/redis/go-redis/v9`

#### AI服务集成
- **AI服务接入**：通过HTTP API调用
- **支持的AI服务**：
  - OpenAI GPT-4/GPT-3.5
  - Claude (Anthropic)
  - 通过配置api-key动态切换
- **HTTP客户端**：标准库 `net/http`

### 1.2 前端技术栈

#### 核心框架
- **基础框架**：React 18+ + TypeScript 5+
- **构建工具**：Vite
- **状态管理**：Zustand
- **UI组件库**：Ant Design 5+

#### 特殊组件
- **代码编辑器**：Monaco Editor (支持PUML语法高亮)
- **图表渲染**：PlantUML Server + 自定义渲染器
- **Markdown渲染**：react-markdown
- **文件处理**：FileSaver.js

## 2. 系统架构设计

### 2.1 整体架构模式

采用**分层架构 + 微服务理念**：

```
前端层 (React)
    ↓
API网关层 (Go Router)
    ↓
业务服务层 (Go Services)
    ↓
数据访问层 (MySQL + Redis)
    ↓
外部服务层 (AI APIs)
```

### 2.2 核心模块划分

#### 2.2.1 前端模块
- **用户界面模块** (`src/components/`)
- **PUML编辑器模块** (`src/editor/`)
- **对话界面模块** (`src/chat/`)
- **项目管理模块** (`src/project/`)
- **状态管理模块** (`src/store/`)

#### 2.2.2 后端模块
- **API路由模块** (`internal/api/`)
- **业务逻辑模块** (`internal/service/`)
- **数据访问模块** (`internal/repository/`)
- **AI服务模块** (`internal/ai/`)
- **工具模块** (`internal/utils/`)

### 2.3 数据流设计

#### 用户请求流程
```
用户操作 → React组件 → API调用 → Go路由器 → 业务服务 → 数据库/AI服务 → 响应返回
```

#### AI处理流程
```
用户输入 → 需求分析服务 → AI API → 结果处理 → PUML生成 → 数据库存储 → 前端展示
```

## 3. 编码规范

### 3.1 Go后端编码规范

#### 文件组织结构
```
cmd/
├── server/
│   └── main.go              # 程序入口
internal/
├── api/
│   ├── router.go            # 路由定义
│   ├── middleware.go        # 中间件
│   └── handlers/            # 处理器
├── service/
│   ├── requirement.go       # 需求分析服务
│   ├── puml.go             # PUML生成服务
│   ├── document.go         # 文档生成服务
│   └── module.go           # 模块管理服务
├── repository/
│   ├── project.go          # 项目数据访问
│   ├── conversation.go     # 对话数据访问
│   └── module.go           # 模块数据访问
├── ai/
│   ├── client.go           # AI客户端
│   ├── openai.go           # OpenAI实现
│   └── claude.go           # Claude实现
├── model/
│   └── *.go                # 数据模型
├── config/
│   └── config.go           # 配置管理
└── utils/
    └── *.go                # 工具函数
```

#### 命名约定
- **包名**：小写，简短，描述性
- **文件名**：小写+下划线，如 `user_service.go`
- **函数名**：驼峰命名，公有函数首字母大写
- **变量名**：驼峰命名，私有变量首字母小写
- **常量名**：全大写+下划线，如 `MAX_RETRY_COUNT`

#### 代码风格
```go
// 标准的处理器函数签名
func (h *Handler) HandleMethod(w http.ResponseWriter, r *http.Request) {
    // 1. 参数验证
    // 2. 业务逻辑调用
    // 3. 响应处理
}

// 统一的错误响应格式
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// 统一的成功响应格式
type SuccessResponse struct {
    Code int         `json:"code"`
    Data interface{} `json:"data"`
    Meta *Meta       `json:"meta,omitempty"`
}
```

#### 错误处理规范
```go
// 自定义错误类型
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return e.Message
}

// 错误包装
func wrapError(err error, message string) error {
    return &AppError{
        Code:    500,
        Message: message,
        Err:     err,
    }
}
```

### 3.2 TypeScript前端编码规范

#### 文件组织结构
```
src/
├── components/
│   ├── common/              # 通用组件
│   ├── layout/              # 布局组件
│   └── business/            # 业务组件
├── pages/
│   ├── Project/             # 项目页面
│   ├── Chat/                # 对话页面
│   └── Editor/              # 编辑器页面
├── hooks/
│   └── *.ts                 # 自定义Hooks
├── services/
│   └── api.ts               # API调用服务
├── store/
│   └── *.ts                 # 状态管理
├── types/
│   └── *.ts                 # 类型定义
└── utils/
    └── *.ts                 # 工具函数
```

#### 类型定义规范
```typescript
// API响应类型
interface ApiResponse<T> {
  code: number;
  data: T;
  message?: string;
}

// 业务实体类型
interface Project {
  id: string;
  name: string;
  description: string;
  status: ProjectStatus;
  createdAt: string;
  updatedAt: string;
}

// 组件Props类型
interface ProjectCardProps {
  project: Project;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}
```

## 4. 开发环境配置

### 4.1 Go开发环境

#### 必需工具
- Go 1.21+
- MySQL 8.0+
- Redis 7.0+
- Git

#### 环境变量配置
```bash
# .env文件
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ai_dev_platform
DB_USER=root
DB_PASSWORD=password

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

OPENAI_API_KEY=your_openai_key
CLAUDE_API_KEY=your_claude_key
AI_PROVIDER=openai  # openai | claude

JWT_SECRET=your_jwt_secret
CORS_ORIGINS=http://localhost:3000
```

### 4.2 前端开发环境

#### 必需工具
- Node.js 18+
- npm/yarn
- VS Code + 相关插件

#### 环境变量配置
```bash
# .env.local文件
VITE_API_BASE_URL=http://localhost:8080/api
VITE_PLANTUML_SERVER=http://localhost:8080/plantuml
```

## 5. 部署架构

### 5.1 开发环境部署
- **后端**：本地Go服务 (端口8080)
- **前端**：Vite开发服务器 (端口3000)
- **数据库**：本地MySQL实例
- **缓存**：本地Redis实例

### 5.2 生产环境部署
- **容器化**：Docker + Docker Compose
- **反向代理**：Nginx
- **数据库**：MySQL主从部署
- **缓存**：Redis集群
- **监控**：Prometheus + Grafana

## 6. 安全规范

### 6.1 API安全
- **身份认证**：JWT Token
- **权限控制**：基于角色的访问控制(RBAC)
- **输入验证**：严格的参数校验
- **SQL注入防护**：使用参数化查询
- **XSS防护**：输出编码

### 6.2 数据安全
- **敏感数据加密**：密码使用bcrypt哈希
- **数据传输加密**：HTTPS强制
- **API密钥保护**：环境变量存储
- **用户数据隔离**：租户级别数据隔离

## 7. 性能规范

### 7.1 后端性能
- **数据库连接池**：最大100个连接
- **Redis连接池**：最大50个连接
- **API响应时间**：< 200ms (普通接口)，< 2s (AI接口)
- **并发处理**：使用goroutine处理并发请求

### 7.2 前端性能
- **代码分割**：按页面懒加载
- **静态资源优化**：压缩、缓存
- **API调用优化**：防抖、缓存
- **首屏加载时间**：< 3s

## 8. 测试规范

### 8.1 后端测试
- **单元测试**：覆盖率 > 80%
- **集成测试**：API接口测试
- **压力测试**：并发用户测试

### 8.2 前端测试
- **组件测试**：React Testing Library
- **E2E测试**：Playwright
- **类型检查**：TypeScript严格模式 