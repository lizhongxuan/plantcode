import React, { useState, useEffect } from 'react';
import api from '@/services/api';
import Button from '@/components/ui/Button';
import { 
  BookOpen, 
  FileText, 
  Download, 
  Edit3, 
  Eye, 
  RefreshCw,
  AlertCircle,
  Plus,
  Save
} from 'lucide-react';
import type { Project } from '@/types';

interface TechnicalDocumentsProps {
  projectId: string;
  project: Project;
}

interface Document {
  document_id: string;
  document_type: string;
  document_name: string;
  content: string;
  format: string;
  version: number;
  generated_at: string;
  is_final: boolean;
}

const TechnicalDocuments: React.FC<TechnicalDocumentsProps> = ({ projectId, project }) => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [selectedDocument, setSelectedDocument] = useState<Document | null>(null);
  const [isGenerating, setIsGenerating] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [editContent, setEditContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadDocuments();
  }, [projectId]);

  const loadDocuments = async () => {
    try {
      setLoading(true);
      const response = await api.get(`/ai/document/project/${projectId}`);
      if (response.data.success) {
        setDocuments(response.data.data || []);
        if (response.data.data && response.data.data.length > 0) {
          setSelectedDocument(response.data.data[0]);
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载文档失败');
    } finally {
      setLoading(false);
    }
  };

  const generateDocument = async () => {
    try {
      setIsGenerating(true);
      setError('');

      // 检查是否有需求分析
      const analysisResponse = await api.get(`/ai/analysis/project/${projectId}`);
      if (!analysisResponse.data.success || !analysisResponse.data.data.length) {
        setError('请先完成需求分析，然后再生成技术文档');
        return;
      }

      const analysisId = analysisResponse.data.data[0].analysis_id;

      // 生成技术文档
      const response = await api.post('/ai/document/generate', {
        analysis_id: analysisId,
        project_id: projectId,
        document_type: 'development_plan'
      });

      if (response.data.success) {
        await loadDocuments();
      }
      
    } catch (err: any) {
      setError(err.response?.data?.error || '生成技术文档失败');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleEdit = () => {
    if (selectedDocument) {
      // 如果是JSON格式，格式化显示
      let content = selectedDocument.content;
      if (selectedDocument.format === 'json') {
        try {
          const parsed = JSON.parse(content);
          content = JSON.stringify(parsed, null, 2);
        } catch (e) {
          // 如果解析失败，使用原始内容
        }
      }
      setEditContent(content);
      setEditMode(true);
    }
  };

  const handleSave = async () => {
    if (!selectedDocument) return;

    try {
      await api.put(`/ai/document/${selectedDocument.document_id}`, {
        content: editContent,
        title: selectedDocument.document_name
      });

      // 更新本地状态
      const updatedDocuments = documents.map(d => 
        d.document_id === selectedDocument.document_id 
          ? { ...d, content: editContent, version: d.version + 1 }
          : d
      );
      setDocuments(updatedDocuments);
      setSelectedDocument({ 
        ...selectedDocument, 
        content: editContent, 
        version: selectedDocument.version + 1 
      });
      setEditMode(false);
    } catch (err: any) {
      setError(err.response?.data?.error || '保存失败');
    }
  };

  const downloadDocument = () => {
    if (!selectedDocument) return;

    const blob = new Blob([selectedDocument.content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${selectedDocument.document_name}.${selectedDocument.format}`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const getDocumentTypeLabel = (type: string) => {
    switch (type) {
      case 'development_plan': return '开发计划文档';
      case 'api_design': return 'API设计文档';
      case 'database_design': return '数据库设计文档';
      case 'technical_spec': return '技术规格文档';
      default: return type;
    }
  };

  const formatContent = (content: string, format: string) => {
    if (format === 'json') {
      try {
        const parsed = JSON.parse(content);
        return JSON.stringify(parsed, null, 2);
      } catch (e) {
        return content;
      }
    }
    return content;
  };

  if (loading) {
    return (
      <div className="p-6 text-center">
        <RefreshCw className="h-8 w-8 animate-spin text-blue-600 mx-auto mb-4" />
        <p className="text-gray-600">正在加载技术文档...</p>
      </div>
    );
  }

  if (documents.length === 0) {
    return (
      <div className="p-6">
        <div className="text-center py-12">
          <BookOpen className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">暂无技术文档</h3>
          <p className="text-gray-500 mb-6">
            基于需求分析生成开发文档和技术规范
          </p>
          
          {error && (
            <div className="max-w-md mx-auto mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-600 text-sm">{error}</p>
            </div>
          )}

          <Button
            onClick={generateDocument}
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
                <BookOpen className="h-4 w-4 mr-2" />
                基于需求分析生成文档
              </>
            )}
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full">
      {/* 左侧文档列表 */}
      <div className="w-80 border-r border-gray-200 bg-gray-50">
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h3 className="font-medium text-gray-900">技术文档</h3>
            <Button
              size="sm"
              onClick={generateDocument}
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
          {documents.map((document) => (
            <div
              key={document.document_id}
              onClick={() => setSelectedDocument(document)}
              className={`p-3 rounded-lg cursor-pointer transition-colors ${
                selectedDocument?.document_id === document.document_id
                  ? 'bg-blue-50 border border-blue-200'
                  : 'bg-white border border-gray-200 hover:bg-gray-50'
              }`}
            >
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium text-sm text-gray-900">
                    {getDocumentTypeLabel(document.document_type)}
                  </h4>
                  <p className="text-xs text-gray-500 mt-1">
                    v{document.version} • {new Date(document.generated_at).toLocaleDateString()}
                  </p>
                </div>
                <div className="flex items-center space-x-1">
                  {document.is_final && (
                    <span className="bg-green-100 text-green-800 text-xs px-2 py-0.5 rounded-full">
                      最终版
                    </span>
                  )}
                  <FileText className="h-4 w-4 text-gray-400" />
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* 右侧文档内容 */}
      <div className="flex-1 flex flex-col">
        {selectedDocument ? (
          <>
            {/* 工具栏 */}
            <div className="p-4 border-b border-gray-200 bg-white">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-lg font-semibold text-gray-900">
                    {getDocumentTypeLabel(selectedDocument.document_type)}
                  </h2>
                  <p className="text-sm text-gray-500">
                    版本 {selectedDocument.version} • {new Date(selectedDocument.generated_at).toLocaleString()}
                  </p>
                </div>
                
                <div className="flex items-center space-x-2">
                  {editMode ? (
                    <>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setEditMode(false)}
                      >
                        取消
                      </Button>
                      <Button size="sm" onClick={handleSave}>
                        <Save className="h-4 w-4 mr-1" />
                        保存修改
                      </Button>
                    </>
                  ) : (
                    <>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={downloadDocument}
                      >
                        <Download className="h-4 w-4 mr-1" />
                        下载
                      </Button>
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
            </div>

            {/* 内容区域 */}
            <div className="flex-1 p-4">
              {editMode ? (
                <div className="h-full">
                  <textarea
                    value={editContent}
                    onChange={(e) => setEditContent(e.target.value)}
                    className="w-full h-full p-4 border border-gray-300 rounded-lg font-mono text-sm resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="编辑文档内容..."
                  />
                </div>
              ) : (
                <div className="h-full bg-gray-50 rounded-lg p-4 overflow-auto">
                  <pre className="whitespace-pre-wrap font-mono text-sm text-gray-800">
                    {formatContent(selectedDocument.content, selectedDocument.format)}
                  </pre>
                </div>
              )}
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <Eye className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500">选择一个文档查看内容</p>
            </div>
          </div>
        )}
      </div>

      {/* 错误提示 */}
      {error && (
        <div className="absolute bottom-4 right-4 max-w-sm bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-center">
            <AlertCircle className="h-4 w-4 text-red-600 mr-2" />
            <p className="text-red-600 text-sm">{error}</p>
          </div>
        </div>
      )}
    </div>
  );
};

export default TechnicalDocuments; 