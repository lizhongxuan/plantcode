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
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        user: state.user,
        token: state.token,
      }),
    }
  )
);

// 创建项目状态store
export const useProjectStore = create<ProjectState>((set) => ({
  projects: [],
  currentProject: null,
  isLoading: false,
  error: null,
  setProjects: (projects) => set({ projects }),
  setCurrentProject: (project) => set({ currentProject: project }),
  addProject: (project) =>
    set((state) => ({
      projects: [project, ...state.projects],
    })),
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