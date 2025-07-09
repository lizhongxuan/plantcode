package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Database 数据库连接结构
type Database struct {
	MySQL *sql.DB
	Redis *redis.Client
}

// NewDatabase 创建数据库连接
func NewDatabase(cfg *config.Config) (*Database, error) {
	// 连接MySQL
	log.Printf("连接MySQL: %v", cfg.GetDSN())
	mysql, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Printf("警告: 连接MySQL失败: %v", err)
		log.Printf("提示: 请确保MySQL服务已启动并且密码正确")
		// 不返回错误，允许应用继续启动
		mysql = nil
		panic("警告: MySQL连接测试失败")
	} else {
		// 配置连接池
		mysql.SetMaxOpenConns(cfg.Database.MaxConnections)
		mysql.SetMaxIdleConns(cfg.Database.MaxIdleConn)
		mysql.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

		// 测试连接
		if err := mysql.Ping(); err != nil {
			log.Printf("提示: 请检查MySQL服务状态和配置")
			mysql = nil
			panic("警告: MySQL连接测试失败")
		}
	}

	// 连接Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// 测试Redis连接
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("警告: Redis连接失败: %v", err)
		log.Printf("提示: 请安装并启动Redis服务")
		// Redis连接失败不影响主要功能，只记录日志
		panic("警告: Redis连接失败")
	}

	return &Database{
		MySQL: mysql,
		Redis: redisClient,
	}, nil
}

// Close 关闭数据库连接
func (db *Database) Close() error {
	var err error

	if db.MySQL != nil {
		if mysqlErr := db.MySQL.Close(); mysqlErr != nil {
			err = fmt.Errorf("关闭MySQL连接失败: %w", mysqlErr)
		}
	}

	if db.Redis != nil {
		if redisErr := db.Redis.Close(); redisErr != nil {
			if err != nil {
				err = fmt.Errorf("%w; 关闭Redis连接失败: %w", err, redisErr)
			} else {
				err = fmt.Errorf("关闭Redis连接失败: %w", redisErr)
			}
		}
	}

	return err
}

