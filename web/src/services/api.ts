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
  timeout: 10000,  // 默认10秒超时
  headers: {
    'Content-Type': 'application/json',
  },
});

// 创建长时间任务的axios实例
const longTaskApi = axios.create({
  baseURL: '/api',
  timeout: 60000,  // 60秒超时，用于AI生成等长时间任务
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

// 为长时间任务API添加相同的拦截器
longTaskApi.interceptors.request.use(
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

// 为长时间任务API添加相同的响应拦截器
longTaskApi.interceptors.response.use(
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

// AI相关API
export const aiApi = {
  // 分阶段生成文档 - 使用长超时
  generateStageDocuments: async (projectId: string, stage: number): Promise<any> => {
    const response = await longTaskApi.post('/ai/generate-stage-documents', {
      project_id: projectId,
      stage: stage
    });
    return response.data;
  },

  // AI需求分析 - 使用长超时
  analyzeRequirement: async (projectId: string, requirement: string): Promise<any> => {
    const response = await longTaskApi.post('/ai/analyze', {
      project_id: projectId,
      requirement: requirement
    });
    return response.data;
  },

  // 获取项目的需求分析结果
  getProjectAnalysis: async (projectId: string): Promise<any> => {
    const response = await api.get(`/ai/analysis/project/${projectId}`);
    return response.data;
  },

  // 项目AI对话
  projectChat: async (projectId: string, message: string, context?: string): Promise<any> => {
    const response = await api.post('/ai/chat', {
      project_id: projectId,
      message: message,
      context: context
    });
    return response.data;
  },

  // 根据需求分析生成阶段文档列表
  generateStageDocumentList: async (projectId: string, stage: number): Promise<any> => {
    const response = await api.post('/ai/generate-document-list', {
      project_id: projectId,
      stage: stage
    });
    return response.data;
  },

  // 获取AI配置
  getConfig: async (): Promise<any> => {
    const response = await api.get('/ai/config');
    return response.data;
  },

  // 更新AI配置
  updateConfig: async (config: any): Promise<any> => {
    const response = await api.put('/ai/config', config);
    return response.data;
  },

  // 测试AI连接
  testConnection: async (): Promise<any> => {
    const response = await api.post('/ai/test-connection');
    return response.data;
  },
};

// PUML相关API
export const pumlApi = {
  // 获取项目PUML图表列表
  getProjectPUMLs: async (projectId: string, stage?: number): Promise<any> => {
    const params = stage ? { stage } : {};
    const response = await api.get(`/puml/project/${projectId}`, { params });
    return response.data;
  },

  // 创建PUML图表
  createPUML: async (data: {
    project_id: string;
    stage: number;
    diagram_type: string;
    diagram_name: string;
    puml_content: string;
  }): Promise<any> => {
    const response = await api.post('/puml/create', data);
    return response.data;
  },

  // 更新PUML图表
  updatePUMLDiagram: async (pumlId: string, data: {
    diagram_name?: string;
    puml_content?: string;
  }): Promise<any> => {
    const response = await api.put(`/puml/${pumlId}`, data);
    return response.data;
  },

  // 删除PUML图表
  deletePUML: async (pumlId: string): Promise<any> => {
    const response = await api.delete(`/puml/${pumlId}`);
    return response.data;
  },

  // 验证PUML语法
  validatePUML: async (content: string): Promise<any> => {
    const response = await api.post('/puml/validate', { puml_content: content });
    return response.data;
  },

  // 生成PUML图片
  generateImage: async (content: string): Promise<any> => {
    const response = await api.post('/puml/generate-image', { puml_content: content });
    return response.data;
  },

  // 渲染PUML为图片
  renderPUML: async (content: string): Promise<string> => {
    const response = await api.post('/puml/render', { puml_content: content });
    if (response.data.success) {
      return response.data.data.url;
    }
    throw new Error('渲染失败');
  },
};

// 异步任务相关API
export const asyncTaskApi = {
  // 启动阶段文档生成任务
  startStageDocumentGeneration: async (projectId: string, stage: number): Promise<any> => {
    const response = await api.post('/async/stage-documents', {
      project_id: projectId,
      stage: stage
    });
    return response.data;
  },

  // 获取任务状态
  getTaskStatus: async (taskId: string): Promise<any> => {
    const response = await api.get(`/async/tasks/${taskId}/status`);
    return response.data;
  },

  // 轮询任务状态（支持长轮询）
  pollTaskStatus: async (taskId: string, timeout: number = 30): Promise<any> => {
    const response = await api.get(`/async/tasks/${taskId}/poll?timeout=${timeout}`);
    return response.data;
  },

  // 获取项目阶段进度
  getStageProgress: async (projectId: string): Promise<any> => {
    const response = await api.get(`/async/projects/${projectId}/progress`);
    return response.data;
  },

  // 便捷方法：直接使用api.get调用
  get: async (path: string): Promise<any> => {
    const response = await api.get(`/async${path}`);
    return response.data;
  },

  // 获取阶段文档列表
  getStageDocuments: async (projectId: string, stage: number): Promise<any> => {
    const response = await api.get(`/async/projects/${projectId}/stages/${stage}/documents`);
    return response.data;
  },

  // 启动完整项目文档生成任务
  startCompleteProjectDocumentGeneration: async (projectId: string): Promise<any> => {
    const response = await api.post('/async/complete-project-documents', {
      project_id: projectId
    });
    return response.data;
  },
};

export default api; 