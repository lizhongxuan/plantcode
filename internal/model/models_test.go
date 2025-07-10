package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelsTestSuite struct {
	suite.Suite
}

func (suite *ModelsTestSuite) TestUser_JSONSerialization() {
	// Arrange
	user := &User{
		UserID:      uuid.New(),
		Username:    "testuser",
		Email:       "test@example.com",
		PasswordHash: "hashed_password",
		FullName:    "Test User",
		Status:      "active",
		Preferences: `{"theme": "dark"}`,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act - 序列化
	jsonData, err := json.Marshal(user)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedUser User
	err = json.Unmarshal(jsonData, &deserializedUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.UserID, deserializedUser.UserID)
	assert.Equal(suite.T(), user.Username, deserializedUser.Username)
	assert.Equal(suite.T(), user.Email, deserializedUser.Email)
	// PasswordHash 字段有 json:"-" 标签，不会被序列化，所以反序列化后应该是空字符串
	assert.Empty(suite.T(), deserializedUser.PasswordHash)
	assert.Equal(suite.T(), user.FullName, deserializedUser.FullName)
	assert.Equal(suite.T(), user.Status, deserializedUser.Status)
	assert.Equal(suite.T(), user.Preferences, deserializedUser.Preferences)
}

func (suite *ModelsTestSuite) TestProject_JSONSerialization() {
	// Arrange
	project := &Project{
		ProjectID:            uuid.New(),
		UserID:               uuid.New(),
		ProjectName:          "测试项目",
		Description:          "这是一个测试项目",
		ProjectType:          "web",
		Status:               "planning",
		CompletionPercentage: 25,
		Settings:             `{"autoSave": true}`,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Act - 序列化
	jsonData, err := json.Marshal(project)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedProject Project
	err = json.Unmarshal(jsonData, &deserializedProject)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), project.ProjectID, deserializedProject.ProjectID)
	assert.Equal(suite.T(), project.UserID, deserializedProject.UserID)
	assert.Equal(suite.T(), project.ProjectName, deserializedProject.ProjectName)
	assert.Equal(suite.T(), project.Description, deserializedProject.Description)
	assert.Equal(suite.T(), project.ProjectType, deserializedProject.ProjectType)
	assert.Equal(suite.T(), project.Status, deserializedProject.Status)
	assert.Equal(suite.T(), project.CompletionPercentage, deserializedProject.CompletionPercentage)
	assert.Equal(suite.T(), project.Settings, deserializedProject.Settings)
}

func (suite *ModelsTestSuite) TestCreateUserRequest_Fields() {
	// Arrange
	req := CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	// Act & Assert - 只验证字段是否存在和可访问
	assert.Equal(suite.T(), "testuser", req.Username)
	assert.Equal(suite.T(), "test@example.com", req.Email)
	assert.Equal(suite.T(), "password123", req.Password)
	assert.Equal(suite.T(), "Test User", req.FullName)
}

func (suite *ModelsTestSuite) TestLoginRequest_Fields() {
	// Arrange
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act & Assert - 只验证字段是否存在和可访问
	assert.Equal(suite.T(), "test@example.com", req.Email)
	assert.Equal(suite.T(), "password123", req.Password)
}

func (suite *ModelsTestSuite) TestCreateProjectRequest_Fields() {
	// Arrange
	req := CreateProjectRequest{
		ProjectName: "测试项目",
		Description: "这是一个测试项目",
		ProjectType: "web",
	}

	// Act & Assert - 只验证字段是否存在和可访问
	assert.Equal(suite.T(), "测试项目", req.ProjectName)
	assert.Equal(suite.T(), "这是一个测试项目", req.Description)
	assert.Equal(suite.T(), "web", req.ProjectType)
}

func (suite *ModelsTestSuite) TestUserStatus_Constants() {
	// Assert
	assert.Equal(suite.T(), "active", UserStatusActive)
	assert.Equal(suite.T(), "inactive", UserStatusInactive)
	assert.Equal(suite.T(), "suspended", UserStatusSuspended)
}

func (suite *ModelsTestSuite) TestProjectStatus_Constants() {
	// Assert - 使用实际存在的常量
	assert.Equal(suite.T(), "draft", ProjectStatusDraft)
	assert.Equal(suite.T(), "analyzing", ProjectStatusAnalyzing)
	assert.Equal(suite.T(), "completed", ProjectStatusCompleted)
	assert.Equal(suite.T(), "archived", ProjectStatusArchived)
}

func (suite *ModelsTestSuite) TestDiagramType_Constants() {
	// Assert
	assert.Equal(suite.T(), "business_flow", DiagramTypeBusinessFlow)
	assert.Equal(suite.T(), "architecture", DiagramTypeArchitecture)
	assert.Equal(suite.T(), "data_model", DiagramTypeDataModel)
	assert.Equal(suite.T(), "sequence", DiagramTypeSequence)
}

