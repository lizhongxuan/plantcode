-- 清空所有表数据但保留表结构
-- 禁用外键约束
SET FOREIGN_KEY_CHECKS = 0;

-- 清空用户表
TRUNCATE TABLE users;

-- 清空项目表
TRUNCATE TABLE projects;

-- 清空需求分析表（如果存在）
TRUNCATE TABLE requirements;

-- 清空聊天会话表（如果存在）
TRUNCATE TABLE chat_sessions;

-- 清空聊天消息表（如果存在）
TRUNCATE TABLE chat_messages;

-- 清空问题表（如果存在）
TRUNCATE TABLE questions;

-- 清空PUML图表表（如果存在）
TRUNCATE TABLE puml_diagrams;

-- 清空文档表（如果存在）
TRUNCATE TABLE documents;

-- 清空业务模块表（如果存在）
TRUNCATE TABLE business_modules;

-- 清空通用模块表（如果存在）
TRUNCATE TABLE common_modules;

-- 清空异步任务表（如果存在）
TRUNCATE TABLE async_tasks;

-- 清空阶段进度表（如果存在）
TRUNCATE TABLE stage_progress;

-- 清空项目文件夹表（如果存在）
TRUNCATE TABLE project_folders;

-- 清空项目文档表（如果存在）
TRUNCATE TABLE project_documents;

-- 清空文档版本表（如果存在）
TRUNCATE TABLE document_versions;

-- 清空AI对话表（如果存在）
TRUNCATE TABLE ai_conversations;

-- 清空AI消息表（如果存在）
TRUNCATE TABLE ai_messages;

-- 清空文档变更表（如果存在）
TRUNCATE TABLE document_changes;

-- 清空用户AI配置表（如果存在）
TRUNCATE TABLE user_ai_configs;

-- 清空项目规范表（如果存在）
TRUNCATE TABLE project_specs;

-- 清空用户故事表（如果存在）
TRUNCATE TABLE user_stories;

-- 清空需求文档表（如果存在）
TRUNCATE TABLE requirements_docs;

-- 启用外键约束
SET FOREIGN_KEY_CHECKS = 1;