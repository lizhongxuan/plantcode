import React, { useState, useEffect, useRef } from 'react';
import api from '@/services/api';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { 
  Target, 
  CheckCircle, 
  Loader2, 
  RefreshCw, 
  Play, 
  MessageCircle, 
  Send, 
  Bot, 
  User,
  Lightbulb,
  FileText,
  Users,
  Workflow,
  Database
} from 'lucide-react';
import type { Project } from '@/types';

interface RequirementAnalysisProps {
  projectId: string;
  project: Project;
  onAnalysisComplete?: () => void;
}

interface AnalysisResult {
  analysis_id: string;
  core_functions: string[];
  roles: string[];
  business_processes: string[];
  data_entities: string[];
  completeness_score: number;
  created_at: string;
}

interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: string;
}

const RequirementAnalysis: React.FC<RequirementAnalysisProps> = ({ projectId, project, onAnalysisComplete }) => {
  const [step, setStep] = useState<'input' | 'analyzing' | 'interactive'>('input');
  const [requirement, setRequirement] = useState('');
  const [analysisResult, setAnalysisResult] = useState<AnalysisResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  
  // AI对话相关状态
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [currentMessage, setCurrentMessage] = useState('');
  const [isChatLoading, setIsChatLoading] = useState(false);
  const chatContainerRef = useRef<HTMLDivElement>(null);

  // 加载现有分析结果
  useEffect(() => {
    loadExistingAnalysis();
  }, [projectId]);

  // 自动滚动到聊天底部
  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [chatMessages]);

  const loadExistingAnalysis = async () => {
    try {
      const response = await api.get(`/ai/analysis/project/${projectId}`);
      if (response.data.success && response.data.data.length > 0) {
        const latestAnalysis = response.data.data[0];
        setAnalysisResult(latestAnalysis);
        setStep('interactive');
        if (onAnalysisComplete) {
          onAnalysisComplete();
        }
        // 初始化对话
        initializeChat();
      }
    } catch (err) {
      // 没有现有分析，保持输入状态
    }
  };

  const initializeChat = () => {
    const welcomeMessage: ChatMessage = {
      id: Date.now().toString(),
      role: 'assistant',
      content: '您好！我是您的AI助手。您的项目需求分析已完成，我可以帮您进一步优化需求分析、回答相关问题或者协助您完善项目细节。请告诉我您想了解什么？',
      timestamp: new Date().toISOString()
    };
    setChatMessages([welcomeMessage]);
  };

  const handleInitialAnalysis = async () => {
    if (!requirement.trim()) {
      setError('请输入项目需求描述');
      return;
    }

    try {
      setIsLoading(true);
      setError('');
      setStep('analyzing');

      const response = await api.post('/ai/analyze', {
        project_id: projectId,
        requirement: requirement
      });

      if (response.data.success) {
        setAnalysisResult(response.data.data);
        setStep('interactive');
        if (onAnalysisComplete) {
          onAnalysisComplete();
        }
        // 初始化对话
        initializeChat();
      } else {
        setError(response.data.error || '分析失败');
        setStep('input');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '分析失败');
      setStep('input');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChatMessage = async () => {
    if (!currentMessage.trim() || isChatLoading) return;

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      role: 'user',
      content: currentMessage,
      timestamp: new Date().toISOString()
    };

    setChatMessages(prev => [...prev, userMessage]);
    setCurrentMessage('');
    setIsChatLoading(true);

    try {
      // 这里调用项目上下文的AI对话API
      const response = await api.post('/ai/chat', {
        project_id: projectId,
        message: currentMessage,
        context: 'requirement_analysis'
      });

      if (response.data.success) {
        const assistantMessage: ChatMessage = {
          id: (Date.now() + 1).toString(),
          role: 'assistant',
          content: response.data.data.message,
          timestamp: new Date().toISOString()
        };
        setChatMessages(prev => [...prev, assistantMessage]);

        // 如果AI返回了更新的分析结果，更新分析数据
        if (response.data.data.updated_analysis) {
          setAnalysisResult(response.data.data.updated_analysis);
        }
      }
    } catch (err: any) {
      const errorMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: '抱歉，我遇到了一些问题。请稍后再试。',
        timestamp: new Date().toISOString()
      };
      setChatMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsChatLoading(false);
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600';
    if (score >= 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleChatMessage();
    }
  };

  // 初始输入阶段
  if (step === 'input') {
    return (
      <div className="p-6">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-8">
            <Target className="h-12 w-12 text-blue-600 mx-auto mb-4" />
            <h2 className="text-2xl font-semibold text-gray-900 mb-2">
              智能需求分析
            </h2>
            <p className="text-gray-600">
              详细描述 <strong>{project.project_name}</strong> 的功能需求，AI将帮您生成结构化的需求文档，并可通过对话进一步优化
            </p>
          </div>

          <div className="bg-white border border-gray-200 rounded-lg p-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              项目需求描述 <span className="text-red-500">*</span>
            </label>
            <textarea
              value={requirement}
              onChange={(e) => setRequirement(e.target.value)}
              placeholder={`请详细描述 ${project.project_name} 的功能需求，例如：
• 主要功能模块
• 用户角色和权限
• 业务流程
• 数据管理需求
• 技术要求等...`}
              className="w-full h-64 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            />
            <div className="flex justify-between items-center mt-4">
              <p className="text-sm text-gray-500">
                建议输入200字以上的详细描述，分析完成后可通过AI对话进一步优化
              </p>
              <span className="text-sm text-gray-500">
                {requirement.length} 字符
              </span>
            </div>
          </div>

          {error && (
            <div className="mt-4 bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-600 text-sm">{error}</p>
            </div>
          )}

          <div className="mt-6 flex justify-center">
            <Button
              onClick={handleInitialAnalysis}
              disabled={isLoading}
              className="px-8 py-3"
            >
              {isLoading ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  分析中...
                </>
              ) : (
                <>
                  <Play className="h-4 w-4 mr-2" />
                  开始AI分析
                </>
              )}
            </Button>
          </div>
        </div>
      </div>
    );
  }

  // 分析中阶段
  if (step === 'analyzing') {
    return (
      <div className="p-6">
        <div className="max-w-2xl mx-auto text-center">
          <Loader2 className="h-16 w-16 text-blue-600 animate-spin mx-auto mb-6" />
          <h2 className="text-xl font-semibold text-gray-900 mb-2">AI正在分析您的需求...</h2>
          <p className="text-gray-600 mb-6">这可能需要几十秒的时间，分析完成后您可以通过AI对话进一步优化需求</p>
        </div>
      </div>
    );
  }

  // 交互式分析阶段
  if (step === 'interactive' && analysisResult) {
    return (
      <div className="h-[800px] flex">
        {/* 左侧 - 需求分析结果 */}
        <div className="w-1/2 p-6 overflow-y-auto border-r border-gray-200">
          <div className="mb-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-3">
                <CheckCircle className="h-6 w-6 text-green-600" />
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">需求分析结果</h3>
                  <p className="text-sm text-gray-600">
                    {new Date(analysisResult.created_at).toLocaleString()}
                  </p>
                </div>
              </div>
              <div className="px-3 py-1 rounded-lg bg-blue-50 border border-blue-200">
                <span className="text-sm font-medium text-gray-600">完整度</span>
                <span className={`ml-2 text-sm font-bold ${getScoreColor(analysisResult.completeness_score)}`}>
                  {analysisResult.completeness_score}%
                </span>
              </div>
            </div>
          </div>

          <div className="space-y-6">
            {/* 核心功能 */}
            <div className="bg-blue-50 rounded-lg p-4">
              <div className="flex items-center mb-3">
                <Lightbulb className="h-5 w-5 text-blue-600 mr-2" />
                <h4 className="font-medium text-gray-900">核心功能</h4>
              </div>
              <ul className="space-y-2">
                {analysisResult.core_functions.map((func, index) => (
                  <li key={index} className="flex items-start">
                    <span className="block w-2 h-2 bg-blue-500 rounded-full mt-2 mr-3 flex-shrink-0"></span>
                    <span className="text-sm text-gray-700">{func}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* 用户角色 */}
            <div className="bg-green-50 rounded-lg p-4">
              <div className="flex items-center mb-3">
                <Users className="h-5 w-5 text-green-600 mr-2" />
                <h4 className="font-medium text-gray-900">用户角色</h4>
              </div>
              <div className="flex flex-wrap gap-2">
                {analysisResult.roles.map((role, index) => (
                  <span key={index} className="px-3 py-1 bg-green-100 text-green-800 text-sm rounded-full">
                    {role}
                  </span>
                ))}
              </div>
            </div>

            {/* 业务流程 */}
            <div className="bg-purple-50 rounded-lg p-4">
              <div className="flex items-center mb-3">
                <Workflow className="h-5 w-5 text-purple-600 mr-2" />
                <h4 className="font-medium text-gray-900">业务流程</h4>
              </div>
              <ul className="space-y-2">
                {analysisResult.business_processes.map((process, index) => (
                  <li key={index} className="flex items-start">
                    <span className="block w-2 h-2 bg-purple-500 rounded-full mt-2 mr-3 flex-shrink-0"></span>
                    <span className="text-sm text-gray-700">{process}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* 数据实体 */}
            <div className="bg-orange-50 rounded-lg p-4">
              <div className="flex items-center mb-3">
                <Database className="h-5 w-5 text-orange-600 mr-2" />
                <h4 className="font-medium text-gray-900">数据实体</h4>
              </div>
              <div className="flex flex-wrap gap-2">
                {analysisResult.data_entities.map((entity, index) => (
                  <span key={index} className="px-3 py-1 bg-orange-100 text-orange-800 text-sm rounded-full">
                    {entity}
                  </span>
                ))}
              </div>
            </div>
          </div>

          <div className="mt-6 pt-6 border-t border-gray-200">
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setStep('input');
                setRequirement('');
                setAnalysisResult(null);
                setChatMessages([]);
              }}
              className="w-full"
            >
              <RefreshCw className="h-4 w-4 mr-2" />
              重新分析
            </Button>
          </div>
        </div>

        {/* 右侧 - AI对话 */}
        <div className="w-1/2 flex flex-col">
          {/* 对话头部 */}
          <div className="p-4 border-b border-gray-200 bg-gray-50">
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-blue-100 rounded-full">
                <Bot className="h-5 w-5 text-blue-600" />
              </div>
              <div>
                <h3 className="font-medium text-gray-900">AI需求分析助手</h3>
                <p className="text-sm text-gray-600">优化需求分析，完善项目细节</p>
              </div>
            </div>
          </div>

          {/* 对话内容 */}
          <div className="flex-1 p-4 overflow-y-auto" ref={chatContainerRef}>
            <div className="space-y-4">
              {chatMessages.map((message) => (
                <div key={message.id} className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                  <div className={`flex items-start space-x-2 max-w-[80%] ${message.role === 'user' ? 'flex-row-reverse space-x-reverse' : ''}`}>
                    <div className={`p-2 rounded-full ${message.role === 'user' ? 'bg-blue-100' : 'bg-gray-100'}`}>
                      {message.role === 'user' ? (
                        <User className="h-4 w-4 text-blue-600" />
                      ) : (
                        <Bot className="h-4 w-4 text-gray-600" />
                      )}
                    </div>
                    <div className={`p-3 rounded-lg ${
                      message.role === 'user' 
                        ? 'bg-blue-600 text-white' 
                        : 'bg-gray-100 text-gray-900'
                    }`}>
                      <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                      <p className={`text-xs mt-1 ${message.role === 'user' ? 'text-blue-100' : 'text-gray-500'}`}>
                        {new Date(message.timestamp).toLocaleTimeString()}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
              {isChatLoading && (
                <div className="flex justify-start">
                  <div className="flex items-start space-x-2">
                    <div className="p-2 bg-gray-100 rounded-full">
                      <Bot className="h-4 w-4 text-gray-600" />
                    </div>
                    <div className="p-3 bg-gray-100 rounded-lg">
                      <div className="flex space-x-1">
                        <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                        <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                        <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* 对话输入 */}
          <div className="p-4 border-t border-gray-200 bg-gray-50">
            <div className="flex space-x-2">
              <Input
                value={currentMessage}
                onChange={(e) => setCurrentMessage(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="询问AI助手，例如：'这个功能还需要考虑什么？'、'用户权限如何设计？'"
                disabled={isChatLoading}
                className="flex-1"
              />
              <Button
                onClick={handleChatMessage}
                disabled={isChatLoading || !currentMessage.trim()}
                size="sm"
              >
                <Send className="h-4 w-4" />
              </Button>
            </div>
            <p className="text-xs text-gray-500 mt-2">
              💡 通过对话可以优化需求分析，完善项目细节，并自动更新PUML图表和技术文档
            </p>
          </div>
        </div>
      </div>
    );
  }

  return null;
};

export default RequirementAnalysis; 