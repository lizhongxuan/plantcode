import React, { useState } from 'react';
import { Card, Button, Typography, Space, Tag, Divider } from 'antd';
import { EditOutlined, DownloadOutlined } from '@ant-design/icons';
import { ProjectDocument } from '../types';
import ReactMarkdown from 'react-markdown';

const { Title, Text } = Typography;

interface DocumentViewerProps {
  document: ProjectDocument;
  onEdit?: (document: ProjectDocument) => void;
}

const DocumentViewer: React.FC<DocumentViewerProps> = ({
  document,
  onEdit
}) => {
  const [viewMode, setViewMode] = useState<'markdown' | 'raw'>('markdown');

  const getDocumentTypeColor = (type: string) => {
    switch (type) {
      case 'requirements_doc':
        return 'blue';
      case 'design_doc':
        return 'green';
      case 'task_list':
        return 'orange';
      case 'puml_diagram':
        return 'purple';
      default:
        return 'default';
    }
  };

  const getDocumentTypeLabel = (type: string) => {
    switch (type) {
      case 'requirements_doc':
        return 'Requirements Document';
      case 'design_doc':
        return 'Design Document';
      case 'task_list':
        return 'Task List';
      case 'puml_diagram':
        return 'PUML Diagram';
      default:
        return 'General Document';
    }
  };

  const downloadDocument = () => {
    const blob = new Blob([document.content], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = window.document.createElement('a');
    a.href = url;
    a.download = `${document.document_name}.md`;
    window.document.body.appendChild(a);
    a.click();
    window.document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="document-viewer">
      {/* Document Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-3">
          <Title level={4} className="mb-0">{document.document_name}</Title>
          <Tag color={getDocumentTypeColor(document.document_type)}>
            {getDocumentTypeLabel(document.document_type)}
          </Tag>
          <Tag>v{document.version}</Tag>
        </div>
        
        <Space>
          <Button
            type={viewMode === 'raw' ? 'primary' : 'default'}
            size="small"
            onClick={() => setViewMode(viewMode === 'raw' ? 'markdown' : 'raw')}
          >
            {viewMode === 'raw' ? 'View Rendered' : 'View Source'}
          </Button>
          <Button
            icon={<DownloadOutlined />}
            size="small"
            onClick={downloadDocument}
          >
            Download
          </Button>
          {onEdit && (
            <Button
              type="primary"
              icon={<EditOutlined />}
              size="small"
              onClick={() => onEdit(document)}
            >
              Edit
            </Button>
          )}
        </Space>
      </div>

      {/* Document Metadata */}
      <Card size="small" className="mb-4">
        <Space direction="vertical" size="small" style={{ width: '100%' }}>
          <div className="flex justify-between">
            <Text type="secondary">Created:</Text>
            <Text>{new Date(document.created_at).toLocaleString()}</Text>
          </div>
          <div className="flex justify-between">
            <Text type="secondary">Last Updated:</Text>
            <Text>{new Date(document.updated_at).toLocaleString()}</Text>
          </div>
          <div className="flex justify-between">
            <Text type="secondary">Version:</Text>
            <Text>{document.version}</Text>
          </div>
          {document.is_template && (
            <div className="flex justify-between">
              <Text type="secondary">Template:</Text>
              <Tag color="cyan">Template Document</Tag>
            </div>
          )}
        </Space>
      </Card>

      <Divider />

      {/* Document Content */}
      <Card>
        <div className="document-content">
          {viewMode === 'markdown' ? (
            <div className="prose max-w-none">
              <ReactMarkdown>{document.content || '*No content*'}</ReactMarkdown>
            </div>
          ) : (
            <pre className="whitespace-pre-wrap font-mono text-sm bg-gray-50 p-4 rounded border">
              {document.content || 'No content'}
            </pre>
          )}
        </div>
      </Card>
    </div>
  );
};

export default DocumentViewer;