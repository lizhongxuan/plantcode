import React, { useState, useEffect, useRef } from 'react';
import { Card, Input, Button, Progress, Space, Typography, Divider, Tag, List, message, Alert } from 'antd';
import { 
  SendOutlined, 
  UserOutlined, 
  RobotOutlined,
  CheckCircleOutlined,
  LoadingOutlined,
  BulbOutlined,
  TeamOutlined,
  PartitionOutlined,
  DatabaseOutlined,
  FileTextOutlined
} from '@ant-design/icons';
import { aiApi, asyncTaskApi } from '@/services/api';

const { TextArea } = Input;
const { Title, Text, Paragraph } = Typography;

interface RequirementAnalysisProps {
  projectId: string;
  onAnalysisUpdate?: () => void;
}

// 需求分析结果结构
interface AnalysisResult {
  analysis_id: string;
  core_functions: string[];
  roles: string[];
  business_processes: Array<{ name: string; steps: string[]; actors: string[]; description: string }>;
  data_entities: Array<{ name: string; attributes: any[]; description: string }>;
  completeness_score: number;
  created_at: string;
}

interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: string;
  analysisUpdate?: AnalysisResult; // 如果消息包含分析更新
}

const RequirementAnalysis: React.FC<RequirementAnalysisProps> = ({ 
  projectId, 
  onAnalysisUpdate 
}) => {
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

  // 转换后端数据为前端使用的格式
  const transformAnalysisData = (rawData: any): AnalysisResult => {
    let structuredData;
    try {
      structuredData = typeof rawData.structured_requirement === 'string' 
        ? JSON.parse(rawData.structured_requirement)
        : rawData.structured_requirement;
    } catch (error) {
      console.error('解析structured_requirement失败:', error);
      structuredData = {
        core_functions: [],
        roles: [],
        business_processes: [],
        data_entities: []
      };
    }

    return {
      analysis_id: rawData.requirement_id || rawData.analysis_id,
      core_functions: structuredData.core_functions || [],
      roles: structuredData.roles || [],
      business_processes: structuredData.business_processes || [],
      data_entities: structuredData.data_entities || [],
      completeness_score: rawData.completeness_score || 0,
      created_at: rawData.created_at || new Date().toISOString()
    };
  };

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
      setIsLoading(true);
      const response = await aiApi.getProjectAnalysis(projectId);
      if (response.success && response.data.length > 0) {
        const latestAnalysis = response.data[0];
        setAnalysisResult(transformAnalysisData(latestAnalysis));
        setStep('interactive');
        // 初始化对话
        initializeChat();
      }
    } catch (err) {
      // 没有现有分析，保持输入状态
      console.log('无现有分析结果，等待用户输入');
    } finally {
      setIsLoading(false);
    }
  };

  const initializeChat = () => {
    const welcomeMessage: ChatMessage = {
      id: Date.now().toString(),
      role: 'assistant',
      content: '您好！我是您的AI需求分析助手。我已经完成了初步的需求分析。您可以：\n\n1. 询问我关于项目的任何细节问题\n2. 要求我补充或修改需求分析\n3. 让我解释某个业务流程或数据实体\n\n请告诉我您想了解或改进什么？',
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

      const response = await aiApi.analyzeRequirement(projectId, requirement);

      if (response.success) {
        setAnalysisResult(transformAnalysisData(response.data));
        setStep('interactive');
        onAnalysisUpdate?.();
        // 初始化对话
        initializeChat();
      } else {
        throw new Error(response.message || '分析失败');
      }
    } catch (err: any) {
      console.error('需求分析失败:', err);
      setError(err.message || '需求分析失败，请稍后重试');
      setStep('input');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChatMessage = async () => {
    if (!currentMessage.trim()) return;

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
      // 使用项目上下文AI对话
      const response = await aiApi.projectChat(projectId, currentMessage, 'requirement_analysis');
      
      if (response.success) {
        const aiMessage: ChatMessage = {
          id: (Date.now() + 1).toString(),
          role: 'assistant',
          content: response.data.response,
          timestamp: new Date().toISOString()
        };

        // 检查是否包含分析更新
        if (response.data.analysis_update) {
          aiMessage.analysisUpdate = transformAnalysisData(response.data.analysis_update);
          setAnalysisResult(aiMessage.analysisUpdate);
          onAnalysisUpdate?.();
        }

        setChatMessages(prev => [...prev, aiMessage]);
      } else {
        throw new Error(response.message || '对话失败');
      }
    } catch (err: any) {
      console.error('AI对话失败:', err);
      const errorMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: `抱歉，我遇到了一些问题：${err.message || '请稍后重试'}`,
        timestamp: new Date().toISOString()
      };
      setChatMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsChatLoading(false);
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 0.8) return '#52c41a';
    if (score >= 0.6) return '#faad14';
    return '#ff4d4f';
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && e.ctrlKey) {
      e.preventDefault();
      if (step === 'input') {
        handleInitialAnalysis();
      } else {
        handleChatMessage();
      }
    }
  };

  // 渲染分析结果
  const renderAnalysisResult = () => {
    if (!analysisResult) return null;

    return (
      <Space direction="vertical" style={{ width: '100%' }} size="large">
        {/* 完整度评分 */}
        <Card size="small" title="需求完整度评估">
          <Progress
            percent={analysisResult.completeness_score * 100}
            strokeColor={getScoreColor(analysisResult.completeness_score)}
            format={percent => `${percent}%`}
          />
          <Text type="secondary" style={{ marginTop: 8, display: 'block' }}>
            当前需求分析完整度为 {(analysisResult.completeness_score * 100).toFixed(1)}%
          </Text>
        </Card>

        {/* 核心功能 */}
        {analysisResult.core_functions.length > 0 && (
          <Card size="small" title={<><BulbOutlined /> 核心功能</>}>
            <Space wrap>
              {analysisResult.core_functions.map((func, index) => (
                <Tag key={index} color="blue">{func}</Tag>
              ))}
            </Space>
          </Card>
        )}

        {/* 用户角色 */}
        {analysisResult.roles.length > 0 && (
          <Card size="small" title={<><TeamOutlined /> 用户角色</>}>
            <Space wrap>
              {analysisResult.roles.map((role, index) => (
                <Tag key={index} color="green">{role}</Tag>
              ))}
            </Space>
          </Card>
        )}

        {/* 业务流程 */}
        {analysisResult.business_processes.length > 0 && (
          <Card size="small" title={<><PartitionOutlined /> 业务流程</>}>
            <List
              size="small"
              dataSource={analysisResult.business_processes}
              renderItem={(process) => (
                <List.Item>
                  <div style={{ width: '100%' }}>
                    <Title level={5}>{process.name}</Title>
                    <Paragraph style={{ margin: 0 }}>
                      <Text type="secondary">{process.description}</Text>
                    </Paragraph>
                    {process.actors && process.actors.length > 0 && (
                      <div style={{ marginTop: 4 }}>
                        <Text strong>参与者：</Text>
                        <Space wrap>
                          {process.actors.map((actor, index) => (
                            <Tag key={index} size="small">{actor}</Tag>
                          ))}
                        </Space>
                      </div>
                    )}
                  </div>
                </List.Item>
              )}
            />
          </Card>
        )}

        {/* 数据实体 */}
        {analysisResult.data_entities.length > 0 && (
          <Card size="small" title={<><DatabaseOutlined /> 数据实体</>}>
            <List
              size="small"
              dataSource={analysisResult.data_entities}
              renderItem={(entity) => (
                <List.Item>
                  <div style={{ width: '100%' }}>
                    <Title level={5}>{entity.name}</Title>
                    <Paragraph style={{ margin: 0 }}>
                      <Text type="secondary">{entity.description}</Text>
                    </Paragraph>
                  </div>
                </List.Item>
              )}
            />
          </Card>
        )}
      </Space>
    );
  };

  // 渲染对话界面
  const renderChatInterface = () => (
    <Card title="AI需求助手对话" size="small">
      <div
        ref={chatContainerRef}
        style={{
          height: '400px',
          overflowY: 'auto',
          border: '1px solid #f0f0f0',
          borderRadius: '6px',
          padding: '12px',
          marginBottom: '12px',
          backgroundColor: '#fafafa'
        }}
      >
        {chatMessages.map((msg) => (
          <div
            key={msg.id}
            style={{
              display: 'flex',
              justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
              marginBottom: '12px'
            }}
          >
            <div
              style={{
                maxWidth: '80%',
                padding: '8px 12px',
                borderRadius: '8px',
                backgroundColor: msg.role === 'user' ? '#1890ff' : '#fff',
                color: msg.role === 'user' ? '#fff' : '#000',
                border: msg.role === 'assistant' ? '1px solid #d9d9d9' : 'none',
                position: 'relative'
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', marginBottom: '4px' }}>
                {msg.role === 'user' ? <UserOutlined /> : <RobotOutlined />}
                <Text
                  style={{
                    fontSize: '12px',
                    marginLeft: '6px',
                    color: msg.role === 'user' ? 'rgba(255,255,255,0.8)' : '#999'
                  }}
                >
                  {msg.role === 'user' ? '您' : 'AI助手'}
                </Text>
              </div>
              <div style={{ whiteSpace: 'pre-line' }}>{msg.content}</div>
              
              {/* 如果消息包含分析更新，显示提示 */}
              {msg.analysisUpdate && (
                <Alert
                  message="需求分析已更新"
                  type="success"
                  size="small"
                  style={{ marginTop: '8px' }}
                  showIcon
                />
              )}
            </div>
          </div>
        ))}
        
        {isChatLoading && (
          <div style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: '12px' }}>
            <div
              style={{
                padding: '8px 12px',
                borderRadius: '8px',
                backgroundColor: '#fff',
                border: '1px solid #d9d9d9'
              }}
            >
              <LoadingOutlined /> AI正在思考...
            </div>
          </div>
        )}
      </div>

      <Space.Compact style={{ width: '100%' }}>
        <TextArea
          value={currentMessage}
          onChange={(e) => setCurrentMessage(e.target.value)}
          onKeyDown={handleKeyPress}
          placeholder="输入您的问题或要求..."
          autoSize={{ minRows: 1, maxRows: 3 }}
          disabled={isChatLoading}
        />
        <Button
          type="primary"
          icon={<SendOutlined />}
          onClick={handleChatMessage}
          disabled={isChatLoading || !currentMessage.trim()}
        >
          发送
        </Button>
      </Space.Compact>
      
      <Text type="secondary" style={{ fontSize: '12px', display: 'block', marginTop: '8px' }}>
        提示：按 Ctrl+Enter 快速发送消息
      </Text>
    </Card>
  );

  if (isLoading && step === 'input') {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <LoadingOutlined style={{ fontSize: '24px' }} />
        <div style={{ marginTop: '16px' }}>正在加载需求分析...</div>
      </div>
    );
  }

  return (
    <div>
      {step === 'input' && (
        <Card title="需求输入" size="small">
          <Space direction="vertical" style={{ width: '100%' }}>
            <TextArea
              value={requirement}
              onChange={(e) => setRequirement(e.target.value)}
              onKeyDown={handleKeyPress}
              placeholder="请详细描述您的项目需求..."
              rows={6}
              disabled={isLoading}
            />
            
            {error && (
              <Alert message={error} type="error" showIcon />
            )}
            
            <Button
              type="primary"
              onClick={handleInitialAnalysis}
              loading={isLoading}
              disabled={!requirement.trim()}
              block
            >
              开始智能分析
            </Button>
            
            <Text type="secondary" style={{ fontSize: '12px' }}>
              提示：按 Ctrl+Enter 快速开始分析
            </Text>
          </Space>
        </Card>
      )}

      {step === 'analyzing' && (
        <Card title="正在分析" size="small">
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <LoadingOutlined style={{ fontSize: '32px', color: '#1890ff' }} />
            <div style={{ marginTop: '16px', fontSize: '16px' }}>
              AI正在分析您的需求...
            </div>
            <Text type="secondary">
              这可能需要几秒钟时间，请耐心等待
            </Text>
          </div>
        </Card>
      )}

      {step === 'interactive' && (
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          {renderAnalysisResult()}
          <Divider />
          {renderChatInterface()}
        </Space>
      )}
    </div>
  );
};

export default RequirementAnalysis;