import React, { useState, useEffect, useCallback } from 'react';
import api from '@/services/api';
import Button from '@/components/ui/Button';
import Modal from '@/components/ui/Modal';
import Input from '@/components/ui/Input';
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
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [newDiagramName, setNewDiagramName] = useState('');
  const [newDiagramCode, setNewDiagramCode] = useState('@startuml\n\n@enduml');
  const [isCreating, setIsCreating] = useState(false);
  const [editorMode, setEditorMode] = useState<'edit' | 'preview' | 'split'>('split');

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

  // 移除旧的renderPreview和selectedDiagramImage逻辑

  useEffect(() => {
    // Cleanup function to revoke URL on component unmount
    return () => {
      if (selectedDiagramImage) {
        URL.revokeObjectURL(selectedDiagramImage);
      }
    };
  }, [selectedDiagramImage]);

  
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

  const handleOpenCreateModal = () => {
    setNewDiagramName('');
    setNewDiagramCode('@startuml\n\n@enduml');
    setIsCreateModalOpen(true);
  };
  const handleCloseCreateModal = () => {
    setIsCreateModalOpen(false);
  };
  const handleCreateDiagram = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDiagramName.trim() || !newDiagramCode.trim()) return;
    setIsCreating(true);
    try {
      const response = await api.post(`/api/puml/create`, {
        project_id: projectId,
        diagram_name: newDiagramName,
        puml_content: newDiagramCode,
      });
      if (response.data.success) {
        setIsCreateModalOpen(false);
        loadDiagrams();
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '新建图表失败');
    } finally {
      setIsCreating(false);
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
        <Button onClick={handleOpenCreateModal} variant="primary">新建PUML图表</Button>
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
                  <OnlinePUMLEditor value={selectedDiagram.puml_content} onChange={() => {}} readOnly />
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

      {isCreateModalOpen && (
        <Modal visible={isCreateModalOpen} onClose={handleCloseCreateModal} width={900} bodyStyle={{ padding: 0, background: '#fff', borderRadius: 12, minHeight: 600, maxHeight: '80vh', overflow: 'auto', display: 'flex', flexDirection: 'column' }}>
          <form onSubmit={handleCreateDiagram} style={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: 0 }}>
            <div style={{ padding: 32, flex: '0 0 auto' }}>
              {/* 表单区 */}
              <div style={{ marginBottom: 24 }}>
                <label style={{ fontWeight: 600, fontSize: 16, display: 'block', marginBottom: 8 }}>图表名称</label>
                <Input value={newDiagramName} onChange={e => setNewDiagramName(e.target.value)} placeholder="请输入图表名称" required style={{ width: '100%', height: 40, borderRadius: 8, fontSize: 16, marginBottom: 16 }} />
              </div>
              {/* 显示模式切换区 */}
              <div style={{ marginBottom: 24, display: 'flex', alignItems: 'center', gap: 16 }}>
                <span style={{ fontWeight: 500, fontSize: 15 }}>显示模式：</span>
                <div style={{ display: 'flex', gap: 8 }}>
                  <Button type={editorMode === 'edit' ? 'primary' : 'default'} onClick={() => setEditorMode('edit')} style={{ borderRadius: 20, fontWeight: 500 }}>纯编辑</Button>
                  <Button type={editorMode === 'preview' ? 'primary' : 'default'} onClick={() => setEditorMode('preview')} style={{ borderRadius: 20, fontWeight: 500 }}>纯预览</Button>
                  <Button type={editorMode === 'split' ? 'primary' : 'default'} onClick={() => setEditorMode('split')} style={{ borderRadius: 20, fontWeight: 500 }}>分屏显示</Button>
                </div>
              </div>
            </div>
            {/* 编辑/预览区，flex:1填满剩余空间 */}
            <div style={{ flex: 1, minHeight: 0, minWidth: 0, display: 'flex', flexDirection: 'column', padding: '0 32px' }}>
              <OnlinePUMLEditor value={newDiagramCode} onChange={setNewDiagramCode} mode={editorMode} style={{ flex: 1, minHeight: 0, minWidth: 0, height: '100%' }} />
            </div>
            {/* 按钮区 */}
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 24, margin: '24px 32px 32px 32px' }}>
              <Button type="button" onClick={handleCloseCreateModal} variant="outline" style={{ borderRadius: 8, fontWeight: 500 }}>取消</Button>
              <Button type="submit" disabled={isCreating} variant="primary" style={{ borderRadius: 8, fontWeight: 500 }}>{isCreating ? '保存中...' : '保存'}</Button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  );
};

export default PUMLDiagrams; 