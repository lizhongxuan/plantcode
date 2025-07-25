import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Spin } from 'antd';

const ProjectDetail: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();

  useEffect(() => {
    if (projectId) {
      // 直接跳转到spec页面
      navigate(`/projects/${projectId}/spec`);
    }
  }, [projectId, navigate]);

  // 在跳转期间显示加载状态
  return (
    <div style={{ 
      height: '100vh', 
      display: 'flex', 
      alignItems: 'center', 
      justifyContent: 'center',
      flexDirection: 'column'
    }}>
      <Spin size="large" />
      <div style={{ marginTop: '16px', color: '#666' }}>
        正在跳转到项目工作区...
      </div>
    </div>
  );
};

export default ProjectDetail;