import React, { useState, useEffect, useRef } from 'react';
import api from '@/services/api';
import Button from '@/components/ui/Button';
import { 
  Workflow, 
  FileCode, 
  Download, 
  Edit3, 
  Eye, 
  RefreshCw,
  AlertCircle,
  Plus,
  Check,
  X,
  ZoomIn,
  ZoomOut,
  RotateCcw,
  Save,
  FileImage,
  FileText,
  Database,
  GitBranch,
  Settings
} from 'lucide-react';
import type { Project } from '@/types';

interface PUMLDiagramsProps {
  projectId: string;
  project: Project;
}

interface PUMLDiagram {
  diagram_id: string;
  diagram_type: string;
  diagram_name: string;
  puml_content: string;
  generated_at: string;
  is_manual: boolean;
  version?: number;
}

interface ValidationResult {
  is_valid: boolean;
  errors: string[];
  warnings: string[];
}

const PUMLDiagrams: React.FC<PUMLDiagramsProps> = ({ projectId, project }) => {
  const [diagrams, setDiagrams] = useState<PUMLDiagram[]>([]);
  const [selectedDiagram, setSelectedDiagram] = useState<PUMLDiagram | null>(null);
  const [isGenerating, setIsGenerating] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [previewMode, setPreviewMode] = useState(false);
  const [editContent, setEditContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [validationResult, setValidationResult] = useState<ValidationResult | null>(null);
  const [previewUrl, setPreviewUrl] = useState('');
  const [isValidating, setIsValidating] = useState(false);
  const [isPreviewing, setIsPreviewing] = useState(false);
  const [activeTab, setActiveTab] = useState<'edit' | 'preview'>('edit');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    loadDiagrams();
  }, [projectId]);

  useEffect(() => {
    if (editMode && editContent) {
      validatePUML(editContent);
    }
  }, [editContent, editMode]);

  const loadDiagrams = async () => {
    try {
      setLoading(true);
      const response = await api.get(`/ai/puml/project/${projectId}`);
      if (response.data.success) {
        setDiagrams(response.data.data || []);
        if (response.data.data && response.data.data.length > 0) {
          setSelectedDiagram(response.data.data[0]);
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载图表失败');
    } finally {
      setLoading(false);
    }
  };

  const generatePUML = async () => {
    try {
      setIsGenerating(true);
      setError('');

      const analysisResponse = await api.get(`/ai/analysis/project/${projectId}`);
      if (!analysisResponse.data.success || !analysisResponse.data.data.length) {
        setError('请先完成需求分析，然后再生成PUML图表');
        return;
      }

      const analysisId = analysisResponse.data.data[0].analysis_id;

      // 生成多种类型的图表
      const diagramTypes = ['architecture', 'business_flow', 'data_model', 'sequence'];
      
      for (const diagramType of diagramTypes) {
        try {
          await api.post('/ai/puml/generate', {
            analysis_id: analysisId,
            diagram_type: diagramType,
            project_id: projectId
          });
        } catch (err) {
          console.log(`生成${diagramType}图表时出错:`, err);
        }
      }

      await loadDiagrams();
      
    } catch (err: any) {
      setError(err.response?.data?.error || '生成PUML图表失败');
    } finally {
      setIsGenerating(false);
    }
  };

  const validatePUML = async (content: string) => {
    if (!content.trim()) {
      setValidationResult(null);
      return;
    }

    try {
      setIsValidating(true);
      const response = await api.post('/puml/validate', {
        puml_code: content
      });
      
      if (response.data.success) {
        setValidationResult(response.data.data);
      }
    } catch (err: any) {
      console.error('验证失败:', err);
    } finally {
      setIsValidating(false);
    }
  };

  const previewPUML = async (content: string) => {
    if (!content.trim()) return;

    try {
      setIsPreviewing(true);
      const response = await api.post('/puml/preview', {
        puml_code: content,
        format: 'png',
        return_type: 'base64'
      });
      
      if (response.data.success) {
        setPreviewUrl(response.data.data);
      }
    } catch (err: any) {
      setError('预览生成失败: ' + (err.response?.data?.error || err.message));
    } finally {
      setIsPreviewing(false);
    }
  };

  const handleEdit = () => {
    if (selectedDiagram) {
      setEditContent(selectedDiagram.puml_content);
      setEditMode(true);
      setActiveTab('edit');
    }
  };

  const handleSave = async () => {
    if (!selectedDiagram) return;

    try {
      await api.put(`/ai/puml/${selectedDiagram.diagram_id}`, {
        content: editContent,
        title: selectedDiagram.diagram_name
      });

      const updatedDiagrams = diagrams.map(d => 
        d.diagram_id === selectedDiagram.diagram_id 
          ? { ...d, puml_content: editContent, is_manual: true }
          : d
      );
      setDiagrams(updatedDiagrams);
      setSelectedDiagram({ ...selectedDiagram, puml_content: editContent, is_manual: true });
      setEditMode(false);
      setPreviewUrl('');
    } catch (err: any) {
      setError(err.response?.data?.error || '保存失败');
    }
  };

  const exportPUML = async (format: string) => {
    if (!selectedDiagram) return;

    try {
      const response = await api.post('/puml/export', {
        puml_code: selectedDiagram.puml_content,
        format: format,
        filename: `${selectedDiagram.diagram_name}.${format}`
      }, {
        responseType: 'blob'
      });

      const blob = new Blob([response.data]);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${selectedDiagram.diagram_name}.${format}`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err: any) {
      setError('导出失败: ' + (err.response?.data?.error || err.message));
    }
  };

  const getDiagramTypeLabel = (type: string) => {
    switch (type) {
      case 'architecture': return '系统架构图';
      case 'business_flow': return '业务流程图';
      case 'sequence': return '时序图';
      case 'data_model': return '数据模型图';
      case 'class': return '类图';
      default: return type;
    }
  };

  const getDiagramIcon = (type: string) => {
    switch (type) {
      case 'architecture': return <Settings className="h-4 w-4" />;
      case 'business_flow': return <GitBranch className="h-4 w-4" />;
      case 'sequence': return <Workflow className="h-4 w-4" />;
      case 'data_model': return <Database className="h-4 w-4" />;
      case 'class': return <FileCode className="h-4 w-4" />;
      default: return <FileCode className="h-4 w-4" />;
    }
  };

  const handlePreview = () => {
    const content = editMode ? editContent : selectedDiagram?.puml_content || '';
    previewPUML(content);
    setActiveTab('preview');
  };

  if (loading) {
    return (
      <div className="p-6 text-center">
        <RefreshCw className="h-8 w-8 animate-spin text-blue-600 mx-auto mb-4" />
        <p className="text-gray-600">正在加载PUML图表...</p>
      </div>
    );
  }

  if (diagrams.length === 0) {
    return (
      <div className="p-6">
        <div className="text-center py-12">
          <Workflow className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">暂无PUML图表</h3>
          <p className="text-gray-500 mb-6">
            基于需求分析生成架构图、流程图、数据模型图等，或手动创建PUML图表
          </p>
          
          {error && (
            <div className="max-w-md mx-auto mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-600 text-sm">{error}</p>
            </div>
          )}

          <div className="flex justify-center space-x-4">
            <Button
              onClick={generatePUML}
              disabled={isGenerating}
              className="px-6 py-3"
            >
              {isGenerating ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  生成中...
                </>
              ) : (
                <>
                  <Workflow className="h-4 w-4 mr-2" />
                  基于需求分析生成
                </>
              )}
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full">
      {/* 左侧图表列表 */}
      <div className="w-80 border-r border-gray-200 bg-gray-50">
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h3 className="font-medium text-gray-900">PUML图表</h3>
            <Button
              size="sm"
              onClick={generatePUML}
              disabled={isGenerating}
            >
              {isGenerating ? (
                <RefreshCw className="h-3 w-3 animate-spin" />
              ) : (
                <Plus className="h-3 w-3" />
              )}
            </Button>
          </div>
        </div>
        
        <div className="p-4 space-y-2">
          {diagrams.map((diagram) => (
            <div
              key={diagram.diagram_id}
              onClick={() => {
                setSelectedDiagram(diagram);
                setEditMode(false);
                setPreviewUrl('');
                setActiveTab('edit');
              }}
              className={`p-3 rounded-lg cursor-pointer transition-colors ${
                selectedDiagram?.diagram_id === diagram.diagram_id
                  ? 'bg-blue-50 border border-blue-200'
                  : 'bg-white border border-gray-200 hover:bg-gray-50'
              }`}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  {getDiagramIcon(diagram.diagram_type)}
                  <div>
                    <h4 className="font-medium text-sm text-gray-900">
                      {getDiagramTypeLabel(diagram.diagram_type)}
                    </h4>
                    <p className="text-xs text-gray-500 mt-1">
                      {new Date(diagram.generated_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-1">
                  {diagram.is_manual && (
                    <span className="bg-green-100 text-green-800 text-xs px-2 py-0.5 rounded-full">
                      手动
                    </span>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* 右侧图表内容 */}
      <div className="flex-1 flex flex-col">
        {selectedDiagram ? (
          <>
            {/* 工具栏 */}
            <div className="p-4 border-b border-gray-200 bg-white">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-lg font-semibold text-gray-900">
                    {getDiagramTypeLabel(selectedDiagram.diagram_type)}
                  </h2>
                  <p className="text-sm text-gray-500">
                    生成时间：{new Date(selectedDiagram.generated_at).toLocaleString()}
                  </p>
                </div>
                
                <div className="flex items-center space-x-2">
                  {editMode ? (
                    <>
                      {/* 编辑模式下的标签页 */}
                      <div className="flex border border-gray-300 rounded-lg">
                        <button
                          onClick={() => setActiveTab('edit')}
                          className={`px-3 py-1 text-sm rounded-l-lg ${
                            activeTab === 'edit' 
                              ? 'bg-blue-500 text-white' 
                              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                          }`}
                        >
                          <Edit3 className="h-3 w-3 mr-1 inline" />
                          编辑
                        </button>
                        <button
                          onClick={handlePreview}
                          disabled={isPreviewing}
                          className={`px-3 py-1 text-sm rounded-r-lg ${
                            activeTab === 'preview' 
                              ? 'bg-blue-500 text-white' 
                              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                          }`}
                        >
                          {isPreviewing ? (
                            <RefreshCw className="h-3 w-3 mr-1 inline animate-spin" />
                          ) : (
                            <Eye className="h-3 w-3 mr-1 inline" />
                          )}
                          预览
                        </button>
                      </div>
                      
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => {
                          setEditMode(false);
                          setPreviewUrl('');
                          setActiveTab('edit');
                        }}
                      >
                        <X className="h-4 w-4 mr-1" />
                        取消
                      </Button>
                      <Button size="sm" onClick={handleSave}>
                        <Save className="h-4 w-4 mr-1" />
                        保存
                      </Button>
                    </>
                  ) : (
                    <>
                      <div className="flex border border-gray-300 rounded-lg">
                        <button
                          onClick={() => exportPUML('puml')}
                          className="px-3 py-1 text-sm bg-gray-100 text-gray-700 hover:bg-gray-200 rounded-l-lg"
                        >
                          <FileText className="h-3 w-3 mr-1 inline" />
                          .puml
                        </button>
                        <button
                          onClick={() => exportPUML('png')}
                          className="px-3 py-1 text-sm bg-gray-100 text-gray-700 hover:bg-gray-200"
                        >
                          <FileImage className="h-3 w-3 mr-1 inline" />
                          .png
                        </button>
                        <button
                          onClick={() => exportPUML('svg')}
                          className="px-3 py-1 text-sm bg-gray-100 text-gray-700 hover:bg-gray-200 rounded-r-lg"
                        >
                          <FileImage className="h-3 w-3 mr-1 inline" />
                          .svg
                        </button>
                      </div>
                      
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleEdit}
                      >
                        <Edit3 className="h-4 w-4 mr-1" />
                        编辑
                      </Button>
                    </>
                  )}
                </div>
              </div>

              {/* 验证结果 */}
              {editMode && validationResult && (
                <div className="mt-3">
                  {validationResult.is_valid ? (
                    <div className="flex items-center text-green-600 text-sm">
                      <Check className="h-4 w-4 mr-1" />
                      语法验证通过
                    </div>
                  ) : (
                    <div className="space-y-1">
                      {validationResult.errors.map((error, index) => (
                        <div key={index} className="flex items-center text-red-600 text-sm">
                          <X className="h-4 w-4 mr-1" />
                          {error}
                        </div>
                      ))}
                      {validationResult.warnings.map((warning, index) => (
                        <div key={index} className="flex items-center text-yellow-600 text-sm">
                          <AlertCircle className="h-4 w-4 mr-1" />
                          {warning}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* 内容区域 */}
            <div className="flex-1 p-4">
              {editMode ? (
                <div className="h-full">
                  {activeTab === 'edit' ? (
                    <textarea
                      ref={textareaRef}
                      value={editContent}
                      onChange={(e) => setEditContent(e.target.value)}
                      className="w-full h-full p-4 border border-gray-300 rounded-lg font-mono text-sm resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      placeholder="输入PUML代码..."
                      spellCheck={false}
                    />
                  ) : (
                    <div className="h-full bg-white rounded-lg border border-gray-300 p-4 flex items-center justify-center">
                      {previewUrl ? (
                        <img 
                          src={previewUrl} 
                          alt="PUML预览" 
                          className="max-w-full max-h-full object-contain"
                        />
                      ) : (
                        <div className="text-center text-gray-500">
                          <Eye className="h-12 w-12 mx-auto mb-2" />
                          <p>点击预览按钮生成图表预览</p>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              ) : (
                <div className="h-full bg-gray-50 rounded-lg p-4 overflow-auto">
                  <pre className="whitespace-pre-wrap font-mono text-sm text-gray-800">
                    {selectedDiagram.puml_content}
                  </pre>
                </div>
              )}
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <Eye className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500">选择一个图表查看内容</p>
            </div>
          </div>
        )}
      </div>

      {/* 错误提示 */}
      {error && (
        <div className="absolute bottom-4 right-4 max-w-sm bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-center">
            <AlertCircle className="h-4 w-4 text-red-600 mr-2 flex-shrink-0" />
            <p className="text-red-600 text-sm">{error}</p>
            <button
              onClick={() => setError('')}
              className="ml-2 text-red-600 hover:text-red-800"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default PUMLDiagrams; 