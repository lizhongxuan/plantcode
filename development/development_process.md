# AI辅助项目开发平台 - 开发流程文档

## 1. 项目开发概述

### 1.1 开发目标
- 构建一个完整的AI辅助项目开发平台
- 实现需求分析、PUML建模、文档生成、代码生成的完整流程
- 提供可视化的项目管理和模块化复用能力

### 1.2 开发原则
- **敏捷开发**：快速迭代，持续集成
- **模块化设计**：松耦合，高内聚
- **测试驱动**：先写测试，后写实现
- **文档同步**：代码与文档同步更新
- **用户体验优先**：界面友好，操作流畅

## 2. 开发环境搭建

### 2.1 基础环境要求
```bash
# Go开发环境
go version go1.21+ 

# Node.js环境
node --version  # 18+
npm --version   # 9+

# 数据库环境
mysql --version  # 8.0+
redis-server --version  # 7.0+

# 开发工具
git --version
docker --version
```

### 2.2 项目初始化
```bash
# 1. 克隆项目
git clone https://github.com/your-org/ai-dev-platform.git
cd ai-dev-platform

# 2. 创建项目目录结构
mkdir -p {cmd/server,internal/{api,service,repository,ai,model,config,utils},web/{src,public}}
mkdir -p {scripts,docs,tests,deployments}

# 3. 初始化Go模块
go mod init ai-dev-platform

# 4. 初始化前端项目
cd web
npm create vite@latest . -- --template react-ts
npm install

# 5. 安装Go依赖
cd ..
go get github.com/go-sql-driver/mysql
go get github.com/redis/go-redis/v9
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
```

### 2.3 环境配置
```bash
# 复制环境配置文件
cp .env.example .env
cp web/.env.example web/.env.local

# 配置数据库
mysql -u root -p < scripts/init_database.sql

# 启动Redis
redis-server --daemonize yes
```

## 3. 详细开发计划

### 3.1 第一阶段：核心基础功能 (2-3周)

#### 里程碑1：基础架构搭建 (第1周)
**优先级**：🔴 高

**后端任务**：
- [√] 搭建Go Web服务器基础架构
- [√] 实现路由系统和中间件
- [√] 配置数据库连接池
- [√] 实现用户认证系统 (JWT)
- [√] 设计基础API响应格式

**前端任务**：
- [√] 搭建React + TypeScript项目
- [√] 配置路由系统 (React Router)
- [√] 设计基础UI组件库
- [√] 实现登录/注册页面
- [√] 配置状态管理 (Zustand)

#### 里程碑2：用户管理模块 (第2周)
**优先级**：🔴 高

**开发任务**：
- [√] 用户注册/登录API
- [√] 用户信息管理
- [√] 权限控制中间件
- [√] 用户偏好设置
- [√] 前端用户界面

**API接口实现**：
```go
// internal/api/handlers/auth.go
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. 参数验证
    // 2. 密码加密
    // 3. 创建用户
    // 4. 返回JWT Token
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    // 1. 验证用户名密码
    // 2. 生成JWT Token
    // 3. 更新登录时间
    // 4. 返回用户信息
}
```

#### 里程碑3：项目管理核心功能 (第3周)
**优先级**：🔴 高

**开发任务**：
- [√] 项目CRUD操作
- [√] 项目状态管理
- [√] 项目配置管理
- [√] 项目列表和搜索
- [√] 前端项目管理界面

### 3.2 第二阶段：AI集成和需求分析 (3-4周)

#### 里程碑4：AI服务集成 (第4周)
**优先级**：🟡 中

**开发任务**：
- [√] AI客户端抽象层设计
- [√] OpenAI API集成
- [√] Claude API集成
- [√] 多AI服务切换机制
- [√] AI响应缓存系统

**AI服务实现**：
```go
// internal/ai/client.go
type AIClient interface {
    AnalyzeRequirement(ctx context.Context, req string) (*AnalysisResult, error)
    GeneratePUML(ctx context.Context, analysis *AnalysisResult) (*PUMLDiagram, error)
    GenerateQuestions(ctx context.Context, gaps []string) ([]Question, error)
}

type OpenAIClient struct {
    apiKey string
    client *http.Client
}

type ClaudeClient struct {
    apiKey string
    client *http.Client
}
```

