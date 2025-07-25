package model

import (
	"github.com/google/uuid"
	"time"
)

// ChatSession 对话会话模型
type ChatSession struct {
	SessionID   uuid.UUID  `json:"session_id" gorm:"type:char(36);primaryKey;column:session_id" db:"session_id"`
	ProjectID   uuid.UUID  `json:"project_id" gorm:"type:char(36);not null;index;column:project_id" db:"project_id"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:char(36);not null;index;column:user_id" db:"user_id"`
	SessionType string     `json:"session_type" gorm:"type:varchar(50);not null;column:session_type" db:"session_type"`
	StartedAt   time.Time  `json:"started_at" gorm:"autoCreateTime;column:started_at" db:"started_at"`
	EndedAt     *time.Time `json:"ended_at" gorm:"column:ended_at" db:"ended_at"`
	Status      string     `json:"status" gorm:"type:varchar(20);default:'active';column:status" db:"status"`
	Context     string     `json:"context" gorm:"type:json;column:context" db:"context"` // JSON字符串
}

// TableName 指定表名
func (ChatSession) TableName() string {
	return "chat_sessions"
}
