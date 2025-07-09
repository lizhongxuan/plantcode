import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Layout, 
  Menu, 
  Card, 
  Button, 
  Typography, 
  Space, 
  Spin, 
  Alert, 
  message,
  Modal,
  Drawer,
  Tabs,
  Divider
} from 'antd';
import { 
  ArrowLeftOutlined, 
  RobotOutlined, 
  FileTextOutlined,
  BarChartOutlined,
  SettingOutlined,
  PlayCircleOutlined,
  EyeOutlined
} from '@ant-design/icons';
import { projectApi, asyncTaskApi, pumlApi } from '@/services/api';
import RequirementAnalysis from '@/pages/components/RequirementAnalysis';
import AsyncTaskProgress from '@/components/AsyncTaskProgress';
import StageDocuments from '@/pages/components/StageDocuments';
import PUMLEditor from '@/pages/components/PUMLEditor';

const { Sider, Content } = Layout;
const { Title, Text } = Typography;
const { TabPane } = Tabs;

interface Project {
  project_id: string;
  project_name: string;
  description: string;
  project_type: string;
  status: string;
  completion_percentage: number;
  created_at: string;
  updated_at: string;
}

const ProjectDetail: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedMenu, setSelectedMenu] = useState('requirement-analysis');
  const [currentTaskId, setCurrentTaskId] = useState<string | null>(null);
  const [progressVisible, setProgressVisible] = useState(false);
  const [stageModalVisible, setStageModalVisible] = useState(false);
  const [selectedStage, setSelectedStage] = useState<number>(1);

  // 获取项目详情
  const fetchProject = async () => {
    if (!projectId) return;
    
    try {
      setLoading(true);
      // projectApi.getProject 直接返回 Project 对象，不是包含 success 字段的响应
      const projectData = await projectApi.getProject(projectId);
      setProject(projectData);
    } catch (error: any) {
      console.error('获取项目详情失败:', error);
      // 根据错误类型显示不同的消息
      if (error.message && error.message.includes('项目不存在')) {
        message.error('项目不存在或无权访问');
      } else {
        message.error(error.message || '获取项目信息失败');
      }
      // 确保设置 project 为 null，这样会显示 "项目不存在" 的 Alert
      setProject(null);
    } finally {
      setLoading(false);
    }
  };

  // 启动阶段文档生成
  const handleStartStageGeneration = async (stage: number) => {
    if (!projectId) return;
    
    try {
      const response = await asyncTaskApi.startStageDocumentGeneration(projectId, stage);
      if (response.success) {
        setCurrentTaskId(response.data.task_id);
        setProgressVisible(true);
        setStageModalVisible(false);
        message.success(`阶段${stage}文档生成任务已启动`);
      } else {
        message.error('启动文档生成失败');
      }
    } catch (error: any) {
      message.error(error.message || '启动文档生成失败');
    }
  };

  // 任务完成回调
  const handleTaskComplete = async (taskResult?: any) => {
    message.success('文档生成完成！');
    setCurrentTaskId(null);
    // 可以在这里刷新页面数据或导航到结果页面
  };

  // 任务失败回调
  const handleTaskFailed = (error?: string) => {
    message.error(`文档生成失败: ${error || '未知错误'}`);
    setCurrentTaskId(null);
  };

  // 启动完整项目文档生成
  const handleStartCompleteGeneration = async () => {
    if (!projectId) return;
    
    try {
      const response = await asyncTaskApi.startCompleteProjectDocumentGeneration(projectId);
      if (response.success) {
        setCurrentTaskId(response.data.task_id);
        setProgressVisible(true);
        message.success('完整项目文档生成任务已启动');
      } else {
        message.error('启动文档生成失败');
      }
    } catch (error: any) {
      message.error(error.message || '启动文档生成失败');
    }
  };

  // 菜单项配置
  const menuItems = [
    {
      key: 'requirement-analysis',
      icon: <RobotOutlined />,
      label: '智能需求分析',
      children: [
        {
          key: 'requirement-chat',
          label: 'AI需求助手',
        },
        {
          key: 'stage-1',
          label: '阶段一：需求梳理',
          children: [
            { key: 'stage-1-documents', label: '需求文档' },
            { key: 'stage-1-puml', label: 'PUML图表' },
          ]
        },
        {
          key: 'stage-2', 
          label: '阶段二：技术设计',
          children: [
            { key: 'stage-2-documents', label: '技术文档' },
            { key: 'stage-2-puml', label: 'PUML图表' },
          ]
        },
        {
          key: 'stage-3',
          label: '阶段三：实施计划', 
          children: [
            { key: 'stage-3-documents', label: '实施文档' },
            { key: 'stage-3-puml', label: 'PUML图表' },
          ]
        }
      ]
    },
    {
      key: 'progress-monitor',
      icon: <BarChartOutlined />,
      label: '进度监控',
    },
    {
      key: 'project-settings',
      icon: <SettingOutlined />,
      label: '项目设置',
    }
  ];

  // 阶段生成按钮
  const StageGenerationCard = ({ stage, title, description }: { stage: number, title: string, description: string }) => (
    <Card 
      size="small" 
      title={title}
      extra={
        <Button 
          type="primary" 
          icon={<PlayCircleOutlined />}
          onClick={() => {
            setSelectedStage(stage);
            setStageModalVisible(true);
          }}
        >
          生成
        </Button>
      }
    >
      <Text type="secondary">{description}</Text>
    </Card>
  );

  // 渲染内容区域
  const renderContent = () => {
    if (!project) return null;

    switch (selectedMenu) {
      case 'requirement-chat':
        return <RequirementAnalysis projectId={project.project_id} />;
      
      case 'stage-1':
      case 'stage-2': 
      case 'stage-3':
        const stage = parseInt(selectedMenu.split('-')[1]);
        return (
          <Space direction="vertical" style={{ width: '100%' }}>
            <StageGenerationCard 
              stage={stage}
              title={`阶段${stage}文档生成`}
              description={`生成阶段${stage}的相关文档和PUML图表`}
            />
          </Space>
        );
      
      case 'stage-1-documents':
      case 'stage-2-documents':
      case 'stage-3-documents':
        const docStage = parseInt(selectedMenu.split('-')[1]);
        return <StageDocuments projectId={project.project_id} stage={docStage} />;
        
      case 'stage-1-puml':
      case 'stage-2-puml': 
      case 'stage-3-puml':
        const pumlStage = parseInt(selectedMenu.split('-')[1]);
        return <PUMLEditor projectId={project.project_id} stage={pumlStage} />;
        
      case 'progress-monitor':
        return (
          <AsyncTaskProgress 
            projectId={project.project_id}
            onTaskComplete={handleTaskComplete}
            onTaskFailed={handleTaskFailed}
          />
        );
        
      case 'project-settings':
        return <div>项目设置 - 开发中</div>;
        
      default:
        return <RequirementAnalysis projectId={project.project_id} />;
    }
  };

  useEffect(() => {
    fetchProject();
  }, [projectId]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!project) {
    return (
      <Alert
        message="项目不存在"
        description="请检查项目ID是否正确"
        type="error"
        showIcon
      />
    );
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={280} theme="light" style={{ borderRight: '1px solid #f0f0f0' }}>
        <div style={{ padding: '16px' }}>
          <Space direction="vertical" style={{ width: '100%' }}>
            <Button 
              icon={<ArrowLeftOutlined />} 
              onClick={() => navigate('/projects')}
              style={{ marginBottom: '8px' }}
            >
              返回项目列表
            </Button>
            
            <Title level={4} ellipsis={{ tooltip: project.project_name }}>
              {project.project_name}
            </Title>
            
            <Text type="secondary" style={{ fontSize: '12px' }}>
              {project.description}
            </Text>

            <Divider style={{ margin: '12px 0' }} />
            
            <Button 
              type="primary" 
              icon={<PlayCircleOutlined />}
              onClick={handleStartCompleteGeneration}
              style={{ width: '100%', marginBottom: '12px' }}
              loading={currentTaskId !== null}
            >
              一键生成项目文档
            </Button>
            
            <Text type="secondary" style={{ fontSize: '11px', lineHeight: '1.4' }}>
              包含10份文档：需求、架构、流程、API、数据库、开发、测试、部署等完整文档
            </Text>
            
            <Divider style={{ margin: '12px 0' }} />
            
            {/* 进度快速查看按钮 */}
            <Button 
              type="dashed" 
              block
              icon={<EyeOutlined />}
              onClick={() => setProgressVisible(true)}
            >
              查看进度
            </Button>
          </Space>
        </div>
        
        <Menu
          mode="inline"
          selectedKeys={[selectedMenu]}
          items={menuItems}
          onClick={({ key }) => setSelectedMenu(key)}
          style={{ borderRight: 0, height: 'calc(100vh - 200px)', overflowY: 'auto' }}
        />
      </Sider>
      
      <Layout>
        <Content style={{ padding: '24px', background: '#fff' }}>
          {renderContent()}
        </Content>
      </Layout>

      {/* 进度监控抽屉 */}
      <Drawer
        title="任务进度监控"
        placement="right"
        width={500}
        onClose={() => setProgressVisible(false)}
        open={progressVisible}
      >
        <AsyncTaskProgress 
          taskId={currentTaskId || undefined}
          projectId={project.project_id}
          onTaskComplete={handleTaskComplete}
          onTaskFailed={handleTaskFailed}
        />
      </Drawer>

      {/* 阶段生成确认对话框 */}
      <Modal
        title={`确认生成阶段${selectedStage}文档`}
        open={stageModalVisible}
        onOk={() => handleStartStageGeneration(selectedStage)}
        onCancel={() => setStageModalVisible(false)}
        okText="开始生成"
        cancelText="取消"
      >
        <Space direction="vertical">
          <Text>即将为项目生成阶段{selectedStage}的文档和图表，这可能需要几分钟时间。</Text>
          <Text type="secondary">
            生成完成后，您可以在对应的菜单中查看和编辑生成的内容。
          </Text>
        </Space>
      </Modal>
    </Layout>
  );
};

export default ProjectDetail; 