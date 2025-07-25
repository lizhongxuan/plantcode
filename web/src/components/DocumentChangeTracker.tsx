import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Timeline, 
  Button, 
  Modal, 
  Typography, 
  Space, 
  Tag, 
  Spin,
  message,
  Descriptions,
  Divider
} from 'antd';
import { 
  HistoryOutlined, 
  EyeOutlined, 
  RollbackOutlined,
  UserOutlined,
  RobotOutlined,
  ClockCircleOutlined
} from '@ant-design/icons';
import { 
  DocumentChangesResponse, 
  DocumentChange, 
  DocumentVersion,
  RevertDocumentRequest 
} from '../types';

const { Text } = Typography;

interface DocumentChangeTrackerProps {
  documentId: string;
}

interface VersionModalProps {
  visible: boolean;
  version?: DocumentVersion;
  onClose: () => void;
  onRevert?: (version: DocumentVersion) => void;
}

const VersionModal: React.FC<VersionModalProps> = ({
  visible,
  version,
  onClose,
  onRevert
}) => {
  const [showDiff, setShowDiff] = useState(false);

  if (!version) return null;

  return (
    <Modal
      title={`Version ${version.version_number} Details`}
      open={visible}
      onCancel={onClose}
      width={800}
      footer={[
        <Button key="close" onClick={onClose}>
          Close
        </Button>,
        <Button 
          key="diff" 
          onClick={() => setShowDiff(!showDiff)}
        >
          {showDiff ? 'Hide Diff' : 'Show Diff'}
        </Button>,
        onRevert && (
          <Button 
            key="revert" 
            type="primary" 
            danger
            icon={<RollbackOutlined />}
            onClick={() => {
              Modal.confirm({
                title: 'Revert to this version?',
                content: 'This will create a new version with the content from this version. The current content will be preserved in history.',
                onOk: () => {
                  onRevert(version);
                  onClose();
                }
              });
            }}
          >
            Revert to This Version
          </Button>
        )
      ]}
    >
      <Descriptions column={2} size="small">
        <Descriptions.Item label="Version">{version.version_number}</Descriptions.Item>
        <Descriptions.Item label="Changed By">
          {version.change_type === 'ai_generated' ? (
            <Space>
              <RobotOutlined />
              AI Assistant
            </Space>
          ) : (
            <Space>
              <UserOutlined />
              User
            </Space>
          )}
        </Descriptions.Item>
        <Descriptions.Item label="Date">
          {new Date(version.created_at).toLocaleString()}
        </Descriptions.Item>
        <Descriptions.Item label="Type">
          <Tag color={version.change_type === 'ai_generated' ? 'green' : 'blue'}>
            {version.change_type}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="Change Note" span={2}>
          {version.change_note || 'No note provided'}
        </Descriptions.Item>
      </Descriptions>

      <Divider />

      <div style={{ maxHeight: '400px', overflow: 'auto' }}>
        <Text strong>Content:</Text>
        <pre className="whitespace-pre-wrap font-mono text-sm bg-gray-50 p-4 rounded border mt-2">
          {version.content}
        </pre>
      </div>
    </Modal>
  );
};

