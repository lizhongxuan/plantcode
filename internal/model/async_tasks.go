package model

import (
	"github.com/google/uuid"
	"time"
)

// AsyncTask 异步任务表
type AsyncTask struct {
	TaskID       uuid.UUID  `json:"task_id" gorm:"type:char(36);primaryKey;column:task_id" db:"task_id"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:char(36);not null;index;column:user_id" db:"user_id"`
	ProjectID    uuid.UUID  `json:"project_id" gorm:"type:char(36);not null;index;column:project_id" db:"project_id"`
	TaskType     string     `json:"task_type" gorm:"type:varchar(50);not null;column:task_type" db:"task_type"` // stage_document_generation, puml_generation, document_generation
	TaskName     string     `json:"task_name" gorm:"type:varchar(200);not null;column:task_name" db:"task_name"`
	Status       string     `json:"status" gorm:"type:varchar(20);default:'pending';column:status" db:"status"` // pending, running, completed, failed
	Progress     int        `json:"progress" gorm:"default:0;column:progress" db:"progress"`                    // 0-100
	ResultData   string     `json:"result_data,omitempty" gorm:"type:json;column:result_data" db:"result_data"` // JSON格式的结果数据
	ErrorMessage string     `json:"error_message,omitempty" gorm:"type:text;column:error_message" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	StartedAt    *time.Time `json:"started_at,omitempty" gorm:"column:started_at" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" gorm:"column:completed_at" db:"completed_at"`
	Metadata     string     `json:"metadata,omitempty" gorm:"type:json;column:metadata" db:"metadata"` // JSON格式的任务元数据
}

// TableName 指定表名
func (AsyncTask) TableName() string {
	return "async_tasks"
}
