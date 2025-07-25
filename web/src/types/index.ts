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

// 文件夹结构相关类型
export interface ProjectFolder {
  folder_id: string;
  project_id: string;
  folder_name: string; // 'requirements' | 'design' | 'tasks'
  folder_type: string;
  parent_id?: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface ProjectDocument {
  document_id: string;
  project_id: string;
  folder_id: string;
  document_name: string;
  document_type: string; // 'requirements_doc' | 'design_doc' | 'task_list' | 'puml_diagram' | 'general'
  content: string;
  version: number;
  is_template: boolean;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ProjectFolderWithDocs {
  folder_id: string;
  project_id: string;
  folder_name: string;
  folder_type: string;
  parent_id?: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
  documents: ProjectDocument[];
}

export interface ProjectStructureResponse {
  project_id: string;
  folders: ProjectFolderWithDocs[];
}

// AI对话相关类型
export interface AIConversation {
  conversation_id: string;
  project_id: string;
  user_id: string;
  title: string;
  context: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AIMessage {
  message_id: string;
  conversation_id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  message_type: 'text' | 'document_change' | 'file_operation';
  metadata?: string;
  created_at: string;
}

export interface AIConversationResponse {
  conversation_id: string;
  messages: AIMessage[];
  context: string;
}

export interface StartAIConversationRequest {
  project_id: string;
  title: string;
}

export interface SendAIMessageRequest {
  conversation_id: string;
  content: string;
  message_type?: string;
}

// 文档变更相关类型
export interface DocumentChange {
  change_id: string;
  document_id: string;
  message_id?: string;
  change_type: 'create' | 'update' | 'delete';
  old_content: string;
  new_content: string;
  change_summary: string;
  is_ai_generated: boolean;
  changed_by: string;
  created_at: string;
}

export interface DocumentVersion {
  version_id: string;
  document_id: string;
  version_number: number;
  content: string;
  change_note: string;
  changed_by: string;
  change_type: string;
  created_at: string;
}

export interface DocumentChangesResponse {
  document_id: string;
  changes: DocumentChange[];
  versions: DocumentVersion[];
}

// 请求类型
export interface CreateDocumentRequest {
  project_id: string;
  folder_id: string;
  document_name: string;
  document_type: string;
  content: string;
  is_template?: boolean;
}

export interface UpdateDocumentRequest {
  document_id: string;
  content: string;
  change_note?: string;
}

export interface RevertDocumentRequest {
  document_id: string;
  version_number: number;
  reason_note?: string;
}

// 阶段导航相关类型
export interface StageInfo {
  name: string;
  label: string;
  description: string;
  folder_id?: string;
  completed: boolean;
  progress: number;
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

// Spec 工作流相关类型
export interface UserStory {
  id: string;
  title: string;
  description: string;
  acceptance_criteria: string[];
  priority: 'high' | 'medium' | 'low';
  story_points?: number;
  created_at: string;
}

export interface RequirementsDoc {
  id: string;
  project_id: string;
  content: string; // markdown 格式的完整需求文档
  user_stories: UserStory[];
  assumptions: string[];
  edge_cases: string[];
  functional_requirements: string[];
  non_functional_requirements: string[];
  created_at: string;
  updated_at: string;
  version: number;
}

export interface PUMLDiagram {
  id: string;
  title: string;
  type: 'sequence' | 'class' | 'activity' | 'component' | 'deployment' | 'use_case';
  code: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface TypeScriptInterface {
  id: string;
  name: string;
  code: string;
  description?: string;
}

export interface APIEndpoint {
  id: string;
  path: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  description: string;
  request_body?: any;
  response_body?: any;
  headers?: Record<string, string>;
}

export interface DesignDoc {
  id: string;
  project_id: string;
  content: string; // markdown 格式的设计文档
  puml_diagrams: PUMLDiagram[];
  interfaces: TypeScriptInterface[];
  api_endpoints: APIEndpoint[];
  database_schema?: string;
  architecture_notes: string[];
  created_at: string;
  updated_at: string;
  version: number;
}

export interface DevelopmentTask {
  id: string;
  title: string;
  description: string;
  type: 'feature' | 'bug' | 'refactor' | 'test' | 'docs';
  priority: 'high' | 'medium' | 'low';
  status: 'todo' | 'in_progress' | 'review' | 'done';
  estimated_hours: number;
  actual_hours?: number;
  assignee?: string;
  dependencies: string[]; // task ids
  user_story_id?: string;
  created_at: string;
  updated_at: string;
}

export interface TestCase {
  id: string;
  title: string;
  description: string;
  type: 'unit' | 'integration' | 'e2e' | 'api';
  steps: string[];
  expected_result: string;
  task_id?: string;
  created_at: string;
}

export interface TaskListDoc {
  id: string;
  project_id: string;
  content: string; // markdown 格式的任务文档
  tasks: DevelopmentTask[];
  test_cases: TestCase[];
  estimated_total_hours: number;
  milestones: string[];
  created_at: string;
  updated_at: string;
  version: number;
}

export interface ProjectSpec {
  id: string;
  project_id: string;
  requirements?: RequirementsDoc;
  design?: DesignDoc;
  tasks?: TaskListDoc;
  current_stage: 'requirements' | 'design' | 'tasks' | 'implementation';
  created_at: string;
  updated_at: string;
}

// Spec 工作流阶段枚举
export const SpecStages = {
  REQUIREMENTS: 'requirements',
  DESIGN: 'design', 
  TASKS: 'tasks',
  IMPLEMENTATION: 'implementation',
} as const;

export type SpecStage = typeof SpecStages[keyof typeof SpecStages];

// AI 生成请求类型
export interface GenerateRequirementsRequest {
  project_id: string;
  initial_prompt: string;
  project_type: string;
  target_audience?: string;
  business_goals?: string[];
}

export interface GenerateDesignRequest {
  project_id: string;
  requirements_id: string;
  focus_areas?: string[];
  architecture_style?: 'monolith' | 'microservices' | 'serverless';
}

export interface GenerateTasksRequest {
  project_id: string;
  requirements_id: string;
  design_id: string;
  team_size?: number;
  sprint_duration?: number;
} 