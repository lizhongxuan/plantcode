import React, { useState, useEffect, useRef } from 'react';
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
  Row,
  Col,
  Divider,
  Tag,
  Popconfirm,
  Tabs,
  Select,
  Drawer
} from 'antd';
import { 
  CodeOutlined, 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined,
  EyeOutlined,
  DownloadOutlined,
  SaveOutlined,
  FullscreenOutlined,
  CopyOutlined,
  DragOutlined,
  PictureOutlined,
  FileImageOutlined,
  RobotOutlined
} from '@ant-design/icons';
import { pumlApi } from '@/services/api';
import AIAssistant from '@/components/AIAssistant';
import OnlinePUMLEditor from './OnlinePUMLEditor';

const { Title, Text } = Typography;
const { TextArea } = Input;
const { TabPane } = Tabs;
const { Option } = Select;

interface PUMLEditorProps {
  projectId: string;
  stage: number;
}

interface PUMLDiagram {
  diagram_id: string;
  diagram_name: string;
  diagram_type: string;
  puml_content: string;
  stage: number;
  version: number;
  is_validated: boolean;
  generated_at: string;
  updated_at: string;
}

const PUMLEditor: React.FC<PUMLEditorProps> = ({ projectId, stage }) => {
  const [diagrams, setDiagrams] = useState<PUMLDiagram[]>([]);
  const [loading, setLoading] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [previewVisible, setPreviewVisible] = useState(false);
  const [showAIAssistant, setShowAIAssistant] = useState(false);
  const [currentDiagram, setCurrentDiagram] = useState<PUMLDiagram | null>(null);
  const [pumlContent, setPumlContent] = useState('');
  const [diagramName, setDiagramName] = useState('');
  const [diagramType, setDiagramType] = useState('sequence');
  const [previewUrl, setPreviewUrl] = useState('');
  
  // 编辑器状态
  const [editorMode, setEditorMode] = useState<'edit' | 'preview' | 'split'>('split');
  const textareaRef = useRef<any>(null);

  // PUML图表类型选项
  const diagramTypes = [
    { value: 'sequence', label: '序列图' },
    { value: 'class', label: '类图' },
    { value: 'activity', label: '活动图' },
    { value: 'component', label: '组件图' },
    { value: 'usecase', label: '用例图' },
    { value: 'state', label: '状态图' },
    { value: 'deployment', label: '部署图' },
    { value: 'architecture', label: '架构图' }
  ];

  // 加载PUML图表列表
  const loadDiagrams = async () => {
    setLoading(true);
    try {
      const response = await pumlApi.getProjectPUMLs(projectId, stage);
      if (response.success) {
        setDiagrams(response.data || []);
      } else {
        setDiagrams([]);
      }
    } catch (error: any) {
      console.error('加载PUML图表列表失败:', error);
      setDiagrams([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadDiagrams();
  }, [projectId, stage]);

  // 新建PUML图表
  const handleNewDiagram = () => {
    setCurrentDiagram(null);
    setDiagramName('');
    setDiagramType('sequence');
    setPumlContent(getDefaultPUML('sequence'));
    setEditModalVisible(true);
  };

  // 编辑PUML图表
  const handleEditDiagram = (diagram: PUMLDiagram) => {
    setCurrentDiagram(diagram);
    setDiagramName(diagram.diagram_name);
    setDiagramType(diagram.diagram_type);
    setPumlContent(diagram.puml_content);
    setEditModalVisible(true);
  };

  // 查看PUML图表
  const handleViewDiagram = (diagram: PUMLDiagram) => {
    setCurrentDiagram(diagram);
    setPumlContent(diagram.puml_content);
    generatePreview(diagram.puml_content);
    setPreviewVisible(true);
  };

  // 保存PUML图表
  const handleSaveDiagram = async () => {
    if (!diagramName.trim() || !pumlContent.trim()) {
      message.error('请填写图表名称和PUML内容');
      return;
    }

    try {
      if (currentDiagram) {
        // 更新现有图表
        await pumlApi.updatePUMLDiagram(currentDiagram.diagram_id, {
          diagram_name: diagramName,
          puml_content: pumlContent
        });
        message.success('PUML图表更新成功！');
      } else {
        // 创建新图表
        await pumlApi.createPUML({
          project_id: projectId,
          stage: stage,
          diagram_type: diagramType,
          diagram_name: diagramName,
          puml_content: pumlContent
        });
        message.success('PUML图表创建成功！');
      }
      
      setEditModalVisible(false);
      loadDiagrams();
    } catch (error: any) {
      message.error(error.message || '保存失败');
    }
  };

  // 删除PUML图表
  const handleDeleteDiagram = async (diagramId: string) => {
    try {
      await pumlApi.deletePUML(diagramId);
      message.success('PUML图表删除成功！');
      loadDiagrams();
    } catch (error: any) {
      message.error(error.message || '删除失败');
    }
  };

  // 生成PUML预览
  const generatePreview = async (content: string) => {
    try {
      const response = await pumlApi.generateImage(content);
      if (response.success) {
        setPreviewUrl(response.data.url);
      } else {
        throw new Error(response.message);
      }
    } catch (error: any) {
      console.error('生成预览失败:', error);
      message.error('生成预览失败: ' + error.message);
      setPreviewUrl('');
    }
  };

  // 获取默认PUML内容
  const getDefaultPUML = (type: string): string => {
    const templates = {
      sequence: `@startuml 序列图
!theme plain
actor User
participant "Web界面" as UI
participant "API服务" as API
participant "数据库" as DB

User -> UI: 发起请求
UI -> API: API调用
API -> DB: 查询数据
DB --> API: 返回结果
API --> UI: 响应数据
UI --> User: 显示结果
@enduml`,

      class: `@startuml 类图
!theme plain
class User {
  -id: Long
  -username: String
  -email: String
  +login()
  +logout()
}

class Project {
  -id: Long
  -name: String
  -description: String
  +create()
  +update()
}

User ||--o{ Project : owns
@enduml`,

      activity: `@startuml 活动图
!theme plain
start
:用户登录;
if (验证通过?) then (是)
  :进入主界面;
  :选择功能;
  :执行操作;
else (否)
  :显示错误信息;
endif
stop
@enduml`,

      component: `@startuml 组件图
!theme plain
package "前端" {
  [用户界面]
  [状态管理]
}

package "后端" {
  [API网关]
  [业务服务]
  [数据访问]
}

package "数据层" {
  database "MySQL"
}

[用户界面] --> [API网关]
[API网关] --> [业务服务]
[业务服务] --> [数据访问]
[数据访问] --> MySQL
@enduml`,

      architecture: `@startuml 架构图
!theme plain
package "表示层" {
  [Web前端] as web
  [移动端] as mobile
}

package "服务层" {
  [API网关] as gateway
  [用户服务] as user
  [业务服务] as business
}

package "数据层" {
  database "主数据库" as maindb
  database "缓存" as cache
}

web --> gateway
mobile --> gateway
gateway --> user
gateway --> business
user --> maindb
business --> maindb
business --> cache
@enduml`
    };

    return templates[type as keyof typeof templates] || templates.sequence;
  };

  // 复制PUML内容
  const handleCopyContent = () => {
    navigator.clipboard.writeText(pumlContent);
    message.success('PUML内容已复制到剪贴板');
  };

  // 下载PUML文件
  const handleDownloadPUML = (diagram: PUMLDiagram) => {
    const blob = new Blob([diagram.puml_content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${diagram.diagram_name}.puml`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // 处理AI助手生成的内容
  const handleAIGenerated = (type: 'document' | 'puml', content: any) => {
    if (type === 'puml') {
      setDiagramName(content.name || '新图表');
      setDiagramType(content.type || 'sequence');
      setPumlContent(content.content || '');
      setEditModalVisible(true);
    }
  };

  // 处理图表类型变化
  const handleTypeChange = (newType: string) => {
    setDiagramType(newType);
    if (!pumlContent || pumlContent === getDefaultPUML(diagramType)) {
      setPumlContent(getDefaultPUML(newType));
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
        {/* PUML图表列表区域 */}
        <Col span={showAIAssistant ? 16 : 24}>
          <Card 
            title={
              <Space>
                <CodeOutlined />
                <span>{getStageTitle(stage)} - PUML图表</span>
              </Space>
            }
            extra={
              <Space>
                <Button 
                  type="primary" 
                  icon={<PlusOutlined />}
                  onClick={handleNewDiagram}
                >
                  新建图表
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
              {diagrams.length === 0 ? (
                <Empty 
                  description={
                    <div>
                      <p>暂无PUML图表</p>
                      <p>点击"新建图表"创建你的第一个PUML图表，或使用AI助手生成</p>
                    </div>
                  }
                  image={Empty.PRESENTED_IMAGE_SIMPLE}
                />
              ) : (
                <List
                  itemLayout="horizontal"
                  dataSource={diagrams}
                  renderItem={(diagram) => (
                    <List.Item
                      actions={[
                        <Button 
                          key="view"
                          type="link" 
                          icon={<EyeOutlined />}
                          onClick={() => handleViewDiagram(diagram)}
                        >
                          预览
                        </Button>,
                        <Button 
                          key="edit"
                          type="link" 
                          icon={<EditOutlined />}
                          onClick={() => handleEditDiagram(diagram)}
                        >
                          编辑
                        </Button>,
                        <Button 
                          key="download"
                          type="link" 
                          icon={<DownloadOutlined />}
                          onClick={() => handleDownloadPUML(diagram)}
                        >
                          下载
                        </Button>,
                        <Popconfirm 
                          key="delete"
                          title="确定要删除这个图表吗？" 
                          onConfirm={() => handleDeleteDiagram(diagram.diagram_id)}
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
                        avatar={<CodeOutlined style={{ fontSize: '24px', color: '#52c41a' }} />}
                        title={
                          <Space>
                            <Text strong>{diagram.diagram_name}</Text>
                            <Tag color="blue">
                              {diagramTypes.find(t => t.value === diagram.diagram_type)?.label || diagram.diagram_type}
                            </Tag>
                            <Tag color={diagram.is_validated ? 'green' : 'orange'}>
                              {diagram.is_validated ? '已验证' : '未验证'}
                            </Tag>
                            <Tag color="default">v{diagram.version}</Tag>
                          </Space>
                        }
                        description={
                          <div>
                            <Text type="secondary">
                              更新时间: {new Date(diagram.updated_at).toLocaleString()}
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
              context="puml"
              stage={stage}
              onGenerate={handleAIGenerated}
            />
          </Col>
        )}
      </Row>

      {/* PUML编辑弹窗 */}
      <Modal
        title={currentDiagram ? `编辑 ${currentDiagram.diagram_name}` : '新建PUML图表'}
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        width={1200}
        footer={[
          <Button key="cancel" onClick={() => setEditModalVisible(false)}>
            取消
          </Button>,
          <Button key="copy" icon={<CopyOutlined />} onClick={handleCopyContent}>
            复制内容
          </Button>,
          <Button key="save" type="primary" icon={<SaveOutlined />} onClick={handleSaveDiagram}>
            保存
          </Button>,
        ]}
        style={{ top: 20 }}
      >
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          {/* 图表信息 */}
          <Row gutter={16}>
            <Col span={12}>
              <Input
                placeholder="图表名称"
                value={diagramName}
                onChange={(e) => setDiagramName(e.target.value)}
                prefix="名称："
              />
            </Col>
            <Col span={12}>
              <Select
                value={diagramType}
                onChange={handleTypeChange}
                style={{ width: '100%' }}
                placeholder="选择图表类型"
              >
                {diagramTypes.map(type => (
                  <Option key={type.value} value={type.value}>
                    {type.label}
                  </Option>
                ))}
              </Select>
            </Col>
          </Row>

          {/* 编辑器模式切换 */}
          <Space>
            <Text>显示模式：</Text>
            <Button.Group>
              <Button 
                type={editorMode === 'edit' ? 'primary' : 'default'}
                onClick={() => setEditorMode('edit')}
              >
                纯编辑
              </Button>
              <Button 
                type={editorMode === 'preview' ? 'primary' : 'default'}
                onClick={() => setEditorMode('preview')}
              >
                纯预览
              </Button>
              <Button 
                type={editorMode === 'split' ? 'primary' : 'default'}
                onClick={() => setEditorMode('split')}
              >
                分屏显示
              </Button>
            </Button.Group>
            <Button 
              icon={<PictureOutlined />}
              onClick={() => generatePreview(pumlContent)}
            >
              生成预览
            </Button>
          </Space>

          {/* 编辑器区域 */}
          <div style={{ height: '500px', border: '1px solid #d9d9d9', borderRadius: '6px', overflow: 'hidden' }}>
            <OnlinePUMLEditor
              value={pumlContent}
              onChange={setPumlContent}
              mode={editorMode}
              readOnly={editorMode === 'preview'}
            />
          </div>
        </Space>
      </Modal>

      {/* PUML预览弹窗 */}
      <Modal
        title={currentDiagram ? `预览：${currentDiagram.diagram_name}` : 'PUML预览'}
        open={previewVisible}
        onCancel={() => setPreviewVisible(false)}
        width={900}
        footer={[
          <Button key="close" onClick={() => setPreviewVisible(false)}>
            关闭
          </Button>
        ]}
      >
        <div style={{ height: '500px', border: '1px solid #d9d9d9', borderRadius: '6px', overflow: 'hidden' }}>
          <OnlinePUMLEditor
            value={pumlContent}
            onChange={() => {}}
            mode="preview"
            readOnly
          />
        </div>
      </Modal>
    </div>
  );
};

export default PUMLEditor; 