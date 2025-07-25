package model

import (
	"github.com/google/uuid"
	"time"
)

// StageProgress 阶段进度跟踪
type StageProgress struct {
	ProgressID     uuid.UUID  `json:"progress_id" gorm:"type:char(36);primaryKey;column:progress_id" db:"progress_id"`
	ProjectID      uuid.UUID  `json:"project_id" gorm:"type:char(36);not null;index;column:project_id" db:"project_id"`
	Stage          int        `json:"stage" gorm:"not null;column:stage" db:"stage"`                                  // 1, 2, 3
	Status         string     `json:"status" gorm:"type:varchar(20);default:'not_started';column:status" db:"status"` // not_started, in_progress, completed, failed
	CompletionRate int        `json:"completion_rate" gorm:"default:0;column:completion_rate" db:"completion_rate"`   // 0-100
	StartedAt      *time.Time `json:"started_at,omitempty" gorm:"column:started_at" db:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty" gorm:"column:completed_at" db:"completed_at"`
	DocumentCount  int        `json:"document_count" gorm:"default:0;column:document_count" db:"document_count"`
	PUMLCount      int        `json:"puml_count" gorm:"default:0;column:puml_count" db:"puml_count"`
	LastTaskID     *uuid.UUID `json:"last_task_id,omitempty" gorm:"type:char(36);column:last_task_id" db:"last_task_id"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime;column:updated_at" db:"updated_at"`
}

// TableName 指定表名
func (StageProgress) TableName() string {
	return "stage_progress"
}
