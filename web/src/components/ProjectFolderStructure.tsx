import React, { useState, useEffect } from 'react';
import { 
  Tree, 
  Card, 
  Button, 
  Modal, 
  Input, 
  Select, 
  message, 
  Popconfirm,
  Typography,
  Space,
  Dropdown,
  MenuProps
} from 'antd';
import { 
  FolderOutlined, 
  FileOutlined, 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined,
  MoreOutlined,
  HistoryOutlined
} from '@ant-design/icons';
import type { DataNode, TreeProps } from 'antd/es/tree';
import {
  ProjectStructureResponse,
  ProjectFolderWithDocs,
  ProjectDocument,
  CreateDocumentRequest,
  UpdateDocumentRequest,
  Project
} from '../types';

const { TextArea } = Input;
const { Option } = Select;
const { Title, Text } = Typography;

interface ProjectFolderStructureProps {
  project: Project;
  onDocumentSelect?: (document: ProjectDocument) => void;
  onStructureChange?: () => void;
}

interface DocumentModalProps {
  visible: boolean;
  document?: ProjectDocument;
  folderId?: string;
  mode: 'create' | 'edit';
  onSave: (data: any) => void;
  onCancel: () => void;
}

const DocumentModal: React.FC<DocumentModalProps> = ({
  visible,
  document,
  folderId,
  mode,
  onSave,
  onCancel
}) => {
  const [form, setForm] = useState({
    document_name: '',
    document_type: 'general',
    content: ''
  });

  useEffect(() => {
    if (document && mode === 'edit') {
      setForm({
        document_name: document.document_name,
        document_type: document.document_type,
        content: document.content
      });
    } else {
      setForm({
        document_name: '',
        document_type: 'general',
        content: ''
      });
    }
  }, [document, mode, visible]);

  const handleSave = () => {
    if (!form.document_name.trim()) {
      message.error('Document name is required');
      return;
    }

    const data = mode === 'create' 
      ? { ...form, folder_id: folderId }
      : { ...form, document_id: document?.document_id };

    onSave(data);
  };

  return (
    <Modal
      title={mode === 'create' ? 'Create Document' : 'Edit Document'}
      open={visible}
      onOk={handleSave}
      onCancel={onCancel}
      width={600}
      okText={mode === 'create' ? 'Create' : 'Save'}
    >
      <Space direction="vertical" style={{ width: '100%' }} size="middle">
        <div>
          <Text strong>Document Name</Text>
          <Input
            value={form.document_name}
            onChange={(e) => setForm({ ...form, document_name: e.target.value })}
            placeholder="Enter document name"
          />
        </div>

        <div>
          <Text strong>Document Type</Text>
          <Select
            value={form.document_type}
            onChange={(value) => setForm({ ...form, document_type: value })}
            style={{ width: '100%' }}
          >
            <Option value="general">General Document</Option>
            <Option value="requirements_doc">Requirements Document</Option>
            <Option value="design_doc">Design Document</Option>
            <Option value="task_list">Task List</Option>
            <Option value="puml_diagram">PUML Diagram</Option>
          </Select>
        </div>

        <div>
          <Text strong>Content</Text>
          <TextArea
            value={form.content}
            onChange={(e) => setForm({ ...form, content: e.target.value })}
            placeholder="Enter document content (markdown supported)"
            autoSize={{ minRows: 6, maxRows: 12 }}
          />
        </div>
      </Space>
    </Modal>
  );
};

