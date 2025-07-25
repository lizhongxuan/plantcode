package model

import (
	"github.com/google/uuid"
	"time"
)

// BusinessModule 业务模块模型
type BusinessModule struct {
	ModuleID        uuid.UUID `json:"module_id" gorm:"type:char(36);primaryKey;column:module_id" db:"module_id"`
	ProjectID       uuid.UUID `json:"project_id" gorm:"type:char(36);not null;index;column:project_id" db:"project_id"`
	ModuleName      string    `json:"module_name" gorm:"type:varchar(100);not null;column:module_name" db:"module_name"`
	Description     string    `json:"description" gorm:"type:text;column:description" db:"description"`
	ModuleType      string    `json:"module_type" gorm:"type:varchar(50);column:module_type" db:"module_type"`
	ComplexityLevel string    `json:"complexity_level" gorm:"type:varchar(20);default:'medium';column:complexity_level" db:"complexity_level"`
	BusinessLogic   string    `json:"business_logic" gorm:"type:json;column:business_logic" db:"business_logic"` // JSON
	Interfaces      string    `json:"interfaces" gorm:"type:json;column:interfaces" db:"interfaces"`             // JSON
	Dependencies    string    `json:"dependencies" gorm:"type:json;column:dependencies" db:"dependencies"`       // JSON
	IsReusable      bool      `json:"is_reusable" gorm:"default:false;column:is_reusable" db:"is_reusable"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
}

// TableName 指定表名
func (BusinessModule) TableName() string {
	return "business_modules"
}
