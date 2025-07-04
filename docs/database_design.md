# AI辅助项目开发平台 - 数据库设计

## 1. 数据库设计原则

### 1.1 设计原则
- **规范化设计**：遵循第三范式，避免数据冗余
- **性能优化**：合理设计索引，优化查询性能
- **可扩展性**：预留扩展字段，支持业务增长
- **数据完整性**：使用外键约束保证数据一致性
- **安全性**：敏感数据加密存储

### 1.2 命名规范
- **表名**：小写+下划线，复数形式，如 `users`、`projects`
- **字段名**：小写+下划线，描述性命名，如 `user_id`、`created_at`
- **索引名**：`idx_表名_字段名`，如 `idx_users_email`
- **外键名**：`fk_表名_引用表名`，如 `fk_projects_users`

### 1.3 字段设计规范
- **主键**：使用UUID，保证全局唯一
- **时间戳**：统一使用 `TIMESTAMP` 类型，UTC时区
- **枚举值**：使用 `VARCHAR` 存储，便于扩展
- **JSON字段**：存储非结构化数据，提高灵活性
- **软删除**：重要数据使用 `deleted_at` 字段标记删除

## 2. 数据库选型

### 2.1 MySQL 8.0+
- **主数据库**：存储结构化业务数据
- **特性**：ACID事务、JSON字段支持、全文索引
- **配置**：UTF8MB4字符集、InnoDB存储引擎

### 2.2 Redis 7.0+
- **缓存数据库**：存储会话、缓存、临时数据
- **用途**：
  - 用户会话管理
  - API响应缓存
  - AI生成结果临时存储
  - 实时对话上下文

## 3. 核心表结构设计

### 3.1 用户表 (users)

