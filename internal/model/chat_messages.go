package model

import (
	"github.com/google/uuid"
	"time"
)

// ChatMessage 对话消息模型
type ChatMessage struct {
	MessageID      uuid.UUID `json:"message_id" gorm:"type:char(36);primaryKey;column:message_id" db:"message_id"`
	SessionID      uuid.UUID `json:"session_id" gorm:"type:char(36);not null;index;column:session_id" db:"session_id"`
	SenderType     string    `json:"sender_type" gorm:"type:varchar(20);not null;column:sender_type" db:"sender_type"` // user, ai, system
	MessageContent string    `json:"message_content" gorm:"type:text;not null;column:message_content" db:"message_content"`
	MessageType    string    `json:"message_type" gorm:"type:varchar(50);default:'text';column:message_type" db:"message_type"` // text, question, answer
	Metadata       string    `json:"metadata" gorm:"type:json;column:metadata" db:"metadata"`                                   // JSON字符串
	Timestamp      time.Time `json:"timestamp" gorm:"autoCreateTime;column:timestamp" db:"timestamp"`
	Processed      bool      `json:"processed" gorm:"default:false;column:processed" db:"processed"`
}

// TableName 指定表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}