func (suite *ModelsTestSuite) TestSessionType_Constants() {
	// Assert
	assert.Equal(suite.T(), "requirement_analysis", SessionTypeRequirementAnalysis)
	assert.Equal(suite.T(), "document_review", SessionTypeDocumentReview)
	assert.Equal(suite.T(), "puml_editing", SessionTypePUMLEditing)
}

func (suite *ModelsTestSuite) TestSenderType_Constants() {
	// Assert
	assert.Equal(suite.T(), "user", SenderTypeUser)
	assert.Equal(suite.T(), "ai", SenderTypeAI)
	assert.Equal(suite.T(), "system", SenderTypeSystem)
}

func (suite *ModelsTestSuite) TestMessageType_Constants() {
	// Assert
	assert.Equal(suite.T(), "text", MessageTypeText)
	assert.Equal(suite.T(), "question", MessageTypeQuestion)
	assert.Equal(suite.T(), "answer", MessageTypeAnswer)
	assert.Equal(suite.T(), "command", MessageTypeCommand)
}

func (suite *ModelsTestSuite) TestRequirement_JSONSerialization() {
	// Arrange
	requirement := &Requirement{
		RequirementID:         uuid.New(),
		ProjectID:            uuid.New(),
		RawRequirement:       "用户需要能够注册和登录",
		StructuredRequirement: `{"functions": ["register", "login"]}`,
		CompletenessScore:    0.8,
		MissingInfoTypes:     `["password policy"]`,
		AnalysisStatus:       "completed",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Act - 序列化
	jsonData, err := json.Marshal(requirement)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedRequirement Requirement
	err = json.Unmarshal(jsonData, &deserializedRequirement)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), requirement.RequirementID, deserializedRequirement.RequirementID)
	assert.Equal(suite.T(), requirement.ProjectID, deserializedRequirement.ProjectID)
	assert.Equal(suite.T(), requirement.RawRequirement, deserializedRequirement.RawRequirement)
	assert.Equal(suite.T(), requirement.StructuredRequirement, deserializedRequirement.StructuredRequirement)
	assert.Equal(suite.T(), requirement.CompletenessScore, deserializedRequirement.CompletenessScore)
	assert.Equal(suite.T(), requirement.MissingInfoTypes, deserializedRequirement.MissingInfoTypes)
	assert.Equal(suite.T(), requirement.AnalysisStatus, deserializedRequirement.AnalysisStatus)
}

func (suite *ModelsTestSuite) TestAsyncTask_JSONSerialization() {
	// Arrange
	task := &AsyncTask{
		TaskID:      uuid.New(),
		UserID:      uuid.New(),
		ProjectID:   uuid.New(),
		TaskType:    "generate_document",
		TaskName:    "生成需求文档",
		Status:      "running",
		Progress:    50,
		ResultData:  `{"document_id": "123"}`,
		ErrorMessage: "",
		CreatedAt:   time.Now(),
		StartedAt:   nil,
		CompletedAt: nil,
		Metadata:    `{"stage": 1}`,
	}

	// Act - 序列化
	jsonData, err := json.Marshal(task)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedTask AsyncTask
	err = json.Unmarshal(jsonData, &deserializedTask)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), task.TaskID, deserializedTask.TaskID)
	assert.Equal(suite.T(), task.UserID, deserializedTask.UserID)
	assert.Equal(suite.T(), task.ProjectID, deserializedTask.ProjectID)
	assert.Equal(suite.T(), task.TaskType, deserializedTask.TaskType)
	assert.Equal(suite.T(), task.TaskName, deserializedTask.TaskName)
	assert.Equal(suite.T(), task.Status, deserializedTask.Status)
	assert.Equal(suite.T(), task.Progress, deserializedTask.Progress)
	assert.Equal(suite.T(), task.ResultData, deserializedTask.ResultData)
	assert.Equal(suite.T(), task.ErrorMessage, deserializedTask.ErrorMessage)
	assert.Equal(suite.T(), task.Metadata, deserializedTask.Metadata)
	assert.Nil(suite.T(), deserializedTask.StartedAt)
	assert.Nil(suite.T(), deserializedTask.CompletedAt)
}

func (suite *ModelsTestSuite) TestPUMLDiagram_JSONSerialization() {
	// Arrange
	diagram := &PUMLDiagram{
		DiagramID:          uuid.New(),
		ProjectID:          uuid.New(),
		DiagramType:        DiagramTypeBusinessFlow,
		DiagramName:        "业务流程图",
		PUMLContent:        "@startuml\n...\n@enduml",
		RenderedURL:        "http://example.com/diagram.png",
		Version:            1,
		Stage:              1,
		TaskID:             nil,
		IsValidated:        true,
		ValidationFeedback: "图表生成成功",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Act - 序列化
	jsonData, err := json.Marshal(diagram)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedDiagram PUMLDiagram
	err = json.Unmarshal(jsonData, &deserializedDiagram)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), diagram.DiagramID, deserializedDiagram.DiagramID)
	assert.Equal(suite.T(), diagram.ProjectID, deserializedDiagram.ProjectID)
	assert.Equal(suite.T(), diagram.DiagramType, deserializedDiagram.DiagramType)
	assert.Equal(suite.T(), diagram.DiagramName, deserializedDiagram.DiagramName)
	assert.Equal(suite.T(), diagram.PUMLContent, deserializedDiagram.PUMLContent)
	assert.Equal(suite.T(), diagram.RenderedURL, deserializedDiagram.RenderedURL)
	assert.Equal(suite.T(), diagram.Version, deserializedDiagram.Version)
	assert.Equal(suite.T(), diagram.Stage, deserializedDiagram.Stage)
	assert.Equal(suite.T(), diagram.IsValidated, deserializedDiagram.IsValidated)
	assert.Equal(suite.T(), diagram.ValidationFeedback, deserializedDiagram.ValidationFeedback)
	assert.Nil(suite.T(), deserializedDiagram.TaskID)
}

