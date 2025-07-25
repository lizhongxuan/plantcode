import React, { useState } from 'react';
import { ListTodo, Wand2, ArrowLeft, CheckCircle, Clock, AlertCircle, Calendar, User, BarChart3 } from 'lucide-react';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import type { 
  GenerateTasksRequest, 
  TaskListDoc,
  DevelopmentTask,
  TestCase 
} from '@/types';

interface TasksStepProps {
  projectId: string;
  requirementsId: string;
  designId: string;
  onComplete: (taskListId: string) => void;
  onBack?: () => void;
}

const TasksStep: React.FC<TasksStepProps> = ({ 
  projectId, 
  requirementsId,
  designId,
  onComplete, 
  onBack 
}) => {
  const [formData, setFormData] = useState<Omit<GenerateTasksRequest, 'project_id' | 'requirements_id' | 'design_id'>>({
    team_size: 3,
    sprint_duration: 2
  });

  const [isGenerating, setIsGenerating] = useState(false);
  const [taskList, setTaskList] = useState<TaskListDoc | null>(null);
  const [error, setError] = useState<string>('');

  const handleGenerate = async () => {
    setIsGenerating(true);
    setError('');

    try {
      const requestData: GenerateTasksRequest = {
        project_id: projectId,
        requirements_id: requirementsId,
        design_id: designId,
        team_size: formData.team_size,
        sprint_duration: formData.sprint_duration
      };

      const response = await fetch(`/api/v1/projects/${projectId}/spec/tasks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(requestData)
      });

      if (!response.ok) {
        throw new Error('生成任务列表失败');
      }

      const result = await response.json();
      
      if (result.success) {
        setTaskList(result.data);
      } else {
        throw new Error(result.message || '生成失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '生成任务列表时发生错误');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleComplete = () => {
    if (taskList) {
      onComplete(taskList.id);
    }
  };

  const handleEditTasks = () => {
    setTaskList(null);
  };

  const getTaskTypeColor = (type: string) => {
    switch (type) {
      case 'feature': return 'bg-blue-100 text-blue-700';
      case 'bug': return 'bg-red-100 text-red-700';
      case 'refactor': return 'bg-yellow-100 text-yellow-700';
      case 'test': return 'bg-green-100 text-green-700';
      case 'docs': return 'bg-purple-100 text-purple-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'bg-red-100 text-red-700';
      case 'medium': return 'bg-yellow-100 text-yellow-700';
      case 'low': return 'bg-green-100 text-green-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getTestTypeColor = (type: string) => {
    switch (type) {
      case 'unit': return 'bg-blue-100 text-blue-700';
      case 'integration': return 'bg-green-100 text-green-700';
      case 'e2e': return 'bg-purple-100 text-purple-700';
      case 'api': return 'bg-orange-100 text-orange-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  return (
    <div className="space-y-6">
      {/* 步骤标题 */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full">
            <span className="text-sm font-medium">3</span>
          </div>
          <div>
            <h2 className="text-xl font-semibold text-gray-900">任务分解 (Tasks)</h2>
            <p className="text-sm text-gray-600">生成详细的开发任务和测试用例</p>
          </div>
        </div>
        
        <div className="flex items-center space-x-2">
          {taskList && (
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

      {!taskList ? (
        /* 任务生成配置 */
        <div className="bg-white rounded-lg shadow p-6">
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Input
                  label="团队规模"
                  type="number"
                  value={formData.team_size?.toString() || ''}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    team_size: parseInt(e.target.value) || 1
                  }))}
                  placeholder="团队人数"
                  min="1"
                  max="20"
                />
                <p className="mt-1 text-xs text-gray-500">
                  AI 将根据团队规模估算任务工作量
                </p>
              </div>

              <div>
                <Input
                  label="Sprint 周期（周）"
                  type="number"
                  value={formData.sprint_duration?.toString() || ''}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    sprint_duration: parseInt(e.target.value) || 1
                  }))}
                  placeholder="Sprint 持续时间"
                  min="1"
                  max="8"
                />
                <p className="mt-1 text-xs text-gray-500">
                  用于合理分配任务到不同的 Sprint
                </p>
              </div>
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
                className="flex items-center"
              >
                <Wand2 className="h-4 w-4 mr-2" />
                {isGenerating ? '正在生成任务列表...' : '生成任务列表'}
              </Button>
            </div>
          </div>
        </div>
      ) : (
        /* 任务列表预览 */
        <div className="space-y-6">
          {/* 任务概览 */}
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-2">
                <BarChart3 className="h-5 w-5 text-green-600" />
                <h3 className="text-lg font-medium text-gray-900">项目概览</h3>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleEditTasks}
              >
                重新生成
              </Button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="bg-blue-50 rounded-lg p-4">
                <div className="flex items-center">
                  <ListTodo className="h-8 w-8 text-blue-600" />
                  <div className="ml-3">
                    <p className="text-sm font-medium text-blue-900">总任务数</p>
                    <p className="text-2xl font-bold text-blue-600">{taskList.tasks?.length || 0}</p>
                  </div>
                </div>
              </div>

              <div className="bg-green-50 rounded-lg p-4">
                <div className="flex items-center">
                  <CheckCircle className="h-8 w-8 text-green-600" />
                  <div className="ml-3">
                    <p className="text-sm font-medium text-green-900">测试用例</p>
                    <p className="text-2xl font-bold text-green-600">{taskList.test_cases?.length || 0}</p>
                  </div>
                </div>
              </div>

              <div className="bg-yellow-50 rounded-lg p-4">
                <div className="flex items-center">
                  <Clock className="h-8 w-8 text-yellow-600" />
                  <div className="ml-3">
                    <p className="text-sm font-medium text-yellow-900">预计工时</p>
                    <p className="text-2xl font-bold text-yellow-600">{taskList.estimated_total_hours}h</p>
                  </div>
                </div>
              </div>

              <div className="bg-purple-50 rounded-lg p-4">
                <div className="flex items-center">
                  <Calendar className="h-8 w-8 text-purple-600" />
                  <div className="ml-3">
                    <p className="text-sm font-medium text-purple-900">预计周期</p>
                    <p className="text-2xl font-bold text-purple-600">
                      {Math.ceil(taskList.estimated_total_hours / (formData.team_size || 1) / 40)}周
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* 任务详情文档 */}
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200">
              <h3 className="text-lg font-medium text-gray-900">任务规划文档</h3>
            </div>
            <div className="p-6">
              <div className="prose max-w-none">
                <div
                  dangerouslySetInnerHTML={{
                    __html: taskList.content.replace(/\n/g, '<br/>')
                  }}
                />
              </div>
            </div>
          </div>

          {/* 开发任务列表 */}
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200">
              <h3 className="text-lg font-medium text-gray-900">开发任务</h3>
            </div>
            <div className="p-6">
              <div className="space-y-4">
                {taskList.tasks?.map((task, index) => (
                  <div
                    key={task.id}
                    className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
                  >
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex-1">
                        <div className="flex items-center space-x-2 mb-2">
                          <h4 className="font-medium text-gray-900">{task.title}</h4>
                          <span className={`px-2 py-1 rounded text-xs font-medium ${getTaskTypeColor(task.type)}`}>
                            {task.type}
                          </span>
                          <span className={`px-2 py-1 rounded text-xs font-medium ${getPriorityColor(task.priority)}`}>
                            {task.priority}
                          </span>
                        </div>
                        <p className="text-sm text-gray-700 mb-2">{task.description}</p>
                        <div className="flex items-center space-x-4 text-xs text-gray-500">
                          <div className="flex items-center">
                            <Clock className="h-3 w-3 mr-1" />
                            <span>{task.estimated_hours}h</span>
                          </div>
                          {task.assignee && (
                            <div className="flex items-center">
                              <User className="h-3 w-3 mr-1" />
                              <span>{task.assignee}</span>
                            </div>
                          )}
                          <div className="flex items-center">
                            <span className="inline-block w-2 h-2 bg-gray-400 rounded-full mr-1"></span>
                            <span>{task.status}</span>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* 任务依赖 */}
                    {task.dependencies && task.dependencies.length > 0 && (
                      <div className="mt-2 p-2 bg-yellow-50 rounded-md">
                        <p className="text-xs text-yellow-700">
                          <strong>依赖任务：</strong> {task.dependencies.join(', ')}
                        </p>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 测试用例 */}
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200">
              <h3 className="text-lg font-medium text-gray-900">测试用例</h3>
            </div>
            <div className="p-6">
              <div className="space-y-4">
                {taskList.test_cases?.map((testCase, index) => (
                  <div
                    key={testCase.id}
                    className="border border-gray-200 rounded-lg p-4"
                  >
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex-1">
                        <div className="flex items-center space-x-2 mb-2">
                          <h4 className="font-medium text-gray-900">{testCase.title}</h4>
                          <span className={`px-2 py-1 rounded text-xs font-medium ${getTestTypeColor(testCase.type)}`}>
                            {testCase.type}
                          </span>
                        </div>
                        <p className="text-sm text-gray-700 mb-3">{testCase.description}</p>
                        
                        {/* 测试步骤 */}
                        {testCase.steps && testCase.steps.length > 0 && (
                          <div className="mb-3">
                            <p className="text-xs font-medium text-gray-600 mb-1">测试步骤：</p>
                            <ol className="text-xs text-gray-600 space-y-1">
                              {testCase.steps.map((step, stepIndex) => (
                                <li key={stepIndex} className="flex">
                                  <span className="mr-2">{stepIndex + 1}.</span>
                                  <span>{step}</span>
                                </li>
                              ))}
                            </ol>
                          </div>
                        )}
                        
                        {/* 预期结果 */}
                        <div className="bg-green-50 rounded-md p-2">
                          <p className="text-xs font-medium text-green-700">预期结果：</p>
                          <p className="text-xs text-green-600">{testCase.expected_result}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 里程碑 */}
          {taskList.milestones && taskList.milestones.length > 0 && (
            <div className="bg-white rounded-lg shadow">
              <div className="px-6 py-4 border-b border-gray-200">
                <h3 className="text-lg font-medium text-gray-900">项目里程碑</h3>
              </div>
              <div className="p-6">
                <div className="space-y-3">
                  {taskList.milestones.map((milestone, index) => (
                    <div
                      key={index}
                      className="flex items-center p-3 bg-indigo-50 rounded-lg"
                    >
                      <div className="flex items-center justify-center w-6 h-6 bg-indigo-100 text-indigo-600 rounded-full text-xs font-medium mr-3">
                        {index + 1}
                      </div>
                      <span className="text-sm text-indigo-700">{milestone}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}
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
              上一步
            </Button>
          )}
        </div>

        <div>
          {taskList && (
            <Button
              onClick={handleComplete}
              className="flex items-center"
            >
              <CheckCircle className="h-4 w-4 mr-2" />
              完成 Spec 工作流
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default TasksStep;