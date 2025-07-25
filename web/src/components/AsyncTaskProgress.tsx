import React, { useState, useEffect, useCallback } from 'react';
import { Card, Progress, Button, Alert, Space, Spin, Typography, Tag } from 'antd';
import { ReloadOutlined, CheckCircleOutlined, CloseCircleOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { asyncTaskApi } from '@/services/api';

const { Text } = Typography;

interface AsyncTaskProgressProps {
  taskId?: string;
  projectId: string;
  onTaskComplete?: (taskResult?: any) => void;
  onTaskFailed?: (error?: string) => void;
}

interface TaskStatus {
  task_id: string;
  status: string;
  progress: number;
  message?: string;
}

interface StageProgress {
  progress_id: string;
  project_id: string;
  stage: number;
  status: string;
  completion_rate: number;
  document_count: number;
  puml_count: number;
  started_at?: string;
  completed_at?: string;
}

interface StageProgressResponse {
  project_id: string;
  stages: StageProgress[];
  overall: {
    completion_rate: number;
    status: string;
  };
}

const AsyncTaskProgress: React.FC<AsyncTaskProgressProps> = ({
  taskId,
  projectId,
  onTaskComplete,
  onTaskFailed,
}) => {
  const [taskStatus, setTaskStatus] = useState<TaskStatus | null>(null);
  const [stageProgress, setStageProgress] = useState<StageProgressResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [isPolling, setIsPolling] = useState<boolean>(false);

  // 获取任务状态
  const fetchTaskStatus = useCallback(async (id: string) => {
    try {
      const response = await asyncTaskApi.getTaskStatus(id);
      if (response.success) {
        const status = response.data;
        setTaskStatus(status);
        
        // 根据状态调用回调函数
        if (status.status === 'completed') {
          setIsPolling(false);
          onTaskComplete?.(status);
        } else if (status.status === 'failed') {
          setIsPolling(false);
          onTaskFailed?.(status.message);
        }
        
        return status;
      }
    } catch (err: any) {
      setError(err.message || '获取任务状态失败');
      setIsPolling(false);
    }
    return null;
  }, [onTaskComplete, onTaskFailed]);

  // 获取阶段进度
  const fetchStageProgress = useCallback(async () => {
    try {
      const response = await asyncTaskApi.getStageProgress(projectId);
      if (response.success) {
        setStageProgress(response.data);
      }
    } catch (err: any) {
      console.error('获取阶段进度失败:', err);
    }
  }, [projectId]);

  // 轮询任务状态
  const startPolling = useCallback(async (id: string) => {
    if (isPolling) return;
    
    setIsPolling(true);
    setError(null);

    const poll = async () => {
      const status = await fetchTaskStatus(id);
      
      if (status && (status.status === 'pending' || status.status === 'running')) {
        // 继续轮询
        setTimeout(poll, 2000); // 每2秒轮询一次
      } else {
        setIsPolling(false);
      }
    };

    poll();
  }, [isPolling, fetchTaskStatus]);

  // 停止轮询
  const stopPolling = useCallback(() => {
    setIsPolling(false);
  }, []);

  // 手动刷新
  const handleRefresh = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      await Promise.all([
        taskId && fetchTaskStatus(taskId),
        fetchStageProgress(),
      ]);
    } finally {
      setLoading(false);
    }
  }, [taskId, fetchTaskStatus, fetchStageProgress]);

  // 组件挂载时获取初始数据
  useEffect(() => {
    const initData = async () => {
      setLoading(true);
      try {
        await Promise.all([
          taskId && fetchTaskStatus(taskId),
          fetchStageProgress(),
        ]);
        
        // 如果有任务ID且任务正在进行中，开始轮询
        if (taskId && taskStatus && (taskStatus.status === 'pending' || taskStatus.status === 'running')) {
          startPolling(taskId);
        }
      } finally {
        setLoading(false);
      }
    };

    initData();
  }, [taskId, projectId]);

  // 当taskId改变时重新开始轮询
  useEffect(() => {
    if (taskId && taskStatus && (taskStatus.status === 'pending' || taskStatus.status === 'running')) {
      startPolling(taskId);
    }
    
    return () => {
      setIsPolling(false);
    };
  }, [taskId, taskStatus?.status, startPolling]);

  // 获取状态标签
  const getStatusTag = (status: string) => {
    switch (status) {
      case 'pending':
        return <Tag icon={<ClockCircleOutlined />} color="default">等待中</Tag>;
      case 'running':
        return <Tag icon={<Spin size="small" />} color="processing">执行中</Tag>;
      case 'completed':
        return <Tag icon={<CheckCircleOutlined />} color="success">已完成</Tag>;
      case 'failed':
        return <Tag icon={<CloseCircleOutlined />} color="error">失败</Tag>;
      case 'not_started':
        return <Tag color="default">未开始</Tag>;
      case 'in_progress':
        return <Tag color="processing">进行中</Tag>;
      default:
        return <Tag>{status}</Tag>;
    }
  };

  // 获取进度颜色
  const getProgressColor = (status: string) => {
    switch (status) {
      case 'completed':
        return '#52c41a';
      case 'failed':
        return '#ff4d4f';
      case 'running':
      case 'in_progress':
        return '#1890ff';
      default:
        return '#d9d9d9';
    }
  };

  return (
    <div>
      {/* 当前任务状态 */}
      {taskId && taskStatus && (
        <Card 
          title="当前任务状态" 
          size="small"
          extra={
            <Space>
              {getStatusTag(taskStatus.status)}
              <Button 
                icon={<ReloadOutlined />} 
                size="small"
                onClick={handleRefresh}
                loading={loading}
              >
                刷新
              </Button>
            </Space>
          }
          style={{ marginBottom: 16 }}
        >
          <Space direction="vertical" style={{ width: '100%' }}>
            <Progress 
              percent={taskStatus.progress} 
              status={taskStatus.status === 'failed' ? 'exception' : undefined}
              strokeColor={getProgressColor(taskStatus.status)}
            />
            {taskStatus.message && (
              <Text type="secondary">{taskStatus.message}</Text>
            )}
            {isPolling && (
              <Alert
                message="正在实时更新任务状态..."
                type="info"
                showIcon
                action={
                  <Button size="small" onClick={stopPolling}>
                    停止更新
                  </Button>
                }
              />
            )}
          </Space>
        </Card>
      )}

      {/* 项目整体进度 */}
      {stageProgress && (
        <Card 
          title="项目进度" 
          size="small"
          extra={getStatusTag(stageProgress.overall.status)}
        >
          <Space direction="vertical" style={{ width: '100%' }}>
            <div>
              <Text strong>整体进度: </Text>
              <Progress 
                percent={stageProgress.overall.completion_rate} 
                strokeColor={getProgressColor(stageProgress.overall.status)}
              />
            </div>
            
            {/* 各阶段详细进度 */}
            {stageProgress.stages.map((stage) => (
              <Card 
                key={stage.progress_id} 
                size="small" 
                title={`阶段 ${stage.stage}`}
                extra={getStatusTag(stage.status)}
              >
                <Space direction="vertical" style={{ width: '100%' }}>
                  <Progress 
                    percent={stage.completion_rate} 
                    strokeColor={getProgressColor(stage.status)}
                    size="small"
                  />
                  <div>
                    <Text type="secondary">
                      文档: {stage.document_count} 个 | PUML图表: {stage.puml_count} 个
                    </Text>
                  </div>
                  {stage.started_at && (
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                      开始时间: {new Date(stage.started_at).toLocaleString()}
                    </Text>
                  )}
                  {stage.completed_at && (
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                      完成时间: {new Date(stage.completed_at).toLocaleString()}
                    </Text>
                  )}
                </Space>
              </Card>
            ))}
          </Space>
        </Card>
      )}

      {/* 错误提示 */}
      {error && (
        <Alert
          message="错误"
          description={error}
          type="error"
          showIcon
          closable
          onClose={() => setError(null)}
          style={{ marginTop: 16 }}
        />
      )}
    </div>
  );
};

export default AsyncTaskProgress; 