# AI辅助项目开发平台 - API接口设计

## 1. API设计原则

### 1.1 RESTful设计
- 使用标准HTTP方法：GET、POST、PUT、DELETE
- 资源导向的URL设计
- 使用HTTP状态码表示操作结果
- 统一的响应格式

### 1.2 版本管理
- API版本前缀：`/api/v1/`
- 向后兼容原则
- 废弃API的渐进式迁移

### 1.3 安全设计
- JWT Token认证
- API密钥验证
- 请求频率限制
- 输入参数校验

## 2. 统一响应格式

### 2.1 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "v1"
  }
}
```

### 2.2 错误响应
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "INVALID_PARAMETER",
  "details": "project_name不能为空",
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "v1"
  }
}
```

### 2.3 分页响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

## 3. 用户认证接口

### 3.1 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "full_name": "测试用户"
}
```

**响应**：
```json
{
  "code": 201,
  "message": "用户注册成功",
  "data": {
    "user_id": "uuid-here",
    "username": "testuser",
    "email": "test@example.com",
    "full_name": "测试用户",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3.2 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**响应**：
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "jwt-token-here",
    "expires_in": 3600,
    "user": {
      "user_id": "uuid-here",
      "username": "testuser",
      "email": "test@example.com",
      "full_name": "测试用户"
    }
  }
}
```

### 3.3 Token刷新
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

## 4. 项目管理接口

### 4.1 创建项目
```http
POST /api/v1/projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "project_name": "图书管理系统",
  "description": "一个完整的图书管理系统",
  "project_type": "web_application"
}
```

**响应**：
```json
{
  "code": 201,
  "message": "项目创建成功",
  "data": {
    "project_id": "uuid-here",
    "project_name": "图书管理系统",
    "description": "一个完整的图书管理系统",
    "project_type": "web_application",
    "status": "draft",
    "completion_percentage": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 4.2 获取项目列表
```http
GET /api/v1/projects?page=1&page_size=20&status=draft
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "project_id": "uuid-here",
        "project_name": "图书管理系统",
        "description": "一个完整的图书管理系统",
        "status": "draft",
        "completion_percentage": 0,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

### 4.3 获取项目详情
```http
GET /api/v1/projects/{project_id}
Authorization: Bearer <token>
```

### 4.4 更新项目
```http
PUT /api/v1/projects/{project_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "project_name": "图书管理系统-更新版",
  "description": "一个更完整的图书管理系统"
}
```

### 4.5 删除项目
```http
DELETE /api/v1/projects/{project_id}
Authorization: Bearer <token>
```

## 5. 需求分析接口

### 5.1 提交需求分析
```http
POST /api/v1/projects/{project_id}/requirements/analyze
Authorization: Bearer <token>
Content-Type: application/json

{
  "raw_requirement": "我想开发一个图书管理系统，用户可以借书、还书，管理员可以管理图书信息"
}
```

**响应**：
```json
{
  "code": 200,
  "message": "需求分析完成",
  "data": {
    "requirement_id": "uuid-here",
    "structured_requirement": {
      "core_functions": ["借书", "还书", "图书管理"],
      "user_roles": ["普通用户", "管理员"],
      "business_entities": ["用户", "图书", "借阅记录"],
      "main_processes": ["借书流程", "还书流程", "图书管理流程"]
    },
    "completeness_score": 0.6,
    "analysis_status": "questions_generated",
    "missing_info_types": ["business_rules", "data_validation", "exception_handling"]
  }
}
```

### 5.2 获取补充问题
```http
GET /api/v1/requirements/{requirement_id}/questions
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "questions": [
      {
        "question_id": "uuid-here",
        "question_text": "用户借书是否有期限限制？逾期如何处理？",
        "question_category": "business_rules",
        "priority_level": 1,
        "answer_status": "pending"
      },
      {
        "question_id": "uuid-here-2",
        "question_text": "图书如何分类？是否需要搜索功能？",
        "question_category": "data_structure",
        "priority_level": 2,
        "answer_status": "pending"
      }
    ]
  }
}
```

### 5.3 回答补充问题
```http
POST /api/v1/requirements/{requirement_id}/questions/{question_id}/answer
Authorization: Bearer <token>
Content-Type: application/json

