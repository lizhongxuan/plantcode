package model

import (
	"github.com/google/uuid"
	"time"
)

// Question 补充问题模型
type Question struct {
	QuestionID       uuid.UUID  `json:"question_id" gorm:"type:char(36);primaryKey;column:question_id" db:"question_id"`
	RequirementID    uuid.UUID  `json:"requirement_id" gorm:"type:char(36);not null;index;column:requirement_id" db:"requirement_id"`
	QuestionText     string     `json:"question_text" gorm:"type:text;not null;column:question_text" db:"question_text"`
	QuestionCategory string     `json:"question_category" gorm:"type:varchar(50);column:question_category" db:"question_category"`
	PriorityLevel    int        `json:"priority_level" gorm:"default:1;column:priority_level" db:"priority_level"`
	AnswerText       string     `json:"answer_text" gorm:"type:text;column:answer_text" db:"answer_text"`
	AnswerStatus     string     `json:"answer_status" gorm:"type:varchar(20);default:'pending';column:answer_status" db:"answer_status"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	AnsweredAt       *time.Time `json:"answered_at" gorm:"column:answered_at" db:"answered_at"`
}

// TableName 指定表名
func (Question) TableName() string {
	return "questions"
}
