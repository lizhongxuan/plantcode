import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Layout, 
  Row, 
  Col, 
  Card, 
  Button, 
  Typography, 
  Space, 
  Breadcrumb,
  FloatButton,
  Badge,
  Progress,
  Tabs,
  message,
  Modal
} from 'antd';
import { 
  ArrowLeftOutlined, 
  MessageOutlined, 
  HistoryOutlined,
  EyeOutlined,
  EditOutlined
} from '@ant-design/icons';
import { 
  Project, 
  ProjectDocument, 
  ProjectStructureResponse,
  DocumentChangesResponse 
} from '../types';
import ProjectFolderStructure from '../components/ProjectFolderStructure';
import PersistentAIChat from '../components/PersistentAIChat';
import DocumentViewer from '../components/DocumentViewer';
import DocumentChangeTracker from '../components/DocumentChangeTracker';

const { Header, Content, Sider } = Layout;
const { Title, Text, Paragraph } = Typography;
const { TabPane } = Tabs;

interface StageNavigation {
  name: string;
  label: string;
  description: string;
  icon: string;
  color: string;
  progress: number;
}

export const ProjectDetailPage: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();

  const [project, setProject] = useState<Project | null>(null);
  const [projectStructure, setProjectStructure] = useState<ProjectStructureResponse | null>(null);
  const [selectedDocument, setSelectedDocument] = useState<ProjectDocument | null>(null);
  const [aiChatVisible, setAIChatVisible] = useState(false);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('structure');

  // 三个阶段的导航信息
  const stageNavigation: StageNavigation[] = [
    {
      name: 'requirements',
      label: 'Requirements',
      description: 'Define what needs to be built',
      icon: '📋',
      color: '#1890ff',
      progress: 0
    },
    {
      name: 'design',
      label: 'Design',
      description: 'Plan how it will be built',
      icon: '🎨',
      color: '#52c41a',
      progress: 0
    },
    {
      name: 'tasks',
      label: 'Tasks',
      description: 'Break down into actionable items',
      icon: '✅',
      color: '#faad14',
      progress: 0
    }
  ];

  const [stageProgress, setStageProgress] = useState<StageNavigation[]>(stageNavigation);

  useEffect(() => {
    if (projectId) {
      loadProjectData();
    }
  }, [projectId]);

  const loadProjectData = async () => {
    setLoading(true);
    try {
      // 加载项目信息
      await Promise.all([
        loadProject(),
        loadProjectStructure()
      ]);
    } catch (error) {
      console.error('Failed to load project data:', error);
      message.error('Failed to load project data');
    } finally {
      setLoading(false);
    }
  };

  const loadProject = async () => {
    try {
      const response = await fetch(`/api/v1/projects/${projectId}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (response.ok) {
        const projectData: Project = await response.json();
        setProject(projectData);
      } else {
        throw new Error('Failed to load project');
      }
    } catch (error) {
      console.error('Failed to load project:', error);
      throw error;
    }
  };

  const loadProjectStructure = async () => {
    try {
      const response = await fetch(`/api/v1/projects/${projectId}/structure`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (response.ok) {
        const structureData: ProjectStructureResponse = await response.json();
        setProjectStructure(structureData);
        calculateStageProgress(structureData);
      } else {
        throw new Error('Failed to load project structure');
      }
    } catch (error) {
      console.error('Failed to load structure:', error);
      throw error;
    }
  };

  const calculateStageProgress = (structure: ProjectStructureResponse) => {
    const updatedStages = stageProgress.map(stage => {
      const folder = structure.folders.find(f => f.folder_name === stage.name);
      const documentCount = folder?.documents.length || 0;
      const progress = Math.min(documentCount * 25, 100); // 每个文档贡献25%进度，最多100%
      
      return {
        ...stage,
        progress
      };
    });

    setStageProgress(updatedStages);
  };

  const handleDocumentSelect = (document: ProjectDocument) => {
    setSelectedDocument(document);
    setActiveTab('viewer');
  };

  const handleDocumentChange = async (documentId: string, changes: any) => {
    // 处理AI对话引起的文档变更
    message.success('Document updated by AI');
    await loadProjectStructure(); // 重新加载结构以反映变更
    
    // 如果当前选中的文档被修改，重新加载该文档
    if (selectedDocument?.document_id === documentId) {
      // 重新加载选中的文档内容
      try {
        const response = await fetch(`/api/documents/${documentId}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        if (response.ok) {
          const updatedDoc: ProjectDocument = await response.json();
          setSelectedDocument(updatedDoc);
        }
      } catch (error) {
        console.error('Failed to reload document:', error);
      }
    }
  };

  const handleStructureChange = () => {
    loadProjectStructure();
  };

  const getOverallProgress = () => {
    const totalProgress = stageProgress.reduce((sum, stage) => sum + stage.progress, 0);
    return Math.round(totalProgress / stageProgress.length);
  };

  if (loading || !project) {
    return (
      <Layout className="min-h-screen">
        <Content className="p-6">
          <div className="flex justify-center items-center h-64">
            <div>Loading project...</div>
          </div>
        </Content>
      </Layout>
    );
  }

  return (
    <Layout className="min-h-screen bg-gray-50">
      <Header className="bg-white shadow-sm px-6">
        <div className="flex items-center justify-between h-full">
          <div className="flex items-center">
            <Button 
              type="text" 
              icon={<ArrowLeftOutlined />}
              onClick={() => navigate('/projects')}
              className="mr-4"
            />
            <Breadcrumb>
              <Breadcrumb.Item>Projects</Breadcrumb.Item>
              <Breadcrumb.Item>{project.project_name}</Breadcrumb.Item>
            </Breadcrumb>
          </div>
          <Space>
            {/* Settings button removed */}
          </Space>
        </div>
      </Header>

      <Content className="p-6">
        {/* Project Header */}
        <Card className="mb-6" size="small">
          <Row gutter={[24, 16]} align="middle">
            <Col span={16}>
              <div>
                <Title level={3} className="mb-2">{project.project_name}</Title>
                <Paragraph className="text-gray-600 mb-3">
                  {project.description}
                </Paragraph>
                <Space>
                  <Badge color={project.status === 'draft' ? 'orange' : 'green'} text={project.status} />
                  <Text type="secondary">Type: {project.project_type}</Text>
                  <Text type="secondary">Created: {new Date(project.created_at).toLocaleDateString()}</Text>
                </Space>
              </div>
            </Col>
            <Col span={8}>
              <div className="text-right">
                <div className="mb-2">
                  <Text strong>Overall Progress</Text>
                </div>
                <Progress 
                  percent={getOverallProgress()} 
                  strokeColor={{
                    '0%': '#108ee9',
                    '100%': '#87d068',
                  }}
                  className="mb-4"
                />
                
                {/* Stage Progress */}
                <div className="space-y-2">
                  {stageProgress.map(stage => (
                    <div key={stage.name} className="flex items-center justify-between">
                      <span className="text-sm">
                        {stage.icon} {stage.label}
                      </span>
                      <Progress 
                        percent={stage.progress} 
                        size="small" 
                        strokeColor={stage.color}
                        className="flex-1 mx-2"
                      />
                      <Text className="text-xs">{stage.progress}%</Text>
                    </div>
                  ))}
                </div>
              </div>
            </Col>
          </Row>
        </Card>

        {/* Stage Navigation */}
        <Row gutter={16} className="mb-6">
          {stageProgress.map(stage => (
            <Col span={8} key={stage.name}>
              <Card 
                size="small" 
                className="text-center cursor-pointer hover:shadow-md transition-shadow"
                style={{ borderColor: stage.color }}
                onClick={() => {
                  // 可以添加阶段切换逻辑
                  message.info(`Navigating to ${stage.label} stage`);
                }}
              >
                <div className="py-2">
                  <div className="text-2xl mb-2">{stage.icon}</div>
                  <Title level={5} className="mb-1">{stage.label}</Title>
                  <Text type="secondary" className="text-xs">{stage.description}</Text>
                  <Progress 
                    percent={stage.progress} 
                    strokeColor={stage.color} 
                    className="mt-2" 
                    size="small"
                  />
                </div>
              </Card>
            </Col>
          ))}
        </Row>

        {/* Main Content */}
        <Row gutter={16}>
          <Col span={6}>
            <ProjectFolderStructure
              project={project}
              onDocumentSelect={handleDocumentSelect}
              onStructureChange={handleStructureChange}
            />
          </Col>
          
          <Col span={18}>
            <Card>
              <Tabs activeKey={activeTab} onChange={setActiveTab}>
                <TabPane tab="Document Structure" key="structure">
                  <div className="p-4 text-center text-gray-500">
                    <FileOutlined style={{ fontSize: '48px' }} className="mb-4" />
                    <div>Select a document from the project structure to view its content</div>
                  </div>
                </TabPane>
                
                <TabPane 
                  tab={
                    <span>
                      <EyeOutlined />
                      Document Viewer
                      {selectedDocument && <Badge dot className="ml-2" />}
                    </span>
                  } 
                  key="viewer"
                  disabled={!selectedDocument}
                >
                  {selectedDocument && (
                    <DocumentViewer
                      document={selectedDocument}
                      onEdit={(doc) => {
                        // 实现编辑功能
                        message.info('Edit mode coming soon');
                      }}
                    />
                  )}
                </TabPane>

                <TabPane 
                  tab={
                    <span>
                      <HistoryOutlined />
                      Change History
                    </span>
                  } 
                  key="history"
                  disabled={!selectedDocument}
                >
                  {selectedDocument && (
                    <DocumentChangeTracker documentId={selectedDocument.document_id} />
                  )}
                </TabPane>
              </Tabs>
            </Card>
          </Col>
        </Row>
      </Content>

      {/* Floating AI Chat Button */}
      <FloatButton
        icon={<MessageOutlined />}
        type="primary"
        style={{ right: 24, bottom: 24 }}
        onClick={() => setAIChatVisible(true)}
        badge={{ count: aiChatVisible ? 0 : '!' }}
      />

      {/* AI Chat Drawer */}
      <PersistentAIChat
        project={project}
        projectStructure={projectStructure}
        visible={aiChatVisible}
        onClose={() => setAIChatVisible(false)}
        onDocumentChange={handleDocumentChange}
      />
    </Layout>
  );
};

export default ProjectDetailPage;