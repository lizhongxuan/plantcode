import React, { useState, useEffect, useRef } from 'react';
import api from '@/services/api';
import Button from '@/components/ui/Button';
import { 
  Send, 
  Bot, 
  User, 
  Loader2, 
  MessageCircle,
  RefreshCw
} from 'lucide-react';
import type { Project } from '@/types';

interface ProjectChatProps {
  projectId: string;
  project: Project;
}

interface ChatMessage {
  message_id: string;
  session_id: string;
  sender_type: 'user' | 'assistant' | 'system';
  message_content: string;
  message_type: string;
  timestamp: string;
  processed: boolean;
}

interface ChatSession {
  session_id: string;
  project_id: string;
  user_id: string;
  session_type: string;
  started_at: string;
  status: string;
}

const ProjectChat: React.FC<ProjectChatProps> = ({ projectId, project }) => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [session, setSession] = useState<ChatSession | null>(null);
  const [inputMessage, setInputMessage] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    initializeChat();
  }, [projectId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const initializeChat = async () => {
    try {
      setIsLoading(true);
      
      // 创建或获取聊天会话
      const sessionResponse = await api.post('/ai/chat/session', {
        project_id: projectId,
        title: `${project.project_name} - AI助手对话`,
      });

      if (sessionResponse.data.success) {
        setSession(sessionResponse.data.data);
        
        // 加载对话历史
        loadChatHistory(sessionResponse.data.data.session_id);
      } else {
        setError(sessionResponse.data.error || '创建聊天会话失败');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '初始化聊天失败');
    } finally {
      setIsLoading(false);
    }
  };

  const loadChatHistory = async (sessionId: string) => {
    try {
      const response = await api.get(`/ai/chat/session/${sessionId}/messages`);
      if (response.data.success) {
        setMessages(response.data.data || []);
      }
    } catch (err) {
      console.error('加载聊天历史失败:', err);
    }
  };

  const sendMessage = async () => {
    if (!inputMessage.trim() || !session) return;

    const messageContent = inputMessage.trim();
    setInputMessage('');
    setIsSending(true);
    setError('');

    try {
      // 添加用户消息到界面
      const userMessage: ChatMessage = {
        message_id: `temp-${Date.now()}`,
        session_id: session.session_id,
        sender_type: 'user',
        message_content: messageContent,
        message_type: 'text',
        timestamp: new Date().toISOString(),
        processed: true,
      };
      
      setMessages(prev => [...prev, userMessage]);

      // 发送消息到服务器
      const response = await api.post('/ai/chat/message', {
        session_id: session.session_id,
        content: messageContent,
        role: 'user',
      });

      if (response.data.success) {
        // 等待2秒后刷新消息列表以获取AI回复
        setTimeout(() => {
          loadChatHistory(session.session_id);
        }, 2000);
      } else {
        setError(response.data.error || '发送消息失败');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '发送消息失败');
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

  const formatTime = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-96">
        <div className="flex flex-col items-center space-y-4">
          <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          <p className="text-gray-600">正在初始化AI助手...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-[600px]">
      {/* 头部 */}
      <div className="bg-white border-b border-gray-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <MessageCircle className="h-6 w-6 text-blue-600" />
            <div>
              <h3 className="text-lg font-semibold text-gray-900">AI项目助手</h3>
              <p className="text-sm text-gray-500">
                针对 <strong>{project.project_name}</strong> 的智能对话助手
              </p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-green-500 rounded-full"></div>
            <span className="text-sm text-gray-600">在线</span>
          </div>
        </div>
      </div>

      {/* 消息列表 */}
      <div className="flex-1 overflow-y-auto px-6 py-4 space-y-4 bg-gray-50">
        {messages.length === 0 ? (
          <div className="text-center py-12">
            <Bot className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500 mb-2">
              欢迎使用AI项目助手！我可以帮您：
            </p>
            <ul className="text-sm text-gray-500 space-y-1">
              <li>• 解答项目相关问题</li>
              <li>• 分析需求和技术方案</li>
              <li>• 提供开发建议</li>
              <li>• 协助解决技术难题</li>
            </ul>
          </div>
        ) : (
          messages.map((message, index) => (
            <div
              key={message.message_id || index}
              className={`flex ${message.sender_type === 'user' ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-[70%] rounded-lg px-4 py-2 ${
                  message.sender_type === 'user'
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-200 text-gray-900'
                }`}
              >
                <div className="flex items-start space-x-2">
                  {message.sender_type !== 'user' && (
                    <Bot className="h-5 w-5 text-blue-600 flex-shrink-0 mt-0.5" />
                  )}
                  <div className="flex-1">
                    <p className="text-sm leading-relaxed whitespace-pre-wrap">
                      {message.message_content}
                    </p>
                    <p className={`text-xs mt-1 ${
                      message.sender_type === 'user' ? 'text-blue-200' : 'text-gray-500'
                    }`}>
                      {formatTime(message.timestamp)}
                    </p>
                  </div>
                  {message.sender_type === 'user' && (
                    <User className="h-5 w-5 text-blue-200 flex-shrink-0 mt-0.5" />
                  )}
                </div>
              </div>
            </div>
          ))
        )}
        
        {isSending && (
          <div className="flex justify-start">
            <div className="bg-white border border-gray-200 rounded-lg px-4 py-2">
              <div className="flex items-center space-x-2">
                <Bot className="h-5 w-5 text-blue-600" />
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                </div>
              </div>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>

      {/* 错误提示 */}
      {error && (
        <div className="px-6 py-2">
          <div className="bg-red-50 border border-red-200 rounded-lg px-4 py-2">
            <p className="text-sm text-red-600">{error}</p>
          </div>
        </div>
      )}

      {/* 输入框 */}
      <div className="bg-white border-t border-gray-200 px-6 py-4">
        <div className="flex items-end space-x-4">
          <div className="flex-1">
            <textarea
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder={`向AI助手提问关于 ${project.project_name} 的问题...`}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
              rows={1}
              style={{ minHeight: '44px', maxHeight: '120px' }}
              disabled={isSending}
            />
          </div>
          <Button
            onClick={sendMessage}
            disabled={!inputMessage.trim() || isSending}
            className="flex items-center px-6 py-3"
          >
            {isSending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Send className="h-4 w-4" />
            )}
            <span className="ml-2">发送</span>
          </Button>
        </div>
        <p className="text-xs text-gray-500 mt-2">
          按 Enter 发送消息，Shift + Enter 换行
        </p>
      </div>
    </div>
  );
};

export default ProjectChat; 