export const DocumentChangeTracker: React.FC<DocumentChangeTrackerProps> = ({
  documentId
}) => {
  const [changesData, setChangesData] = useState<DocumentChangesResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedVersion, setSelectedVersion] = useState<DocumentVersion | null>(null);
  const [versionModalVisible, setVersionModalVisible] = useState(false);
  const [diffModalVisible, setDiffModalVisible] = useState(false);
  const [diffData, setDiffData] = useState<{
    oldContent: string;
    newContent: string;
    title: string;
  } | null>(null);

  useEffect(() => {
    loadChanges();
  }, [documentId]);

  const loadChanges = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/api/documents/${documentId}/changes`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (response.ok) {
        const data: DocumentChangesResponse = await response.json();
        setChangesData(data);
      } else {
        throw new Error('Failed to load changes');
      }
    } catch (error) {
      console.error('Failed to load changes:', error);
      message.error('Failed to load document changes');
    } finally {
      setLoading(false);
    }
  };

  const handleRevertVersion = async (version: DocumentVersion) => {
    try {
      const request: RevertDocumentRequest = {
        document_id: documentId,
        version_number: version.version_number,
        reason_note: `Reverted to version ${version.version_number}`
      };

      const response = await fetch(`/api/documents/${documentId}/revert`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(request)
      });

      if (response.ok) {
        message.success(`Successfully reverted to version ${version.version_number}`);
        loadChanges(); // 重新加载变更历史
      } else {
        throw new Error('Failed to revert document');
      }
    } catch (error) {
      console.error('Failed to revert:', error);
      message.error('Failed to revert document');
    }
  };

  const showVersionDetails = (version: DocumentVersion) => {
    setSelectedVersion(version);
    setVersionModalVisible(true);
  };

  const showDiff = (change: DocumentChange) => {
    setDiffData({
      oldContent: change.old_content,
      newContent: change.new_content,
      title: change.change_summary
    });
    setDiffModalVisible(true);
  };

  const getChangeIcon = (change: DocumentChange) => {
    if (change.is_ai_generated) {
      return <RobotOutlined style={{ color: '#52c41a' }} />;
    }
    return <UserOutlined style={{ color: '#1890ff' }} />;
  };

  const getChangeColor = (changeType: string) => {
    switch (changeType) {
      case 'create':
        return 'green';
      case 'update':
        return 'blue';
      case 'delete':
        return 'red';
      default:
        return 'default';
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-32">
        <Spin size="large" />
      </div>
    );
  }

  if (!changesData) {
    return (
      <div className="text-center text-gray-500 p-8">
        <HistoryOutlined style={{ fontSize: '48px' }} className="mb-4" />
        <div>No change history available</div>
      </div>
    );
  }

  return (
    <div className="document-change-tracker">
      <div className="flex items-center justify-between mb-4">
        <Text strong>Document Change History</Text>
        <Button 
          icon={<ClockCircleOutlined />} 
          onClick={loadChanges}
          size="small"
        >
          Refresh
        </Button>
      </div>

      {/* Recent Changes Timeline */}
      <Card title="Recent Changes" size="small" className="mb-4">
        <Timeline mode="left">
          {changesData.changes.map((change, _index) => (
            <Timeline.Item
              key={change.change_id}
              dot={getChangeIcon(change)}
              color={getChangeColor(change.change_type)}
            >
              <div className="flex items-center justify-between">
                <div>
                  <Space>
                    <Tag color={getChangeColor(change.change_type)}>
                      {change.change_type}
                    </Tag>
                    <Text strong>{change.change_summary}</Text>
                  </Space>
                  <div className="text-xs text-gray-500 mt-1">
                    {new Date(change.created_at).toLocaleString()} • 
                    {change.is_ai_generated ? ' AI Generated' : ' Manual Change'}
                  </div>
                </div>
                <Space>
                  <Button 
                    size="small" 
                    icon={<EyeOutlined />}
                    onClick={() => showDiff(change)}
                  >
                    Diff
                  </Button>
                </Space>
              </div>
            </Timeline.Item>
          ))}
        </Timeline>
      </Card>

      {/* Version History */}
      <Card title="Version History" size="small">
        <div className="space-y-2">
          {changesData.versions.map((version) => (
            <div 
              key={version.version_id}
              className="flex items-center justify-between p-3 border rounded hover:bg-gray-50"
            >
              <div>
                <Space>
                  <Text strong>Version {version.version_number}</Text>
                  <Tag color={version.change_type === 'ai_generated' ? 'green' : 'blue'}>
                    {version.change_type}
                  </Tag>
                </Space>
                <div className="text-sm text-gray-600 mt-1">
                  {version.change_note || 'No note provided'}
                </div>
                <div className="text-xs text-gray-500">
                  {new Date(version.created_at).toLocaleString()}
                </div>
              </div>
              <Space>
                <Button 
                  size="small" 
                  icon={<EyeOutlined />}
                  onClick={() => showVersionDetails(version)}
                >
                  View
                </Button>
                {version.version_number > 1 && (
                  <Button 
                    size="small" 
                    icon={<RollbackOutlined />}
                    onClick={() => handleRevertVersion(version)}
                  >
                    Revert
                  </Button>
                )}
              </Space>
            </div>
          ))}
        </div>
      </Card>

      {/* Version Details Modal */}
      <VersionModal
        visible={versionModalVisible}
        version={selectedVersion || undefined}
        onClose={() => setVersionModalVisible(false)}
        onRevert={handleRevertVersion}
      />

      {/* Diff Modal */}
      <Modal
        title={`Changes: ${diffData?.title}`}
        open={diffModalVisible}
        onCancel={() => setDiffModalVisible(false)}
        width={1000}
        footer={[
          <Button key="close" onClick={() => setDiffModalVisible(false)}>
            Close
          </Button>
        ]}
      >
        {diffData && (
          <div style={{ maxHeight: '500px', overflow: 'auto' }}>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Text strong>Previous Version</Text>
                <pre className="whitespace-pre-wrap font-mono text-sm bg-red-50 p-4 rounded border mt-2">
                  {diffData.oldContent}
                </pre>
              </div>
              <div>
                <Text strong>Current Version</Text>
                <pre className="whitespace-pre-wrap font-mono text-sm bg-green-50 p-4 rounded border mt-2">
                  {diffData.newContent}
                </pre>
              </div>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default DocumentChangeTracker;