func (suite *ModelsTestSuite) TestDocument_JSONSerialization() {
	// Arrange
	document := &Document{
		DocumentID:   uuid.New(),
		ProjectID:    uuid.New(),
		DocumentType: "requirement_analysis",
		DocumentName: "需求分析文档",
		Content:      "这是需求分析文档的内容",
		Format:       "markdown",
		FilePath:     "/documents/requirement.md",
		Version:      1,
		Stage:        1,
		TaskID:       nil,
		GeneratedAt:  time.Now(),
		IsFinal:      true,
	}

	// Act - 序列化
	jsonData, err := json.Marshal(document)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedDocument Document
	err = json.Unmarshal(jsonData, &deserializedDocument)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), document.DocumentID, deserializedDocument.DocumentID)
	assert.Equal(suite.T(), document.ProjectID, deserializedDocument.ProjectID)
	assert.Equal(suite.T(), document.DocumentType, deserializedDocument.DocumentType)
	assert.Equal(suite.T(), document.DocumentName, deserializedDocument.DocumentName)
	assert.Equal(suite.T(), document.Content, deserializedDocument.Content)
	assert.Equal(suite.T(), document.Format, deserializedDocument.Format)
	assert.Equal(suite.T(), document.FilePath, deserializedDocument.FilePath)
	assert.Equal(suite.T(), document.Version, deserializedDocument.Version)
	assert.Equal(suite.T(), document.Stage, deserializedDocument.Stage)
	assert.Equal(suite.T(), document.IsFinal, deserializedDocument.IsFinal)
	assert.Nil(suite.T(), deserializedDocument.TaskID)
}

func (suite *ModelsTestSuite) TestChatSession_JSONSerialization() {
	// Arrange
	session := &ChatSession{
		SessionID:   uuid.New(),
		ProjectID:   uuid.New(),
		UserID:      uuid.New(),
		SessionType: SessionTypeRequirementAnalysis,
		StartedAt:   time.Now(),
		EndedAt:     nil,
		Status:      "active",
		Context:     `{"project_stage": 1}`,
	}

	// Act - 序列化
	jsonData, err := json.Marshal(session)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedSession ChatSession
	err = json.Unmarshal(jsonData, &deserializedSession)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), session.SessionID, deserializedSession.SessionID)
	assert.Equal(suite.T(), session.ProjectID, deserializedSession.ProjectID)
	assert.Equal(suite.T(), session.UserID, deserializedSession.UserID)
	assert.Equal(suite.T(), session.SessionType, deserializedSession.SessionType)
	assert.Equal(suite.T(), session.Status, deserializedSession.Status)
	assert.Equal(suite.T(), session.Context, deserializedSession.Context)
	assert.Nil(suite.T(), deserializedSession.EndedAt)
}

func (suite *ModelsTestSuite) TestChatMessage_JSONSerialization() {
	// Arrange
	message := &ChatMessage{
		MessageID:      uuid.New(),
		SessionID:      uuid.New(),
		SenderType:     SenderTypeUser,
		MessageContent: "这是一条测试消息",
		MessageType:    MessageTypeText,
		Metadata:       `{"source": "web"}`,
		Timestamp:      time.Now(),
		Processed:      true,
	}

	// Act - 序列化
	jsonData, err := json.Marshal(message)
	assert.NoError(suite.T(), err)

	// Act - 反序列化
	var deserializedMessage ChatMessage
	err = json.Unmarshal(jsonData, &deserializedMessage)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), message.MessageID, deserializedMessage.MessageID)
	assert.Equal(suite.T(), message.SessionID, deserializedMessage.SessionID)
	assert.Equal(suite.T(), message.SenderType, deserializedMessage.SenderType)
	assert.Equal(suite.T(), message.MessageContent, deserializedMessage.MessageContent)
	assert.Equal(suite.T(), message.MessageType, deserializedMessage.MessageType)
	assert.Equal(suite.T(), message.Metadata, deserializedMessage.Metadata)
	assert.Equal(suite.T(), message.Processed, deserializedMessage.Processed)
}

func TestModelsTestSuite(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
} 