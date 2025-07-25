# AI对话助手工作流设计

> **最后更新：2025-01-24**  
> **文档版本：v1.1**
> **状态：基于当前实现的架构分析**

## 当前实现状态 ✅

基于代码分析，AI对话助手已经具备以下核心功能：

### 已实现功能
- **Spec工作流**：三阶段开发流程完整实现 ✅
- **多AI提供商支持**：OpenAI、Gemini集成完成，Claude待完善
- **用户个性化配置**：API密钥、模型选择、提供商配置
- **项目上下文管理**：项目级别的AI对话和上下文保持
- **在线PUML编辑**：实时预览、双向同步、AI辅助修改 ✅
- **需求分析工作流**：基于EARS语法的需求提取和文档生成 ✅
- **文档管理系统**：结构化项目文档组织和版本控制 ✅
- **路由重构**：项目详情直接进入Spec工作流 ✅

### 架构设计

```
[前端组件] 
    ↓
[AIAssistant.tsx] - React对话界面
    ↓
[API接口层] - /api/ai/* 路由
    ↓  
[AIController] - 请求处理和验证
    ↓
[AIService] - 业务逻辑处理  
    ↓
[AIManager] - 多提供商统一管理
    ↓
[具体AI客户端] - OpenAI/Claude/Gemini
    ↓
[外部AI服务] - 各提供商API
```

---

## 工作流程详解

### 工作流1：Spec工作流集成对话 ✅ **已实现**

**实现状态**：Spec工作流已全面实现,包含三阶段开发流程和AI对话集成

```
阶段1：需求确定 ✅
[ProjectSpecWorkflow] → [PersistentAIChat]
用户通过AI对话实时生成和更新需求文档
支持EARS语法标准的需求结构化

阶段2：设计生成 ✅  
[PUML图表管理] → [OnlinePUMLEditor] → [AI对话修改]
基于需求文档生成项目架构PUML图
在线编辑器支持实时预览和AI辅助修改

阶段3：任务分解 ✅
[TODO文档生成] → [AI对话优化]
基于需求和设计生成开发任务列表和执行步骤
```

**核心组件**：
- ✅ ProjectSpecWorkflow: 主工作流页面，三阶段管理
- ✅ PersistentAIChat: 右侧常驻AI对话助手  
- ✅ ProjectFolderStructure: 左侧文件目录管理
- ✅ OnlinePUMLEditor: 在线PUML编辑和预览
- ✅ API集成: projectApi.generateDesign/generateTodos

**页面布局**：
- **左侧(320px)**: 项目信息 + 阶段进度 + 文件目录树  
- **中间(flex)**: 当前阶段内容 + 文档查看标签页
- **右侧(400px)**: AI对话助手 (可收起/展开)

---

## 待完善的功能

### 1. Claude API集成 ⚠️ **高优先级**

**问题**：`/internal/ai/claude_client.go` 中实现返回"暂时不可用"

**需要完成**：
```go
func (c *ClaudeClient) Chat(messages []Message, options *ChatOptions) (*ChatResponse, error) {
    // TODO: 实现Claude API调用逻辑
    // 1. 构建Claude API请求格式
    // 2. 发送HTTP请求到 https://api.anthropic.com
    // 3. 解析响应并转换为统一格式
    // 4. 错误处理和重试逻辑
}
```

### 2. 知识库RAG集成 📚 **核心特性**

**当前缺失**：README承诺的知识库驱动功能未实现

**需要设计**：
```
[用户问题] → [RAG检索系统]
    ↓
[向量化搜索] → [知识库文档]
    ↓  
[相关文档] → [AI上下文增强]
    ↓
[基于知识库的回答] → [用户]
```

**实现计划**：
1. 设计KnowledgeBase数据模型
2. 集成向量数据库(如Milvus/Pinecone)
3. 文档向量化和存储
4. RAG检索逻辑实现
5. AI提示词增强

### 3. 上下文工程优化 🧠 **体验提升**

**当前状况**：基础的项目级别对话历史

**改进方向**：
- 长对话的智能摘要和压缩  
- 关键信息的持久化存储
- 多轮对话的意图理解
- 上下文相关性评分

---

## 技术架构优化建议

### 数据流优化
```
当前：[前端] → [API] → [Service] → [AI] → [响应]
优化：[前端] → [API] → [缓存层] → [Service] → [AI] → [异步处理] → [WebSocket推送]
```

### 缓存策略
- AI响应缓存：相似问题复用结果
- 用户配置缓存：减少数据库查询
- PUML渲染缓存：提高图表加载速度

### 错误处理
- AI服务降级：主要AI不可用时切换备用
- 重试机制：网络异常和API限流处理
- 用户友好错误提示：技术错误转换为用户理解的描述

---

*本文档基于 `requirements_todo.md` 中的分析结果编写，反映了项目当前真实的实现状态和待完善功能。*