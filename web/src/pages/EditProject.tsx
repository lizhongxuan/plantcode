import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useProjectStore } from '@/store';
import { projectApi } from '@/services/api';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { ArrowLeft, Save, RefreshCw } from 'lucide-react';
import type { Project, ProjectUpdateRequest } from '@/types';
import { ProjectTypes } from '@/types';

const EditProject: React.FC = () => {
  const navigate = useNavigate();
  const { projectId } = useParams<{ projectId: string }>();
  const { updateProject } = useProjectStore();
  
  const [project, setProject] = useState<Project | null>(null);
  const [formData, setFormData] = useState<ProjectUpdateRequest>({
    project_name: '',
    description: '',
    project_type: ProjectTypes.WEB_APPLICATION,
    status: 'active'
  });
  
  const [errors, setErrors] = useState<Partial<ProjectUpdateRequest>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(true);
  const [apiError, setApiError] = useState<string>('');

  // 获取项目详情
  useEffect(() => {
    const fetchProject = async () => {
      if (!projectId) return;
      
      try {
        setIsFetching(true);
        const projectData = await projectApi.getProject(projectId);
        setProject(projectData);
        setFormData({
          project_name: projectData.project_name,
          description: projectData.description,
          project_type: projectData.project_type,
          status: projectData.status
        });
      } catch (err: any) {
        setApiError(err.message || '获取项目详情失败');
      } finally {
        setIsFetching(false);
      }
    };

    fetchProject();
  }, [projectId]);

  const validateForm = (): boolean => {
    const newErrors: Partial<ProjectUpdateRequest> = {};
    
    if (!formData.project_name?.trim()) {
      newErrors.project_name = '请输入项目名称';
    } else if (formData.project_name.length < 3) {
      newErrors.project_name = '项目名称至少需要3个字符';
    }
    
    if (!formData.description?.trim()) {
      newErrors.description = '请输入项目描述';
    } else if (formData.description.length < 10) {
      newErrors.description = '项目描述至少需要10个字符';
    }
    
    if (!formData.project_type) {
      newErrors.project_type = '请选择项目类型';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!projectId) return;
    
    setApiError('');
    
    if (!validateForm()) return;
    
    setIsLoading(true);
    
    try {
      const updatedProject = await projectApi.updateProject(projectId, formData);
      updateProject(updatedProject);
      navigate(`/projects/${projectId}`);
    } catch (error) {
      setApiError(error instanceof Error ? error.message : '更新项目失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // 清除对应字段的错误
    if (errors[name as keyof ProjectUpdateRequest]) {
      setErrors(prev => ({ ...prev, [name]: undefined }));
    }
  };

  const projectTypeOptions = [
    { value: ProjectTypes.WEB_APPLICATION, label: 'Web应用程序' },
    { value: ProjectTypes.MOBILE_APP, label: '移动应用' },
    { value: ProjectTypes.API_SERVICE, label: 'API服务' },
    { value: ProjectTypes.DATA_ANALYSIS, label: '数据分析' },
    { value: ProjectTypes.MACHINE_LEARNING, label: '机器学习' },
    { value: ProjectTypes.OTHER, label: '其他' },
  ];

  const statusOptions = [
    { value: 'draft', label: '草稿' },
    { value: 'active', label: '进行中' },
    { value: 'completed', label: '已完成' },
    { value: 'archived', label: '已归档' },
  ];

  if (isFetching) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="flex flex-col items-center space-y-4">
          <RefreshCw className="h-8 w-8 animate-spin text-blue-600" />
          <p className="text-gray-600">正在加载项目信息...</p>
        </div>
      </div>
    );
  }

  if (apiError && !project) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <p className="text-red-600 mb-4">{apiError}</p>
          <Button onClick={() => navigate('/projects')}>
            返回项目列表
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 页面标题 */}
      <div className="flex items-center space-x-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate(`/projects/${projectId}`)}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900">编辑项目</h1>
          <p className="text-gray-600 mt-1">
            {project?.project_name || '加载中...'}
          </p>
        </div>
      </div>

      {/* 表单 */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-medium text-gray-900">项目信息</h2>
        </div>
        
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="md:col-span-2">
              <Input
                label="项目名称"
                type="text"
                name="project_name"
                value={formData.project_name}
                onChange={handleChange}
                error={errors.project_name}
                placeholder="输入项目名称"
                required
              />
            </div>
            
            <div>
              <label className="form-label">项目类型</label>
              <select
                name="project_type"
                value={formData.project_type}
                onChange={handleChange}
                className="form-input"
                required
              >
                {projectTypeOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
              {errors.project_type && (
                <p className="mt-1 text-sm text-red-600">{errors.project_type}</p>
              )}
            </div>

            <div>
              <label className="form-label">项目状态</label>
              <select
                name="status"
                value={formData.status}
                onChange={handleChange}
                className="form-input"
                required
              >
                {statusOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>
            
            <div className="md:col-span-2">
              <label className="form-label">项目描述</label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleChange}
                rows={4}
                className="form-input"
                placeholder="详细描述您的项目需求、目标和功能..."
                required
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description}</p>
              )}
            </div>
          </div>

          {apiError && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-600 text-sm">{apiError}</p>
            </div>
          )}

          <div className="flex items-center justify-between pt-6 border-t border-gray-200">
            <Button
              type="button"
              variant="ghost"
              onClick={() => navigate(`/projects/${projectId}`)}
            >
              取消
            </Button>
            <Button
              type="submit"
              disabled={isLoading}
              className="flex items-center"
            >
              {isLoading ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  保存中...
                </>
              ) : (
                <>
                  <Save className="h-4 w-4 mr-2" />
                  保存更改
                </>
              )}
            </Button>
          </div>
        </form>
      </div>

      {/* 项目创建信息 */}
      {project && (
        <div className="bg-gray-50 rounded-lg p-4">
          <h3 className="text-sm font-medium text-gray-900 mb-2">项目信息</h3>
          <div className="grid grid-cols-2 gap-4 text-sm text-gray-600">
            <div>
              <span className="font-medium">创建时间：</span>
              {new Date(project.created_at).toLocaleString()}
            </div>
            <div>
              <span className="font-medium">最后更新：</span>
              {new Date(project.updated_at).toLocaleString()}
            </div>
            <div>
              <span className="font-medium">项目ID：</span>
              {project.project_id}
            </div>
            <div>
              <span className="font-medium">完成度：</span>
              {project.completion_percentage}%
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EditProject; 