package repository

import (
	"database/sql"
	"testing"
	"time"

	"ai-dev-platform/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// convertNullTime 将 sql.NullTime 转换为 *time.Time
func convertNullTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// convertTimePtr 将 *time.Time 转换为 sql.NullTime
func convertTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

type UserRepositoryTestSuite struct {
	suite.Suite
	db         *sql.DB
	mock       sqlmock.Sqlmock
	repository Repository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	// 这个测试需要重构，使用GORM和真实的测试数据库
	// 暂时注释掉，创建一个SQL DB到GORM的桥接
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: suite.db,
	}), &gorm.Config{})
	assert.NoError(suite.T(), err)
	
	database := &Database{GORM: gormDB}
	suite.repository = NewMySQLRepository(database)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *UserRepositoryTestSuite) TestCreateUser_Success() {
	// Arrange
	user := &model.User{
		UserID:       uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		Status:       model.UserStatusActive,
		Preferences:  `{"theme": "dark"}`,
	}

	suite.mock.ExpectExec("INSERT INTO users").
		WithArgs(
			user.UserID,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.FullName,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			user.Status,
			user.Preferences,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := suite.repository.CreateUser(user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestCreateUser_DatabaseError() {
	// Arrange
	user := &model.User{
		UserID:       uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		Status:       model.UserStatusActive,
		Preferences:  `{"theme": "dark"}`,
	}

	suite.mock.ExpectExec("INSERT INTO users").
		WithArgs(
			user.UserID,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.FullName,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			user.Status,
			user.Preferences,
		).
		WillReturnError(sql.ErrConnDone)

	// Act
	err := suite.repository.CreateUser(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "创建用户失败")
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_Success() {
	// Arrange
	userID := uuid.New()
	email := "test@example.com"
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"user_id", "username", "email", "password_hash", "full_name",
		"created_at", "updated_at", "last_login", "status", "preferences",
	}).AddRow(
		userID, "testuser", email, "hashedpassword", "Test User",
		now, now, now, model.UserStatusActive, `{"theme": "dark"}`,
	)

	suite.mock.ExpectQuery("SELECT .+ FROM users WHERE email = \\? AND status != 'deleted'").
		WithArgs(email).
		WillReturnRows(rows)

	// Act
	user, err := suite.repository.GetUserByEmail(email)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.UserID)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.Equal(suite.T(), email, user.Email)
	assert.Equal(suite.T(), "hashedpassword", user.PasswordHash)
	assert.Equal(suite.T(), "Test User", user.FullName)
	assert.Equal(suite.T(), model.UserStatusActive, user.Status)
	assert.Equal(suite.T(), `{"theme": "dark"}`, user.Preferences)
	assert.NotNil(suite.T(), user.LastLogin)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_NotFound() {
	// Arrange
	email := "nonexistent@example.com"

	suite.mock.ExpectQuery("SELECT .+ FROM users WHERE email = \\? AND status != 'deleted'").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	// Act
	user, err := suite.repository.GetUserByEmail(email)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "用户不存在")
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_WithNullValues() {
	// Arrange
	userID := uuid.New()
	email := "test@example.com"
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"user_id", "username", "email", "password_hash", "full_name",
		"created_at", "updated_at", "last_login", "status", "preferences",
	}).AddRow(
		userID, "testuser", email, "hashedpassword", "Test User",
		now, now, nil, model.UserStatusActive, nil, // null last_login and preferences
	)

	suite.mock.ExpectQuery("SELECT .+ FROM users WHERE email = \\? AND status != 'deleted'").
		WithArgs(email).
		WillReturnRows(rows)

	// Act
	user, err := suite.repository.GetUserByEmail(email)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.UserID)
	assert.Nil(suite.T(), user.LastLogin)
	assert.Equal(suite.T(), "", user.Preferences)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByID_Success() {
	// Arrange
	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"user_id", "username", "email", "password_hash", "full_name",
		"created_at", "updated_at", "last_login", "status", "preferences",
	}).AddRow(
		userID, "testuser", "test@example.com", "hashedpassword", "Test User",
		now, now, now, model.UserStatusActive, `{"theme": "dark"}`,
	)

	suite.mock.ExpectQuery("SELECT .+ FROM users WHERE user_id = \\? AND status != 'deleted'").
		WithArgs(userID).
		WillReturnRows(rows)

	// Act
	user, err := suite.repository.GetUserByID(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.UserID)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByID_NotFound() {
	// Arrange
	userID := uuid.New()

	suite.mock.ExpectQuery("SELECT .+ FROM users WHERE user_id = \\? AND status != 'deleted'").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	// Act
	user, err := suite.repository.GetUserByID(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "用户不存在")
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_Success() {
	// Arrange
	user := &model.User{
		UserID:      uuid.New(),
		Username:    "updateduser",
		Email:       "updated@example.com",
		FullName:    "Updated User",
		Preferences: `{"theme": "light"}`,
	}

	suite.mock.ExpectExec("UPDATE users SET username = \\?, email = \\?, full_name = \\?, updated_at = \\?, preferences = \\? WHERE user_id = \\?").
		WithArgs(
			user.Username,
			user.Email,
			user.FullName,
			sqlmock.AnyArg(), // updated_at
			user.Preferences,
			user.UserID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act
	err := suite.repository.UpdateUser(user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_NotFound() {
	// Arrange
	user := &model.User{
		UserID:      uuid.New(),
		Username:    "updateduser",
		Email:       "updated@example.com",
		FullName:    "Updated User",
		Preferences: `{"theme": "light"}`,
	}

	suite.mock.ExpectExec("UPDATE users SET username = \\?, email = \\?, full_name = \\?, updated_at = \\?, preferences = \\? WHERE user_id = \\?").
		WithArgs(
			user.Username,
			user.Email,
			user.FullName,
			sqlmock.AnyArg(),
			user.Preferences,
			user.UserID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected

	// Act
	err := suite.repository.UpdateUser(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "用户不存在或未更新")
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_DatabaseError() {
	// Arrange
	user := &model.User{
		UserID:      uuid.New(),
		Username:    "updateduser",
		Email:       "updated@example.com",
		FullName:    "Updated User",
		Preferences: `{"theme": "light"}`,
	}

	suite.mock.ExpectExec("UPDATE users SET username = \\?, email = \\?, full_name = \\?, updated_at = \\?, preferences = \\? WHERE user_id = \\?").
		WithArgs(
			user.Username,
			user.Email,
			user.FullName,
			sqlmock.AnyArg(),
			user.Preferences,
			user.UserID,
		).
		WillReturnError(sql.ErrConnDone)

	// Act
	err := suite.repository.UpdateUser(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "更新用户失败")
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateUserLastLogin_Success() {
	// Arrange
	userID := uuid.New()

	suite.mock.ExpectExec("UPDATE users SET last_login = \\?, updated_at = \\? WHERE user_id = \\?").
		WithArgs(
			sqlmock.AnyArg(), // last_login
			sqlmock.AnyArg(), // updated_at
			userID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act
	err := suite.repository.UpdateUserLastLogin(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateUserLastLogin_NotFound() {
	// Arrange
	userID := uuid.New()

	suite.mock.ExpectExec("UPDATE users SET last_login = \\?, updated_at = \\? WHERE user_id = \\?").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			userID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected

	// Act
	err := suite.repository.UpdateUserLastLogin(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "用户不存在")
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateUserLastLogin_DatabaseError() {
	// Arrange
	userID := uuid.New()

	suite.mock.ExpectExec("UPDATE users SET last_login = \\?, updated_at = \\? WHERE user_id = \\?").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			userID,
		).
		WillReturnError(sql.ErrConnDone)

	// Act
	err := suite.repository.UpdateUserLastLogin(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "更新用户登录时间失败")
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// 测试数据库为空的情况
func (suite *UserRepositoryTestSuite) TestCreateUser_NilDatabase() {
	// Arrange
	repository := NewMySQLRepository(&Database{GORM: nil})
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Act
	err := repository.CreateUser(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "数据库连接不可用")
}

// 辅助函数测试
func (suite *UserRepositoryTestSuite) TestConvertNullTime() {
	// Null time
	nullTime := sql.NullTime{Valid: false}
	result := convertNullTime(nullTime)
	assert.Nil(suite.T(), result)

	// Valid time
	now := time.Now()
	validTime := sql.NullTime{Time: now, Valid: true}
	result = convertNullTime(validTime)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), now, *result)
}

func (suite *UserRepositoryTestSuite) TestConvertTimePtr() {
	// Nil pointer
	result := convertTimePtr(nil)
	assert.False(suite.T(), result.Valid)

	// Valid pointer
	now := time.Now()
	result = convertTimePtr(&now)
	assert.True(suite.T(), result.Valid)
	assert.Equal(suite.T(), now, result.Time)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
} 