package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"ai-dev-platform/internal/ai"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/repository"

	"github.com/google/uuid"
)

// AsyncTaskService 异步任务管理服务
type AsyncTaskService struct {
	repo         repository.Repository
	aiService    *AIService
	executors    map[string]TaskExecutor
	mu           sync.RWMutex
	activeAIManager *ai.AIManager // 默认AI管理器
}

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	Execute(ctx context.Context, task *model.AsyncTask) error
}

// NewAsyncTaskService 创建异步任务管理服务
func NewAsyncTaskService(repo repository.Repository, aiService *AIService, aiManager *ai.AIManager) *AsyncTaskService {
	service := &AsyncTaskService{
		repo:         repo,
		aiService:    aiService,
		executors:    make(map[string]TaskExecutor),
		activeAIManager: aiManager,
	}

	// 注册任务执行器
	service.registerExecutors()

	return service
}

// registerExecutors 注册任务执行器
func (s *AsyncTaskService) registerExecutors() {
	s.executors[model.TaskTypeStageDocuments] = &StageDocumentExecutor{
		service:   s,
		aiService: s.aiService,
	}
	s.executors[model.TaskTypePUMLGeneration] = &PUMLGenerationExecutor{
		service:   s,
		aiService: s.aiService,
	}
	s.executors[model.TaskTypeDocumentGeneration] = &DocumentGenerationExecutor{
		service:   s,
		aiService: s.aiService,
	}
	s.executors[model.TaskTypeRequirementAnalysis] = &RequirementAnalysisExecutor{
		service:   s,
		aiService: s.aiService,
	}
	s.executors[model.TaskTypeCompleteProjectDocuments] = &CompleteProjectDocumentsExecutor{
		service:   s,
		aiService: s.aiService,
	}
}

// StartStageDocumentGeneration 启动阶段文档生成任务
func (s *AsyncTaskService) StartStageDocumentGeneration(projectID uuid.UUID, userID uuid.UUID, stage int) (*model.AsyncTaskResponse, error) {
	// 创建任务
	task := &model.AsyncTask{
		TaskID:    uuid.New(),
		UserID:    userID,
		ProjectID: projectID,
		TaskType:  model.TaskTypeStageDocuments,
		TaskName:  fmt.Sprintf("阶段%d文档生成", stage),
		Status:    model.TaskStatusPending,
		Progress:  0,
		CreatedAt: time.Now(),
		Metadata:  fmt.Sprintf(`{"stage": %d}`, stage),
	}

	// 保存任务到数据库
	if err := s.repo.CreateAsyncTask(task); err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	// 初始化或更新阶段进度
	if err := s.initializeStageProgress(projectID, stage, task.TaskID); err != nil {
		log.Printf("初始化阶段进度失败: %v", err)
	}

	// 异步执行任务
	go s.executeTask(context.Background(), task)

	return &model.AsyncTaskResponse{
		TaskID:   task.TaskID,
		Status:   task.Status,
		Progress: task.Progress,
		Message:  "任务已启动",
	}, nil
}

// StartCompleteProjectDocumentGeneration 启动完整项目文档生成任务
func (s *AsyncTaskService) StartCompleteProjectDocumentGeneration(projectID uuid.UUID, userID uuid.UUID) (*model.AsyncTaskResponse, error) {
	// 创建任务
	task := &model.AsyncTask{
		TaskID:    uuid.New(),
		UserID:    userID,
		ProjectID: projectID,
		TaskType:  model.TaskTypeCompleteProjectDocuments,
		TaskName:  "完整项目文档生成",
		Status:    model.TaskStatusPending,
		Progress:  0,
		CreatedAt: time.Now(),
		Metadata:  `{"stages": [1, 2, 3]}`,
	}

	// 保存任务到数据库
	if err := s.repo.CreateAsyncTask(task); err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	// 初始化所有阶段进度
	for stage := 1; stage <= 3; stage++ {
		if err := s.initializeStageProgress(projectID, stage, task.TaskID); err != nil {
			log.Printf("初始化阶段%d进度失败: %v", stage, err)
		}
	}

	// 异步执行任务
	go s.executeTask(context.Background(), task)

	return &model.AsyncTaskResponse{
		TaskID:   task.TaskID,
		Status:   task.Status,
		Progress: task.Progress,
		Message:  "完整项目文档生成任务已启动",
	}, nil
}

