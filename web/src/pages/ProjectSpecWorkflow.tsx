import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Layout, 
  Card, 
  Button, 
  Typography, 
  Steps, 
  Space, 
  message, 
  Spin, 
  Alert,
  Tree,
  Input,
  Modal
} from 'antd';
import { 
  ArrowLeftOutlined,
  RobotOutlined,
  FileTextOutlined,
  CheckSquareOutlined,
  FolderOutlined,
  FileOutlined,
  ZoomInOutlined,
  ZoomOutOutlined,
  BorderOuterOutlined,
  FullscreenOutlined,
  FullscreenExitOutlined
} from '@ant-design/icons';
import { projectApi, aiApi } from '@/services/api';
import type { Project, ProjectDocument } from '@/types';

const { Content } = Layout;
const { Text } = Typography;
const { TextArea } = Input;

enum SpecStage {
  REQUIREMENTS = 1,
  DESIGN = 2,
  TASKS = 3
}

const ProjectSpecWorkflow: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  
  const [project, setProject] = useState<Project | null>(null);
  const [currentStage, setCurrentStage] = useState<SpecStage>(SpecStage.REQUIREMENTS);
  const [selectedDocument, setSelectedDocument] = useState<ProjectDocument | null>(null);
  const [documentContent, setDocumentContent] = useState<string>('');
  const [editMode, setEditMode] = useState(false);
  const [isPumlFile, setIsPumlFile] = useState(false);
  
  // PUML预览相关状态
  const [pumlZoom, setPumlZoom] = useState(100); // 缩放百分比
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  // AI聊天状态
  const [chatMessages, setChatMessages] = useState<any[]>([]);
  const [inputText, setInputText] = useState('');
  const chatContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (projectId) {
      loadProjectData();
    }
  }, [projectId]);

  // 禁用页面滚动
  useEffect(() => {
    // 禁用页面滚动
    document.body.style.overflow = 'hidden';
    document.documentElement.style.overflow = 'hidden';
    
    return () => {
      // 恢复页面滚动
      document.body.style.overflow = '';
      document.documentElement.style.overflow = '';
    };
  }, []);

  // ESC键退出全屏
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isFullscreen) {
        setIsFullscreen(false);
      }
    };
    
    document.addEventListener('keydown', handleKeyPress);
    return () => {
      document.removeEventListener('keydown', handleKeyPress);
    };
  }, [isFullscreen]);

  // 自动滚动到最新消息
  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [chatMessages]);

  const loadProjectData = async () => {
    try {
      setLoading(true);
      setError('');
      
      if (!projectId) {
        throw new Error('项目ID不能为空');
      }
      
      console.log('Loading project:', projectId);
      
      // 获取项目信息
      const projectData = await projectApi.getProject(projectId);
      console.log('Project data loaded:', projectData);
      setProject(projectData);

      // 模拟项目结构数据
      // setProjectStructure({
      //   folders: [
      //     {
      //       folder_id: '1',
      //       folder_name: 'requirements',
      //       documents: []
      //     },
      //     {
      //       folder_id: '2', 
      //       folder_name: 'design',
      //       documents: []
      //     },
      //     {
      //       folder_id: '3',
      //       folder_name: 'tasks', 
      //       documents: []
      //     }
      //   ]
      // } as ProjectStructureResponse);
      
    } catch (err) {
      console.error('Error loading project data:', err);
      const errorMessage = err instanceof Error ? err.message : '加载项目数据失败';
      setError(errorMessage);
      console.log('Setting error:', errorMessage);
    } finally {
      setLoading(false);
    }
  };


  const handleGenerateTodos = async () => {
    try {
      await projectApi.generateTodos(projectId!);
      message.success('TODO文档生成中，请稍候...');
      setCurrentStage(SpecStage.TASKS);
    } catch (error: any) {
      message.error(error.message || '生成TODO文档失败');
    }
  };

  // AI对话处理 - 支持真实的文档生成
  const handleSendMessage = async () => {
    if (!inputText.trim()) return;
    
    const newMessage = {
      id: Date.now(),
      role: 'user',
      content: inputText,
      timestamp: new Date().toLocaleTimeString()
    };
    
    setChatMessages(prev => [...prev, newMessage]);
    const messageToSend = inputText;
    setInputText('');
    
    try {
      // 调用真实的AI API
      const response = await aiApi.projectChat(
        projectId!, 
        messageToSend, 
        `当前阶段: ${currentStage === 1 ? 'requirements' : currentStage === 2 ? 'design' : 'tasks'}`
      );
      
      const aiReply = {
        id: Date.now() + 1,
        role: 'assistant',
        content: response.message || '我已经处理了您的请求。',
        timestamp: new Date().toLocaleTimeString()
      };
      
      setChatMessages(prev => [...prev, aiReply]);
      
      // 如果AI生成了文档，刷新项目结构
      if (response.documents_generated) {
        await loadProjectData();
        message.success('已生成新的文档内容');
      }
      
    } catch (error: any) {
      const errorReply = {
        id: Date.now() + 1,
        role: 'assistant', 
        content: '抱歉，处理您的请求时出现了错误。请稍后再试。',
        timestamp: new Date().toLocaleTimeString()
      };
      setChatMessages(prev => [...prev, errorReply]);
      console.error('AI chat error:', error);
    }
  };

  // PUML缩放控制函数
  const handleZoomIn = () => {
    setPumlZoom(prev => Math.min(prev + 25, 300)); // 最大300%
  };
  
  const handleZoomOut = () => {
    setPumlZoom(prev => Math.max(prev - 25, 25)); // 最小25%
  };
  
  const handleFitToWindow = () => {
    setPumlZoom(100); // 重置为100%
  };
  
  const handleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };
  
  // 鼠标滚轮缩放
  const handleWheel = (e: React.WheelEvent) => {
    if (e.ctrlKey || e.metaKey) {
      e.preventDefault();
      const delta = e.deltaY > 0 ? -10 : 10;
      setPumlZoom(prev => Math.max(25, Math.min(300, prev + delta)));
    }
  };

  // PUML预览URL生成
  const getPumlPreviewUrl = (content: string) => {
    try {
      // 使用PlantUML的压缩编码
      const encoded = btoa(unescape(encodeURIComponent(content)));
      return `https://www.plantuml.com/plantuml/svg/~1${encoded}`;
    } catch (error) {
      console.error('PUML encoding error:', error);
      return `https://www.plantuml.com/plantuml/svg/~1${btoa('@startuml\n!theme plain\ntitle 编码错误\n[错误] --> [请检查PUML语法]\n@enduml')}`;
    }
  };

  // 保存文档内容
  const handleSaveDocument = async () => {
    if (!selectedDocument || !projectId) return;
    
    try {
      // 这里应该调用API保存文档内容
      // await documentApi.saveDocument(selectedDocument.document_id, documentContent);
      
      // 临时模拟保存成功
      message.success(`${selectedDocument.document_name} 保存成功`);
      
      // 更新选中文档的内容
      setSelectedDocument({
        ...selectedDocument,
        content: documentContent,
        updated_at: new Date().toISOString()
      });
      
    } catch (error: any) {
      console.error('Save document error:', error);
      message.error('保存文档失败');
    }
  };

  // 处理文件选择
  const handleFileSelect = async (selectedKeys: any[]) => {
    const key = selectedKeys[0];
    if (!key || key.includes('folder')) return;
    
    // 检测是否为PUML文件
    const isPuml = key.endsWith('.puml') || key.endsWith('.plantuml');
    setIsPumlFile(isPuml);
    
    // 模拟加载文档内容
    try {
      let content = '';
      switch(key) {
        case 'requirements.md':
          content = '# 项目需求文档\n\n## 功能需求\n\n请通过AI对话来完善此文档的内容。';
          break;
        case 'architecture.puml':
          content = '@startuml\n!theme plain\ntitle 系统架构图\n\n[用户] --> [前端应用]\n[前端应用] --> [后端API]\n[后端API] --> [数据库]\n\n@enduml';
          break;
        case 'sequence.puml':
          content = '@startuml\n!theme plain\ntitle 时序图\n\nactor 用户\nparticipant 前端\nparticipant 后端\ndatabase 数据库\n\n用户 -> 前端: 发送请求\n前端 -> 后端: API调用\n后端 -> 数据库: 查询数据\n数据库 --> 后端: 返回结果\n后端 --> 前端: 返回响应\n前端 --> 用户: 显示结果\n\n@enduml';
          break;
        case 'class-diagram.puml':
          content = '@startuml\n!theme plain\ntitle 类图\n\nclass User {\n  +String name\n  +String email\n  +login()\n  +logout()\n}\n\nclass Project {\n  +String name\n  +String description\n  +create()\n  +update()\n}\n\nUser ||--o{ Project : owns\n\n@enduml';
          break;
        case 'task-list.md':
          content = '# 项目任务列表\n\n## 开发任务\n\n- [ ] 任务1\n- [ ] 任务2\n\n请通过AI对话来生成详细的任务列表。';
          break;
        default:
          if (isPuml) {
            content = '@startuml\n!theme plain\ntitle 图表标题\n\n[组件A] --> [组件B]\n\n@enduml';
          } else {
            content = '文件内容正在加载...';
          }
      }
      
      setDocumentContent(content);
      setSelectedDocument({
        document_id: key,
        document_name: key,
        document_type: isPuml ? 'puml' : 'markdown',
        content: content,
        version: 1,
        project_id: projectId || '',
        folder_id: '',
        is_template: false,
        created_by: '',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      });
      
      // 重置缩放和全屏状态
      setPumlZoom(100);
      setIsFullscreen(false);
      
    } catch (error) {
      console.error('Error loading document:', error);
      message.error('加载文档失败');
    }
  };

  const getStageContent = () => {
    // 如果选择了文档，显示编辑器
    if (selectedDocument) {
      // PUML文件显示在线编辑器
      if (isPumlFile) {
        return (
          <Card 
            title={
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span>{selectedDocument.document_name}</span>
                <Space>
                  <Button 
                    type="text" 
                    onClick={() => setEditMode(!editMode)}
                  >
                    {editMode ? '预览模式' : '编辑模式'}
                  </Button>
                  <Button type="primary" size="small" onClick={handleSaveDocument}>
                    保存
                  </Button>
                </Space>
              </div>
            } 
            className="h-full"
            bodyStyle={{ padding: 0, height: 'calc(100vh - 4rem - 48px - 200px)', overflow: 'hidden' }}
          >
            <div style={{ display: 'flex', height: '100%' }}>
              {/* 左侧编辑区 */}
              <div style={{ 
                width: editMode ? '50%' : '100%', 
                borderRight: editMode ? '1px solid #e8e8e8' : 'none',
                display: 'flex',
                flexDirection: 'column'
              }}>
                <div style={{ 
                  padding: '8px 12px', 
                  background: '#f5f5f5', 
                  borderBottom: '1px solid #e8e8e8',
                  fontSize: '12px',
                  color: '#666'
                }}>
                  PUML编辑器
                </div>
                <TextArea
                  value={documentContent}
                  onChange={(e) => setDocumentContent(e.target.value)}
                  onKeyDown={(e) => {
                    // Ctrl+S 保存
                    if (e.ctrlKey && e.key === 's') {
                      e.preventDefault();
                      handleSaveDocument();
                    }
                  }}
                  style={{ 
                    flex: 1,
                    border: 'none',
                    borderRadius: 0,
                    fontFamily: 'Monaco, Consolas, "Courier New", monospace',
                    fontSize: '13px',
                    resize: 'none'
                  }}
                  placeholder="输入PlantUML代码...&#10;&#10;提示：&#10;- 使用 @startuml 开始，@enduml 结束&#10;- Ctrl+S 快速保存&#10;- 切换到编辑模式查看实时预览&#10;- Ctrl+滚轮缩放预览"
                />
              </div>
              
              {/* 右侧预览区 */}
              {editMode && (
                <div style={{ 
                  width: '50%',
                  display: 'flex',
                  flexDirection: 'column',
                  background: '#fff'
                }}>
                  <div style={{ 
                    padding: '8px 12px', 
                    background: '#f5f5f5', 
                    borderBottom: '1px solid #e8e8e8',
                    fontSize: '12px',
                    color: '#666',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center'
                  }}>
                    <span>实时预览 ({pumlZoom}%)</span>
                    <Space size="small">
                      <Button 
                        type="text" 
                        size="small" 
                        icon={<ZoomOutOutlined />}
                        onClick={handleZoomOut}
                        disabled={pumlZoom <= 25}
                        title="缩小"
                      />
                      <Button 
                        type="text" 
                        size="small" 
                        icon={<BorderOuterOutlined />}
                        onClick={handleFitToWindow}
                        title="适应窗口"
                      />
                      <Button 
                        type="text" 
                        size="small" 
                        icon={<ZoomInOutlined />}
                        onClick={handleZoomIn}
                        disabled={pumlZoom >= 300}
                        title="放大"
                      />
                      <Button 
                        type="text" 
                        size="small" 
                        icon={<FullscreenOutlined />}
                        onClick={handleFullscreen}
                        title="全屏预览"
                      />
                    </Space>
                  </div>
                  <div 
                    style={{ 
                      flex: 1, 
                      padding: '16px',
                      overflow: 'auto',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      background: '#fafafa'
                    }}
                    onWheel={handleWheel}
                  >
                    <div style={{ 
                      transform: `scale(${pumlZoom / 100})`,
                      transformOrigin: 'center center',
                      transition: 'transform 0.2s ease'
                    }}>
                      <iframe
                        src={getPumlPreviewUrl(documentContent)}
                        style={{ 
                          width: '600px', 
                          height: '400px', 
                          border: '1px solid #e8e8e8',
                          borderRadius: '4px',
                          background: '#fff'
                        }}
                        title="PlantUML预览"
                        onError={(e) => {
                          console.error('PUML preview error:', e);
                        }}
                      />
                    </div>
                  </div>
                </div>
              )}
            </div>
          </Card>
        );
      }
      
      // 普通文档显示编辑器
      return (
        <Card 
          title={
            <div style={{ display: 'flex', justifyContent: 'between', alignItems: 'center' }}>
              <span>{selectedDocument.document_name}</span>
              <Button 
                type="text" 
                onClick={() => setEditMode(!editMode)}
                style={{ marginLeft: 'auto' }}
              >
                {editMode ? '预览' : '编辑'}
              </Button>
            </div>
          } 
          className="h-full"
          bodyStyle={{ padding: '0', height: 'calc(100vh - 4rem - 48px - 200px)', overflow: 'hidden', display: 'flex', flexDirection: 'column' }}
        >
          {editMode ? (
            <div style={{ flex: 1, padding: '16px', overflow: 'hidden' }}>
              <TextArea
                value={documentContent}
                onChange={(e) => setDocumentContent(e.target.value)}
                style={{ 
                  height: '100%', 
                  fontFamily: 'Monaco, Consolas, monospace',
                  fontSize: '13px',
                  border: 'none',
                  resize: 'none'
                }}
                placeholder="编辑文档内容..."
              />
            </div>
          ) : (
            <div style={{ 
              flex: 1,
              margin: '16px',
              background: '#fafafa',
              border: '1px solid #e8e8e8',
              borderRadius: '4px',
              fontFamily: 'Monaco, Consolas, monospace',
              fontSize: '13px',
              whiteSpace: 'pre-wrap',
              overflow: 'auto',
              padding: '16px'
            }}>
              {documentContent || '暂无内容，请通过AI对话生成内容'}
            </div>
          )}
        </Card>
      );
    }

    // 默认阶段内容
    switch (currentStage) {
      case SpecStage.REQUIREMENTS:
        return (
          <Card 
            title="第一阶段：需求确定" 
            className="h-full"
            bodyStyle={{ height: 'calc(100vh - 4rem - 48px - 200px)', overflow: 'auto' }}
          >
            <Space direction="vertical" style={{ width: '100%' }} size="middle">
              <Text>
                通过与AI对话确定项目需求，生成需求文档。AI会实时更新需求文档的内容。
              </Text>
              
            </Space>
          </Card>
        );

      case SpecStage.DESIGN:
        return (
          <Card 
            title="第二阶段：项目设计" 
            className="h-full"
            bodyStyle={{ height: 'calc(100vh - 4rem - 48px - 200px)', overflow: 'auto' }}
          >
            <Space direction="vertical" style={{ width: '100%' }} size="middle">
              <Text>
                基于需求文档生成项目架构PUML图和项目逻辑交互PUML图。
              </Text>
              
              <Button 
                type="primary"
                icon={<CheckSquareOutlined />}
                onClick={handleGenerateTodos}
                size="large"
              >
                生成TODO文档
              </Button>
            </Space>
          </Card>
        );

      case SpecStage.TASKS:
        return (
          <Card 
            title="第三阶段：任务分解" 
            className="h-full"
            bodyStyle={{ height: 'calc(100vh - 4rem - 48px - 200px)', overflow: 'auto' }}
          >
            <Space direction="vertical" style={{ width: '100%' }} size="middle">
              <Text>
                根据前两个阶段的文档生成项目TODO执行步骤列表文档。
              </Text>
            </Space>
          </Card>
        );

      default:
        return null;
    }
  };

  const treeData = [
    {
      title: '📋 Requirements',
      key: 'requirements',
      icon: <FolderOutlined />,
      children: [
        {
          title: 'requirements.md',
          key: 'requirements.md',
          icon: <FileOutlined />
        }
      ]
    },
    {
      title: '🎨 Design', 
      key: 'design',
      icon: <FolderOutlined />,
      children: [
        {
          title: 'architecture.puml',
          key: 'architecture.puml', 
          icon: <FileOutlined />
        },
        {
          title: 'sequence.puml',
          key: 'sequence.puml',
          icon: <FileOutlined />
        },
        {
          title: 'class-diagram.puml',
          key: 'class-diagram.puml',
          icon: <FileOutlined />
        }
      ]
    },
    {
      title: '✅ Tasks',
      key: 'tasks',
      icon: <FolderOutlined />,
      children: [
        {
          title: 'task-list.md',
          key: 'task-list.md',
          icon: <FileOutlined />
        }
      ]
    }
  ];

  if (loading) {
    return (
      <div style={{ 
        height: '100vh', 
        background: '#1e1e1e', 
        color: '#d4d4d4',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        flexDirection: 'column'
      }}>
        <Spin size="large" />
        <div style={{ marginTop: '16px', fontSize: '14px' }}>
          正在加载项目数据...
        </div>
      </div>
    );
  }

  if (error || !project) {
    return (
      <div style={{ 
        height: '100vh', 
        background: '#1e1e1e', 
        color: '#d4d4d4',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '20px'
      }}>
        <Alert
          message="加载失败"
          description={error || '项目不存在'}
          type="error"
          showIcon
          style={{
            background: '#2d2d30',
            border: '1px solid #3e3e42',
            color: '#d4d4d4'
          }}
          action={
            <Button onClick={() => navigate('/projects')}>
              返回项目列表
            </Button>
          }
        />
      </div>
    );
  }

  return (
    <>
      <Layout style={{ 
        height: 'calc(100vh - 4rem - 48px)', // 减去Header 4rem + padding 48px (p-6 = 24px * 2)
        width: '100%',
        background: '#f5f5f5', 
        overflow: 'hidden'
      }}>
      {/* 左侧文件树 */}
      <Layout.Sider 
        width={280} 
        theme="light"
        style={{ 
          background: '#fafafa',
          borderRight: '1px solid #e8e8e8',
          display: 'flex',
          flexDirection: 'column',
          height: 'calc(100vh - 4rem - 48px)',
          overflow: 'hidden'
        }}
      >
        {/* 项目头部 */}
        <div style={{ 
          padding: '12px 16px', 
          borderBottom: '1px solid #e8e8e8',
          background: '#fff'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
            <Button 
              type="text" 
              icon={<ArrowLeftOutlined />} 
              onClick={() => navigate('/projects')}
              style={{ marginRight: '8px' }}
            />
            <Text style={{ fontSize: '14px', fontWeight: 500 }}>
              {project.project_name}
            </Text>
          </div>
          
          {/* 阶段进度 */}
          <Steps 
            direction="vertical" 
            size="small" 
            current={currentStage - 1}
            style={{ fontSize: '12px' }}
            items={[
              {
                title: <span style={{ fontSize: '12px' }}>需求确定</span>,
                icon: <FileTextOutlined style={{ color: currentStage >= 1 ? '#1890ff' : '#999' }} />,
              },
              {
                title: <span style={{ fontSize: '12px' }}>项目设计</span>,
                icon: <FileTextOutlined style={{ color: currentStage >= 2 ? '#1890ff' : '#999' }} />,
              },
              {
                title: <span style={{ fontSize: '12px' }}>任务分解</span>,
                icon: <CheckSquareOutlined style={{ color: currentStage >= 3 ? '#1890ff' : '#999' }} />,
              }
            ]}
          />
        </div>

        {/* 文件资源管理器 */}
        <div style={{ flex: 1, padding: '8px', overflow: 'auto' }}>
          <div style={{ 
            color: '#666', 
            fontSize: '11px', 
            fontWeight: 600, 
            marginBottom: '8px',
            textTransform: 'uppercase',
            letterSpacing: '0.5px'
          }}>
            EXPLORER
          </div>
          <Tree
            showIcon
            defaultExpandAll
            treeData={treeData}
            onSelect={handleFileSelect}
            style={{ 
              background: 'transparent'
            }}
          />
        </div>
      </Layout.Sider>

      <Layout style={{ 
        display: 'flex', 
        flexDirection: 'column', 
        flex: 1,
        height: 'calc(100vh - 4rem - 48px)',
        overflow: 'hidden'
      }}>
        {/* 标签页头部 */}
        <div style={{ 
          height: '40px', 
          background: '#fff', 
          borderBottom: '1px solid #e8e8e8',
          display: 'flex',
          alignItems: 'center',
          padding: '0 12px'
        }}>
          <div style={{ 
            background: '#f0f0f0',
            color: '#333',
            padding: '6px 16px',
            fontSize: '13px',
            borderRadius: '4px',
            marginRight: '2px'
          }}>
            {currentStage === 1 ? 'Requirements' : 
             currentStage === 2 ? 'Design' : 'Tasks'}
          </div>
        </div>

        {/* 主要内容 */}
        <Content style={{ 
          background: '#fff',
          padding: '24px',
          margin: '0',
          flex: 1,
          overflow: 'hidden'
        }}>
          {getStageContent()}
        </Content>
      </Layout>

      {/* 右侧AI助手面板 */}
      <div
        style={{ 
          width: '350px',
          background: '#fafafa',
          borderLeft: '1px solid #e8e8e8',
          display: 'flex',
          flexDirection: 'column',
          height: 'calc(100vh - 4rem - 48px)',
          overflow: 'hidden'
        }}
      >
        {/* AI助手头部 */}
        <div style={{ 
          height: '40px',
          background: '#fff', 
          borderBottom: '1px solid #e8e8e8',
          display: 'flex',
          alignItems: 'center',
          padding: '0 12px'
        }}>
          <RobotOutlined style={{ color: '#1890ff', marginRight: '8px' }} />
          <Text style={{ fontSize: '13px', fontWeight: 500 }}>
            Kiro AI Assistant
          </Text>
        </div>

        {/* 对话历史 */}
        <div 
          ref={chatContainerRef}
          style={{ 
            flex: 1,
            padding: '12px',
            overflow: 'auto',
            minHeight: 0
          }}
        >
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            gap: '12px',
            minHeight: '100%',
            justifyContent: chatMessages.length === 0 ? 'center' : 'flex-start'
          }}>
          {chatMessages.length === 0 ? (
            <div style={{ 
              textAlign: 'center', 
              color: '#999', 
              fontSize: '13px',
              marginTop: '50px'
            }}>
              <RobotOutlined style={{ fontSize: '32px', marginBottom: '12px', color: '#1890ff' }} />
              <div>Hi! I'm Kiro, your AI assistant.</div>
              <div>How can I help you with your project today?</div>
            </div>
          ) : (
            chatMessages.map((msg) => (
              <div key={msg.id} style={{ 
                display: 'flex', 
                flexDirection: 'column',
                alignItems: msg.role === 'user' ? 'flex-end' : 'flex-start'
              }}>
                <div style={{
                  background: msg.role === 'user' ? '#1890ff' : '#f5f5f5',
                  color: msg.role === 'user' ? '#fff' : '#333',
                  padding: '8px 12px',
                  borderRadius: '8px',
                  maxWidth: '85%',
                  fontSize: '13px',
                  lineHeight: '1.4',
                  border: msg.role === 'assistant' ? '1px solid #e8e8e8' : 'none',
                  wordWrap: 'break-word',
                  whiteSpace: 'pre-wrap',
                  overflowWrap: 'break-word'
                }}>
                  {msg.content}
                </div>
                <Text style={{ 
                  fontSize: '11px', 
                  color: '#999', 
                  marginTop: '4px',
                  alignSelf: msg.role === 'user' ? 'flex-end' : 'flex-start'
                }}>
                  {msg.timestamp}
                </Text>
              </div>
            ))
          )}
          </div>
        </div>

        {/* 输入区域 */}
        <div style={{ 
          padding: '12px',
          borderTop: '1px solid #e8e8e8',
          background: '#fff',
          flexShrink: 0
        }}>
          <div style={{ display: 'flex', gap: '8px', alignItems: 'flex-end' }}>
            <Input.TextArea
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder="Type your message..."
              autoSize={{ minRows: 1, maxRows: 3 }}
              onPressEnter={(e) => {
                if (!e.shiftKey) {
                  e.preventDefault();
                  handleSendMessage();
                }
              }}
              style={{ flex: 1 }}
            />
            <Button 
              type="primary"
              onClick={handleSendMessage}
              disabled={!inputText.trim()}
              style={{ flexShrink: 0 }}
            >
              Send
            </Button>
          </div>
        </div>
      </div>
    </Layout>
    
    {/* 全屏PUML预览Modal */}
    {isFullscreen && isPumlFile && selectedDocument && (
      <Modal
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>全屏预览 - {selectedDocument.document_name}</span>
            <Space>
              <span style={{ fontSize: '12px', color: '#666' }}>缩放: {pumlZoom}%</span>
              <Button 
                type="text" 
                size="small" 
                icon={<ZoomOutOutlined />}
                onClick={handleZoomOut}
                disabled={pumlZoom <= 25}
              />
              <Button 
                type="text" 
                size="small" 
                icon={<BorderOuterOutlined />}
                onClick={handleFitToWindow}
              />
              <Button 
                type="text" 
                size="small" 
                icon={<ZoomInOutlined />}
                onClick={handleZoomIn}
                disabled={pumlZoom >= 300}
              />
            </Space>
          </div>
        }
        open={isFullscreen}
        onCancel={handleFullscreen}
        footer={[
          <Button key="close" onClick={handleFullscreen} icon={<FullscreenExitOutlined />}>
            退出全屏
          </Button>
        ]}
        width="90vw"
        style={{ top: 20 }}
        bodyStyle={{ 
          height: '80vh', 
          overflow: 'auto',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: '#fafafa'
        }}
      >
        <div 
          style={{ 
            transform: `scale(${pumlZoom / 100})`,
            transformOrigin: 'center center',
            transition: 'transform 0.2s ease'
          }}
          onWheel={handleWheel}
        >
          <iframe
            src={getPumlPreviewUrl(documentContent)}
            style={{ 
              width: '800px', 
              height: '600px', 
              border: '1px solid #e8e8e8',
              borderRadius: '4px',
              background: '#fff'
            }}
            title="PlantUML全屏预览"
          />
        </div>
      </Modal>
    )}
    </>
  );
};

export default ProjectSpecWorkflow;