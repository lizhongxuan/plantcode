import React, { useState, useEffect, useRef } from 'react';
import { Button, Input, Card, Avatar, Typography, Spin, message, Drawer } from 'antd';
import { SendOutlined, RobotOutlined, UserOutlined, CloseOutlined } from '@ant-design/icons';
import { 
  AIConversation, 
  AIMessage, 
  AIConversationResponse, 
  StartAIConversationRequest,
  SendAIMessageRequest,
  Project,
  ProjectStructureResponse
} from '../types';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;

interface PersistentAIChatProps {
  project: Project;
  projectStructure?: ProjectStructureResponse;
  visible: boolean;
  onClose: () => void;
  onDocumentChange?: (documentId: string, changes: any) => void;
}

export const PersistentAIChat: React.FC<PersistentAIChatProps> = ({
  project,
  projectStructure,
  visible,
  onClose,
  onDocumentChange
}) => {
  const [conversation, setConversation] = useState<AIConversation | null>(null);
  const [messages, setMessages] = useState<AIMessage[]>([]);
  const [inputText, setInputText] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSending, setIsSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // 滚动到消息底部
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // 初始化或获取活跃对话
  useEffect(() => {
    if (visible && project) {
      initializeConversation();
    }
  }, [visible, project.project_id]);

  const initializeConversation = async () => {
    setIsLoading(true);
    try {
      // 先尝试获取活跃对话
      const activeResponse = await fetch(`/api/ai/conversations/active/${project.project_id}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (activeResponse.ok) {
        const activeConv: AIConversation = await activeResponse.json();
        setConversation(activeConv);
        await loadConversationMessages(activeConv.conversation_id);
      } else {
        // 如果没有活跃对话，创建新的
        await createNewConversation();
      }
    } catch (error) {
      console.error('Failed to initialize conversation:', error);
      message.error('Failed to initialize AI chat');
    } finally {
      setIsLoading(false);
    }
  };

  const createNewConversation = async () => {
    try {
      const request: StartAIConversationRequest = {
        project_id: project.project_id,
        title: `AI Assistant for ${project.project_name}`
      };

      const response = await fetch('/api/ai/conversations/start', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(request)
      });

      if (response.ok) {
        const newConversation: AIConversation = await response.json();
        setConversation(newConversation);
        await loadConversationMessages(newConversation.conversation_id);
      } else {
        throw new Error('Failed to create conversation');
      }
    } catch (error) {
      console.error('Failed to create conversation:', error);
      message.error('Failed to start AI chat');
    }
  };

  const loadConversationMessages = async (conversationId: string) => {
    try {
      const response = await fetch(`/api/ai/conversations/${conversationId}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (response.ok) {
        const convResponse: AIConversationResponse = await response.json();
        setMessages(convResponse.messages);
      }
    } catch (error) {
      console.error('Failed to load messages:', error);
    }
  };

  const sendMessage = async () => {
    if (!inputText.trim() || !conversation || isSending) return;

    setIsSending(true);
    const userMessage: AIMessage = {
      message_id: `temp-${Date.now()}`,
      conversation_id: conversation.conversation_id,
      role: 'user',
      content: inputText,
      message_type: 'text',
      created_at: new Date().toISOString()
    };

    // 立即添加用户消息到UI
    setMessages(prev => [...prev, userMessage]);
    setInputText('');

    try {
      const request: SendAIMessageRequest = {
        conversation_id: conversation.conversation_id,
        content: inputText,
        message_type: 'text'
      };

      const response = await fetch('/api/ai/conversations/message', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(request)
      });

      if (response.ok) {
        const aiMessage: AIMessage = await response.json();
        
        // 更新消息列表，替换临时用户消息并添加AI回复
        setMessages(prev => {
          const withoutTemp = prev.filter(msg => msg.message_id !== userMessage.message_id);
          return [...withoutTemp, { ...userMessage, message_id: `user-${Date.now()}` }, aiMessage];
        });

        // 检查AI回复是否包含文档操作
        if (aiMessage.message_type === 'document_change' && onDocumentChange) {
          try {
            const metadata = JSON.parse(aiMessage.metadata || '{}');
            onDocumentChange(metadata.document_id, metadata.changes);
          } catch (error) {
            console.error('Failed to parse document change metadata:', error);
          }
        }
      } else {
        throw new Error('Failed to send message');
      }
    } catch (error) {
      console.error('Failed to send message:', error);
      message.error('Failed to send message');
      // 移除失败的用户消息
      setMessages(prev => prev.filter(msg => msg.message_id !== userMessage.message_id));
    } finally {
      setIsSending(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const getContextInfo = () => {
    if (!projectStructure) return '';
    
    const folderSummary = projectStructure.folders.map(folder => 
      `${folder.folder_name}: ${folder.documents.length} documents`
    ).join(', ');

    return `Project: ${project.project_name} | Folders: ${folderSummary}`;
  };

  const renderMessage = (msg: AIMessage) => {
    const isUser = msg.role === 'user';
    const isSystem = msg.role === 'system';
    
    return (
      <div
        key={msg.message_id}
        className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-4`}
      >
        <div className={`flex max-w-[80%] ${isUser ? 'flex-row-reverse' : 'flex-row'}`}>
          <Avatar
            className={`${isUser ? 'ml-2' : 'mr-2'}`}
            style={{ 
              backgroundColor: isSystem ? '#f56a00' : isUser ? '#1890ff' : '#52c41a' 
            }}
            icon={isUser ? <UserOutlined /> : <RobotOutlined />}
          />
          <Card
            size="small"
            className={`${isUser ? 'bg-blue-50' : isSystem ? 'bg-orange-50' : 'bg-green-50'}`}
            bodyStyle={{ padding: '8px 12px' }}
          >
            <div className="text-xs text-gray-500 mb-1">
              {new Date(msg.created_at).toLocaleTimeString()}
            </div>
            <Paragraph 
              className="mb-0 whitespace-pre-wrap"
              style={{ fontSize: '14px' }}
            >
              {msg.content}
            </Paragraph>
            {msg.message_type === 'document_change' && (
              <div className="mt-2 p-2 bg-yellow-100 rounded text-xs">
                <Text type="secondary">📝 Document modified</Text>
              </div>
            )}
          </Card>
        </div>
      </div>
    );
  };

  return (
    <Drawer
      title={
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <RobotOutlined className="mr-2" />
            <span>AI Assistant</span>
          </div>
          <Button 
            type="text" 
            icon={<CloseOutlined />} 
            onClick={onClose}
            size="small"
          />
        </div>
      }
      placement="right"
      width={400}
      open={visible}
      onClose={onClose}
      closable={false}
      bodyStyle={{ padding: 0 }}
    >
      <div className="flex flex-col h-full">
        {/* Context Info */}
        <div className="p-3 bg-gray-50 border-b text-xs text-gray-600">
          {getContextInfo()}
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto p-4">
          {isLoading ? (
            <div className="flex justify-center items-center h-32">
              <Spin size="large" />
            </div>
          ) : (
            <>
              {messages.map(renderMessage)}
              {isSending && (
                <div className="flex justify-start mb-4">
                  <div className="flex">
                    <Avatar
                      className="mr-2"
                      style={{ backgroundColor: '#52c41a' }}
                      icon={<RobotOutlined />}
                    />
                    <Card size="small" className="bg-green-50">
                      <Spin size="small" />
                      <Text className="ml-2 text-xs text-gray-500">AI is thinking...</Text>
                    </Card>
                  </div>
                </div>
              )}
              <div ref={messagesEndRef} />
            </>
          )}
        </div>

        {/* Input Area */}
        <div className="border-t p-4">
          <div className="flex gap-2">
            <TextArea
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Ask AI to help with requirements, design, or tasks..."
              autoSize={{ minRows: 1, maxRows: 3 }}
              disabled={!conversation || isLoading}
            />
            <Button
              type="primary"
              icon={<SendOutlined />}
              onClick={sendMessage}
              disabled={!inputText.trim() || !conversation || isSending || isLoading}
            />
          </div>
          <div className="text-xs text-gray-500 mt-2">
            Press Enter to send, Shift+Enter for new line
          </div>
        </div>
      </div>
    </Drawer>
  );
};

export default PersistentAIChat;