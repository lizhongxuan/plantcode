package model

import (
	"github.com/google/uuid"
	"time"
)

// CommonModule 通用模块库模型
type CommonModule struct {
	CommonModuleID uuid.UUID `json:"common_module_id" gorm:"type:char(36);primaryKey;column:common_module_id" db:"common_module_id"`
	ModuleName     string    `json:"module_name" gorm:"type:varchar(100);not null;column:module_name" db:"module_name"`
	Category       string    `json:"category" gorm:"type:varchar(50);not null;index;column:category" db:"category"`
	Description    string    `json:"description" gorm:"type:text;column:description" db:"description"`
	Functionality  string    `json:"functionality" gorm:"type:json;column:functionality" db:"functionality"`    // JSON
	InterfaceSpec  string    `json:"interface_spec" gorm:"type:json;column:interface_spec" db:"interface_spec"` // JSON
	CodeTemplate   string    `json:"code_template" gorm:"type:text;column:code_template" db:"code_template"`
	UsageExamples  string    `json:"usage_examples" gorm:"type:json;column:usage_examples" db:"usage_examples"` // JSON
	Version        string    `json:"version" gorm:"type:varchar(20);default:'1.0.0';column:version" db:"version"`
	DownloadsCount int       `json:"downloads_count" gorm:"default:0;column:downloads_count" db:"downloads_count"`
	Rating         float64   `json:"rating" gorm:"type:decimal(3,2);default:0;column:rating" db:"rating"`
	Tags           string    `json:"tags" gorm:"type:json;column:tags" db:"tags"` // JSON
	CreatedBy      uuid.UUID `json:"created_by" gorm:"type:char(36);not null;column:created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at" db:"updated_at"`
}

// TableName 指定表名
func (CommonModule) TableName() string {
	return "common_modules"
}
