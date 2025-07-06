import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore, useProjectStore } from '@/store';
import { projectApi } from '@/services/api';
import Button from '@/components/ui/Button';
import { Plus, FolderOpen, Activity, Clock, CheckCircle } from 'lucide-react';
import type { Project } from '@/types';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const { projects, setProjects, isLoading, setLoading, setError } = useProjectStore();
  const [recentProjects, setRecentProjects] = useState<Project[]>([]);

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        setLoading(true);
        const response = await projectApi.getProjects(1, 5);
        setProjects(response.data);
        setRecentProjects(response.data.slice(0, 3));
      } catch (error) {
        setError(error instanceof Error ? error.message : '获取项目失败');
      } finally {
        setLoading(false);
      }
    };

    if (user) {
      fetchProjects();
    }
  }, [user, setProjects, setLoading, setError]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'completed':
        return 'bg-blue-100 text-blue-800';
      case 'draft':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'active':
        return <Activity className="h-4 w-4" />;
      case 'completed':
        return <CheckCircle className="h-4 w-4" />;
      case 'draft':
        return <Clock className="h-4 w-4" />;
      default:
        return <FolderOpen className="h-4 w-4" />;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active':
        return '进行中';
      case 'completed':
        return '已完成';
      case 'draft':
        return '草稿';
      default:
        return '未知';
    }
  };

  return (
    <div className="space-y-6">
      {/* 欢迎区域 */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              欢迎回来，{user?.full_name}！
            </h1>
            <p className="text-gray-600 mt-1">
              在这里管理您的AI开发项目
            </p>
          </div>
          <Button
            onClick={() => navigate('/projects/new')}
            className="flex items-center"
          >
            <Plus className="h-4 w-4 mr-2" />
            新建项目
          </Button>
        </div>
      </div>

      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-3 bg-blue-100 rounded-full">
              <FolderOpen className="h-6 w-6 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">总项目数</p>
              <p className="text-2xl font-bold text-gray-900">{projects.length}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-3 bg-green-100 rounded-full">
              <Activity className="h-6 w-6 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">进行中项目</p>
              <p className="text-2xl font-bold text-gray-900">
                {projects.filter(p => p.status === 'active').length}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-3 bg-purple-100 rounded-full">
              <CheckCircle className="h-6 w-6 text-purple-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">已完成项目</p>
              <p className="text-2xl font-bold text-gray-900">
                {projects.filter(p => p.status === 'completed').length}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* 最近项目 */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-medium text-gray-900">最近项目</h2>
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate('/projects')}
            >
              查看全部
            </Button>
          </div>
        </div>
        
        <div className="p-6">
          {isLoading ? (
            <div className="flex justify-center items-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
            </div>
          ) : recentProjects.length > 0 ? (
            <div className="space-y-4">
              {recentProjects.map((project) => (
                <div
                  key={project.project_id}
                  className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                  onClick={() => navigate(`/projects/${project.project_id}`)}
                >
                  <div className="flex items-center space-x-4">
                    <div className="p-2 bg-primary-100 rounded-lg">
                      <FolderOpen className="h-5 w-5 text-primary-600" />
                    </div>
                    <div>
                      <h3 className="font-medium text-gray-900">{project.project_name}</h3>
                      <p className="text-sm text-gray-500">{project.description}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="text-right">
                      <p className="text-sm text-gray-500">完成度</p>
                      <p className="text-sm font-medium">{project.completion_percentage}%</p>
                    </div>
                    <div className={`flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(project.status)}`}>
                      {getStatusIcon(project.status)}
                      <span className="ml-1">{getStatusText(project.status)}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <FolderOpen className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500">暂无项目，创建您的第一个项目吧！</p>
              <Button
                onClick={() => navigate('/projects/new')}
                className="mt-4"
                variant="outline"
              >
                新建项目
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard; 