{
  "answer_text": "用户借书期限为30天，逾期每天收取0.5元罚金，超过7天自动续借一次"
}
```

## 6. 对话交互接口

### 6.1 创建对话会话
```http
POST /api/v1/projects/{project_id}/chat/sessions
Authorization: Bearer <token>
Content-Type: application/json

{
  "session_type": "requirement_analysis"
}
```

### 6.2 发送消息
```http
POST /api/v1/chat/sessions/{session_id}/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "message_content": "请详细描述借书流程的业务规则",
  "message_type": "question"
}
```

**响应**：
```json
{
  "code": 200,
  "message": "消息发送成功",
  "data": {
    "message_id": "uuid-here",
    "ai_response": {
      "message_id": "uuid-ai",
      "message_content": "借书流程的业务规则包括以下几个方面：\n1. 用户身份验证...",
      "message_type": "answer",
      "timestamp": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 6.3 获取对话历史
```http
GET /api/v1/chat/sessions/{session_id}/messages?page=1&page_size=50
Authorization: Bearer <token>
```

## 7. PUML图表接口

### 7.1 生成PUML图表
```http
POST /api/v1/projects/{project_id}/puml/generate
Authorization: Bearer <token>
Content-Type: application/json

{
  "diagram_types": ["business_flow", "architecture", "data_model"]
}
```

**响应**：
```json
{
  "code": 200,
  "message": "PUML图表生成成功",
  "data": {
    "diagrams": [
      {
        "diagram_id": "uuid-here",
        "diagram_type": "business_flow",
        "diagram_name": "图书借阅业务流程图",
        "puml_content": "@startuml\n...\n@enduml",
        "rendered_url": "/api/v1/puml/render/{diagram_id}",
        "version": 1,
        "is_validated": true
      }
    ]
  }
}
```

### 7.2 获取项目PUML图表
```http
GET /api/v1/projects/{project_id}/puml?type=business_flow
Authorization: Bearer <token>
```

### 7.3 更新PUML图表
```http
PUT /api/v1/puml/{diagram_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "puml_content": "@startuml\n更新后的PUML内容\n@enduml",
  "diagram_name": "更新后的图表名称"
}
```

### 7.4 验证PUML语法
```http
POST /api/v1/puml/validate
Authorization: Bearer <token>
Content-Type: application/json

{
  "puml_content": "@startuml\n...\n@enduml"
}
```

**响应**：
```json
{
  "code": 200,
  "message": "PUML语法验证成功",
  "data": {
    "is_valid": true,
    "validation_errors": [],
    "rendered_url": "/api/v1/puml/preview/{hash}"
  }
}
```

### 7.5 渲染PUML图表
```http
GET /api/v1/puml/render/{diagram_id}
Authorization: Bearer <token>
```

**响应**：返回PNG/SVG图片

## 8. 文档生成接口

### 8.1 生成项目文档
```http
POST /api/v1/projects/{project_id}/documents/generate
Authorization: Bearer <token>
Content-Type: application/json

{
  "document_types": ["requirement_doc", "tech_spec", "api_doc", "test_case"]
}
```

**响应**：
```json
{
  "code": 200,
  "message": "文档生成成功",
  "data": {
    "documents": [
      {
        "document_id": "uuid-here",
        "document_type": "requirement_doc",
        "document_name": "图书管理系统需求规格说明书",
        "format": "markdown",
        "file_path": "/documents/requirement_doc_uuid.md",
        "version": 1,
        "generated_at": "2024-01-01T00:00:00Z",
        "is_final": false
      }
    ]
  }
}
```

### 8.2 获取项目文档列表
```http
GET /api/v1/projects/{project_id}/documents?type=requirement_doc
Authorization: Bearer <token>
```

### 8.3 获取文档内容
```http
GET /api/v1/documents/{document_id}
Authorization: Bearer <token>
```

### 8.4 更新文档
```http
PUT /api/v1/documents/{document_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "更新后的文档内容",
  "document_name": "更新后的文档名称"
}
```

### 8.5 下载文档
```http
GET /api/v1/documents/{document_id}/download?format=pdf
Authorization: Bearer <token>
```

## 9. 模块管理接口

### 9.1 获取业务模块
```http
GET /api/v1/projects/{project_id}/modules
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "modules": [
      {
        "module_id": "uuid-here",
        "module_name": "用户认证模块",
        "description": "负责用户注册、登录、权限验证",
        "module_type": "authentication",
        "complexity_level": "medium",
        "is_reusable": true,
        "interfaces": {
          "input": ["用户信息", "认证请求"],
          "output": ["认证结果", "用户会话"]
        },
        "dependencies": ["数据库模块", "缓存模块"]
      }
    ]
  }
}
```

### 9.2 获取通用模块库
```http
GET /api/v1/modules/common?category=authentication&page=1&page_size=20
Authorization: Bearer <token>
```

### 9.3 搜索通用模块
```http
GET /api/v1/modules/common/search?q=用户认证&category=authentication
Authorization: Bearer <token>
```

### 9.4 使用通用模块
```http
POST /api/v1/projects/{project_id}/modules/use
Authorization: Bearer <token>
Content-Type: application/json

{
  "common_module_id": "uuid-here",
  "usage_type": "modified_use",
  "customization": {
    "auth_method": "jwt",
    "session_timeout": 3600
  }
}
```

## 10. 代码生成接口

### 10.1 生成项目代码
```http
POST /api/v1/projects/{project_id}/code/generate
Authorization: Bearer <token>
Content-Type: application/json

{
  "language": "go",
  "framework": "gin",
  "modules": ["authentication", "user_management", "book_management"]
}
```

### 10.2 获取生成的代码文件
```http
GET /api/v1/projects/{project_id}/code
Authorization: Bearer <token>
```

### 10.3 下载项目代码包
```http
GET /api/v1/projects/{project_id}/export
Authorization: Bearer <token>
```

**响应**：返回ZIP文件

## 11. 配置管理接口

### 11.1 获取项目配置
```http
GET /api/v1/projects/{project_id}/config
Authorization: Bearer <token>
```

### 11.2 更新项目配置
```http
PUT /api/v1/projects/{project_id}/config
Authorization: Bearer <token>
Content-Type: application/json

{
  "tech_stack": {
    "backend": "go",
    "frontend": "react",
    "database": "mysql"
  },
  "database_config": {
    "charset": "utf8mb4",
    "timezone": "Asia/Shanghai"
  }
}
```

## 12. 错误码定义

| 错误码 | HTTP状态码 | 错误信息 | 描述 |
|--------|------------|----------|------|
| 1001 | 400 | 参数错误 | 请求参数格式错误或缺失 |
| 1002 | 401 | 未授权 | Token无效或已过期 |
| 1003 | 403 | 无权限 | 没有访问该资源的权限 |
| 1004 | 404 | 资源不存在 | 请求的资源不存在 |
| 1005 | 409 | 资源冲突 | 资源已存在或状态冲突 |
| 1006 | 422 | 业务逻辑错误 | 业务规则验证失败 |
| 1007 | 429 | 请求过于频繁 | 超过API调用频率限制 |
| 1008 | 500 | 服务器内部错误 | 服务器处理请求时发生错误 |
| 1009 | 502 | 外部服务错误 | AI服务或其他外部服务不可用 |
| 1010 | 503 | 服务不可用 | 服务正在维护或临时不可用 |

## 13. API限流规则

- **认证接口**：每分钟最多5次请求
- **AI接口**：每分钟最多10次请求  
- **普通业务接口**：每分钟最多100次请求
- **文件上传接口**：每小时最多50次请求 