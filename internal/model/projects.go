package model

import (
	"github.com/google/uuid"
	"time"
)

// Project 项目模型
type Project struct {
	ProjectID            uuid.UUID `json:"project_id" gorm:"type:char(36);primaryKey;column:project_id" db:"project_id"`
	UserID               uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index;column:user_id" db:"user_id"`
	ProjectName          string    `json:"project_name" gorm:"type:varchar(100);not null;column:project_name" db:"project_name"`
	Description          string    `json:"description" gorm:"type:text;column:description" db:"description"`
	ProjectType          string    `json:"project_type" gorm:"type:varchar(20);default:'web';column:project_type" db:"project_type"`
	Status               string    `json:"status" gorm:"type:varchar(20);default:'planning';column:status" db:"status"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at" db:"updated_at"`
	CompletionPercentage int       `json:"completion_percentage" gorm:"default:0;column:completion_percentage" db:"completion_percentage"`
	Settings             string    `json:"settings" gorm:"type:json;column:settings" db:"settings"` // JSON字符串

	// GORM 关联 - User 是被引用的，不应该有外键
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}