#### 里程碑5：需求分析模块 (第5周)
**优先级**：🟡 中

**开发任务**：
- [√] 原始需求解析
- [√] 结构化需求提取
- [√] 缺失信息识别
- [√] 补充问题生成
- [√] 需求完整性评分

#### 里程碑6：对话交互系统 (第6-7周)
**优先级**：🟡 中

**开发任务**：
- [ ] 实时对话接口
- [ ] 对话上下文管理
- [ ] 多轮对话逻辑
- [ ] 对话历史存储
- [ ] 前端聊天界面

### 3.3 第三阶段：PUML生成和可视化 (2-3周) ✅ **已完成**

#### 里程碑7：PUML图表生成 (第8周) ✅ **已完成**
**优先级**：🟡 中

**开发任务**：
- [x] 业务流程图生成
- [x] 系统架构图生成  
- [x] 数据模型图生成
- [x] 交互流程图生成
- [x] PUML语法验证

**完成详情**：
- ✅ 在AI客户端中添加了数据模型图(data_model)类型支持
- ✅ 完善了OpenAI客户端的PUML生成逻辑，支持5种图表类型
- ✅ 实现了完整的PUML语法验证功能，包括基本语法检查、括号匹配、标记验证等

#### 里程碑8：图表渲染和编辑 (第9周) ✅ **已完成**  
**优先级**：🟡 中

**开发任务**：
- [x] PlantUML渲染服务
- [x] 在线PUML编辑器
- [x] 图表版本管理
- [x] 图表导出功能
- [x] 图表预览功能

**完成详情**：
- ✅ 创建了完整的PUMLService，支持在线PlantUML渲染
- ✅ 实现了PUML代码的zlib压缩和base64编码
- ✅ 增强前端PUMLDiagrams组件，添加了预览、编辑、导出功能
- ✅ 支持多种格式导出：PUML源码、PNG图片、SVG矢量图
- ✅ 实现了实时语法验证和错误提示
- ✅ 添加了图表缓存机制，提升渲染性能

**第三阶段API接口**：
- POST `/api/puml/render` - PUML代码渲染
- POST `/api/puml/validate` - PUML语法验证
- POST `/api/puml/preview` - 图表预览
- POST `/api/puml/export` - 图表导出
- GET `/api/puml/stats` - 服务统计
- POST `/api/puml/cache/clear` - 清空缓存

### 3.4 第四阶段：模块管理 (3-4周)

#### 里程碑9：业务模块管理 (第10周)
**优先级**：🟢 低

**开发任务**：
- [ ] 业务模块识别
- [ ] 模块依赖分析
- [ ] 通用模块库
- [ ] 模块搜索和推荐

#### 里程碑10：文档生成系统 (第11周)
**优先级**：🟢 低

**开发任务**：
- [ ] 需求文档生成
- [ ] 技术规范生成
- [ ] API文档生成
- [ ] 测试用例生成
- [ ] 文档模板系统

## 5. 开发最佳实践

### 5.1 后端开发规范

#### 错误处理
```go
// 统一错误处理
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *AppError) Error() string {
    return e.Message
}

// 错误处理中间件
func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // 记录错误日志
                log.Printf("Panic: %v", err)
                // 返回500错误
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

#### 日志记录
```go
import "log/slog"

// 结构化日志
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

logger.Info("User login", 
    slog.String("user_id", userID),
    slog.String("ip", clientIP),
    slog.Duration("duration", time.Since(start)))