// GetTaskStatus 获取任务状态
func (s *AsyncTaskService) GetTaskStatus(taskID uuid.UUID) (*model.AsyncTaskResponse, error) {
	task, err := s.repo.GetAsyncTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}

	response := &model.AsyncTaskResponse{
		TaskID:   task.TaskID,
		Status:   task.Status,
		Progress: task.Progress,
	}

	if task.Status == model.TaskStatusCompleted && task.ResultData != "" {
		response.Message = "任务完成"
	} else if task.Status == model.TaskStatusFailed && task.ErrorMessage != "" {
		response.Message = task.ErrorMessage
	} else if task.Status == model.TaskStatusRunning {
		response.Message = "任务正在执行中..."
	}

	return response, nil
}

// GetStageProgress 获取阶段进度
func (s *AsyncTaskService) GetStageProgress(projectID uuid.UUID) (*model.StageProgressResponse, error) {
	progresses, err := s.repo.GetStageProgress(projectID)
	if err != nil {
		return nil, fmt.Errorf("获取阶段进度失败: %w", err)
	}

	// 如果没有进度记录，初始化三个阶段
	if len(progresses) == 0 {
		for stage := 1; stage <= 3; stage++ {
			progress := &model.StageProgress{
				ProgressID:     uuid.New(),
				ProjectID:      projectID,
				Stage:          stage,
				Status:         model.StageStatusNotStarted,
				CompletionRate: 0,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := s.repo.CreateStageProgress(progress); err != nil {
				log.Printf("初始化阶段%d进度失败: %v", stage, err)
			} else {
				progresses = append(progresses, progress)
			}
		}
	}

	response := &model.StageProgressResponse{
		ProjectID: projectID,
		Stages:    progresses,
	}

	// 计算总体进度
	totalCompletion := 0
	allCompleted := true
	anyStarted := false

	for _, stage := range progresses {
		totalCompletion += stage.CompletionRate
		if stage.Status != model.StageStatusCompleted {
			allCompleted = false
		}
		if stage.Status != model.StageStatusNotStarted {
			anyStarted = true
		}
	}

	response.Overall.CompletionRate = totalCompletion / len(progresses)
	if allCompleted {
		response.Overall.Status = model.StageStatusCompleted
	} else if anyStarted {
		response.Overall.Status = model.StageStatusInProgress
	} else {
		response.Overall.Status = model.StageStatusNotStarted
	}

	return response, nil
}

// executeTask 执行任务
func (s *AsyncTaskService) executeTask(ctx context.Context, task *model.AsyncTask) {
	// 更新任务状态为运行中
	task.Status = model.TaskStatusRunning
	task.Progress = 0
	now := time.Now()
	task.StartedAt = &now

	if err := s.repo.UpdateAsyncTask(task); err != nil {
		log.Printf("更新任务状态失败: %v", err)
		return
	}

	// 获取对应的执行器
	executor, exists := s.executors[task.TaskType]
	if !exists {
		s.markTaskFailed(task, fmt.Sprintf("不支持的任务类型: %s", task.TaskType))
		return
	}

	// 执行任务
	if err := executor.Execute(ctx, task); err != nil {
		s.markTaskFailed(task, err.Error())
		return
	}

	// 标记任务完成
	s.markTaskCompleted(task)
}

// markTaskFailed 标记任务失败
func (s *AsyncTaskService) markTaskFailed(task *model.AsyncTask, errorMsg string) {
	task.Status = model.TaskStatusFailed
	task.Progress = 0
	task.ErrorMessage = errorMsg
	now := time.Now()
	task.CompletedAt = &now

	if err := s.repo.UpdateAsyncTask(task); err != nil {
		log.Printf("更新失败任务状态失败: %v", err)
	}

	// 更新阶段进度
	s.updateStageProgressOnFailure(task)
}

// markTaskCompleted 标记任务完成
func (s *AsyncTaskService) markTaskCompleted(task *model.AsyncTask) {
	task.Status = model.TaskStatusCompleted
	task.Progress = 100
	now := time.Now()
	task.CompletedAt = &now

	if err := s.repo.UpdateAsyncTask(task); err != nil {
		log.Printf("更新完成任务状态失败: %v", err)
	}

	// 更新阶段进度
	s.updateStageProgressOnCompletion(task)
}

// initializeStageProgress 初始化阶段进度
func (s *AsyncTaskService) initializeStageProgress(projectID uuid.UUID, stage int, taskID uuid.UUID) error {
	// 检查是否已存在
	existing, err := s.repo.GetStageProgressByStage(projectID, stage)
	if err == nil {
		// 更新现有进度
		existing.Status = model.StageStatusInProgress
		existing.LastTaskID = &taskID
		now := time.Now()
		existing.StartedAt = &now
		existing.UpdatedAt = now
		return s.repo.UpdateStageProgress(existing)
	}

	// 创建新的进度记录
	progress := &model.StageProgress{
		ProgressID:     uuid.New(),
		ProjectID:      projectID,
		Stage:          stage,
		Status:         model.StageStatusInProgress,
		CompletionRate: 0,
		LastTaskID:     &taskID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	now := time.Now()
	progress.StartedAt = &now

	return s.repo.CreateStageProgress(progress)
}

// updateStageProgressOnCompletion 任务完成时更新阶段进度
func (s *AsyncTaskService) updateStageProgressOnCompletion(task *model.AsyncTask) {
	// 解析stage
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(task.Metadata), &metadata); err != nil {
		log.Printf("解析任务元数据失败: %v", err)
		return
	}

	stage, ok := metadata["stage"].(float64)
	if !ok {
		log.Printf("无法从任务元数据中获取阶段信息")
		return
	}

	progress, err := s.repo.GetStageProgressByStage(task.ProjectID, int(stage))
	if err != nil {
		log.Printf("获取阶段进度失败: %v", err)
		return
	}

	// 更新进度
	progress.Status = model.StageStatusCompleted
	progress.CompletionRate = 100
	progress.LastTaskID = &task.TaskID
	now := time.Now()
	progress.CompletedAt = &now
	progress.UpdatedAt = now

	// 统计生成的文档和图表数量
	s.updateDocumentCounts(progress)

	if err := s.repo.UpdateStageProgress(progress); err != nil {
		log.Printf("更新阶段进度失败: %v", err)
	}
}

// updateStageProgressOnFailure 任务失败时更新阶段进度
func (s *AsyncTaskService) updateStageProgressOnFailure(task *model.AsyncTask) {
	// 解析stage
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(task.Metadata), &metadata); err != nil {
		log.Printf("解析任务元数据失败: %v", err)
		return
	}

	stage, ok := metadata["stage"].(float64)
	if !ok {
		log.Printf("无法从任务元数据中获取阶段信息")
		return
	}

	progress, err := s.repo.GetStageProgressByStage(task.ProjectID, int(stage))
	if err != nil {
		log.Printf("获取阶段进度失败: %v", err)
		return
	}

	// 更新进度
	progress.Status = model.StageStatusFailed
	progress.CompletionRate = 0
	progress.LastTaskID = &task.TaskID
	progress.UpdatedAt = time.Now()

	if err := s.repo.UpdateStageProgress(progress); err != nil {
		log.Printf("更新阶段进度失败: %v", err)
	}
}

