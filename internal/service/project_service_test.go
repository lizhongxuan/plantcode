package service

import (
	"fmt"
	"testing"

	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProjectServiceTestSuite struct {
	suite.Suite
	mockRepo       *MockRepository
	projectService ProjectService
}

func (suite *ProjectServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepository)
	suite.projectService = NewProjectService(suite.mockRepo)
}

func (suite *ProjectServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestCreateProject_Success() {
	// Arrange
	userID := uuid.New()
	req := &model.CreateProjectRequest{
		ProjectName: "测试项目",
		Description: "这是一个测试项目",
		ProjectType: "web_application",
	}

	suite.mockRepo.On("CreateProject", mock.MatchedBy(func(project *model.Project) bool {
		return project.UserID == userID &&
			project.ProjectName == req.ProjectName &&
			project.Description == req.Description &&
			project.ProjectType == req.ProjectType &&
			project.Status == "draft" &&
			project.CompletionPercentage == 0 &&
			project.Settings == `{}` &&
			project.ProjectID != uuid.Nil
	})).Return(nil)

	// Act
	project, err := suite.projectService.CreateProject(userID, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), project)
	assert.Equal(suite.T(), userID, project.UserID)
	assert.Equal(suite.T(), req.ProjectName, project.ProjectName)
	assert.Equal(suite.T(), req.Description, project.Description)
	assert.Equal(suite.T(), req.ProjectType, project.ProjectType)
	assert.Equal(suite.T(), "draft", project.Status)
	assert.Equal(suite.T(), 0, project.CompletionPercentage)
	assert.Equal(suite.T(), `{}`, project.Settings)
	assert.NotEqual(suite.T(), uuid.Nil, project.ProjectID)
}

func (suite *ProjectServiceTestSuite) TestCreateProject_InvalidInput() {
	testCases := []struct {
		name string
		req  *model.CreateProjectRequest
	}{
		{
			name: "empty project name",
			req: &model.CreateProjectRequest{
				ProjectName: "",
				Description: "描述",
				ProjectType: "web_application",
			},
		},
		{
			name: "empty description",
			req: &model.CreateProjectRequest{
				ProjectName: "项目名",
				Description: "",
				ProjectType: "web_application",
			},
		},
		{
			name: "invalid project type",
			req: &model.CreateProjectRequest{
				ProjectName: "项目名",
				Description: "描述",
				ProjectType: "invalid",
			},
		},
	}

	userID := uuid.New()

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Act
			project, err := suite.projectService.CreateProject(userID, tc.req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, project)
			assert.Contains(t, err.Error(), "输入验证失败")
		})
	}
}

func (suite *ProjectServiceTestSuite) TestCreateProject_RepositoryError() {
	// Arrange
	userID := uuid.New()
	req := &model.CreateProjectRequest{
		ProjectName: "测试项目",
		Description: "这是一个测试项目",
		ProjectType: "web_application",
	}

	suite.mockRepo.On("CreateProject", mock.Anything).Return(fmt.Errorf("数据库错误"))

	// Act
	project, err := suite.projectService.CreateProject(userID, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), project)
	assert.Contains(suite.T(), err.Error(), "创建项目失败")
}

func (suite *ProjectServiceTestSuite) TestGetUserProjects_Success() {
	// Arrange
	userID := uuid.New()
	page := 1
	pageSize := 10

	projects := []*model.Project{
		{
			ProjectID:   uuid.New(),
			UserID:      userID,
			ProjectName: "项目1",
			Description: "项目1描述",
			ProjectType: "web_application",
			Status:      "draft",
		},
		{
			ProjectID:   uuid.New(),
			UserID:      userID,
			ProjectName: "项目2",
			Description: "项目2描述",
			ProjectType: "mobile_app",
			Status:      "in_progress",
		},
	}
	totalCount := int64(2)

	suite.mockRepo.On("GetProjectsByUserID", userID, page, pageSize).Return(projects, totalCount, nil)

	// Act
	result, pagination, err := suite.projectService.GetUserProjects(userID, page, pageSize)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), projects, result)
	assert.Equal(suite.T(), utils.PaginationInfo{
		Page:      page,
		PageSize:  pageSize,
		Total:     totalCount,
		TotalPage: 1,
	}, pagination)
}

