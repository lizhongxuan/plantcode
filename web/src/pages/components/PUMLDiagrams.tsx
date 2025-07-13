import React, { useState, useEffect, useCallback } from 'react';
import api from '@/services/api';
import Button from '@/components/ui/Button';
import { 
  Workflow, 
  FileCode, 
  Edit3, 
  FileImage,
  Database,
  GitBranch,
  Settings,
  RefreshCw,
  AlertCircle,
} from 'lucide-react';
import type { Project } from '@/types';
import OnlinePUMLEditor from './OnlinePUMLEditor';

interface PUMLDiagramsProps {
  projectId: string;
  project: Project;
}

interface PUMLDiagram {
  diagram_id: string;
  diagram_type: string;
  diagram_name: string;
  puml_content: string;
  is_manual: boolean;
}

const PUMLDiagrams: React.FC<PUMLDiagramsProps> = ({ projectId }) => {
  const [diagrams, setDiagrams] = useState<PUMLDiagram[]>([]);
  const [selectedDiagram, setSelectedDiagram] = useState<PUMLDiagram | null>(null);
  const [selectedDiagramImage, setSelectedDiagramImage] = useState<string>('');
  const [isPreviewLoading, setIsPreviewLoading] = useState(false);
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const [diagramToEdit, setDiagramToEdit] = useState<PUMLDiagram | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadDiagrams = useCallback(async () => {
    try {
      setLoading(true);
      setError('');
      const response = await api.get(`/api/puml/project/${projectId}`);
      if (response.data.success) {
        const loadedDiagrams = response.data.data || [];
        setDiagrams(loadedDiagrams);
        if (loadedDiagrams.length > 0) {
          setSelectedDiagram(loadedDiagrams[0]);
        } else {
          setSelectedDiagram(null);
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载图表列表失败');
    } finally {
      setLoading(false);
    }
  }, [projectId]);

  useEffect(() => {
    loadDiagrams();
  }, [loadDiagrams]);

  const renderPreview = useCallback(async (pumlCode: string) => {
    if (!pumlCode.trim()) {
      setSelectedDiagramImage('');
      return;
    }
    setIsPreviewLoading(true);
    try {
      // Use the existing /render endpoint which returns a blob
      const response = await api.post('/api/puml/render', {
        puml_code: pumlCode,
        format: 'png',
      }, { responseType: 'blob' });
      
      const blob = new Blob([response.data], { type: 'image/png' });
      const imageUrl = URL.createObjectURL(blob);
      
      // Clean up previous blob URL
      if (selectedDiagramImage) {
        URL.revokeObjectURL(selectedDiagramImage);
      }
      
      setSelectedDiagramImage(imageUrl);
    } catch (err) {
      console.error('预览渲染失败:', err);
      setSelectedDiagramImage(''); // Clear image on error
    } finally {
      setIsPreviewLoading(false);
    }
  }, [selectedDiagramImage]);

  useEffect(() => {
    if (selectedDiagram) {
      renderPreview(selectedDiagram.puml_content);
    } else {
      setSelectedDiagramImage('');
    }
    // Cleanup function to revoke URL on component unmount
    return () => {
      if (selectedDiagramImage) {
        URL.revokeObjectURL(selectedDiagramImage);
      }
    };
  }, [selectedDiagram, renderPreview]);

  
  const handleOpenOnlineEditor = (diagram: PUMLDiagram) => {
    setDiagramToEdit(diagram);
    setIsEditorOpen(true);
  };

  const handleCloseOnlineEditor = () => {
    setDiagramToEdit(null);
    setIsEditorOpen(false);
  };

  const handleSaveChanges = async (diagramId: string, newCode: string) => {
    const originalDiagrams = [...diagrams];
    const diagramToUpdate = diagrams.find(d => d.diagram_id === diagramId);
    if (!diagramToUpdate) return;

    // Optimistic UI update
    const updatedDiagrams = diagrams.map(d =>
      d.diagram_id === diagramId ? { ...d, puml_content: newCode, is_manual: true } : d
    );
    setDiagrams(updatedDiagrams);
    if (selectedDiagram?.diagram_id === diagramId) {
       setSelectedDiagram({ ...selectedDiagram, puml_content: newCode, is_manual: true });
    }

    try {
      await api.put(`/api/puml/${diagramId}`, {
        content: newCode,
        title: diagramToUpdate.diagram_name,
      });
    } catch (err: any) {
      setError(err.response?.data?.error || '保存失败，更改已回滚');
      setDiagrams(originalDiagrams); // Rollback on error
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
  
  if (loading) return <div className="p-4 text-center">正在加载PUML图表...</div>;

  if (error) return <div className="p-4 text-red-500 text-center"><AlertCircle className="inline-block mr-2" />{error}</div>;

  return (
    <div className="puml-diagrams-container p-4 h-full flex flex-col">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">PUML 图表</h2>
        {/* Placeholder for future global actions */}
      </div>
      {diagrams.length === 0 ? (
        <div className="text-center py-10 bg-gray-800 rounded-lg">
          <p>项目中还没有PUML图表。</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 flex-grow min-h-0">
          <div className="col-span-1 bg-gray-800 rounded-lg p-3 overflow-y-auto">
            <ul>
              {diagrams.map(d => (
                <li 
                  key={d.diagram_id}
                  className={`p-3 rounded-md cursor-pointer mb-2 transition-colors ${selectedDiagram?.diagram_id === d.diagram_id ? 'bg-blue-600' : 'hover:bg-gray-700'}`}
                  onClick={() => setSelectedDiagram(d)}
                >
                  <div className="flex items-center">
                    {getDiagramIcon(d.diagram_type)}
                    <span className="ml-2 font-medium">{d.diagram_name}</span>
                    {d.is_manual && <Edit3 className="ml-auto h-3 w-3 text-yellow-400" title="已被手动修改" />}
                  </div>
                </li>
              ))}
            </ul>
          </div>

          <div className="col-span-3 bg-gray-800 rounded-lg p-4 flex flex-col">
            {selectedDiagram ? (
              <>
                <div className="flex justify-between items-center mb-3 pb-3 border-b border-gray-700">
                  <h3 className="text-lg font-semibold">{selectedDiagram.diagram_name}</h3>
                  <div className="flex items-center gap-2">
                    <Button onClick={() => handleOpenOnlineEditor(selectedDiagram)} variant="outline">
                      <Edit3 className="mr-2 h-4 w-4" />
                      在线编辑
                    </Button>
                  </div>
                </div>
                <div className="flex-grow bg-white rounded flex items-center justify-center p-2 min-h-0">
                  {isPreviewLoading ? (
                    <div>正在生成预览...</div>
                  ) : selectedDiagramImage ? (
                    <img 
                      src={selectedDiagramImage}
                      alt={selectedDiagram.diagram_name}
                      className="max-w-full max-h-full object-contain"
                    />
                  ) : (
                    <div>无法生成预览。</div>
                  )}
                </div>
              </>
            ) : (
              <div className="flex items-center justify-center h-full text-gray-500">
                请从左侧选择一个图表以预览
              </div>
            )}
          </div>
        </div>
      )}

      {isEditorOpen && diagramToEdit && (
        <OnlinePUMLEditor
          initialCode={diagramToEdit.puml_content}
          onClose={handleCloseOnlineEditor}
          onSave={(newCode) => {
            handleSaveChanges(diagramToEdit.diagram_id, newCode);
          }}
        />
      )}
    </div>
  );
};

export default PUMLDiagrams; 