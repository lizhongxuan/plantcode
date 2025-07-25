import React, { useEffect, useState } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from '@/store';
import { authApi } from '@/services/api';
import PrivateRoute from '@/components/PrivateRoute';
import MainLayout from '@/components/layout/MainLayout';
import Login from '@/pages/Login';
import Register from '@/pages/Register';
import Dashboard from '@/pages/Dashboard';
import ProjectList from '@/pages/ProjectList';
import CreateProject from '@/pages/CreateProject';
import EditProject from '@/pages/EditProject';
import AIAnalysis from '@/pages/AIAnalysis';
import Chat from '@/pages/Chat';
import Settings from '@/pages/Settings';
import ProjectDetail from '@/pages/ProjectDetail';
import ProjectSpecWorkflow from '@/pages/ProjectSpecWorkflow';

const App: React.FC = () => {
  const { isAuthenticated, setAuth, clearAuth } = useAuthStore();
  const [isInitializing, setIsInitializing] = useState(true);

  useEffect(() => {
    const validateAuth = async () => {
      // 先从zustand持久化存储中获取认证状态
      const persistedAuth = useAuthStore.getState();
      
      // 如果已经有认证状态且token存在，直接使用
      if (persistedAuth.isAuthenticated && persistedAuth.token) {
        try {
          // 验证token的有效性
          const user = await authApi.validateToken();
          // token有效，确保认证状态正确
          setAuth(user, persistedAuth.token);
        } catch (error) {
          // token无效，清除认证状态
          console.warn('Token validation failed:', error);
          clearAuth();
        } finally {
          setIsInitializing(false);
        }
        return;
      }

      // 如果没有持久化的认证状态，检查localStorage中的token
      const token = localStorage.getItem('auth_token');
      
      if (!token) {
        // 没有token，清除认证状态
        clearAuth();
        setIsInitializing(false);
        return;
      }

      try {
        // 验证token的有效性
        const user = await authApi.validateToken();
        // token有效，设置认证状态
        setAuth(user, token);
      } catch (error) {
        // token无效，清除认证状态
        console.warn('Token validation failed:', error);
        clearAuth();
        localStorage.removeItem('auth_token');
      } finally {
        setIsInitializing(false);
      }
    };

    validateAuth();
  }, [setAuth, clearAuth]);

  // 如果正在初始化，显示加载状态
  if (isInitializing) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">正在加载...</p>
        </div>
      </div>
    );
  }

  return (
    <Routes>
      {/* 公共路由 */}
      <Route 
        path="/login" 
        element={
          isAuthenticated ? <Navigate to="/dashboard" replace /> : <Login />
        } 
      />
      <Route 
        path="/register" 
        element={
          isAuthenticated ? <Navigate to="/dashboard" replace /> : <Register />
        } 
      />
      
      {/* 私有路由 */}
      <Route 
        path="/" 
        element={
          <PrivateRoute>
            <MainLayout />
          </PrivateRoute>
        }
      >
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="projects" element={<ProjectList />} />
        <Route path="projects/new" element={<CreateProject />} />
        
        {/* 项目功能路由 */}
        <Route path="projects/:projectId" element={<ProjectDetail />} />
        <Route path="projects/:projectId/edit" element={<EditProject />} />
        
        {/* 保留原有AI功能路由（向后兼容） */}
        <Route path="project/:projectId/ai-analysis" element={<AIAnalysis />} />
        <Route path="project/:projectId/chat" element={<Chat />} />
        
        {/* 通用功能路由 */}
        <Route path="chat" element={<Chat />} />
        <Route path="settings" element={<Settings />} />
        <Route path="profile" element={<div className="p-6 bg-white rounded-lg shadow">个人资料页面 - 开发中</div>} />
      </Route>
      
      {/* Spec 工作流路由 - 独立于MainLayout外部 */}
      <Route 
        path="projects/:projectId/spec" 
        element={
          <PrivateRoute>
            <ProjectSpecWorkflow />
          </PrivateRoute>
        } 
      />
      
      {/* 404路由 */}
      <Route path="*" element={
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <div className="text-center">
            <h1 className="text-6xl font-bold text-gray-900">404</h1>
            <p className="text-xl text-gray-600 mt-4">页面未找到</p>
            <button 
              onClick={() => window.history.back()}
              className="mt-4 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
            >
              返回
            </button>
          </div>
        </div>
      } />
    </Routes>
  );
};

export default App; 