func (suite *ProjectServiceTestSuite) TestGetUserProjects_DefaultPagination() {
	// Arrange
	userID := uuid.New()
	page := 0 // 无效页码，应该被修正为1
	pageSize := 0 // 无效页大小，应该被修正为10

	projects := make([]*model.Project, 0) // 使用空slice而不是nil
	totalCount := int64(0)

	suite.mockRepo.On("GetProjectsByUserID", userID, 1, 10).Return(projects, totalCount, nil)

	// Act
	result, pagination, err := suite.projectService.GetUserProjects(userID, page, pageSize)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), projects, result)
	assert.Equal(suite.T(), utils.PaginationInfo{
		Page:      1,
		PageSize:  10,
		Total:     totalCount,
		TotalPage: 0,
	}, pagination)
}

func (suite *ProjectServiceTestSuite) TestGetUserProjects_RepositoryError() {
	// Arrange
	userID := uuid.New()
	page := 1
	pageSize := 10

	suite.mockRepo.On("GetProjectsByUserID", userID, page, pageSize).Return(nil, int64(0), fmt.Errorf("数据库错误"))

	// Act
	result, _, err := suite.projectService.GetUserProjects(userID, page, pageSize)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "获取项目列表失败")
	// 注意：错误时返回的pagination不一定是空的，因为参数已经被修正
}

func (suite *ProjectServiceTestSuite) TestGetProject_Success() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()

	project := &model.Project{
		ProjectID:   projectID,
		UserID:      userID,
		ProjectName: "测试项目",
		Description: "这是一个测试项目",
		ProjectType: "web_application",
		Status:      "draft",
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(project, nil)

	// Act
	result, err := suite.projectService.GetProject(projectID, userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), project, result)
}

func (suite *ProjectServiceTestSuite) TestGetProject_NotFound() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()

	suite.mockRepo.On("GetProjectByID", projectID).Return(nil, fmt.Errorf("项目不存在"))

	// Act
	result, err := suite.projectService.GetProject(projectID, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "项目不存在")
}

func (suite *ProjectServiceTestSuite) TestGetProject_AccessDenied() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New() // 不同的用户ID

	project := &model.Project{
		ProjectID:   projectID,
		UserID:      ownerID, // 项目属于其他用户
		ProjectName: "测试项目",
		Description: "这是一个测试项目",
		ProjectType: "web_application",
		Status:      "draft",
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(project, nil)

	// Act
	result, err := suite.projectService.GetProject(projectID, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "无权访问此项目")
}

func (suite *ProjectServiceTestSuite) TestUpdateProject_Success() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()
	newName := "更新后的项目名"
	newDescription := "更新后的描述"
	newStatus := "in_progress"
	newPercentage := 50

	updates := &ProjectUpdateRequest{
		ProjectName:          &newName,
		Description:          &newDescription,
		Status:               &newStatus,
		CompletionPercentage: &newPercentage,
	}

	existingProject := &model.Project{
		ProjectID:            projectID,
		UserID:               userID,
		ProjectName:          "原项目名",
		Description:          "原描述",
		ProjectType:          "web",
		Status:               "draft",
		CompletionPercentage: 0,
		Settings:             `{}`,
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(existingProject, nil)
	suite.mockRepo.On("UpdateProject", mock.MatchedBy(func(project *model.Project) bool {
		return project.ProjectID == projectID &&
			project.UserID == userID &&
			project.ProjectName == newName &&
			project.Description == newDescription &&
			project.Status == newStatus &&
			project.CompletionPercentage == newPercentage
	})).Return(nil)

	// Act
	result, err := suite.projectService.UpdateProject(projectID, userID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), newName, result.ProjectName)
	assert.Equal(suite.T(), newDescription, result.Description)
	assert.Equal(suite.T(), newStatus, result.Status)
	assert.Equal(suite.T(), newPercentage, result.CompletionPercentage)
}

