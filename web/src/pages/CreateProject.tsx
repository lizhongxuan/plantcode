import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProjectStore } from '@/store';
import { projectApi } from '@/services/api';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { ArrowLeft, FolderPlus } from 'lucide-react';
import type { CreateProjectRequest } from '@/types';
import { ProjectTypes } from '@/types';

const CreateProject: React.FC = () => {
  const navigate = useNavigate();
  const { addProject } = useProjectStore();
  
  const [formData, setFormData] = useState<CreateProjectRequest>({
    project_name: '',
    description: '',
    project_type: ProjectTypes.WEB_APPLICATION,
  });
  
  const [errors, setErrors] = useState<Partial<CreateProjectRequest>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [apiError, setApiError] = useState<string>('');

  const validateForm = (): boolean => {
    const newErrors: Partial<CreateProjectRequest> = {};
    
    if (!formData.project_name.trim()) {
      newErrors.project_name = '请输入项目名称';
    } else if (formData.project_name.length < 3) {
      newErrors.project_name = '项目名称至少需要3个字符';
    }
    
    if (!formData.description.trim()) {
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
    setApiError('');
    
    if (!validateForm()) return;
    
    setIsLoading(true);
    
    try {
      const newProject = await projectApi.createProject(formData);
      addProject(newProject);
      navigate(`/projects/${newProject.project_id}`);
    } catch (error) {
      setApiError(error instanceof Error ? error.message : '创建项目失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // 清除对应字段的错误
    if (errors[name as keyof CreateProjectRequest]) {
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

  return (
    <div className="space-y-6">
      {/* 页面标题 */}
      <div className="flex items-center space-x-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate('/projects')}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900">新建项目</h1>
          <p className="text-gray-600 mt-1">创建一个新的AI开发项目</p>
        </div>
      </div>

      {/* 表单 */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-center space-x-2">
            <FolderPlus className="h-5 w-5 text-primary-600" />
            <h2 className="text-lg font-medium text-gray-900">项目信息</h2>
          </div>
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
            
            <div className="md:col-span-2">
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
            <div className="rounded-md bg-red-50 p-4">
              <p className="text-sm text-red-800">{apiError}</p>
            </div>
          )}

          <div className="flex items-center justify-end space-x-4 pt-4 border-t border-gray-200">
            <Button
              type="button"
              variant="outline"
              onClick={() => navigate('/projects')}
            >
              取消
            </Button>
            <Button
              type="submit"
              isLoading={isLoading}
              className="flex items-center"
            >
              <FolderPlus className="h-4 w-4 mr-2" />
              创建项目
            </Button>
          </div>
        </form>
      </div>

      {/* 项目类型说明 */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">项目类型说明</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-3">
            <div>
              <h4 className="font-medium text-gray-900">Web应用程序</h4>
              <p className="text-sm text-gray-600">基于Web技术的应用程序，包括前端和后端开发</p>
            </div>
            <div>
              <h4 className="font-medium text-gray-900">移动应用</h4>
              <p className="text-sm text-gray-600">iOS、Android或跨平台移动应用开发</p>
            </div>
            <div>
              <h4 className="font-medium text-gray-900">API服务</h4>
              <p className="text-sm text-gray-600">RESTful API、GraphQL或其他API服务开发</p>
            </div>
          </div>
          <div className="space-y-3">
            <div>
              <h4 className="font-medium text-gray-900">数据分析</h4>
              <p className="text-sm text-gray-600">数据处理、分析和可视化项目</p>
            </div>
            <div>
              <h4 className="font-medium text-gray-900">机器学习</h4>
              <p className="text-sm text-gray-600">AI模型训练、部署和推理服务</p>
            </div>
            <div>
              <h4 className="font-medium text-gray-900">其他</h4>
              <p className="text-sm text-gray-600">不属于以上类型的其他项目</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateProject; 