# AI开发平台

## 项目概述

AI开发平台是一个智能的软件开发辅助工具，帮助开发者进行需求分析、系统设计、文档生成和开发指导。

## 功能特性

### 核心功能
- 🧠 **智能需求分析**：AI自动分析用户需求，提取核心功能和业务流程
- 📊 **PUML图表生成**：自动生成系统架构图、流程图等技术文档
- 📝 **分阶段文档生成**：按开发阶段自动生成规范化的技术文档
- 💬 **项目上下文AI对话**：基于项目信息的智能AI助手
- ⚙️ **用户自定义AI配置**：支持OpenAI、Claude、Gemini等多种AI提供商

## 项目解决问题
当遇到大项目,AI会面临以下几个问题：
1. 注意力稀释： 在极长的上下文中，AI很难判断哪些是当前任务最核心的信息，早期的、可能已经过时的指令会干扰它对当前需求的理解。
2. “灾难性遗忘”： 为了处理新的指令，模型可能会“忘记”或忽略掉一些早期的关键细节或约束。
3. 错误累积： 如果AI在早期生成了一个小错误，而您没有发现，在后续的开发中，AI会基于这个错误的认知继续构建，导致问题像滚雪球一样越滚越大。
4. token限制： 即使是拥有超长上下文窗口的模型，也存在物理上限。超出限制后，最早的信息会被截断，导致信息丢失。

要解决这个问题，核心思想是**“化整为零，分而治之” (Divide and Conquer)**。您需要从一个“项目经理”的角色，转变为一个“技术总监”或“架构师”的角色，将与AI的协作模式进行优化。

不要把AI看作一个无所不能、一次性就能完成任务的黑箱，而是把它看作一个能力极强但需要明确指导的初级程序员。

推荐流程： 你像一个建筑师，首先画好了大楼的设计图（规划与分解），然后把任务分包出去。你对负责打地基的团队说：“这是地基图纸，按这个做”（模块化开发）。你规定好水电管道的接口标准（定义接口），并对每一面墙进行质量检测（单元测试），最后才将所有部分组装起来（集成）。整个过程都用文档和照片记录在案（版本控制）。


### AI配置功能 ✨ **NEW**
系统已升级支持用户个性化AI配置：
- **多AI提供商支持**：OpenAI、Claude (Anthropic)、Gemini (Google)
- **个人API密钥配置**：安全存储用户的API密钥
- **模型选择**：支持不同AI模型的选择和配置
- **全功能适配**：所有AI功能都使用用户的个人配置

## 技术栈

### 后端
- **语言**：Go 1.23
- **框架**：标准库 net/http
- **数据库**：MySQL 8.0
- **AI集成**：OpenAI API、Google Gemini API、Claude API

### 前端
- **框架**：React 18 + TypeScript
- **构建工具**：Vite
- **UI组件**：Ant Design
- **状态管理**：Zustand
- **样式**：Tailwind CSS

## 快速开始

### 环境要求
- Go 1.23+
- Node.js 18+
- MySQL 8.0+

### 启动步骤

1. **启动后端服务**
```bash
go run cmd/server/main.go
```

2. **启动前端服务**
```bash
cd web
npm install
npm run dev
```

3. **访问应用**
- 前端地址：http://localhost:3000 (或自动分配的端口)
- 后端API：http://localhost:8080

### AI配置指南

#### 设置个人AI配置
1. 登录系统后，进入「设置」页面
2. 选择「AI配置」选项卡
3. 选择您偏好的AI提供商
4. 输入对应的API密钥
5. 选择默认模型和参数
6. 保存配置

#### 支持的AI提供商

| 提供商 | 支持模型 | 特点 |
|--------|----------|------|
| **OpenAI** | GPT-4, GPT-4 Turbo, GPT-3.5 Turbo | 强大的通用AI模型，适用于各种开发任务 |
| **Claude** | Claude-3 Opus, Sonnet, Haiku | 专注于代码生成和分析的AI助手 |
| **Gemini** | Gemini-2.5 Pro, 1.5 Pro, 1.5 Flash | Google的多模态AI模型，支持文本和图像处理 |

#### 获取API密钥
- **OpenAI**：访问 [OpenAI API Keys](https://platform.openai.com/api-keys)
- **Claude**：访问 [Anthropic Console](https://console.anthropic.com/)
- **Gemini**：访问 [Google AI Studio](https://makersuite.google.com/app/apikey)

## 主要功能说明

### 1. 需求分析
- 输入自然语言需求描述
- AI自动提取核心功能、角色、业务流程
- 生成结构化的需求文档
- 识别缺失信息并生成补充问题

### 2. PUML图表生成
- 支持多种图表类型：序列图、类图、活动图、组件图等
- 在线编辑器支持实时预览
- 拖拽功能和语法验证
- 版本管理和下载功能

### 3. 分阶段文档生成
- **阶段一**：需求梳理 - 项目需求文档、系统架构图、业务流程图
- **阶段二**：技术设计 - 技术规范、API设计、数据库设计
- **阶段三**：开发指导 - 开发流程、测试用例、部署文档

### 4. AI助手对话
- 基于项目上下文的智能对话
- 支持文档生成、PUML编辑、需求分析等场景
- 实时响应和动态操作建议

## API文档

### 主要API端点

#### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

#### AI功能
- `POST /api/ai/chat` - 项目AI对话
- `POST /api/ai/analyze` - 需求分析
- `POST /api/ai/generate-stage-documents` - 分阶段文档生成

#### AI配置
- `GET /api/ai/config` - 获取用户AI配置
- `PUT /api/ai/config` - 更新用户AI配置
- `POST /api/ai/test-connection` - 测试AI连接

#### PUML功能
- `POST /api/puml/render` - 渲染PUML图片
- `POST /api/puml/create` - 创建PUML图表
- `GET /api/puml/project/{projectId}` - 获取项目PUML列表

## 开发指南

### 项目结构
```
plant_code/
├── cmd/server/          # 服务器入口
├── internal/
│   ├── ai/             # AI服务集成
│   ├── api/            # HTTP处理器
│   ├── config/         # 配置管理
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── utils/          # 工具函数
├── web/                # 前端代码
│   ├── src/
│   │   ├── components/ # React组件
│   │   ├── pages/      # 页面组件
│   │   ├── services/   # API服务
│   │   └── types/      # TypeScript类型
└── docs/               # 项目文档
```


## 更新日志

### v1.2.0 (2025-01-08)
- ✨ 新增用户个性化AI配置功能
- 🔧 所有AI功能现在使用用户配置而非系统默认配置
- 🐛 修复PUML预览和AI助手的错误
- 📚 改进API文档和用户指南

### v1.1.0
- 分阶段文档生成功能
- PUML在线编辑器
- 项目上下文AI对话

### v1.0.0
- 基础需求分析功能
- 用户认证系统
- 项目管理功能

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或联系开发团队。 