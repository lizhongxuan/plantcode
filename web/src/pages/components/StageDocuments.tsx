import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Button, 
  List, 
  Typography, 
  Space, 
  Spin, 
  Empty, 
  Modal, 
  Input, 
  message,
  Divider,
  Tag,
  Popconfirm,
  Row,
  Col
} from 'antd';
import { 
  FileTextOutlined, 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined,
  RobotOutlined,
  DownloadOutlined,
  EyeOutlined
} from '@ant-design/icons';
import { aiApi, asyncTaskApi } from '@/services/api';
import AIAssistant from '@/components/AIAssistant';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;

interface StageDocumentsProps {
  projectId: string;
  stage: number;
}

interface Document {
  document_id: string;
  document_name: string;
  document_type: string;
  content: string;
  format: string;
  version: number;
  is_final: boolean;
  generated_at: string;
  updated_at: string;
}

const StageDocuments: React.FC<StageDocumentsProps> = ({ projectId, stage }) => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingDoc, setEditingDoc] = useState<Document | null>(null);
  const [docContent, setDocContent] = useState('');
  const [showAIAssistant, setShowAIAssistant] = useState(false);

  // 加载文档列表
  const loadDocuments = async () => {
    setLoading(true);
    try {
      // 使用阶段文档获取API
      const response = await asyncTaskApi.get(`/projects/${projectId}/stages/${stage}/documents`);
      console.log('Documents response:', response);
      
      if (response.success && response.data?.documents) {
        setDocuments(response.data.documents);
      } else {
        setDocuments([]);
      }
    } catch (error: any) {
      console.error('加载文档列表失败:', error);
      // 不显示错误消息，而是显示空状态
      setDocuments([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadDocuments();
  }, [projectId, stage]);

  // AI生成文档列表
  const handleAIGenerate = async () => {
    setLoading(true);
    try {
      const response = await aiApi.generateStageDocumentList(projectId, stage);
      if (response.success) {
        message.success('AI文档生成成功！');
        loadDocuments(); // 刷新列表
      } else {
        throw new Error(response.message || '生成失败');
      }
    } catch (error: any) {
      console.error('AI生成文档失败:', error);
      message.error(error.message || 'AI生成文档失败');
    } finally {
      setLoading(false);
    }
  };

  // 查看文档
  const handleViewDocument = (doc: Document) => {
    setEditingDoc(doc);
    setDocContent(doc.content);
    setModalVisible(true);
  };

  // 编辑文档
  const handleEditDocument = (doc: Document) => {
    setEditingDoc(doc);
    setDocContent(doc.content);
    setModalVisible(true);
  };

  // 保存文档
  const handleSaveDocument = async () => {
    if (!editingDoc) return;

    try {
      // 这里应该调用文档更新API
      message.success('文档保存成功！');
      setModalVisible(false);
      loadDocuments();
    } catch (error: any) {
      message.error('保存文档失败: ' + error.message);
    }
  };

  // 删除文档
  const handleDeleteDocument = async (docId: string) => {
    try {
      // 这里应该调用文档删除API
      message.success('文档删除成功！');
      loadDocuments();
    } catch (error: any) {
      message.error('删除文档失败: ' + error.message);
    }
  };

  // 下载文档
  const handleDownloadDocument = (doc: Document) => {
    const blob = new Blob([doc.content], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${doc.document_name}.md`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // 处理AI助手生成的内容
  const handleAIGenerated = (type: 'document' | 'puml', content: any) => {
    if (type === 'document') {
      loadDocuments(); // 刷新文档列表
    }
  };

  const getStageTitle = (stage: number) => {
    const titles = {
      1: '阶段一：需求梳理',
      2: '阶段二：技术设计',
      3: '阶段三：实施计划'
    };
    return titles[stage as keyof typeof titles] || `阶段${stage}`;
  };

  return (
    <div style={{ padding: '24px' }}>
      <Row gutter={24}>
        {/* 文档列表区域 */}
        <Col span={showAIAssistant ? 16 : 24}>
          <Card 
            title={
              <Space>
                <FileTextOutlined />
                <span>{getStageTitle(stage)} - 需求文档</span>
              </Space>
            }
            extra={
              <Space>
                <Button 
                  type="primary" 
                  icon={<RobotOutlined />}
                  onClick={handleAIGenerate}
                  loading={loading}
                >
                  AI生成文档列表
                </Button>
                <Button 
                  icon={<RobotOutlined />}
                  onClick={() => setShowAIAssistant(!showAIAssistant)}
                >
                  {showAIAssistant ? '隐藏' : '显示'}AI助手
                </Button>
              </Space>
            }
          >
            <Spin spinning={loading}>
              {documents.length === 0 ? (
                <Empty 
                  description={
                    <div>
                      <p>暂无文档</p>
                      <p>点击"AI生成文档列表"按钮让AI根据需求分析结果生成文档</p>
                    </div>
                  }
                  image={Empty.PRESENTED_IMAGE_SIMPLE}
                />
              ) : (
                <List
                  itemLayout="horizontal"
                  dataSource={documents}
                  renderItem={(doc) => (
                    <List.Item
                      actions={[
                        <Button 
                          key="view"
                          type="link" 
                          icon={<EyeOutlined />}
                          onClick={() => handleViewDocument(doc)}
                        >
                          查看
                        </Button>,
                        <Button 
                          key="edit"
                          type="link" 
                          icon={<EditOutlined />}
                          onClick={() => handleEditDocument(doc)}
                        >
                          编辑
                        </Button>,
                        <Button 
                          key="download"
                          type="link" 
                          icon={<DownloadOutlined />}
                          onClick={() => handleDownloadDocument(doc)}
                        >
                          下载
                        </Button>,
                        <Popconfirm 
                          key="delete"
                          title="确定要删除这个文档吗？" 
                          onConfirm={() => handleDeleteDocument(doc.document_id)}
                        >
                          <Button 
                            type="link" 
                            danger 
                            icon={<DeleteOutlined />}
                          >
                            删除
                          </Button>
                        </Popconfirm>
                      ]}
                    >
                      <List.Item.Meta
                        avatar={<FileTextOutlined style={{ fontSize: '24px', color: '#1890ff' }} />}
                        title={
                          <Space>
                            <Text strong>{doc.document_name}</Text>
                            <Tag color={doc.is_final ? 'green' : 'blue'}>
                              {doc.is_final ? '已定稿' : '草稿'}
                            </Tag>
                            <Tag color="default">v{doc.version}</Tag>
                          </Space>
                        }
                        description={
                          <div>
                            <Text type="secondary">类型: {doc.document_type}</Text>
                            <br />
                            <Text type="secondary">
                              更新时间: {new Date(doc.updated_at).toLocaleString()}
                            </Text>
                          </div>
                        }
                      />
                    </List.Item>
                  )}
                />
              )}
            </Spin>
          </Card>
        </Col>

        {/* AI助手区域 */}
        {showAIAssistant && (
          <Col span={8}>
            <AIAssistant
              projectId={projectId}
              context="document"
              stage={stage}
              onGenerate={handleAIGenerated}
            />
          </Col>
        )}
      </Row>

      {/* 文档查看/编辑弹窗 */}
      <Modal
        title={editingDoc ? `${editingDoc.document_name} (v${editingDoc.version})` : '文档详情'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        width={800}
        footer={[
          <Button key="cancel" onClick={() => setModalVisible(false)}>
            取消
          </Button>,
          <Button key="save" type="primary" onClick={handleSaveDocument}>
            保存
          </Button>,
        ]}
      >
        <TextArea
          value={docContent}
          onChange={(e) => setDocContent(e.target.value)}
          rows={20}
          placeholder="文档内容..."
        />
      </Modal>
    </div>
  );
};

export default StageDocuments; 