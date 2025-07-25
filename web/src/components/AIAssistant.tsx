import React, { useState, useRef, useEffect } from 'react';
import { 
  Card, 
  Input, 
  Button, 
  List, 
  Typography, 
  Space, 
  Spin, 
  Divider,
  Avatar,
  message 
} from 'antd';
import { 
  RobotOutlined, 
  SendOutlined, 
  UserOutlined,
  FileTextOutlined,
  CodeOutlined,
  ClearOutlined
} from '@ant-design/icons';
import { aiApi } from '@/services/api';

const { Text } = Typography;
const { TextArea } = Input;

interface AIAssistantProps {
  projectId: string;
  context: 'document' | 'puml' | 'requirement_analysis';
  stage?: number;
  onGenerate?: (type: 'document' | 'puml', content: any) => void;
}

interface ChatMessage {
  id: string;
  type: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  actions?: Array<{
    type: 'generate_document' | 'generate_puml';
    label: string;
    data: any;
  }>;
}

const AIAssistant: React.FC<AIAssistantProps> = ({ 
  projectId, 
  context, 
  stage, 
  onGenerate 
}) => {
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      id: '1',
      type: 'assistant',
      content: getWelcomeMessage(context),
      timestamp: new Date(),
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // 滚动到底部
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // 获取欢迎消息
  function getWelcomeMessage(context: string): string {
    switch (context) {
      case 'document':
        return `你好！我是AI文档助手。我可以帮你：\n\n• 根据项目需求生成各类文档\n• 优化和改进现有文档内容\n• 回答文档相关的问题\n\n请告诉我你需要什么帮助？`;
      case 'puml':
        return `你好！我是PUML图表助手。我可以帮你：\n\n• 生成各种类型的PUML图表\n• 优化现有图表结构\n• 解答PUML语法问题\n\n你想创建什么类型的图表？`;
      case 'requirement_analysis':
        return `你好！我是需求分析助手。我可以帮你：\n\n• 分析和梳理项目需求\n• 识别缺失的需求信息\n• 优化需求描述\n\n请描述你的项目需求或问题。`;
      default:
        return '你好！我是AI助手，很高兴为你服务！';
    }
  }

  // 发送消息
  const handleSendMessage = async () => {
    if (!inputValue.trim() || loading) return;

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      type: 'user',
      content: inputValue.trim(),
      timestamp: new Date(),
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setLoading(true);

    try {
      // 调用AI对话接口
      const response = await aiApi.projectChat(projectId, userMessage.content, context);
      
      if (response.success) {
        const aiMessage: ChatMessage = {
          id: (Date.now() + 1).toString(),
          type: 'assistant',
          content: response.data.message || response.data.response || '我明白了，让我为你处理。',
          timestamp: new Date(),
          actions: generateActions(response.data, context)
        };

        setMessages(prev => [...prev, aiMessage]);
      } else {
        throw new Error(response.message || 'AI对话失败');
      }
    } catch (error: any) {
      console.error('AI对话失败:', error);
      
      const errorMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        type: 'assistant',
        content: `抱歉，我遇到了一些问题：${error.message}。请稍后再试或换个方式描述你的需求。`,
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setLoading(false);
    }
  };

  // 生成操作按钮
  const generateActions = (_data: any, context: string) => {
    const actions = [];

    // 根据上下文和AI响应生成不同的操作
    if (context === 'document') {
      actions.push({
        type: 'generate_document' as const,
        label: '生成文档',
        data: { stage, type: 'requirement' }
      });
    } else if (context === 'puml') {
      actions.push({
        type: 'generate_puml' as const,
        label: '生成架构图',
        data: { type: 'architecture', stage }
      });
      actions.push({
        type: 'generate_puml' as const,
        label: '生成流程图',
        data: { type: 'business_flow', stage }
      });
    }

    return actions;
  };

  // 处理操作按钮点击
  const handleAction = async (action: any) => {
    try {
      if (action.type === 'generate_document') {
        // 生成文档
        const response = await aiApi.generateStageDocumentList(projectId, action.data.stage || 1);
        if (response.success && onGenerate) {
          onGenerate('document', response.data);
          message.success('文档生成成功！');
        }
      } else if (action.type === 'generate_puml') {
        // 生成PUML图表
        const pumlContent = generateSamplePUML(action.data.type);
        if (onGenerate) {
          onGenerate('puml', {
            type: action.data.type,
            content: pumlContent,
            name: `${action.data.type === 'architecture' ? '系统架构图' : '业务流程图'}`
          });
          message.success('PUML图表生成成功！');
        }
      }
    } catch (error: any) {
      message.error(error.message || '生成失败');
    }
  };

  // 生成示例PUML代码
  const generateSamplePUML = (type: string): string => {
    if (type === 'architecture') {
      return `@startuml 系统架构图
!theme plain

package "前端层" {
  [Web界面] as frontend
  [移动端] as mobile
}

package "服务层" {
  [API网关] as gateway
  [业务服务] as service
  [认证服务] as auth
}

package "数据层" {
  database "MySQL" as db
  database "Redis" as cache
}

frontend --> gateway
mobile --> gateway
gateway --> service
gateway --> auth
service --> db
service --> cache

@enduml`;
    } else {
      return `@startuml 业务流程图
|用户|
start
:登录系统;
:选择功能;

if (权限验证) then (通过)
  :执行操作;
  :返回结果;
else (失败)
  :显示错误信息;
endif

stop
@enduml`;
    }
  };

  // 清空对话
  const handleClearChat = () => {
    setMessages([{
      id: '1',
      type: 'assistant',
      content: getWelcomeMessage(context),
      timestamp: new Date(),
    }]);
  };

  // 按回车发送
  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <Card 
      title={
        <Space>
          <RobotOutlined />
          <span>AI助手</span>
        </Space>
      }
      extra={
        <Button 
          icon={<ClearOutlined />} 
          size="small" 
          onClick={handleClearChat}
        >
          清空对话
        </Button>
      }
      style={{ height: '500px', display: 'flex', flexDirection: 'column' }}
      bodyStyle={{ flex: 1, display: 'flex', flexDirection: 'column', padding: '12px' }}
    >
      {/* 聊天消息区域 */}
      <div style={{ flex: 1, overflowY: 'auto', marginBottom: '12px' }}>
        <List
          itemLayout="horizontal"
          dataSource={messages}
          renderItem={(message) => (
            <List.Item style={{ padding: '8px 0', border: 'none' }}>
              <List.Item.Meta
                avatar={
                  <Avatar 
                    icon={message.type === 'user' ? <UserOutlined /> : <RobotOutlined />}
                    style={{ 
                      backgroundColor: message.type === 'user' ? '#1890ff' : '#52c41a' 
                    }}
                  />
                }
                title={
                  <Space>
                    <Text strong>{message.type === 'user' ? '你' : 'AI助手'}</Text>
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                      {message.timestamp.toLocaleTimeString()}
                    </Text>
                  </Space>
                }
                description={
                  <div>
                    <Text style={{ whiteSpace: 'pre-wrap' }}>{message.content}</Text>
                    {message.actions && message.actions.length > 0 && (
                      <div style={{ marginTop: '8px' }}>
                        <Space>
                          {message.actions.map((action, index) => (
                            <Button
                              key={index}
                              size="small"
                              type="primary"
                              ghost
                              icon={action.type === 'generate_document' ? <FileTextOutlined /> : <CodeOutlined />}
                              onClick={() => handleAction(action)}
                            >
                              {action.label}
                            </Button>
                          ))}
                        </Space>
                      </div>
                    )}
                  </div>
                }
              />
            </List.Item>
          )}
        />
        <div ref={messagesEndRef} />
      </div>

      {/* 输入区域 */}
      <div>
        <Divider style={{ margin: '8px 0' }} />
        <Space.Compact style={{ display: 'flex' }}>
          <TextArea
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder={`输入你的${context === 'document' ? '文档需求' : context === 'puml' ? 'PUML图表需求' : '问题'}...`}
            autoSize={{ minRows: 1, maxRows: 3 }}
            style={{ flex: 1 }}
            disabled={loading}
          />
          <Button 
            type="primary" 
            icon={loading ? <Spin size="small" /> : <SendOutlined />}
            onClick={handleSendMessage}
            disabled={loading || !inputValue.trim()}
          >
            发送
          </Button>
        </Space.Compact>
      </div>
    </Card>
  );
};

export default AIAssistant; 