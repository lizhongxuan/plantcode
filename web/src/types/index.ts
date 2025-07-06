// 用户相关类型
export interface User {
  user_id: string;
  username: string;
  email: string;
  full_name: string;
  created_at: string;
  updated_at: string;
  last_login?: string;
  status: string;
  preferences?: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  full_name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  token: string;
}

export interface UserUpdateRequest {
  username?: string;
  email?: string;
  full_name?: string;
  preferences?: string;
}

// 项目相关类型
export interface Project {
  project_id: string;
  user_id: string;
  project_name: string;
  description: string;
  project_type: string;
  status: string;
  created_at: string;
  updated_at: string;
  completion_percentage: number;
  settings?: string;
}

export interface CreateProjectRequest {
  project_name: string;
  description: string;
  project_type: string;
}

export interface ProjectUpdateRequest {
  project_name?: string;
  description?: string;
  project_type?: string;
  status?: string;
  completion_percentage?: number;
  settings?: string;
}

// API响应类型
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
  code: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  message?: string;
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_page: number;
  };
}

// 通用类型
export interface LoadingState {
  isLoading: boolean;
  error: string | null;
}

export interface FormState<T> extends LoadingState {
  data: T;
  isDirty: boolean;
  isValid: boolean;
}

// 项目类型枚举
export const ProjectTypes = {
  WEB_APPLICATION: 'web_application',
  MOBILE_APP: 'mobile_app',
  API_SERVICE: 'api_service',
  DATA_ANALYSIS: 'data_analysis',
  MACHINE_LEARNING: 'machine_learning',
  OTHER: 'other',
} as const;

export type ProjectType = typeof ProjectTypes[keyof typeof ProjectTypes];

// 项目状态枚举
export const ProjectStatus = {
  DRAFT: 'draft',
  ACTIVE: 'active',
  COMPLETED: 'completed',
  ARCHIVED: 'archived',
  DELETED: 'deleted',
} as const;

export type ProjectStatusType = typeof ProjectStatus[keyof typeof ProjectStatus]; 