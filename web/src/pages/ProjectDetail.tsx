import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store';
import { projectApi } from '@/services/api';
import Button from '@/components/ui/Button';
import { 
  ArrowLeft, 
  FileText, 
  Settings, 
  Eye,
  Edit3,
  Download,
  RefreshCw,
  Target,
  Workflow,
  BookOpen
} from 'lucide-react';
import type { Project } from '@/types';

// 导入子组件
import RequirementAnalysis from './components/RequirementAnalysis';
import PUMLDiagrams from './components/PUMLDiagrams';
import TechnicalDocuments from './components/TechnicalDocuments';

const ProjectDetail: React.FC = () => {
  const navigate = useNavigate();
  const { projectId } = useParams<{ projectId: string }>();
  const { user } = useAuthStore();
  
  const [project, setProject] = useState<Project | null>(null);
  const [activeTab, setActiveTab] = useState('analysis');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasRequirementAnalysis, setHasRequirementAnalysis] = useState(false);

  // 检查项目是否已完成需求分析
  const checkRequirementAnalysisStatus = (project: Project): boolean => {
    // 检查条件：
    // 1. 项目描述不为空且有实际内容
    // 2. 项目状态不是草稿状态
    // 3. 完成度大于0
    return !!(
      project.description && 
      project.description.trim().length > 10 && // 至少有一些实际描述内容
      project.status !== 'draft' &&
      project.completion_percentage > 0
    );
  };

  // 获取项目详情
  useEffect(() => {
    const fetchProject = async () => {
      if (!projectId) return;
      
      try {
        setLoading(true);
        const projectData = await projectApi.getProject(projectId);
        setProject(projectData);
        
        // 检查需求分析状态
        const hasAnalysis = checkRequirementAnalysisStatus(projectData);
        setHasRequirementAnalysis(hasAnalysis);
        
        // 如果没有需求分析，强制切换到分析标签页
        if (!hasAnalysis && activeTab !== 'analysis') {
          setActiveTab('analysis');
        }
      } catch (err: any) {
        setError(err.message || '获取项目详情失败');
      } finally {
        setLoading(false);
      }
    };

    fetchProject();
  }, [projectId]);

  // 监听需求分析完成事件，更新状态
  const handleAnalysisComplete = () => {
    if (project) {
      setHasRequirementAnalysis(checkRequirementAnalysisStatus(project));
    }
  };

  const tabs = [
    {
      id: 'analysis',
      name: '智能需求分析',
      icon: <Target className="h-4 w-4" />,
      description: '通过AI对话深度分析项目需求，生成结构化需求文档',
      primary: true,
      enabled: true, // 需求分析始终可用
      requiresAnalysis: false
    },
    {
      id: 'puml',
      name: 'PUML图表',
      icon: <Workflow className="h-4 w-4" />,
      description: '基于需求分析自动生成架构图和流程图',
      dependsOn: 'analysis',
      enabled: hasRequirementAnalysis,
      requiresAnalysis: true
    },
    {
      id: 'documents',
      name: '技术文档',
      icon: <BookOpen className="h-4 w-4" />,
      description: '基于需求分析自动生成开发文档和技术规范',
      dependsOn: 'analysis',
      enabled: hasRequirementAnalysis,
      requiresAnalysis: true
    }
  ];

  // 处理标签页切换
  const handleTabClick = (tabId: string) => {
    const tab = tabs.find(t => t.id === tabId);
    if (tab && tab.enabled) {
      setActiveTab(tabId);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="flex flex-col items-center space-y-4">
          <RefreshCw className="h-8 w-8 animate-spin text-blue-600" />
          <p className="text-gray-600">正在加载项目详情...</p>
        </div>
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <p className="text-red-600 mb-4">{error || '项目不存在'}</p>
          <Button onClick={() => navigate('/projects')}>
            返回项目列表
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 项目头部 */}
      <div className="bg-white border border-gray-200 rounded-lg">
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => navigate('/projects')}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <div>
                <h1 className="text-2xl font-semibold text-gray-900">
                  {project.project_name}
                </h1>
                <p className="text-sm text-gray-500 mt-1">
                  {project.description || '暂无描述'}
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                project.status === 'active' 
                  ? 'bg-green-100 text-green-800'
                  : project.status === 'completed'
                  ? 'bg-blue-100 text-blue-800'
                  : 'bg-gray-100 text-gray-800'
              }`}>
                {project.status === 'active' ? '进行中' : 
                 project.status === 'completed' ? '已完成' : '暂停'}
              </span>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => navigate(`/projects/${projectId}/edit`)}
              >
                <Settings className="h-4 w-4" />
                项目设置
              </Button>
            </div>
          </div>
        </div>

        {/* 功能选项卡 */}
        <div className="border-t border-gray-200">
          <div className="px-6">
            <nav className="flex space-x-8">
              {tabs.map((tab) => (
                <div key={tab.id} className="relative group">
                  <button
                    onClick={() => handleTabClick(tab.id)}
                    disabled={!tab.enabled}
                    className={`flex items-center space-x-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                      activeTab === tab.id
                        ? 'border-blue-500 text-blue-600'
                        : tab.enabled
                        ? 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                        : 'border-transparent text-gray-400 cursor-not-allowed'
                    }`}
                  >
                    {tab.icon}
                    <span>{tab.name}</span>
                    {tab.primary && (
                      <span className="bg-blue-100 text-blue-800 text-xs px-2 py-0.5 rounded-full ml-2">
                        核心
                      </span>
                    )}
                    {!tab.enabled && tab.requiresAnalysis && (
                      <span className="bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded-full ml-2">
                        需要分析
                      </span>
                    )}
                  </button>
                  
                  {/* 禁用状态的提示信息 */}
                  {!tab.enabled && tab.requiresAnalysis && (
                    <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
                      请先完成需求分析才能使用此功能
                      <div className="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
                    </div>
                  )}
                </div>
              ))}
            </nav>
          </div>
        </div>
      </div>

      {/* 主内容区域 */}
      <div className="bg-white border border-gray-200 rounded-lg">
        {activeTab === 'analysis' && projectId && project && (
          <RequirementAnalysis 
            projectId={projectId} 
            project={project} 
            onAnalysisComplete={handleAnalysisComplete}
          />
        )}
        {activeTab === 'puml' && projectId && project && (
          <PUMLDiagrams projectId={projectId} project={project} />
        )}
        {activeTab === 'documents' && projectId && project && (
          <TechnicalDocuments projectId={projectId} project={project} />
        )}
      </div>
    </div>
  );
};

export default ProjectDetail; 