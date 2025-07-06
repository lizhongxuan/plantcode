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
        
        {/* 占位路由 - 待后续实现 */}
        <Route path="projects/:id" element={<div className="p-6 bg-white rounded-lg shadow">项目详情页面 - 开发中</div>} />
        <Route path="projects/:id/edit" element={<div className="p-6 bg-white rounded-lg shadow">编辑项目页面 - 开发中</div>} />
        <Route path="documents" element={<div className="p-6 bg-white rounded-lg shadow">文档管理页面 - 开发中</div>} />
        <Route path="chat" element={<div className="p-6 bg-white rounded-lg shadow">AI对话页面 - 开发中</div>} />
        <Route path="settings" element={<div className="p-6 bg-white rounded-lg shadow">设置页面 - 开发中</div>} />
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