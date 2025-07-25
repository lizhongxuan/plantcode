import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { FileText, Wand2, ArrowRight, ArrowLeft, CheckCircle, Clock, AlertCircle } from 'lucide-react';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import type { 
  ProjectSpec, 
  GenerateRequirementsRequest, 
  RequirementsDoc,
  SpecStage 
} from '@/types';

interface RequirementsStepProps {
  projectId: string;
  projectType: string;
  onNext: (requirementsId: string) => void;
  onBack?: () => void;
}

const RequirementsStep: React.FC<RequirementsStepProps> = ({ 
  projectId, 
  projectType, 
  onNext, 
  onBack 
}) => {
  const [formData, setFormData] = useState<Omit<GenerateRequirementsRequest, 'project_id'>>({
    initial_prompt: '',
    project_type: projectType,
    target_audience: '',
    business_goals: []
  });

  const [isGenerating, setIsGenerating] = useState(false);
  const [requirements, setRequirements] = useState<RequirementsDoc | null>(null);
  const [error, setError] = useState<string>('');
  const [businessGoalInput, setBusinessGoalInput] = useState('');

  const handleAddBusinessGoal = () => {
    if (businessGoalInput.trim()) {
      setFormData(prev => ({
        ...prev,
        business_goals: [...(prev.business_goals || []), businessGoalInput.trim()]
      }));
      setBusinessGoalInput('');
    }
  };

  const handleRemoveBusinessGoal = (index: number) => {
    setFormData(prev => ({
      ...prev,
      business_goals: prev.business_goals?.filter((_, i) => i !== index) || []
    }));
  };

  const handleGenerate = async () => {
    if (!formData.initial_prompt.trim()) {
      setError('请输入项目需求描述');
      return;
    }

    setIsGenerating(true);
    setError('');

    try {
      const requestData: GenerateRequirementsRequest = {
        project_id: projectId,
        initial_prompt: formData.initial_prompt,
        project_type: formData.project_type,
        target_audience: formData.target_audience || undefined,
        business_goals: formData.business_goals?.length ? formData.business_goals : undefined
      };

      const response = await fetch(`/api/projects/${projectId}/spec/requirements`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(requestData)
      });

      if (!response.ok) {
        throw new Error('生成需求文档失败');
      }

      const result = await response.json();
      
      if (result.success) {
        setRequirements(result.data);
      } else {
        throw new Error(result.message || '生成失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '生成需求文档时发生错误');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleNext = () => {
    if (requirements) {
      onNext(requirements.id);
    }
  };

  const handleEditRequirements = () => {
    setRequirements(null);
  };

  return (
    <div className="space-y-6">
      {/* 步骤标题 */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full">
            <span className="text-sm font-medium">1</span>
          </div>
          <div>
            <h2 className="text-xl font-semibold text-gray-900">需求分析 (Requirements)</h2>
            <p className="text-sm text-gray-600">使用 EARS 语法生成结构化需求文档</p>
          </div>
        </div>
        
        <div className="flex items-center space-x-2">
          {requirements && (
            <div className="flex items-center text-green-600">
              <CheckCircle className="h-4 w-4 mr-1" />
              <span className="text-sm">已完成</span>
            </div>
          )}
          {isGenerating && (
            <div className="flex items-center text-blue-600">
              <Clock className="h-4 w-4 mr-1 animate-spin" />
              <span className="text-sm">生成中...</span>
            </div>
          )}
        </div>
      </div>

      {!requirements ? (
        /* 需求输入表单 */
        <div className="bg-white rounded-lg shadow p-6">
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                项目需求描述 *
              </label>
              <textarea
                value={formData.initial_prompt}
                onChange={(e) => setFormData(prev => ({ ...prev, initial_prompt: e.target.value }))}
                placeholder="请详细描述您的项目需求、目标和核心功能..."
                rows={6}
                className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <p className="mt-1 text-xs text-gray-500">
                AI 将基于此描述生成详细的需求文档，请尽可能描述清楚
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Input
                  label="目标用户群体"
                  type="text"
                  value={formData.target_audience}
                  onChange={(e) => setFormData(prev => ({ ...prev, target_audience: e.target.value }))}
                  placeholder="如：个人用户、企业客户、开发者等"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  项目类型
                </label>
                <select
                  value={formData.project_type}
                  onChange={(e) => setFormData(prev => ({ ...prev, project_type: e.target.value }))}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="web_application">Web应用程序</option>
                  <option value="mobile_app">移动应用</option>
                  <option value="api_service">API服务</option>
                  <option value="data_analysis">数据分析</option>
                  <option value="machine_learning">机器学习</option>
                  <option value="other">其他</option>
                </select>
              </div>
            </div>

            {/* 业务目标输入 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                业务目标
              </label>
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={businessGoalInput}
                  onChange={(e) => setBusinessGoalInput(e.target.value)}
                  placeholder="输入一个业务目标"
                  className="flex-1 border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  onKeyPress={(e) => e.key === 'Enter' && handleAddBusinessGoal()}
                />
                <Button
                  type="button"
                  onClick={handleAddBusinessGoal}
                  size="sm"
                >
                  添加
                </Button>
              </div>

              {/* 已添加的业务目标 */}
              {formData.business_goals && formData.business_goals.length > 0 && (
                <div className="mt-2 space-y-1">
                  {formData.business_goals.map((goal, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between bg-blue-50 px-3 py-2 rounded-md"
                    >
                      <span className="text-sm text-blue-700">{goal}</span>
                      <button
                        onClick={() => handleRemoveBusinessGoal(index)}
                        className="text-blue-600 hover:text-blue-800"
                      >
                        ×
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {error && (
              <div className="flex items-center space-x-2 text-red-600 bg-red-50 p-3 rounded-md">
                <AlertCircle className="h-4 w-4" />
                <span className="text-sm">{error}</span>
              </div>
            )}

            <div className="flex justify-center">
              <Button
                onClick={handleGenerate}
                isLoading={isGenerating}
                disabled={!formData.initial_prompt.trim()}
                className="flex items-center"
              >
                <Wand2 className="h-4 w-4 mr-2" />
                {isGenerating ? '正在生成需求文档...' : '生成需求文档'}
              </Button>
            </div>
          </div>
        </div>
      ) : (
        /* 需求文档预览 */
        <div className="bg-white rounded-lg shadow">
          <div className="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <FileText className="h-5 w-5 text-blue-600" />
              <h3 className="text-lg font-medium text-gray-900">需求文档</h3>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleEditRequirements}
            >
              重新生成
            </Button>
          </div>

          <div className="p-6">
            <div className="prose max-w-none">
              <div
                dangerouslySetInnerHTML={{
                  __html: requirements.content.replace(/\n/g, '<br/>')
                }}
              />
            </div>

            {/* 用户故事预览 */}
            <div className="mt-8">
              <h4 className="text-lg font-medium text-gray-900 mb-4">用户故事</h4>
              <div className="grid gap-4">
                {requirements.user_stories?.map((story, index) => (
                  <div
                    key={story.id}
                    className="border border-gray-200 rounded-lg p-4"
                  >
                    <div className="flex items-center justify-between mb-2">
                      <h5 className="font-medium text-gray-900">{story.title}</h5>
                      <span className={`px-2 py-1 rounded text-xs ${
                        story.priority === 'high' ? 'bg-red-100 text-red-700' :
                        story.priority === 'medium' ? 'bg-yellow-100 text-yellow-700' :
                        'bg-green-100 text-green-700'
                      }`}>
                        {story.priority}
                      </span>
                    </div>
                    <p className="text-gray-700 text-sm mb-2">{story.description}</p>
                    {story.acceptance_criteria && story.acceptance_criteria.length > 0 && (
                      <div>
                        <p className="text-xs font-medium text-gray-600 mb-1">验收标准：</p>
                        <ul className="text-xs text-gray-600 space-y-1">
                          {story.acceptance_criteria.map((criteria, idx) => (
                            <li key={idx}>• {criteria}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* 导航按钮 */}
      <div className="flex items-center justify-between pt-6 border-t border-gray-200">
        <div>
          {onBack && (
            <Button
              variant="outline"
              onClick={onBack}
              className="flex items-center"
            >
              <ArrowLeft className="h-4 w-4 mr-2" />
              返回
            </Button>
          )}
        </div>

        <div>
          {requirements && (
            <Button
              onClick={handleNext}
              className="flex items-center"
            >
              下一步：系统设计
              <ArrowRight className="h-4 w-4 ml-2" />
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default RequirementsStep;