// updateDocumentCounts 更新文档和图表数量
func (s *AsyncTaskService) updateDocumentCounts(progress *model.StageProgress) {
	// 获取阶段相关的文档数量
	documents, err := s.repo.GetDocumentsByProjectID(progress.ProjectID)
	if err != nil {
		log.Printf("获取文档列表失败: %v", err)
		return
	}

	// 获取阶段相关的PUML图表数量
	diagrams, err := s.repo.GetPUMLDiagramsByProjectID(progress.ProjectID)
	if err != nil {
		log.Printf("获取PUML图表列表失败: %v", err)
		return
	}

	// 统计当前阶段的文档和图表
	docCount := 0
	pumlCount := 0

	for _, doc := range documents {
		if doc.Stage == progress.Stage {
			docCount++
		}
	}

	for _, diagram := range diagrams {
		if diagram.Stage == progress.Stage {
			pumlCount++
		}
	}

	progress.DocumentCount = docCount
	progress.PUMLCount = pumlCount
}

// ===== 具体的任务执行器实现 =====

// StageDocumentExecutor 阶段文档生成执行器
type StageDocumentExecutor struct {
	service   *AsyncTaskService
	aiService *AIService
}

func (e *StageDocumentExecutor) Execute(ctx context.Context, task *model.AsyncTask) error {
	// 解析任务元数据
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(task.Metadata), &metadata); err != nil {
		return fmt.Errorf("解析任务元数据失败: %w", err)
	}

	stage, ok := metadata["stage"].(float64)
	if !ok {
		return fmt.Errorf("任务元数据中缺少阶段信息")
	}

	// 创建阶段文档生成请求
	req := &model.GenerateStageDocumentsRequest{
		ProjectID: task.ProjectID,
		Stage:     int(stage),
	}

	// 更新进度到50%
	task.Progress = 50
	if err := e.service.repo.UpdateAsyncTask(task); err != nil {
		log.Printf("更新任务进度失败: %v", err)
	}

	// 调用AI服务生成文档
	result, err := e.aiService.GenerateStageDocuments(ctx, req, task.UserID)
	if err != nil {
		return fmt.Errorf("生成阶段文档失败: %w", err)
	}

	// 为生成的文档和图表添加阶段和任务ID
	for _, doc := range result.Documents {
		doc.Stage = int(stage)
		doc.TaskID = &task.TaskID
		if err := e.service.repo.UpdateDocument(doc); err != nil {
			log.Printf("更新文档阶段信息失败: %v", err)
		}
	}

	for _, diagram := range result.PUMLDiagrams {
		diagram.Stage = int(stage)
		diagram.TaskID = &task.TaskID
		if err := e.service.repo.UpdatePUMLDiagram(diagram); err != nil {
			log.Printf("更新PUML图表阶段信息失败: %v", err)
		}
	}

	// 序列化结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("序列化结果失败: %w", err)
	}

	task.ResultData = string(resultJSON)
	task.Progress = 100

	return nil
}