```sql
CREATE TABLE users (
    user_id CHAR(36) PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    full_name VARCHAR(100) NOT NULL COMMENT '全名',
    avatar_url VARCHAR(500) COMMENT '头像URL',
    status VARCHAR(20) DEFAULT 'active' COMMENT '用户状态',
    preferences JSON COMMENT '用户偏好设置',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    
    INDEX idx_users_email (email),
    INDEX idx_users_username (username),
    INDEX idx_users_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

### 3.2 项目表 (projects)

```sql
CREATE TABLE projects (
    project_id CHAR(36) PRIMARY KEY COMMENT '项目ID',
    user_id CHAR(36) NOT NULL COMMENT '创建者ID',
    project_name VARCHAR(100) NOT NULL COMMENT '项目名称',
    description TEXT COMMENT '项目描述',
    project_type VARCHAR(50) NOT NULL COMMENT '项目类型',
    status VARCHAR(20) DEFAULT 'draft' COMMENT '项目状态',
    completion_percentage TINYINT UNSIGNED DEFAULT 0 COMMENT '完成百分比',
    settings JSON COMMENT '项目设置',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    
    FOREIGN KEY fk_projects_users (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    INDEX idx_projects_user_id (user_id),
    INDEX idx_projects_status (status),
    FULLTEXT idx_projects_search (project_name, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';
```

### 3.3 需求分析表 (requirement_analyses)

```sql
CREATE TABLE requirement_analyses (
    requirement_id CHAR(36) PRIMARY KEY COMMENT '需求分析ID',
    project_id CHAR(36) NOT NULL COMMENT '项目ID',
    raw_requirement TEXT NOT NULL COMMENT '原始需求描述',
    structured_requirement JSON COMMENT '结构化需求',
    completeness_score DECIMAL(3,2) DEFAULT 0.00 COMMENT '完整性评分',
    analysis_status VARCHAR(50) DEFAULT 'pending' COMMENT '分析状态',
    missing_info_types JSON COMMENT '缺失信息类型',
    ai_model_used VARCHAR(50) COMMENT '使用的AI模型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY fk_requirements_projects (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    INDEX idx_requirements_project_id (project_id),
    INDEX idx_requirements_status (analysis_status),
    FULLTEXT idx_requirements_content (raw_requirement)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='需求分析表';
```

### 3.4 对话会话表 (chat_sessions)

```sql
CREATE TABLE chat_sessions (
    session_id CHAR(36) PRIMARY KEY COMMENT '会话ID',
    project_id CHAR(36) NOT NULL COMMENT '项目ID',
    user_id CHAR(36) NOT NULL COMMENT '用户ID',
    session_type VARCHAR(50) DEFAULT 'general' COMMENT '会话类型',
    status VARCHAR(20) DEFAULT 'active' COMMENT '会话状态',
    context JSON COMMENT '会话上下文',
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    ended_at TIMESTAMP NULL COMMENT '结束时间',
    message_count INT DEFAULT 0 COMMENT '消息数量',
    
    FOREIGN KEY fk_sessions_projects (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    FOREIGN KEY fk_sessions_users (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    INDEX idx_sessions_project_id (project_id),
    INDEX idx_sessions_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话会话表';
```

### 3.5 对话消息表 (chat_messages)

```sql
CREATE TABLE chat_messages (
    message_id CHAR(36) PRIMARY KEY COMMENT '消息ID',
    session_id CHAR(36) NOT NULL COMMENT '会话ID',
    sender_type VARCHAR(20) NOT NULL COMMENT '发送者类型',
    message_content TEXT NOT NULL COMMENT '消息内容',
    message_type VARCHAR(50) DEFAULT 'question' COMMENT '消息类型',
    metadata JSON COMMENT '消息元数据',
    ai_model_used VARCHAR(50) COMMENT '使用的AI模型',
    processing_time INT COMMENT '处理时间(毫秒)',
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '时间戳',
    
    FOREIGN KEY fk_messages_sessions (session_id) REFERENCES chat_sessions(session_id) ON DELETE CASCADE,
    INDEX idx_messages_session_id (session_id),
    INDEX idx_messages_timestamp (timestamp),
    FULLTEXT idx_messages_content (message_content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话消息表';
```

### 3.6 PUML图表表 (puml_diagrams)

```sql
CREATE TABLE puml_diagrams (
    diagram_id CHAR(36) PRIMARY KEY COMMENT '图表ID',
    project_id CHAR(36) NOT NULL COMMENT '项目ID',
    diagram_type VARCHAR(50) NOT NULL COMMENT '图表类型',
    diagram_name VARCHAR(100) NOT NULL COMMENT '图表名称',
    puml_content LONGTEXT NOT NULL COMMENT 'PUML内容',
    rendered_url VARCHAR(500) COMMENT '渲染后的图片URL',
    version INT DEFAULT 1 COMMENT '版本号',
    is_validated BOOLEAN DEFAULT FALSE COMMENT '是否已验证',
    validation_feedback TEXT COMMENT '验证反馈',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY fk_diagrams_projects (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    INDEX idx_diagrams_project_id (project_id),
    INDEX idx_diagrams_type (diagram_type),
    FULLTEXT idx_diagrams_name (diagram_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='PUML图表表';
```

### 3.7 业务模块表 (business_modules)

```sql
CREATE TABLE business_modules (
    module_id CHAR(36) PRIMARY KEY COMMENT '模块ID',
    project_id CHAR(36) NOT NULL COMMENT '项目ID',
    module_name VARCHAR(100) NOT NULL COMMENT '模块名称',
    description TEXT COMMENT '模块描述',
    module_type VARCHAR(50) NOT NULL COMMENT '模块类型',
    complexity_level VARCHAR(20) DEFAULT 'medium' COMMENT '复杂度级别',
    business_logic JSON COMMENT '业务逻辑',
    interfaces JSON COMMENT '接口定义',
    dependencies JSON COMMENT '依赖关系',
    is_reusable BOOLEAN DEFAULT FALSE COMMENT '是否可复用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY fk_modules_projects (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    INDEX idx_modules_project_id (project_id),
    INDEX idx_modules_type (module_type),
    FULLTEXT idx_modules_search (module_name, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='业务模块表';
```

### 3.8 通用模块库表 (common_module_library)

```sql
CREATE TABLE common_module_library (
    common_module_id CHAR(36) PRIMARY KEY COMMENT '通用模块ID',
    module_name VARCHAR(100) NOT NULL COMMENT '模块名称',
    category VARCHAR(50) NOT NULL COMMENT '模块分类',
    description TEXT COMMENT '模块描述',
    functionality JSON COMMENT '功能说明',
    interface_spec JSON COMMENT '接口规范',
    code_template LONGTEXT COMMENT '代码模板',
    usage_examples JSON COMMENT '使用示例',
    version VARCHAR(20) DEFAULT '1.0.0' COMMENT '版本号',
    downloads_count INT DEFAULT 0 COMMENT '下载次数',
    rating DECIMAL(2,1) DEFAULT 0.0 COMMENT '评分',
    tags JSON COMMENT '标签',
    created_by CHAR(36) COMMENT '创建者ID',
    status VARCHAR(20) DEFAULT 'draft' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY fk_common_modules_users (created_by) REFERENCES users(user_id) ON SET NULL,
    INDEX idx_common_modules_category (category),
    INDEX idx_common_modules_status (status),
    FULLTEXT idx_common_modules_search (module_name, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通用模块库表';
```

### 3.9 生成文档表 (generated_documents)

```sql
CREATE TABLE generated_documents (
    document_id CHAR(36) PRIMARY KEY COMMENT '文档ID',
    project_id CHAR(36) NOT NULL COMMENT '项目ID',
    document_type VARCHAR(50) NOT NULL COMMENT '文档类型',
    document_name VARCHAR(100) NOT NULL COMMENT '文档名称',
    content LONGTEXT NOT NULL COMMENT '文档内容',
    format VARCHAR(20) DEFAULT 'markdown' COMMENT '文档格式',
    file_path VARCHAR(500) COMMENT '文件路径',
    version INT DEFAULT 1 COMMENT '版本号',
    is_final BOOLEAN DEFAULT FALSE COMMENT '是否为最终版本',
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '生成时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY fk_documents_projects (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    INDEX idx_documents_project_id (project_id),
    INDEX idx_documents_type (document_type),
    FULLTEXT idx_documents_search (document_name, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='生成文档表';
```

## 4. Redis缓存设计

### 4.1 用户会话缓存
```
KEY: session:{user_id}
TYPE: Hash
TTL: 7200 seconds (2 hours)
FIELDS:
  - token: JWT token
  - login_time: 登录时间
  - user_info: 用户基本信息JSON
```

### 4.2 AI结果缓存
```
KEY: ai_result:{request_hash}
TYPE: Hash
TTL: 3600 seconds (1 hour)
FIELDS:
  - result: AI生成结果
  - model_used: 使用的AI模型
  - processing_time: 处理时间
```

### 4.3 对话上下文缓存
```
KEY: chat_context:{session_id}
TYPE: List
TTL: 1800 seconds (30 minutes)
VALUE: 对话消息列表（最近50条）
```

### 4.4 PUML渲染缓存
```
KEY: puml_render:{content_hash}
TYPE: String
TTL: 86400 seconds (24 hours)
VALUE: 渲染后的图片URL
```

## 5. 索引优化策略

### 5.1 主要查询索引
```sql
-- 用户登录查询
CREATE INDEX idx_users_login ON users(email, password_hash);

-- 项目列表查询
CREATE INDEX idx_projects_user_status ON projects(user_id, status, created_at);

-- 对话消息查询
CREATE INDEX idx_messages_session_time ON chat_messages(session_id, timestamp);

-- 模块搜索查询
CREATE INDEX idx_modules_search ON business_modules(project_id, module_type, is_reusable);
```

### 5.2 复合索引设计
```sql
-- 项目管理复合索引
CREATE INDEX idx_projects_composite ON projects(user_id, status, project_type, created_at);

-- 消息查询复合索引  
CREATE INDEX idx_messages_composite ON chat_messages(session_id, sender_type, timestamp);
```

## 6. 数据初始化

### 6.1 基础数据
```sql
-- 插入系统管理员用户
INSERT INTO users (user_id, username, email, password_hash, full_name, status) VALUES
('admin-uuid-001', 'admin', 'admin@example.com', '$2a$10$...', '系统管理员', 'active');

-- 插入基础通用模块
INSERT INTO common_module_library (common_module_id, module_name, category, description) VALUES
('auth-module-001', '用户认证模块', 'authentication', 'JWT based authentication module'),
('user-mgmt-001', '用户管理模块', 'user_management', 'Complete user management with CRUD operations'),
('file-upload-001', '文件上传模块', 'file_management', 'Secure file upload with validation');
```

## 7. 备份与恢复策略

### 7.1 备份策略
- **全量备份**：每日凌晨进行全量备份
- **增量备份**：每小时进行增量备份  
- **实时同步**：主从复制保证数据一致性
- **异地备份**：备份数据异地存储

### 7.2 恢复策略
- **故障恢复**：从最近的全量备份+增量备份恢复
- **数据验证**：恢复后进行数据完整性检查
- **回滚机制**：支持快速回滚到指定时间点

## 8. 性能监控与优化

### 8.1 性能监控
- **慢查询监控**：记录执行时间超过1秒的查询
- **连接数监控**：监控数据库连接池使用情况
- **索引效率监控**：定期分析索引使用情况
- **表空间监控**：监控表空间使用率

### 8.2 优化建议
- **定期清理**：删除过期的临时数据和日志
- **索引重建**：定期重建碎片化的索引
- **分区表**：对大表进行分区存储
- **读写分离**：读操作使用从库，写操作使用主库 