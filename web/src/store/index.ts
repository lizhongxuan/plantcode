import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User, Project } from '@/types';

// 认证状态接口
interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
  setAuth: (user: User, token: string) => void;
  clearAuth: () => void;
  updateUser: (user: User) => void;
}

// 项目状态接口
interface ProjectState {
  projects: Project[];
  currentProject: Project | null;
  isLoading: boolean;
  error: string | null;
  setProjects: (projects: Project[]) => void;
  setCurrentProject: (project: Project | null) => void;
  addProject: (project: Project) => void;
  updateProject: (project: Project) => void;
  removeProject: (projectId: string) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

// 创建认证状态store
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      isAuthenticated: false,
      user: null,
      token: null,
      setAuth: (user, token) => {
        set({
          isAuthenticated: true,
          user,
          token,
        });
        // 同时保存到localStorage
        localStorage.setItem('auth_token', token);
      },
      clearAuth: () => {
        set({
          isAuthenticated: false,
          user: null,
          token: null,
        });
        // 清除localStorage
        localStorage.removeItem('auth_token');
        // 清除持久化存储
        localStorage.removeItem('auth-store');
      },
      updateUser: (user) => {
        set((state) => ({
          ...state,
          user,
        }));
      },
    }),
    {
      name: 'auth-store',
      // 存储所有认证相关状态
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        user: state.user,
        token: state.token,
      }),
      // 添加storage选项，确保在SSR环境下正常工作
      storage: {
        getItem: (name) => {
          try {
            const item = localStorage.getItem(name);
            return item ? JSON.parse(item) : null;
          } catch (error) {
            console.warn('Failed to get item from localStorage:', error);
            return null;
          }
        },
        setItem: (name, value) => {
          try {
            localStorage.setItem(name, JSON.stringify(value));
          } catch (error) {
            console.warn('Failed to set item to localStorage:', error);
          }
        },
        removeItem: (name) => {
          try {
            localStorage.removeItem(name);
          } catch (error) {
            console.warn('Failed to remove item from localStorage:', error);
          }
        },
      },
    }
  )
);

// 创建项目状态store
export const useProjectStore = create<ProjectState>((set) => ({
  projects: [],
  currentProject: null,
  isLoading: false,
  error: null,
  setProjects: (projects) => set({ projects: projects || [] }),
  setCurrentProject: (project) => set({ currentProject: project }),
  addProject: (project) =>
    set((state) => {
      // 检查项目是否已存在，避免重复
      const existingIndex = state.projects.findIndex(p => p.project_id === project.project_id);
      if (existingIndex !== -1) {
        // 如果存在，更新现有项目
        const updatedProjects = [...state.projects];
        updatedProjects[existingIndex] = project;
        return { projects: updatedProjects };
      } else {
        // 如果不存在，添加到前面
        return { projects: [project, ...state.projects] };
      }
    }),
  updateProject: (updatedProject) =>
    set((state) => ({
      projects: state.projects.map((p) =>
        p.project_id === updatedProject.project_id ? updatedProject : p
      ),
      currentProject:
        state.currentProject?.project_id === updatedProject.project_id
          ? updatedProject
          : state.currentProject,
    })),
  removeProject: (projectId) =>
    set((state) => ({
      projects: state.projects.filter((p) => p.project_id !== projectId),
      currentProject:
        state.currentProject?.project_id === projectId
          ? null
          : state.currentProject,
    })),
  setLoading: (loading) => set({ isLoading: loading }),
  setError: (error) => set({ error }),
}));

// 通用应用状态接口
interface AppState {
  sidebarCollapsed: boolean;
  theme: 'light' | 'dark';
  setSidebarCollapsed: (collapsed: boolean) => void;
  setTheme: (theme: 'light' | 'dark') => void;
}

// 创建应用状态store
export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      sidebarCollapsed: false,
      theme: 'light',
      setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
      setTheme: (theme) => set({ theme }),
    }),
    {
      name: 'app-store',
    }
  )
); 