package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-dev-platform/internal/ai"
	"ai-dev-platform/internal/api"
	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/repository"
	"ai-dev-platform/internal/service"
)

func main() {
	// 加载配置
	cfg := config.Load()

	log.Printf("启动AI开发平台服务器 [环境: %s] [端口: %s]", cfg.Env, cfg.Port)

	// 初始化数据库
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		}
	}()

	// 在开发环境下创建数据表
	if cfg.IsDevelopment() {
		if err := db.CreateTables(); err != nil {
			log.Fatalf("创建数据表失败: %v", err)
		}
		log.Println("数据表初始化完成")
	}

	// 初始化仓库
	repo := repository.NewMySQLRepository(db)
	aiRepo := repository.NewAIRepository(db.MySQL)

	// 初始化AI管理器
	aiManagerConfig := ai.AIManagerConfig{
		DefaultProvider: ai.AIProvider(cfg.AI.DefaultProvider),
		OpenAIConfig: &ai.OpenAIConfig{
			APIKey:  cfg.AI.OpenAIConfig.APIKey,
			BaseURL: os.Getenv("OPENAI_BASE_URL"),
			Model:   cfg.AI.OpenAIConfig.DefaultModel,
		},
		EnableCache: cfg.AI.EnableCache,
		CacheTTL:    cfg.AI.CacheTTL,
	}

	aiManager, err := ai.NewAIManager(aiManagerConfig)
	if err != nil {
		log.Printf("AI管理器初始化失败，将以有限功能模式运行: %v", err)
		// 创建一个空的AI管理器，确保服务不中断
		aiManager, _ = ai.NewAIManager(ai.AIManagerConfig{
			DefaultProvider: ai.ProviderOpenAI,
			EnableCache:     false,
		})
	}

	// 初始化服务
	userService := service.NewUserService(repo, cfg)
	projectService := service.NewProjectService(repo)
	aiService := service.NewAIService(aiManager, aiRepo, repo.(*repository.MySQLRepository))

	// 初始化PUML渲染服务
	pumlService := service.NewPUMLService(&cfg.PUML)

	// 初始化异步任务服务
	asyncTaskService := service.NewAsyncTaskService(repo, aiService, aiManager)

	// 初始化路由器
	router := api.NewRouter(cfg, userService, projectService, aiService, pumlService, asyncTaskService)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.GetHandler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("服务器启动在端口 :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("服务器强制关闭: %v", err)
	} else {
		log.Println("服务器已优雅关闭")
	}
}