export const ProjectFolderStructure: React.FC<ProjectFolderStructureProps> = ({
  project,
  onDocumentSelect,
  onStructureChange
}) => {
  const [structure, setStructure] = useState<ProjectStructureResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [documentModal, setDocumentModal] = useState<{
    visible: boolean;
    document?: ProjectDocument;
    folderId?: string;
    mode: 'create' | 'edit';
  }>({
    visible: false,
    mode: 'create'
  });

  useEffect(() => {
    loadProjectStructure();
  }, [project.project_id]);

  const loadProjectStructure = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/api/v1/projects/${project.project_id}/structure`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (response.ok) {
        const data: ProjectStructureResponse = await response.json();
        setStructure(data);
      } else {
        throw new Error('Failed to load project structure');
      }
    } catch (error) {
      console.error('Failed to load structure:', error);
      message.error('Failed to load project structure');
    } finally {
      setLoading(false);
    }
  };

  const createDocument = async (data: CreateDocumentRequest) => {
    try {
      const response = await fetch('/api/documents', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          ...data,
          project_id: project.project_id
        })
      });

      if (response.ok) {
        message.success('Document created successfully');
        setDocumentModal({ visible: false, mode: 'create' });
        loadProjectStructure();
        onStructureChange?.();
      } else {
        throw new Error('Failed to create document');
      }
    } catch (error) {
      console.error('Failed to create document:', error);
      message.error('Failed to create document');
    }
  };

  const updateDocument = async (data: UpdateDocumentRequest) => {
    try {
      const response = await fetch(`/api/documents/${data.document_id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          content: data.content,
          change_note: 'Manual update'
        })
      });

      if (response.ok) {
        message.success('Document updated successfully');
        setDocumentModal({ visible: false, mode: 'edit' });
        loadProjectStructure();
        onStructureChange?.();
      } else {
        throw new Error('Failed to update document');
      }
    } catch (error) {
      console.error('Failed to update document:', error);
      message.error('Failed to update document');
    }
  };

  const deleteDocument = async (documentId: string) => {
    try {
      const response = await fetch(`/api/documents/${documentId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (response.ok) {
        message.success('Document deleted successfully');
        loadProjectStructure();
        onStructureChange?.();
      } else {
        throw new Error('Failed to delete document');
      }
    } catch (error) {
      console.error('Failed to delete document:', error);
      message.error('Failed to delete document');
    }
  };

  const getDocumentMenuItems = (document: ProjectDocument): MenuProps['items'] => [
    {
      key: 'edit',
      label: 'Edit',
      icon: <EditOutlined />,
      onClick: () => setDocumentModal({
        visible: true,
        document,
        mode: 'edit'
      })
    },
    {
      key: 'history',
      label: 'View History',
      icon: <HistoryOutlined />,
      onClick: () => {
        // TODO: Implement document history view
        message.info('Document history feature coming soon');
      }
    },
    {
      key: 'delete',
      label: 'Delete',
      icon: <DeleteOutlined />,
      danger: true,
      onClick: () => {
        Modal.confirm({
          title: 'Delete Document',
          content: 'Are you sure you want to delete this document?',
          onOk: () => deleteDocument(document.document_id)
        });
      }
    }
  ];

  const getFolderIcon = (folderName: string) => {
    switch (folderName) {
      case 'requirements':
        return '📋';
      case 'design':
        return '🎨';
      case 'tasks':
        return '✅';
      default:
        return '📁';
    }
  };

  const getDocumentIcon = (documentType: string) => {
    switch (documentType) {
      case 'requirements_doc':
        return '📄';
      case 'design_doc':
        return '📐';
      case 'task_list':
        return '📝';
      case 'puml_diagram':
        return '📊';
      default:
        return '📄';
    }
  };

  const convertToTreeData = (): DataNode[] => {
    if (!structure) return [];

    return structure.folders.map(folder => ({
      title: (
        <div className="flex items-center justify-between group">
          <span className="flex items-center">
            <span className="mr-1">{getFolderIcon(folder.folder_name)}</span>
            <span className="font-medium">{folder.folder_name}</span>
            <span className="ml-2 text-xs text-gray-500">
              ({folder.documents.length} docs)
            </span>
          </span>
          <Button
            type="text"
            size="small"
            icon={<PlusOutlined />}
            className="opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={(e) => {
              e.stopPropagation();
              setDocumentModal({
                visible: true,
                folderId: folder.folder_id,
                mode: 'create'
              });
            }}
          />
        </div>
      ),
      key: folder.folder_id,
      icon: <FolderOutlined />,
      children: folder.documents.map(doc => ({
        title: (
          <div className="flex items-center justify-between group">
            <span 
              className="flex items-center cursor-pointer flex-1"
              onClick={() => onDocumentSelect?.(doc)}
            >
              <span className="mr-1">{getDocumentIcon(doc.document_type)}</span>
              <span>{doc.document_name}</span>
              <span className="ml-2 text-xs text-gray-400">
                v{doc.version}
              </span>
            </span>
            <Dropdown
              menu={{ items: getDocumentMenuItems(doc) }}
              trigger={['click']}
            >
              <Button
                type="text"
                size="small"
                icon={<MoreOutlined />}
                className="opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={(e) => e.stopPropagation()}
              />
            </Dropdown>
          </div>
        ),
        key: doc.document_id,
        icon: <FileOutlined />,
        isLeaf: true
      }))
    }));
  };

  const onTreeSelect: TreeProps['onSelect'] = (selectedKeys, info) => {
    if (info.node.isLeaf && onDocumentSelect) {
      const document = structure?.folders
        .flatMap(f => f.documents)
        .find(d => d.document_id === selectedKeys[0]);
      if (document) {
        onDocumentSelect(document);
      }
    }
  };

  return (
    <Card 
      title={
        <div className="flex items-center justify-between">
          <Title level={5} className="mb-0">Project Structure</Title>
          <Button
            icon={<PlusOutlined />}
            size="small"
            onClick={() => setDocumentModal({
              visible: true,
              mode: 'create'
            })}
          >
            New Doc
          </Button>
        </div>
      }
      size="small"
      loading={loading}
    >
      {structure && (
        <Tree
          showIcon
          defaultExpandAll
          onSelect={onTreeSelect}
          treeData={convertToTreeData()}
          className="project-folder-tree"
        />
      )}

      <DocumentModal
        visible={documentModal.visible}
        document={documentModal.document}
        folderId={documentModal.folderId}
        mode={documentModal.mode}
        onSave={documentModal.mode === 'create' ? createDocument : updateDocument}
        onCancel={() => setDocumentModal({ visible: false, mode: 'create' })}
      />
    </Card>
  );
};

export default ProjectFolderStructure;