```

#### 数据库操作
```go
// 使用事务
func (r *ProjectRepository) CreateProjectWithModules(ctx context.Context, 
    project *Project, modules []Module) error {
    
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 创建项目
    if err := r.createProject(ctx, tx, project); err != nil {
        return err
    }

    // 创建模块
    for _, module := range modules {
        if err := r.createModule(ctx, tx, &module); err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

### 5.2 前端开发规范

#### 组件设计
```typescript
// 组件Props类型定义
interface ProjectCardProps {
  project: Project;
  onEdit?: (id: string) => void;
  onDelete?: (id: string) => void;
  className?: string;
}

// 使用React.memo优化性能
export const ProjectCard = React.memo<ProjectCardProps>(({
  project,
  onEdit,
  onDelete,
  className
}) => {
  const handleEdit = useCallback(() => {
    onEdit?.(project.id);
  }, [project.id, onEdit]);

  return (
    <div className={`project-card ${className}`}>
      {/* 组件内容 */}
    </div>
  );
});
```

#### 状态管理
```typescript
// Zustand状态管理
interface ProjectState {
  projects: Project[];
  loading: boolean;
  error: string | null;
  
  fetchProjects: () => Promise<void>;
  createProject: (project: CreateProjectRequest) => Promise<void>;
  updateProject: (id: string, updates: Partial<Project>) => Promise<void>;
  deleteProject: (id: string) => Promise<void>;
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  projects: [],
  loading: false,
  error: null,

  fetchProjects: async () => {
    set({ loading: true });
    try {
      const projects = await projectApi.getProjects();
      set({ projects, loading: false });
    } catch (error) {
      set({ error: error.message, loading: false });
    }
  },
  // ... 其他方法
}));
```

## 6. 测试策略

### 6.1 后端测试
```go
// 单元测试示例
func TestUserService_CreateUser(t *testing.T) {
    // 准备测试数据
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    user := &User{
        Username: "testuser",
        Email: "test@example.com",
        Password: "password123",
    }

    // 执行测试
    result, err := service.CreateUser(context.Background(), user)

    // 验证结果
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ID)
    assert.Equal(t, user.Username, result.Username)
}

// API集成测试
func TestProjectAPI(t *testing.T) {
    // 设置测试服务器
    server := setupTestServer()
    defer server.Close()

    // 测试创建项目
    resp, err := http.Post(server.URL+"/api/v1/projects", 
        "application/json", 
        strings.NewReader(`{"name":"test project"}`))
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

### 6.2 前端测试
```typescript
// 组件测试
import { render, screen, fireEvent } from '@testing-library/react';
import { ProjectCard } from './ProjectCard';

describe('ProjectCard', () => {
  const mockProject = {
    id: '1',
    name: 'Test Project',
    description: 'Test Description',
    status: 'draft'
  };

  it('renders project information', () => {
    render(<ProjectCard project={mockProject} />);
    
    expect(screen.getByText('Test Project')).toBeInTheDocument();
    expect(screen.getByText('Test Description')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', () => {
    const mockOnEdit = jest.fn();
    render(<ProjectCard project={mockProject} onEdit={mockOnEdit} />);
    
    fireEvent.click(screen.getByRole('button', { name: /edit/i }));
    expect(mockOnEdit).toHaveBeenCalledWith('1');
  });
});
```

## 7. 性能优化策略

### 7.1 后端性能优化
- **数据库优化**：合理使用索引，避免N+1查询
- **缓存策略**：Redis缓存热点数据
- **连接池管理**：数据库和Redis连接池配置
- **异步处理**：使用goroutine处理耗时操作
- **API限流**：防止API被滥用

### 7.2 前端性能优化
- **代码分割**：按路由拆分代码包
- **懒加载**：组件和图片懒加载
- **缓存策略**：API响应缓存
- **虚拟化**：大列表虚拟化渲染
- **图片优化**：图片压缩和格式优化

## 8. 部署准备

### 8.1 生产环境配置
```bash
# 环境变量配置
export GO_ENV=production
export DB_HOST=prod-mysql-host
export REDIS_HOST=prod-redis-host
export JWT_SECRET=super-secret-key
export API_RATE_LIMIT=1000
```

### 8.2 Docker容器化
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 9. 项目交付标准

### 9.1 功能交付标准
- [ ] 所有核心功能正常工作
- [ ] API接口文档完整
- [ ] 前端界面用户友好
- [ ] 数据库设计合理
- [ ] 系统性能满足要求

### 9.2 质量交付标准
- [ ] 代码测试覆盖率 > 80%
- [ ] 没有严重安全漏洞
- [ ] 代码审查通过
- [ ] 文档完整准确
- [ ] 部署流程验证通过

### 9.3 维护交付标准
- [ ] 监控系统配置完成
- [ ] 日志系统正常工作
- [ ] 备份恢复流程验证
- [ ] 运维文档完整
- [ ] 团队技术交接完成 