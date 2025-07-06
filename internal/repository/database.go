package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	mysql, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 配置连接池
	mysql.SetMaxOpenConns(cfg.Database.MaxConnections)
	mysql.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	mysql.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := mysql.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL连接测试失败: %w", err)
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
		log.Printf("Redis连接失败: %v", err)
		// Redis连接失败不影响主要功能，只记录日志
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

// CreateTables 创建数据表（用于开发环境）
func (db *Database) CreateTables() error {
	tables := []string{
		// 用户表
		`CREATE TABLE IF NOT EXISTS users (
			user_id CHAR(36) PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			last_login TIMESTAMP NULL,
			status VARCHAR(20) DEFAULT 'active',
			preferences JSON,
			INDEX idx_username (username),
			INDEX idx_email (email),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 项目表
		`CREATE TABLE IF NOT EXISTS projects (
			project_id CHAR(36) PRIMARY KEY,
			user_id CHAR(36) NOT NULL,
			project_name VARCHAR(100) NOT NULL,
			description TEXT,
			project_type VARCHAR(50) NOT NULL,
			status VARCHAR(20) DEFAULT 'draft',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			completion_percentage INT DEFAULT 0,
			settings JSON,
			INDEX idx_user_id (user_id),
			INDEX idx_status (status),
			INDEX idx_project_type (project_type),
			INDEX idx_created_at (created_at),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 需求分析表
		`CREATE TABLE IF NOT EXISTS requirement_analyses (
			requirement_id CHAR(36) PRIMARY KEY,
			project_id CHAR(36) NOT NULL,
			raw_requirement TEXT NOT NULL,
			structured_requirement JSON,
			completeness_score DECIMAL(3,2) DEFAULT 0.00,
			analysis_status VARCHAR(50) DEFAULT 'pending',
			missing_info_types JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_project_id (project_id),
			INDEX idx_analysis_status (analysis_status),
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 对话会话表
		`CREATE TABLE IF NOT EXISTS chat_sessions (
			session_id CHAR(36) PRIMARY KEY,
			project_id CHAR(36) NOT NULL,
			user_id CHAR(36) NOT NULL,
			session_type VARCHAR(50) NOT NULL,
			started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			ended_at TIMESTAMP NULL,
			status VARCHAR(20) DEFAULT 'active',
			context JSON,
			INDEX idx_project_id (project_id),
			INDEX idx_user_id (user_id),
			INDEX idx_session_type (session_type),
			INDEX idx_status (status),
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 对话消息表
		`CREATE TABLE IF NOT EXISTS chat_messages (
			message_id CHAR(36) PRIMARY KEY,
			session_id CHAR(36) NOT NULL,
			sender_type VARCHAR(20) NOT NULL,
			message_content TEXT NOT NULL,
			message_type VARCHAR(50) NOT NULL,
			metadata JSON,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			processed BOOLEAN DEFAULT FALSE,
			INDEX idx_session_id (session_id),
			INDEX idx_sender_type (sender_type),
			INDEX idx_message_type (message_type),
			INDEX idx_timestamp (timestamp),
			FOREIGN KEY (session_id) REFERENCES chat_sessions(session_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 补充问题表
		`CREATE TABLE IF NOT EXISTS questions (
			question_id CHAR(36) PRIMARY KEY,
			requirement_id CHAR(36) NOT NULL,
			question_text TEXT NOT NULL,
			question_category VARCHAR(50) NOT NULL,
			priority_level INT DEFAULT 1,
			answer_text TEXT,
			answer_status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			answered_at TIMESTAMP NULL,
			INDEX idx_requirement_id (requirement_id),
			INDEX idx_question_category (question_category),
			INDEX idx_priority_level (priority_level),
			INDEX idx_answer_status (answer_status),
			FOREIGN KEY (requirement_id) REFERENCES requirement_analyses(requirement_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// PUML图表表
		`CREATE TABLE IF NOT EXISTS puml_diagrams (
			diagram_id CHAR(36) PRIMARY KEY,
			project_id CHAR(36) NOT NULL,
			diagram_type VARCHAR(50) NOT NULL,
			diagram_name VARCHAR(100) NOT NULL,
			puml_content TEXT NOT NULL,
			rendered_url VARCHAR(500),
			version INT DEFAULT 1,
			is_validated BOOLEAN DEFAULT FALSE,
			validation_feedback TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_project_id (project_id),
			INDEX idx_diagram_type (diagram_type),
			INDEX idx_version (version),
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 业务模块表
		`CREATE TABLE IF NOT EXISTS business_modules (
			module_id CHAR(36) PRIMARY KEY,
			project_id CHAR(36) NOT NULL,
			module_name VARCHAR(100) NOT NULL,
			description TEXT,
			module_type VARCHAR(50) NOT NULL,
			complexity_level VARCHAR(20) DEFAULT 'medium',
			business_logic JSON,
			interfaces JSON,
			dependencies JSON,
			is_reusable BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_project_id (project_id),
			INDEX idx_module_type (module_type),
			INDEX idx_complexity_level (complexity_level),
			INDEX idx_is_reusable (is_reusable),
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 通用模块库表
		`CREATE TABLE IF NOT EXISTS common_module_library (
			common_module_id CHAR(36) PRIMARY KEY,
			module_name VARCHAR(100) NOT NULL,
			category VARCHAR(50) NOT NULL,
			description TEXT,
			functionality JSON,
			interface_spec JSON,
			code_template TEXT,
			usage_examples JSON,
			version VARCHAR(20) DEFAULT '1.0.0',
			downloads_count INT DEFAULT 0,
			rating DECIMAL(2,1) DEFAULT 0.0,
			tags JSON,
			created_by CHAR(36) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_module_name (module_name),
			INDEX idx_category (category),
			INDEX idx_rating (rating),
			INDEX idx_created_by (created_by),
			FOREIGN KEY (created_by) REFERENCES users(user_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// 生成文档表
		`CREATE TABLE IF NOT EXISTS generated_documents (
			document_id CHAR(36) PRIMARY KEY,
			project_id CHAR(36) NOT NULL,
			document_type VARCHAR(50) NOT NULL,
			document_name VARCHAR(100) NOT NULL,
			content TEXT NOT NULL,
			format VARCHAR(20) DEFAULT 'markdown',
			file_path VARCHAR(500),
			version INT DEFAULT 1,
			generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_final BOOLEAN DEFAULT FALSE,
			INDEX idx_project_id (project_id),
			INDEX idx_document_type (document_type),
			INDEX idx_version (version),
			FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for _, table := range tables {
		if _, err := db.MySQL.Exec(table); err != nil {
			return fmt.Errorf("创建表失败: %w", err)
		}
	}

	log.Println("数据库表创建成功")
	return nil
}

// Health 检查数据库健康状态
func (db *Database) Health() error {
	// 检查MySQL
	if err := db.MySQL.Ping(); err != nil {
		return fmt.Errorf("MySQL连接异常: %w", err)
	}

	// 检查Redis
	ctx := context.Background()
	if err := db.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接异常: %w", err)
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