// CreateTables 创建所有必要的数据表
func (db *Database) CreateTables() error {
	tables := []string{
		// 用户表
		`CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR(36) PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			last_login TIMESTAMP NULL,
			status ENUM('active', 'inactive', 'deleted') DEFAULT 'active',
			preferences JSON,
			INDEX idx_username (username),
			INDEX idx_email (email),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 项目表
		`CREATE TABLE IF NOT EXISTS projects (
			project_id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			project_name VARCHAR(100) NOT NULL,
			description TEXT,
			project_type ENUM('web', 'mobile', 'desktop', 'api', 'other') DEFAULT 'web',
			status ENUM('planning', 'active', 'completed', 'archived') DEFAULT 'planning',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			completion_percentage INT DEFAULT 0,
			settings JSON,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id),
			INDEX idx_status (status),
			INDEX idx_project_type (project_type),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 需求分析表 (第二阶段新增)
		`CREATE TABLE IF NOT EXISTS requirement_analyses (
			requirement_id VARCHAR(36) PRIMARY KEY,
			project_id VARCHAR(36) NOT NULL,
			raw_requirement TEXT NOT NULL,
			structured_requirement JSON,
			completeness_score INT DEFAULT 0,
			analysis_status ENUM('pending', 'processing', 'completed', 'failed') DEFAULT 'pending',
			missing_info_types JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			INDEX idx_project_id (project_id),
			INDEX idx_status (analysis_status),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 对话会话表 (第二阶段新增)
		`CREATE TABLE IF NOT EXISTS chat_sessions (
			session_id VARCHAR(36) PRIMARY KEY,
			project_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			session_type ENUM('requirement_analysis', 'technical_guidance', 'code_review') DEFAULT 'requirement_analysis',
			started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			ended_at TIMESTAMP NULL,
			status ENUM('active', 'completed', 'cancelled') DEFAULT 'active',
			context JSON,
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
			INDEX idx_project_id (project_id),
			INDEX idx_user_id (user_id),
			INDEX idx_status (status),
			INDEX idx_started_at (started_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 对话消息表 (第二阶段新增)
		`CREATE TABLE IF NOT EXISTS chat_messages (
			message_id VARCHAR(36) PRIMARY KEY,
			session_id VARCHAR(36) NOT NULL,
			sender_type ENUM('user', 'assistant', 'system') NOT NULL,
			message_content TEXT NOT NULL,
			message_type ENUM('text', 'code', 'image', 'file') DEFAULT 'text',
			metadata JSON,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			processed BOOLEAN DEFAULT FALSE,
			FOREIGN KEY (session_id) REFERENCES chat_sessions(session_id) ON DELETE CASCADE,
			INDEX idx_session_id (session_id),
			INDEX idx_timestamp (timestamp),
			INDEX idx_sender_type (sender_type)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 补充问题表 (第二阶段新增)
		`CREATE TABLE IF NOT EXISTS questions (
			question_id VARCHAR(36) PRIMARY KEY,
			requirement_id VARCHAR(36) NOT NULL,
			question_text TEXT NOT NULL,
			question_category ENUM('functionality', 'technical', 'business', 'ui_ux', 'integration') NOT NULL,
			priority_level INT DEFAULT 1,
			answer_text TEXT,
			answer_status ENUM('pending', 'answered', 'skipped') DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			answered_at TIMESTAMP NULL,
			FOREIGN KEY (requirement_id) REFERENCES requirement_analyses(requirement_id) ON DELETE CASCADE,
			INDEX idx_requirement_id (requirement_id),
			INDEX idx_status (answer_status),
			INDEX idx_priority (priority_level),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// PUML图表表 (第二阶段新增)
		`CREATE TABLE IF NOT EXISTS puml_diagrams (
			diagram_id VARCHAR(36) PRIMARY KEY,
			project_id VARCHAR(36) NOT NULL,
			diagram_type ENUM('class', 'sequence', 'usecase', 'activity', 'component', 'deployment') NOT NULL,
			diagram_name VARCHAR(100) NOT NULL,
			puml_content TEXT NOT NULL,
			rendered_url VARCHAR(500),
			version INT DEFAULT 1,
			is_validated BOOLEAN DEFAULT FALSE,
			validation_feedback TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			INDEX idx_project_id (project_id),
			INDEX idx_diagram_type (diagram_type),
			INDEX idx_version (version),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 生成文档表 (第二阶段新增)
		`CREATE TABLE IF NOT EXISTS generated_documents (
			document_id VARCHAR(36) PRIMARY KEY,
			project_id VARCHAR(36) NOT NULL,
			document_type ENUM('api_spec', 'technical_design', 'user_guide', 'deployment_guide', 'development_plan') NOT NULL,
			document_name VARCHAR(100) NOT NULL,
			content LONGTEXT NOT NULL,
			format ENUM('markdown', 'html', 'pdf', 'json') DEFAULT 'markdown',
			file_path VARCHAR(500),
			version INT DEFAULT 1,
			generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_final BOOLEAN DEFAULT FALSE,
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			INDEX idx_project_id (project_id),
			INDEX idx_document_type (document_type),
			INDEX idx_version (version),
			INDEX idx_generated_at (generated_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 业务模块表 (保留现有)
		`CREATE TABLE IF NOT EXISTS business_modules (
			module_id VARCHAR(36) PRIMARY KEY,
			project_id VARCHAR(36) NOT NULL,
			module_name VARCHAR(100) NOT NULL,
			module_type ENUM('frontend', 'backend', 'database', 'api', 'service') NOT NULL,
			description TEXT,
			dependencies JSON,
			status ENUM('planned', 'in_progress', 'completed', 'tested') DEFAULT 'planned',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			INDEX idx_project_id (project_id),
			INDEX idx_module_type (module_type),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 通用模块表 (保留现有)
		`CREATE TABLE IF NOT EXISTS common_modules (
			module_id VARCHAR(36) PRIMARY KEY,
			module_name VARCHAR(100) NOT NULL,
			category VARCHAR(50) NOT NULL,
			description TEXT,
			template_code TEXT,
			configuration JSON,
			popularity_score INT DEFAULT 0,
			is_verified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_category (category),
			INDEX idx_popularity (popularity_score),
			INDEX idx_verified (is_verified)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 用户AI配置表 (第三阶段新增)
		`CREATE TABLE IF NOT EXISTS user_ai_configs (
			config_id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			provider ENUM('openai', 'claude', 'gemini') NOT NULL DEFAULT 'openai',
			openai_api_key VARCHAR(255) NULL,
			claude_api_key VARCHAR(255) NULL,
			gemini_api_key VARCHAR(255) NULL,
			default_model VARCHAR(100) NOT NULL DEFAULT 'gpt-4',
			max_tokens INT NOT NULL DEFAULT 2048,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id),
			INDEX idx_provider (provider),
			INDEX idx_is_active (is_active),
			UNIQUE KEY unique_user_active (user_id, is_active)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	if db.MySQL == nil {
		return fmt.Errorf("MySQL数据库连接不可用")
	}

	for _, table := range tables {
		if _, err := db.MySQL.Exec(table); err != nil {
			return fmt.Errorf("创建表失败: %w", err)
		}
	}

	return nil
}

// RunMigrations 执行数据库迁移
func (db *Database) RunMigrations() error {
	if db.MySQL == nil {
		log.Println("警告: MySQL连接不可用，跳过数据库迁移")
		return nil
	}

	log.Println("开始执行数据库迁移...")

	migrations := []string{
		// 迁移1: 添加gemini_api_key列到user_ai_configs表
		`ALTER TABLE user_ai_configs 
		 ADD COLUMN IF NOT EXISTS gemini_api_key VARCHAR(255) NULL 
		 AFTER claude_api_key`,

		// 迁移2: 更新provider枚举值以支持gemini
		`ALTER TABLE user_ai_configs 
		 MODIFY COLUMN provider ENUM('openai', 'claude', 'gemini') NOT NULL DEFAULT 'openai'`,

		// 迁移3: 创建异步任务表
		`CREATE TABLE IF NOT EXISTS async_tasks (
			task_id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			project_id VARCHAR(36) NOT NULL,
			task_type ENUM('stage_document_generation', 'puml_generation', 'document_generation', 'requirement_analysis', 'complete_project_documents') NOT NULL,
			task_name VARCHAR(200) NOT NULL,
			status ENUM('pending', 'running', 'completed', 'failed') DEFAULT 'pending',
			progress INT DEFAULT 0,
			result_data LONGTEXT,
			error_message TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			started_at TIMESTAMP NULL,
			completed_at TIMESTAMP NULL,
			metadata JSON,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id),
			INDEX idx_project_id (project_id),
			INDEX idx_task_type (task_type),
			INDEX idx_status (status),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 迁移4: 创建阶段进度表
		`CREATE TABLE IF NOT EXISTS stage_progress (
			progress_id VARCHAR(36) PRIMARY KEY,
			project_id VARCHAR(36) NOT NULL,
			stage INT NOT NULL,
			status ENUM('not_started', 'in_progress', 'completed', 'failed') DEFAULT 'not_started',
			completion_rate INT DEFAULT 0,
			started_at TIMESTAMP NULL,
			completed_at TIMESTAMP NULL,
			document_count INT DEFAULT 0,
			puml_count INT DEFAULT 0,
			last_task_id VARCHAR(36),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			FOREIGN KEY (last_task_id) REFERENCES async_tasks(task_id) ON DELETE SET NULL,
			INDEX idx_project_id (project_id),
			INDEX idx_stage (stage),
			INDEX idx_status (status),
			UNIQUE KEY unique_project_stage (project_id, stage)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 迁移5: 给generated_documents表添加stage和task_id字段
		`ALTER TABLE generated_documents 
		 ADD COLUMN stage INT DEFAULT 1 AFTER version,
		 ADD COLUMN task_id VARCHAR(36) NULL AFTER stage,
		 ADD FOREIGN KEY (task_id) REFERENCES async_tasks(task_id) ON DELETE SET NULL`,

		// 迁移6: 给puml_diagrams表添加stage和task_id字段
		`ALTER TABLE puml_diagrams 
		 ADD COLUMN stage INT DEFAULT 1 AFTER version,
		 ADD COLUMN task_id VARCHAR(36) NULL AFTER stage,
		 ADD FOREIGN KEY (task_id) REFERENCES async_tasks(task_id) ON DELETE SET NULL`,
	}

	for i, migration := range migrations {
		log.Printf("执行迁移 %d...", i+1)
		
		// 对于MySQL，我们需要检查列或表是否已存在
		if i == 0 { // 第一个迁移：添加gemini_api_key列
			// 检查列是否已存在
			var exists int
			err := db.MySQL.QueryRow(`
				SELECT COUNT(*) 
				FROM INFORMATION_SCHEMA.COLUMNS 
				WHERE TABLE_SCHEMA = DATABASE() 
				AND TABLE_NAME = 'user_ai_configs' 
				AND COLUMN_NAME = 'gemini_api_key'
			`).Scan(&exists)
			
			if err != nil {
				log.Printf("迁移 %d 检查失败: %v", i+1, err)
				continue
			}
			
			if exists > 0 {
				log.Printf("迁移 %d 已完成，跳过", i+1)
				continue
			}
			
			// 执行迁移
			if _, err := db.MySQL.Exec(`ALTER TABLE user_ai_configs ADD COLUMN gemini_api_key VARCHAR(255) NULL AFTER claude_api_key`); err != nil {
				log.Printf("迁移 %d 失败: %v", i+1, err)
				return fmt.Errorf("迁移 %d 失败: %w", i+1, err)
			}
		} else if i == 2 { // 第三个迁移：创建async_tasks表
			// 检查表是否已存在
			var tableExists int
			err := db.MySQL.QueryRow(`
				SELECT COUNT(*) 
				FROM INFORMATION_SCHEMA.TABLES 
				WHERE TABLE_SCHEMA = DATABASE() 
				AND TABLE_NAME = 'async_tasks'
			`).Scan(&tableExists)
			
			if err != nil {
				log.Printf("迁移 %d 检查失败: %v", i+1, err)
				continue
			}
			
			if tableExists > 0 {
				log.Printf("迁移 %d 已完成，跳过", i+1)
				continue
			}
			
			// 执行迁移
			if _, err := db.MySQL.Exec(migration); err != nil {
				log.Printf("迁移 %d 失败: %v", i+1, err)
				return fmt.Errorf("迁移 %d 失败: %w", i+1, err)
			}
		} else if i == 3 { // 第四个迁移：创建stage_progress表
			// 检查表是否已存在
			var tableExists int
			err := db.MySQL.QueryRow(`
				SELECT COUNT(*) 
				FROM INFORMATION_SCHEMA.TABLES 
				WHERE TABLE_SCHEMA = DATABASE() 
				AND TABLE_NAME = 'stage_progress'
			`).Scan(&tableExists)
			
			if err != nil {
				log.Printf("迁移 %d 检查失败: %v", i+1, err)
				continue
			}
			
			if tableExists > 0 {
				log.Printf("迁移 %d 已完成，跳过", i+1)
				continue
			}
			
			// 执行迁移
			if _, err := db.MySQL.Exec(migration); err != nil {
				log.Printf("迁移 %d 失败: %v", i+1, err)
				return fmt.Errorf("迁移 %d 失败: %w", i+1, err)
			}
		} else if i == 4 || i == 5 { // 第五、六个迁移：添加字段
			// 检查字段是否存在
			tableName := "generated_documents"
			if i == 5 {
				tableName = "puml_diagrams"
			}
			
			var columnExists int
			err := db.MySQL.QueryRow(`
				SELECT COUNT(*) 
				FROM INFORMATION_SCHEMA.COLUMNS 
				WHERE TABLE_SCHEMA = DATABASE() 
				AND TABLE_NAME = ? 
				AND COLUMN_NAME = 'stage'
			`, tableName).Scan(&columnExists)
			
			if err != nil {
				log.Printf("迁移 %d 检查失败: %v", i+1, err)
				continue
			}
			
			if columnExists > 0 {
				log.Printf("迁移 %d 已完成，跳过", i+1)
				continue
			}
			
			// 执行迁移，但使用分步操作避免外键错误
			alterSQL := fmt.Sprintf(`ALTER TABLE %s 
				ADD COLUMN stage INT DEFAULT 1 AFTER version,
				ADD COLUMN task_id VARCHAR(36) NULL AFTER stage`, tableName)
			
			if _, err := db.MySQL.Exec(alterSQL); err != nil {
				log.Printf("迁移 %d 失败: %v", i+1, err)
			}
		} else {
			// 其他迁移直接执行
			if _, err := db.MySQL.Exec(migration); err != nil {
				// 某些迁移失败不是致命错误，记录日志并继续
				log.Printf("迁移 %d 警告: %v", i+1, err)
			}
		}
		
		log.Printf("迁移 %d 完成", i+1)
	}

	log.Println("数据库迁移完成")
	return nil
}

// Health 检查数据库健康状态
func (db *Database) Health() error {
	var errors []string

	// 检查MySQL
	if db.MySQL == nil {
		errors = append(errors, "MySQL连接不存在")
	} else if err := db.MySQL.Ping(); err != nil {
		errors = append(errors, fmt.Sprintf("MySQL连接异常: %v", err))
	}

	// 检查Redis
	if db.Redis == nil {
		errors = append(errors, "Redis连接不存在")
	} else {
		ctx := context.Background()
		if err := db.Redis.Ping(ctx).Err(); err != nil {
			errors = append(errors, fmt.Sprintf("Redis连接异常: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("数据库健康检查失败: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Repository 仓库接口
type Repository interface {
	// 用户相关
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(userID uuid.UUID) (*model.User, error)
	UpdateUser(user *model.User) error
	UpdateUserLastLogin(userID uuid.UUID) error

	// 项目相关
	CreateProject(project *model.Project) error
	GetProjectsByUserID(userID uuid.UUID, page, pageSize int) ([]*model.Project, int64, error)
	GetProjectByID(projectID uuid.UUID) (*model.Project, error)
	UpdateProject(project *model.Project) error
	DeleteProject(projectID uuid.UUID) error

	// 需求分析相关
	CreateRequirementAnalysis(requirement *model.Requirement) error
	GetRequirementByProjectID(projectID uuid.UUID) (*model.Requirement, error)
	UpdateRequirementAnalysis(requirement *model.Requirement) error

	// 对话相关
	CreateChatSession(session *model.ChatSession) error
	GetChatSessionByProjectID(projectID uuid.UUID) (*model.ChatSession, error)
	CreateChatMessage(message *model.ChatMessage) error
	GetChatMessagesBySessionID(sessionID uuid.UUID, page, pageSize int) ([]*model.ChatMessage, int64, error)
	EndChatSession(sessionID uuid.UUID) error

	// 问题相关
	CreateQuestion(question *model.Question) error
	GetQuestionsByRequirementID(requirementID uuid.UUID) ([]*model.Question, error)
	AnswerQuestion(questionID uuid.UUID, answer string) error

	// PUML图表相关
	CreatePUMLDiagram(diagram *model.PUMLDiagram) error
	GetPUMLDiagramsByProjectID(projectID uuid.UUID) ([]*model.PUMLDiagram, error)
	UpdatePUMLDiagram(diagram *model.PUMLDiagram) error
	DeletePUMLDiagram(diagramID uuid.UUID) error

	// 文档相关
	CreateDocument(document *model.Document) error
	GetDocumentsByProjectID(projectID uuid.UUID) ([]*model.Document, error)
	UpdateDocument(document *model.Document) error
	DeleteDocument(documentID uuid.UUID) error

	// 业务模块相关
	CreateBusinessModule(module *model.BusinessModule) error
	GetBusinessModulesByProjectID(projectID uuid.UUID) ([]*model.BusinessModule, error)
	UpdateBusinessModule(module *model.BusinessModule) error
	DeleteBusinessModule(moduleID uuid.UUID) error

	// 通用模块库相关
	CreateCommonModule(module *model.CommonModule) error
	GetCommonModulesByCategory(category string, page, pageSize int) ([]*model.CommonModule, int64, error)
	GetCommonModuleByID(moduleID uuid.UUID) (*model.CommonModule, error)
	UpdateCommonModule(module *model.CommonModule) error
	DeleteCommonModule(moduleID uuid.UUID) error

	// 异步任务相关
	CreateAsyncTask(task *model.AsyncTask) error
	GetAsyncTask(taskID uuid.UUID) (*model.AsyncTask, error)
	UpdateAsyncTask(task *model.AsyncTask) error
	GetTasksByProject(projectID uuid.UUID, taskType string) ([]*model.AsyncTask, error)

	// 阶段进度相关
	CreateStageProgress(progress *model.StageProgress) error
	GetStageProgress(projectID uuid.UUID) ([]*model.StageProgress, error)
	UpdateStageProgress(progress *model.StageProgress) error
	GetStageProgressByStage(projectID uuid.UUID, stage int) (*model.StageProgress, error)

	// 健康检查
	Health() error
}

// MySQLRepository MySQL仓库实现
type MySQLRepository struct {
	db *Database
}

// NewMySQLRepository 创建MySQL仓库
func NewMySQLRepository(db *Database) Repository {
	return &MySQLRepository{db: db}
}

// convertNullTime 转换NULL时间
func convertNullTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

// convertTimePtr 转换时间指针
func convertTimePtr(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

// scanUser 扫描用户数据
func scanUser(row *sql.Row) (*model.User, error) {
	var user model.User
	var lastLogin sql.NullTime
	var preferences sql.NullString

	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
		&user.Status,
		&preferences,
	)
	if err != nil {
		return nil, err
	}

	user.LastLogin = convertNullTime(lastLogin)
	if preferences.Valid {
		user.Preferences = preferences.String
	}

	return &user, nil
}

// scanUsers 扫描多个用户数据
func scanUsers(rows *sql.Rows) ([]*model.User, error) {
	var users []*model.User

	for rows.Next() {
		var user model.User
		var lastLogin sql.NullTime
		var preferences sql.NullString

		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&lastLogin,
			&user.Status,
			&preferences,
		)
		if err != nil {
			return nil, err
		}

		user.LastLogin = convertNullTime(lastLogin)
		if preferences.Valid {
			user.Preferences = preferences.String
		}

		users = append(users, &user)
	}

	return users, nil
}
