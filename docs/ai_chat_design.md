# AI对话助手工作流

### 根据对话内容生成PUML序列图

[前端] => [后端Websocket/API]
- 发送消息: { "projectId": "xxx", "conversationId": "yyy", "message": "帮我画一个用户登录的序列图" }

[后端] => [AI助手模块]
- 调用函数: assistantModule.HandleMessage(projectId, conversationId, message)

- 目的: AI助手模块是总指挥。它首先加载conversationId对应的历史对话，形成完整的上下文。

[AI助手模块] => [图表生成模块]

- 它识别出用户的意图是“画图”，于是调用图表模块。

- 调用函数: diagramModule.GenerateDiagramFromText(fullContext, userAiConfig)

- 参数 (fullContext): "历史对话... 用户最新消息：帮我画一个用户登录的序列图"

[图表生成模块] 内部调用 => [AI服务代理模块]

- 和工作流1类似，diagramModule构造一个专门用于生成PUML的prompt，然后通过aiProxyModule将请求发给用户选择的AI（例如，这次用户可能配置了Claude，因为它擅长代码生成）。

- Prompt示例: "根据以下描述，生成一个PUML格式的序列图代码：[fullContext]"

- 返回数据流: AI返回PUML代码 => aiProxyModule => diagramModule => assistantModule => [前端]

- 最终通过Websocket推送给前端的JSON:
JSON
{
"type": "puml_diagram",
"content": "@startuml\nUser -> Server: Login Request\nServer -> Database: Verify Credentials\nDatabase --> Server: User Info\nServer --> User: Login Success\n@enduml"
}
- 前端接收到此数据后，调用PUML渲染库，实时显示序列图。