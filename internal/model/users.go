package model

import (
	"github.com/google/uuid"
	"time"
)

// User 用户模型
type User struct {
	UserID       uuid.UUID  `json:"user_id" gorm:"type:char(36);primaryKey;column:user_id" db:"user_id"`
	Username     string     `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;column:username" db:"username"`
	Email        string     `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;column:email" db:"email"`
	PasswordHash string     `json:"-" gorm:"type:varchar(255);not null;column:password_hash" db:"password_hash"` // 不返回给前端
	FullName     string     `json:"full_name" gorm:"type:varchar(100);not null;column:full_name" db:"full_name"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime;column:updated_at" db:"updated_at"`
	LastLogin    *time.Time `json:"last_login" gorm:"column:last_login" db:"last_login"`
	Status       string     `json:"status" gorm:"type:varchar(20);default:'active';column:status" db:"status"`
	Preferences  string     `json:"preferences" gorm:"type:json;column:preferences" db:"preferences"` // JSON字符串
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
