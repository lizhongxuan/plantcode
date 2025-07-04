# AI辅助业务分析与文档生成平台 - 设计文档

## 1. 项目概述

### 1.1 项目背景
传统的软件开发需求分析往往不够深入，容易遗漏重要信息，导致后期开发过程中频繁变更。本项目旨在构建一个AI辅助的业务分析平台，通过智能对话分析业务需求，自动生成开发所需的PUML图表和开发步骤文档。

### 1.2 核心价值
- **智能需求分析**：通过AI对话深度挖掘业务需求，识别缺失信息
- **自动文档生成**：基于分析结果生成标准化的PUML图表和开发文档
- **灵活修改完善**：支持用户直接编辑生成的内容，适应实际业务变化

## 2. 核心功能

### 2.1 对话式业务分析
- **自然语言输入**：用户用日常语言描述业务需求
- **智能信息提取**：AI识别核心功能、参与角色、主要流程、数据实体
- **缺失信息识别**：自动检测业务描述中的不完整部分
- **主动问题生成**：针对缺失信息生成具体的补充问题

### 2.2 PUML图表生成
- **业务流程图**：基于业务逻辑生成活动图或序列图
- **系统架构图**：基于功能需求生成组件架构图
- **实时预览**：支持PUML代码的实时编辑和预览

### 2.3 开发文档生成
- **功能模块清单**：详细的功能分解和描述
- **开发优先级**：基于业务重要性和技术依赖的开发顺序
- **技术选型建议**：根据项目特点推荐合适的技术栈
- **数据库设计**：基于数据实体生成数据表结构
- **API接口设计**：定义系统间的接口规范

## 3. 系统架构

### 3.1 整体架构
系统采用简洁的分层架构：
- **前端界面层**：Web界面 + PUML编辑器
- **核心服务层**：对话分析 + 业务分析 + 文档生成
- **AI服务**：大语言模型
- **数据存储**：MySQL数据库

### 3.2 核心组件

#### 3.2.1 对话分析服务
- **功能**：处理用户输入，管理对话上下文
- **职责**：
  - 接收用户业务描述
  - 维护对话历史
  - 传递信息给业务分析引擎

#### 3.2.2 业务分析引擎
- **功能**：核心的业务逻辑分析组件
- **职责**：
  - 分析业务需求，提取关键信息
  - 识别缺失信息，生成补充问题
  - 协调PUML和文档生成

#### 3.2.3 PUML生成器
- **功能**：生成业务相关的PUML图表
- **输出类型**：
  - 业务流程图（活动图/序列图）
  - 系统架构图（组件图）
- **特性**：支持语法验证和实时预览

#### 3.2.4 文档生成器
- **功能**：生成结构化的开发步骤文档
- **内容**：
  - 功能清单和开发计划
  - 技术选型和架构建议
  - 数据库和API设计

## 4. 技术选型

### 4.1 后端技术栈
- **编程语言**：Go
- **Web框架**：标准库 net/http
- **数据库**：MySQL
- **AI服务**：OpenAI GPT-4

### 4.2 前端技术栈
- **框架**：React + TypeScript
- **UI库**：Ant Design
- **代码编辑器**：Monaco Editor
- **图表渲染**：PlantUML渲染器

## 5. 数据模型

### 5.1 核心实体

#### 5.1.1 项目 (Project)
```sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id INTEGER,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 5.1.2 对话记录 (Conversation)
```sql
CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    user_input TEXT NOT NULL,
    ai_response TEXT,
    analysis_result JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 5.1.3 生成结果 (Generated_Content)
```sql
CREATE TABLE generated_content (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    content_type VARCHAR(50), -- 'business_flow', 'architecture', 'dev_doc'
    puml_content TEXT,
    document_content TEXT,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 6. API接口设计

### 6.1 项目管理
- `POST /api/projects` - 创建项目
- `GET /api/projects` - 获取项目列表
- `GET /api/projects/{id}` - 获取项目详情

### 6.2 对话交互
- `POST /api/conversations` - 发送用户输入
- `GET /api/conversations/{project_id}` - 获取对话历史

### 6.3 文档生成
- `POST /api/generate/analysis` - 触发业务分析
- `POST /api/generate/puml` - 生成PUML图表
- `POST /api/generate/document` - 生成开发文档

### 6.4 内容编辑
- `PUT /api/content/{id}` - 更新生成的内容
- `GET /api/content/{project_id}` - 获取项目所有内容

## 7. 开发计划

### 7.1 第一阶段：核心对话功能
- [ ] 搭建基础框架（Go后端 + React前端）
- [ ] 实现对话分析服务
- [ ] 集成OpenAI API
- [ ] 基础的业务信息提取

### 7.2 第二阶段：文档生成功能
- [ ] 实现业务分析引擎
- [ ] 开发PUML生成器
- [ ] 开发文档生成器
- [ ] 数据库设计和实现

### 7.3 第三阶段：编辑和完善功能
- [ ] 实现PUML编辑器
- [ ] 添加实时预览功能
- [ ] 实现文档编辑功能
- [ ] 完善用户界面

### 7.4 第四阶段：优化和部署
- [ ] 性能优化
- [ ] 错误处理完善
- [ ] 用户体验优化
- [ ] 部署和运维

## 8. 使用流程示例

### 8.1 典型使用场景
**用户输入**：
> "我想做一个在线图书管理系统，用户可以借书还书，管理员可以管理图书信息"

**AI分析结果**：
- 识别角色：用户、管理员
- 核心功能：借书、还书、图书管理
- 缺失信息：借书期限、逾期处理、图书分类等

**AI生成问题**：
1. 用户借书有期限限制吗？逾期如何处理？
2. 图书如何分类？是否需要搜索功能？
3. 是否需要预约功能？
4. 管理员权限包括哪些具体操作？

**用户补充回答后，AI生成**：
- 业务流程图PUML（借书还书流程）
- 系统架构图PUML（用户模块、图书模块、借阅模块等）
- 开发步骤文档（功能清单、数据库设计、API设计等）

## 9. 项目优势

### 9.1 简化开发流程
- 减少需求分析时间
- 自动生成标准化文档
- 降低沟通成本

### 9.2 提高文档质量
- AI确保信息完整性
- 标准化的PUML图表
- 结构化的开发文档

### 9.3 灵活性强
- 支持实时修改
- 适应业务变化
- 可重复使用 