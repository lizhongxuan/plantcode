package model

import (
	"github.com/google/uuid"
	"time"
)

// Document 生成文档模型
type Document struct {
	DocumentID   uuid.UUID  `json:"document_id" gorm:"type:char(36);primaryKey;column:document_id" db:"document_id"`
	ProjectID    uuid.UUID  `json:"project_id" gorm:"type:char(36);not null;index;column:project_id" db:"project_id"`
	DocumentType string     `json:"document_type" gorm:"type:varchar(50);not null;column:document_type" db:"document_type"`
	DocumentName string     `json:"document_name" gorm:"type:varchar(200);not null;column:document_name" db:"document_name"`
	Content      string     `json:"content" gorm:"type:text;not null;column:content" db:"content"`
	Format       string     `json:"format" gorm:"type:varchar(50);default:'markdown';column:format" db:"format"`
	FilePath     string     `json:"file_path" gorm:"type:varchar(255);column:file_path" db:"file_path"`
	Version      int        `json:"version" gorm:"default:1;column:version" db:"version"`
	Stage        int        `json:"stage" gorm:"default:1;column:stage" db:"stage"`                     // 新增：所属阶段 1,2,3
	TaskID       *uuid.UUID `json:"task_id,omitempty" gorm:"type:char(36);column:task_id" db:"task_id"` // 新增：关联的异步任务ID
	GeneratedAt  time.Time  `json:"generated_at" gorm:"autoCreateTime;column:generated_at" db:"generated_at"`
	IsFinal      bool       `json:"is_final" gorm:"default:false;column:is_final" db:"is_final"`
}

// TableName 指定表名
func (Document) TableName() string {
	return "generated_documents"
}