func (suite *ProjectServiceTestSuite) TestUpdateProject_PartialUpdate() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()
	newName := "更新后的项目名"

	updates := &ProjectUpdateRequest{
		ProjectName: &newName,
		// 其他字段为nil，应该保持原值
	}

	existingProject := &model.Project{
		ProjectID:            projectID,
		UserID:               userID,
		ProjectName:          "原项目名",
		Description:          "原描述",
		ProjectType:          "web",
		Status:               "draft",
		CompletionPercentage: 25,
		Settings:             `{"theme": "dark"}`,
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(existingProject, nil)
	suite.mockRepo.On("UpdateProject", mock.MatchedBy(func(project *model.Project) bool {
		return project.ProjectID == projectID &&
			project.UserID == userID &&
			project.ProjectName == newName &&
			project.Description == existingProject.Description && // 应该保持原值
			project.Status == existingProject.Status && // 应该保持原值
			project.CompletionPercentage == existingProject.CompletionPercentage // 应该保持原值
	})).Return(nil)

	// Act
	result, err := suite.projectService.UpdateProject(projectID, userID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), newName, result.ProjectName)
	assert.Equal(suite.T(), existingProject.Description, result.Description)
	assert.Equal(suite.T(), existingProject.Status, result.Status)
	assert.Equal(suite.T(), existingProject.CompletionPercentage, result.CompletionPercentage)
}

func (suite *ProjectServiceTestSuite) TestUpdateProject_ProjectNotFound() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()
	updates := &ProjectUpdateRequest{
		ProjectName: utils.StringPtr("新名称"),
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(nil, fmt.Errorf("项目不存在"))

	// Act
	result, err := suite.projectService.UpdateProject(projectID, userID, updates)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "项目不存在")
}

func (suite *ProjectServiceTestSuite) TestUpdateProject_AccessDenied() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New() // 不同的用户ID
	updates := &ProjectUpdateRequest{
		ProjectName: utils.StringPtr("新名称"),
	}

	existingProject := &model.Project{
		ProjectID: projectID,
		UserID:    ownerID, // 项目属于其他用户
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(existingProject, nil)

	// Act
	result, err := suite.projectService.UpdateProject(projectID, userID, updates)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "无权修改此项目")
}

func (suite *ProjectServiceTestSuite) TestUpdateProject_RepositoryError() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()
	updates := &ProjectUpdateRequest{
		ProjectName: utils.StringPtr("新名称"),
	}

	existingProject := &model.Project{
		ProjectID: projectID,
		UserID:    userID,
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(existingProject, nil)
	suite.mockRepo.On("UpdateProject", mock.Anything).Return(fmt.Errorf("数据库错误"))

	// Act
	result, err := suite.projectService.UpdateProject(projectID, userID, updates)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "更新项目失败")
}

func (suite *ProjectServiceTestSuite) TestDeleteProject_Success() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()

	existingProject := &model.Project{
		ProjectID: projectID,
		UserID:    userID,
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(existingProject, nil)
	suite.mockRepo.On("DeleteProject", projectID).Return(nil)

	// Act
	err := suite.projectService.DeleteProject(projectID, userID)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *ProjectServiceTestSuite) TestDeleteProject_ProjectNotFound() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()

	suite.mockRepo.On("GetProjectByID", projectID).Return(nil, fmt.Errorf("项目不存在"))

	// Act
	err := suite.projectService.DeleteProject(projectID, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "项目不存在")
}

func (suite *ProjectServiceTestSuite) TestDeleteProject_AccessDenied() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New() // 不同的用户ID

	existingProject := &model.Project{
		ProjectID: projectID,
		UserID:    ownerID, // 项目属于其他用户
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(existingProject, nil)

	// Act
	err := suite.projectService.DeleteProject(projectID, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "无权删除此项目")
}

func (suite *ProjectServiceTestSuite) TestDeleteProject_RepositoryError() {
	// Arrange
	projectID := uuid.New()
	userID := uuid.New()

	existingProject := &model.Project{
		ProjectID: projectID,
		UserID:    userID,
	}

	suite.mockRepo.On("GetProjectByID", projectID).Return(existingProject, nil)
	suite.mockRepo.On("DeleteProject", projectID).Return(fmt.Errorf("数据库错误"))

	// Act
	err := suite.projectService.DeleteProject(projectID, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "删除项目失败")
}

func TestProjectServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectServiceTestSuite))
} 