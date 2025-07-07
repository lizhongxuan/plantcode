import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from '@/store';
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

const App: React.FC = () => {
  const { isAuthenticated } = useAuthStore();

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
        
        {/* 项目功能路由 - 新的集成式项目详情页面 */}
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