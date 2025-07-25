-- ====================================================================
-- 数据库完全重建脚本 - 解决外键约束问题
-- ====================================================================

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 删除所有可能存在的表（按依赖关系逆序）
DROP TABLE IF EXISTS `document_changes`;
DROP TABLE IF EXISTS `ai_messages`;
DROP TABLE IF EXISTS `ai_conversations`;
DROP TABLE IF EXISTS `document_versions`;
DROP TABLE IF EXISTS `project_documents`;
DROP TABLE IF EXISTS `project_folders`;
DROP TABLE IF EXISTS `stage_progress`;
DROP TABLE IF EXISTS `user_ai_configs`;
DROP TABLE IF EXISTS `requirements_docs`;
DROP TABLE IF EXISTS `user_stories`;
DROP TABLE IF EXISTS `project_specs`;
DROP TABLE IF EXISTS `async_tasks`;
DROP TABLE IF EXISTS `common_modules`;
DROP TABLE IF EXISTS `business_modules`;
DROP TABLE IF EXISTS `documents`;
DROP TABLE IF EXISTS `puml_diagrams`;
DROP TABLE IF EXISTS `questions`;
DROP TABLE IF EXISTS `chat_messages`;
DROP TABLE IF EXISTS `chat_sessions`;
DROP TABLE IF EXISTS `requirements`;
DROP TABLE IF EXISTS `projects`;
DROP TABLE IF EXISTS `users`;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示成功消息
SELECT '数据库表已全部删除，准备重新创建...' as Status;