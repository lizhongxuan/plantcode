import axios, { AxiosResponse } from 'axios';
import type {
  ApiResponse,
  PaginatedResponse,
  User,
  CreateUserRequest,
  LoginRequest,
  LoginResponse,
  UserUpdateRequest,
  Project,
  CreateProjectRequest,
  ProjectUpdateRequest,
} from '@/types';

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器 - 添加认证token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理通用错误
api.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // 清除本地token并跳转到登录页
      localStorage.removeItem('auth_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 认证相关API
export const authApi = {
  // 用户注册
  register: async (data: CreateUserRequest): Promise<User> => {
    const response = await api.post<ApiResponse<User>>('/auth/register', data);
    if (!response.data.success) {
      throw new Error(response.data.error || '注册失败');
    }
    return response.data.data!;
  },

  // 用户登录
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', data);
    if (!response.data.success) {
      throw new Error(response.data.error || '登录失败');
    }
    const loginData = response.data.data!;
    // 保存token到localStorage
    localStorage.setItem('auth_token', loginData.token);
    return loginData;
  },

  // 退出登录
  logout: () => {
    localStorage.removeItem('auth_token');
  },
};

// 用户相关API
export const userApi = {
  // 获取当前用户信息
  getCurrentUser: async (): Promise<User> => {
    const response = await api.get<ApiResponse<User>>('/user/profile');
    if (!response.data.success) {
      throw new Error(response.data.error || '获取用户信息失败');
    }
    return response.data.data!;
  },

  // 更新用户信息
  updateUser: async (data: UserUpdateRequest): Promise<User> => {
    const response = await api.put<ApiResponse<User>>('/user/profile/update', data);
    if (!response.data.success) {
      throw new Error(response.data.error || '更新用户信息失败');
    }
    return response.data.data!;
  },
};

// 项目相关API
export const projectApi = {
  // 创建项目
  createProject: async (data: CreateProjectRequest): Promise<Project> => {
    const response = await api.post<ApiResponse<Project>>('/projects', data);
    if (!response.data.success) {
      throw new Error(response.data.error || '创建项目失败');
    }
    return response.data.data!;
  },

  // 获取项目列表
  getProjects: async (page: number = 1, pageSize: number = 10): Promise<{ data: PaginatedResponse<Project> }> => {
    const response = await api.get<PaginatedResponse<Project>>('/projects/list', {
      params: { page, page_size: pageSize },
    });
    if (!response.data.success) {
      throw new Error(response.data.message || '获取项目列表失败');
    }
    return response;
  },

  // 获取项目详情
  getProject: async (id: string): Promise<Project> => {
    const response = await api.get<ApiResponse<Project>>(`/projects/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error || '获取项目详情失败');
    }
    return response.data.data!;
  },

  // 更新项目
  updateProject: async (id: string, data: ProjectUpdateRequest): Promise<Project> => {
    const response = await api.put<ApiResponse<Project>>(`/projects/${id}`, data);
    if (!response.data.success) {
      throw new Error(response.data.error || '更新项目失败');
    }
    return response.data.data!;
  },

  // 删除项目
  deleteProject: async (id: string): Promise<void> => {
    const response = await api.delete<ApiResponse>(`/projects/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error || '删除项目失败');
    }
  },
};

// 系统相关API
export const systemApi = {
  // 健康检查
  healthCheck: async (): Promise<any> => {
    const response = await axios.get('/health');
    return response.data;
  },
};

export default api; 