// PUMLGenerationExecutor PUML生成执行器
type PUMLGenerationExecutor struct {
	service   *AsyncTaskService
	aiService *AIService
}

func (e *PUMLGenerationExecutor) Execute(ctx context.Context, task *model.AsyncTask) error {
	// TODO: 实现PUML生成逻辑
	task.Progress = 100
	return nil
}

// DocumentGenerationExecutor 文档生成执行器
type DocumentGenerationExecutor struct {
	service   *AsyncTaskService
	aiService *AIService
}

func (e *DocumentGenerationExecutor) Execute(ctx context.Context, task *model.AsyncTask) error {
	// TODO: 实现文档生成逻辑
	task.Progress = 100
	return nil
}

// RequirementAnalysisExecutor 需求分析执行器
type RequirementAnalysisExecutor struct {
	service   *AsyncTaskService
	aiService *AIService
}

func (e *RequirementAnalysisExecutor) Execute(ctx context.Context, task *model.AsyncTask) error {
	// TODO: 实现需求分析逻辑
	task.Progress = 100
	return nil
}

// CompleteProjectDocumentsExecutor 完整项目文档生成执行器
type CompleteProjectDocumentsExecutor struct {
	service   *AsyncTaskService
	aiService *AIService
}

func (e *CompleteProjectDocumentsExecutor) Execute(ctx context.Context, task *model.AsyncTask) error {
	log.Printf("开始执行完整项目文档生成任务: %s", task.TaskID)
	
	// 第一阶段：项目需求文档 + 系统架构图 + 交互流程图 + 业务流程图 (4份)
	task.Progress = 10
	if err := e.service.repo.UpdateAsyncTask(task); err != nil {
		log.Printf("更新任务进度失败: %v", err)
	}
	
	stage1Req := &model.GenerateStageDocumentsRequest{
		ProjectID: task.ProjectID,
		Stage:     1,
	}
	
	log.Printf("生成第一阶段文档...")
	stage1Result, err := e.aiService.GenerateSpecificStageDocuments(ctx, stage1Req, task.UserID, []string{
		"项目需求文档.md",
		"系统架构图.puml", 
		"交互流程图.puml",
		"业务流程图.puml",
	})
	if err != nil {
		return fmt.Errorf("生成第一阶段文档失败: %w", err)
	}
	
	// 标记文档和图表的阶段信息
	for _, doc := range stage1Result.Documents {
		doc.Stage = 1
		doc.TaskID = &task.TaskID
		if err := e.service.repo.UpdateDocument(doc); err != nil {
			log.Printf("更新第一阶段文档信息失败: %v", err)
		}
	}
	
	for _, diagram := range stage1Result.PUMLDiagrams {
		diagram.Stage = 1
		diagram.TaskID = &task.TaskID
		if err := e.service.repo.UpdatePUMLDiagram(diagram); err != nil {
			log.Printf("更新第一阶段PUML图表信息失败: %v", err)
		}
	}
	
	// 第二阶段：技术规范文档 + API设计 + 数据库设计 (3份)
	task.Progress = 40
	if err := e.service.repo.UpdateAsyncTask(task); err != nil {
		log.Printf("更新任务进度失败: %v", err)
	}
	
	stage2Req := &model.GenerateStageDocumentsRequest{
		ProjectID: task.ProjectID,
		Stage:     2,
	}
	
	log.Printf("生成第二阶段文档...")
	stage2Result, err := e.aiService.GenerateSpecificStageDocuments(ctx, stage2Req, task.UserID, []string{
		"技术规范文档.md",
		"API设计.md",
		"数据库设计.md",
	})
	if err != nil {
		return fmt.Errorf("生成第二阶段文档失败: %w", err)
	}
	
	// 标记文档的阶段信息
	for _, doc := range stage2Result.Documents {
		doc.Stage = 2
		doc.TaskID = &task.TaskID
		if err := e.service.repo.UpdateDocument(doc); err != nil {
			log.Printf("更新第二阶段文档信息失败: %v", err)
		}
	}
	
	// 第三阶段：开发流程文档 + 测试用例文档 + 部署文档 (3份)
	task.Progress = 70
	if err := e.service.repo.UpdateAsyncTask(task); err != nil {
		log.Printf("更新任务进度失败: %v", err)
	}
	
	stage3Req := &model.GenerateStageDocumentsRequest{
		ProjectID: task.ProjectID,
		Stage:     3,
	}
	
	log.Printf("生成第三阶段文档...")
	stage3Result, err := e.aiService.GenerateSpecificStageDocuments(ctx, stage3Req, task.UserID, []string{
		"开发流程文档.md",
		"测试用例文档.md",
		"部署文档.md",
	})
	if err != nil {
		return fmt.Errorf("生成第三阶段文档失败: %w", err)
	}
	
	// 标记文档的阶段信息
	for _, doc := range stage3Result.Documents {
		doc.Stage = 3
		doc.TaskID = &task.TaskID
		if err := e.service.repo.UpdateDocument(doc); err != nil {
			log.Printf("更新第三阶段文档信息失败: %v", err)
		}
	}
	
	// 汇总所有结果
	totalResult := &model.CompleteProjectDocumentsResult{
		ProjectID:    task.ProjectID,
		GeneratedAt:  time.Now(),
		Stage1:       stage1Result,
		Stage2:       stage2Result,
		Stage3:       stage3Result,
		TotalDocuments: len(stage1Result.Documents) + len(stage2Result.Documents) + len(stage3Result.Documents),
		TotalPUMLDiagrams: len(stage1Result.PUMLDiagrams) + len(stage2Result.PUMLDiagrams) + len(stage3Result.PUMLDiagrams),
	}
	
	// 序列化结果
	resultJSON, err := json.Marshal(totalResult)
	if err != nil {
		return fmt.Errorf("序列化结果失败: %w", err)
	}
	
	task.ResultData = string(resultJSON)
	task.Progress = 100
	
	log.Printf("完整项目文档生成任务完成: %s，共生成 %d 份文档，%d 个PUML图表", 
		task.TaskID, totalResult.TotalDocuments, totalResult.TotalPUMLDiagrams)
	
	return nil
} 