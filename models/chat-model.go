package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatSession представляет собой сессию чата между пользователем и техподдержкой
type ChatSession struct {
	ID        uint       `gorm:"primaryKey;autoincrement"`
	ClientID  uint       `gorm:"index"`            // Идентификатор клиента
	SupportID *uint      `gorm:"index;null"`       // Идентификатор поддержки, может быть null до первого ответа
	Status    string     `gorm:"type:varchar(20)"` // Статус сессии (активен, завершен)
	StartedAt time.Time  `gorm:"autoCreateTime"`
	ClosedAt  *time.Time `gorm:"default:null"`
}

// Настройки ChatSession
type ChatSessions []ChatSession

func (ChatSession) TableName() string {
	return "ChatSessions"
}

func (m *ChatSession) BeforeCreate(db *gorm.DB) (err error) {
	if err := CheckSequence(m.TableName(), db); err != nil {
		return err
	}
	return nil
}

// ------------------ CHAT MESSAGES ------------------>

// ChatMessage представляет сообщение в чате
type ChatMessage struct {
	ID        uint      `json:"id"`
	SessionID uint      `json:"sessionId"`
	SenderID  uint      `json:"senderId"`  // ID отправителя (может быть клиент или поддержка)
	IsSupport bool      `json:"isSupport"` // Флаг, который указывает, отправлено ли сообщение поддержкой
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

// Настройки ChatMessage
type ChatMessages []ChatMessage

func (ChatMessage) TableName() string {
	return "ChatMessages"
}

func (m *ChatMessage) BeforeCreate(db *gorm.DB) (err error) {
	if err := CheckSequence(m.TableName(), db); err != nil {
		return err
	}
	return nil
}
