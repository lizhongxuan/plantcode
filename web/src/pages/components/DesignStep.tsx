import React, { useState, useEffect } from 'react';
import { Layers, Wand2, ArrowRight, ArrowLeft, CheckCircle, Clock, AlertCircle, Plus, Edit3 } from 'lucide-react';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import OnlinePUMLEditor from './OnlinePUMLEditor';
import type { 
  GenerateDesignRequest, 
  DesignDoc,
  PUMLDiagram,
  TypeScriptInterface,
  APIEndpoint 
} from '@/types';

interface DesignStepProps {
  projectId: string;
  requirementsId: string;
  onNext: (designId: string) => void;
  onBack?: () => void;
}

const DesignStep: React.FC<DesignStepProps> = ({ 
  projectId, 
  requirementsId,
  onNext, 
  onBack 
}) => {
  const [formData, setFormData] = useState<Omit<GenerateDesignRequest, 'project_id' | 'requirements_id'>>({
    focus_areas: [],
    architecture_style: 'monolith'
  });

  const [isGenerating, setIsGenerating] = useState(false);
  const [design, setDesign] = useState<DesignDoc | null>(null);
  const [error, setError] = useState<string>('');
  const [focusAreaInput, setFocusAreaInput] = useState('');
  
  // PUML 编辑器状态
  const [selectedDiagram, setSelectedDiagram] = useState<PUMLDiagram | null>(null);
  const [showPUMLEditor, setShowPUMLEditor] = useState(false);
  const [editingDiagramIndex, setEditingDiagramIndex] = useState<number>(-1);

  const handleAddFocusArea = () => {
    if (focusAreaInput.trim()) {
      setFormData(prev => ({
        ...prev,
        focus_areas: [...(prev.focus_areas || []), focusAreaInput.trim()]
      }));
      setFocusAreaInput('');
    }
  };

  const handleRemoveFocusArea = (index: number) => {
    setFormData(prev => ({
      ...prev,
      focus_areas: prev.focus_areas?.filter((_, i) => i !== index) || []
    }));
  };

  const handleGenerate = async () => {
    setIsGenerating(true);
    setError('');

    try {
      const requestData: GenerateDesignRequest = {
        project_id: projectId,
        requirements_id: requirementsId,
        focus_areas: formData.focus_areas?.length ? formData.focus_areas : undefined,
        architecture_style: formData.architecture_style
      };

      const response = await fetch(`/api/v1/projects/${projectId}/spec/design`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(requestData)
      });

      if (!response.ok) {
        throw new Error('生成设计文档失败');
      }

      const result = await response.json();
      
      if (result.success) {
        setDesign(result.data);
      } else {
        throw new Error(result.message || '生成失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '生成设计文档时发生错误');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleNext = () => {
    if (design) {
      onNext(design.id);
    }
  };

  const handleEditDesign = () => {
    setDesign(null);
  };

  const handleEditPUMLDiagram = (diagram: PUMLDiagram, index: number) => {
    setSelectedDiagram(diagram);
    setEditingDiagramIndex(index);
    setShowPUMLEditor(true);
  };

  const handleCreateNewDiagram = () => {
    const newDiagram: PUMLDiagram = {
      id: `new-${Date.now()}`,
      title: '新建图表',
      type: 'component',
      code: '@startuml\n\n@enduml',
      description: '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };
    setSelectedDiagram(newDiagram);
    setEditingDiagramIndex(-1);
    setShowPUMLEditor(true);
  };

  const handleSavePUMLDiagram = (updatedCode: string) => {
    if (design && selectedDiagram) {
      const updatedDiagram = {
        ...selectedDiagram,
        code: updatedCode,
        updated_at: new Date().toISOString()
      };

      if (editingDiagramIndex >= 0) {
        // 更新现有图表
        const updatedDiagrams = [...design.puml_diagrams];
        updatedDiagrams[editingDiagramIndex] = updatedDiagram;
        setDesign(prev => prev ? { ...prev, puml_diagrams: updatedDiagrams } : null);
      } else {
        // 添加新图表
        setDesign(prev => prev ? { 
          ...prev, 
          puml_diagrams: [...prev.puml_diagrams, updatedDiagram] 
        } : null);
      }
    }
    setShowPUMLEditor(false);
    setSelectedDiagram(null);
    setEditingDiagramIndex(-1);
  };

  const handleClosePUMLEditor = () => {
    setShowPUMLEditor(false);
    setSelectedDiagram(null);
    setEditingDiagramIndex(-1);
  };

  if (showPUMLEditor && selectedDiagram) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full">
              <span className="text-sm font-medium">2</span>
            </div>
            <div>
              <h2 className="text-xl font-semibold text-gray-900">系统设计 - PUML 编辑</h2>
              <p className="text-sm text-gray-600">编辑 {selectedDiagram.title}</p>
            </div>
          </div>
          
          <Button
            variant="outline"
            onClick={handleClosePUMLEditor}
          >
            返回设计文档
          </Button>
        </div>

        <OnlinePUMLEditor
          value={selectedDiagram.code}
          onSave={handleSavePUMLDiagram}
          onClose={handleClosePUMLEditor}
          mode="split"
        />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 步骤标题 */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full">
            <span className="text-sm font-medium">2</span>
          </div>
          <div>
            <h2 className="text-xl font-semibold text-gray-900">系统设计 (Design)</h2>
            <p className="text-sm text-gray-600">生成技术架构和 PUML 设计图表</p>
          </div>
        </div>
        
        <div className="flex items-center space-x-2">
          {design && (
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

      {!design ? (
        /* 设计参数配置 */
        <div className="bg-white rounded-lg shadow p-6">
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                架构风格
              </label>
              <select
                value={formData.architecture_style}
                onChange={(e) => setFormData(prev => ({ ...prev, architecture_style: e.target.value as 'monolith' | 'microservices' | 'serverless' }))}
                className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="monolith">单体架构</option>
                <option value="microservices">微服务架构</option>
                <option value="serverless">无服务器架构</option>
              </select>
            </div>

            {/* 重点关注领域 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                重点关注领域
              </label>
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={focusAreaInput}
                  onChange={(e) => setFocusAreaInput(e.target.value)}
                  placeholder="如：性能优化、安全性、可扩展性等"
                  className="flex-1 border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  onKeyPress={(e) => e.key === 'Enter' && handleAddFocusArea()}
                />
                <Button
                  type="button"
                  onClick={handleAddFocusArea}
                  size="sm"
                >
                  添加
                </Button>
              </div>

              {/* 已添加的关注领域 */}
              {formData.focus_areas && formData.focus_areas.length > 0 && (
                <div className="mt-2 space-y-1">
                  {formData.focus_areas.map((area, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between bg-purple-50 px-3 py-2 rounded-md"
                    >
                      <span className="text-sm text-purple-700">{area}</span>
                      <button
                        onClick={() => handleRemoveFocusArea(index)}
                        className="text-purple-600 hover:text-purple-800"
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
                className="flex items-center"
              >
                <Wand2 className="h-4 w-4 mr-2" />
                {isGenerating ? '正在生成设计文档...' : '生成设计文档'}
              </Button>
            </div>
          </div>
        </div>
      ) : (
        /* 设计文档预览 */
        <div className="space-y-6">
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <Layers className="h-5 w-5 text-purple-600" />
                <h3 className="text-lg font-medium text-gray-900">设计文档</h3>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleEditDesign}
              >
                重新生成
              </Button>
            </div>

            <div className="p-6">
              <div className="prose max-w-none">
                <div
                  dangerouslySetInnerHTML={{
                    __html: design.content.replace(/\n/g, '<br/>')
                  }}
                />
              </div>
            </div>
          </div>

          {/* PUML 图表区域 */}
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
              <h3 className="text-lg font-medium text-gray-900">PUML 设计图表</h3>
              <Button
                size="sm"
                onClick={handleCreateNewDiagram}
                className="flex items-center"
              >
                <Plus className="h-4 w-4 mr-1" />
                新建图表
              </Button>
            </div>

            <div className="p-6">
              <div className="grid gap-4">
                {design.puml_diagrams?.map((diagram, index) => (
                  <div
                    key={diagram.id}
                    className="border border-gray-200 rounded-lg p-4"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div>
                        <h4 className="font-medium text-gray-900">{diagram.title}</h4>
                        <p className="text-sm text-gray-600 capitalize">{diagram.type} 图</p>
                      </div>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleEditPUMLDiagram(diagram, index)}
                        className="flex items-center"
                      >
                        <Edit3 className="h-3 w-3 mr-1" />
                        编辑
                      </Button>
                    </div>
                    
                    {diagram.description && (
                      <p className="text-sm text-gray-700 mb-3">{diagram.description}</p>
                    )}
                    
                    <div className="bg-gray-50 rounded-md p-3">
                      <pre className="text-xs text-gray-600 overflow-x-auto">
                        {diagram.code}
                      </pre>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* API 端点预览 */}
          {design.api_endpoints && design.api_endpoints.length > 0 && (
            <div className="bg-white rounded-lg shadow">
              <div className="px-6 py-4 border-b border-gray-200">
                <h3 className="text-lg font-medium text-gray-900">API 端点设计</h3>
              </div>
              <div className="p-6">
                <div className="space-y-3">
                  {design.api_endpoints.map((endpoint, index) => (
                    <div
                      key={endpoint.id}
                      className="flex items-center justify-between p-3 border border-gray-200 rounded-md"
                    >
                      <div className="flex items-center space-x-3">
                        <span className={`px-2 py-1 rounded text-xs font-medium ${
                          endpoint.method === 'GET' ? 'bg-blue-100 text-blue-700' :
                          endpoint.method === 'POST' ? 'bg-green-100 text-green-700' :
                          endpoint.method === 'PUT' ? 'bg-yellow-100 text-yellow-700' :
                          endpoint.method === 'DELETE' ? 'bg-red-100 text-red-700' :
                          'bg-gray-100 text-gray-700'
                        }`}>
                          {endpoint.method}
                        </span>
                        <span className="font-mono text-sm">{endpoint.path}</span>
                      </div>
                      <p className="text-sm text-gray-600">{endpoint.description}</p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* TypeScript 接口预览 */}
          {design.interfaces && design.interfaces.length > 0 && (
            <div className="bg-white rounded-lg shadow">
              <div className="px-6 py-4 border-b border-gray-200">
                <h3 className="text-lg font-medium text-gray-900">TypeScript 接口</h3>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  {design.interfaces.map((iface, index) => (
                    <div
                      key={iface.id}
                      className="border border-gray-200 rounded-lg p-4"
                    >
                      <div className="flex items-center justify-between mb-2">
                        <h4 className="font-medium text-gray-900">{iface.name}</h4>
                      </div>
                      {iface.description && (
                        <p className="text-sm text-gray-600 mb-3">{iface.description}</p>
                      )}
                      <div className="bg-gray-900 rounded-md p-3">
                        <pre className="text-sm text-green-400 overflow-x-auto">
                          {iface.code}
                        </pre>
                      </div>
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
          {design && (
            <Button
              onClick={handleNext}
              className="flex items-center"
            >
              下一步：任务分解
              <ArrowRight className="h-4 w-4 ml-2" />
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default DesignStep;