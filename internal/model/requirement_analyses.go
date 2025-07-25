package model

import (
	"github.com/google/uuid"
	"time"
)

// Requirement 需求分析模型
type Requirement struct {
	RequirementID         uuid.UUID `json:"requirement_id" gorm:"type:char(36);primaryKey;column:requirement_id" db:"requirement_id"`
	ProjectID             uuid.UUID `json:"project_id" gorm:"type:char(36);not null;index;column:project_id" db:"project_id"`
	RawRequirement        string    `json:"raw_requirement" gorm:"type:text;not null;column:raw_requirement" db:"raw_requirement"`
	StructuredRequirement string    `json:"structured_requirement" gorm:"type:json;column:structured_requirement" db:"structured_requirement"` // JSON
	CompletenessScore     float64   `json:"completeness_score" gorm:"type:decimal(5,2);default:0;column:completeness_score" db:"completeness_score"`
	AnalysisStatus        string    `json:"analysis_status" gorm:"type:varchar(50);default:'pending';column:analysis_status" db:"analysis_status"`
	MissingInfoTypes      string    `json:"missing_info_types" gorm:"type:json;column:missing_info_types" db:"missing_info_types"` // JSON
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Requirement) TableName() string {
	return "requirement_analyses"
}
