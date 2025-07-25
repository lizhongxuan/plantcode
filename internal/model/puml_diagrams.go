package model

import (
	"github.com/google/uuid"
	"time"
)

// PUMLDiagram PUML图表模型
type PUMLDiagram struct {
	DiagramID          uuid.UUID  `json:"diagram_id" gorm:"type:char(36);primaryKey;column:diagram_id" db:"diagram_id"`
	ProjectID          uuid.UUID  `json:"project_id" gorm:"type:char(36);not null;index;column:project_id" db:"project_id"`
	DiagramType        string     `json:"diagram_type" gorm:"type:varchar(50);not null;column:diagram_type" db:"diagram_type"` // business_flow, architecture, data_model
	DiagramName        string     `json:"diagram_name" gorm:"type:varchar(200);not null;column:diagram_name" db:"diagram_name"`
	PUMLContent        string     `json:"puml_content" gorm:"type:text;not null;column:puml_content" db:"puml_content"`
	RenderedURL        string     `json:"rendered_url" gorm:"type:varchar(255);column:rendered_url" db:"rendered_url"`
	Version            int        `json:"version" gorm:"default:1;column:version" db:"version"`
	Stage              int        `json:"stage" gorm:"default:1;column:stage" db:"stage"`                     // 新增：所属阶段 1,2,3
	TaskID             *uuid.UUID `json:"task_id,omitempty" gorm:"type:char(36);column:task_id" db:"task_id"` // 新增：关联的异步任务ID
	IsValidated        bool       `json:"is_validated" gorm:"default:false;column:is_validated" db:"is_validated"`
	ValidationFeedback string     `json:"validation_feedback" gorm:"type:text;column:validation_feedback" db:"validation_feedback"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime;column:updated_at" db:"updated_at"`
}

// TableName 指定表名
func (PUMLDiagram) TableName() string {
	return "puml_diagrams"
}
