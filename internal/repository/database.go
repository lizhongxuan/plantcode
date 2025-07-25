package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库连接结构
type Database struct {
	GORM  *gorm.DB
	Redis *redis.Client
}

// NewDatabase 创建数据库连接
func NewDatabase(cfg *config.Config) (*Database, error) {
	var gormDB *gorm.DB
	var err error

	// 配置 GORM 日志级别
	var gormLogLevel logger.LogLevel
	if cfg.IsDevelopment() {
		gormLogLevel = logger.Info
	} else {
		gormLogLevel = logger.Error
	}

	// 使用 GORM 连接 MySQL
	log.Printf("连接MySQL (GORM): %v", cfg.GetDSN())
	gormDB, err = gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("GORM MySQL连接失败: %w", err)
	}

	// 获取底层的 *sql.DB 以配置连接池
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取GORM底层SQL连接失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

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
		if cfg.IsProduction() {
			return nil, fmt.Errorf("生产环境: Redis连接失败: %w", err)
		}
		redisClient = nil
	}

	return &Database{
		GORM:  gormDB,
		Redis: redisClient,
	}, nil
}

// Close 关闭数据库连接
func (db *Database) Close() error {
	var err error

	// 关闭 GORM 连接
	if db.GORM != nil {
		sqlDB, gormErr := db.GORM.DB()
		if gormErr == nil {
			if mysqlErr := sqlDB.Close(); mysqlErr != nil {
				err = fmt.Errorf("关闭GORM MySQL连接失败: %w", mysqlErr)
			}
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
	if db.GORM == nil {
		return fmt.Errorf("GORM数据库连接不可用")
	}

	// 使用 GORM 自动迁移创建表
	err := db.GORM.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Requirement{},
		&model.ChatSession{},
		&model.ChatMessage{},
		&model.Question{},
		&model.PUMLDiagram{},
		&model.Document{},
		&model.BusinessModule{},
		&model.CommonModule{},
		&model.StageProgress{},
		&model.UserAIConfig{},
		&model.AsyncTask{},
	)
	if err != nil {
		return fmt.Errorf("GORM 自动迁移失败: %w", err)
	}

	log.Println("GORM 数据表自动迁移完成")
	return nil
}

// Health 检查数据库健康状态
func (db *Database) Health() error {
	var errors []string

	// 检查GORM连接
	if db.GORM == nil {
		errors = append(errors, "GORM连接不存在")
	} else {
		sqlDB, err := db.GORM.DB()
		if err != nil {
			errors = append(errors, fmt.Sprintf("获取GORM底层连接失败: %v", err))
		} else if err := sqlDB.Ping(); err != nil {
			errors = append(errors, fmt.Sprintf("GORM MySQL连接异常: %v", err))
		}
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
		return fmt.Errorf("数据库健康检查失败: %s", fmt.Sprintf("%v", errors))
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

	// 用户AI配置相关
	GetUserAIConfig(userID uuid.UUID) (*model.UserAIConfig, error)
	CreateUserAIConfig(config *model.UserAIConfig) error
	UpdateUserAIConfig(config *model.UserAIConfig) error
	DeleteUserAIConfig(userID uuid.UUID) error

	// 扩展方法（用于兼容性）
	GetRequirementAnalysis(analysisID uuid.UUID) (*model.Requirement, error)
	GetRequirementAnalysesByProject(projectID uuid.UUID) ([]*model.Requirement, error)
	GetChatSession(sessionID uuid.UUID) (*model.ChatSession, error)
	GetChatSessionsByProject(projectID uuid.UUID) ([]*model.ChatSession, error)
	GetChatMessages(sessionID uuid.UUID) ([]*model.ChatMessage, error)
	GetPUMLDiagram(diagramID uuid.UUID) (*model.PUMLDiagram, error)
	GetDocument(documentID uuid.UUID) (*model.Document, error)
	GetQuestions(requirementID uuid.UUID) ([]*model.Question, error)

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
