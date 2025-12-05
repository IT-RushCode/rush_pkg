package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatSession представляет собой сессию чата между пользователем и техподдержкой
type ChatSession struct {
	ID        uint       `gorm:"primaryKey;autoincrement;comment:Первичный ключ с автоинкрементом"`
	ClientID  uint       `gorm:"index;comment:Идентификатор клиента"`
	SupportID *uint      `gorm:"index;null;comment:Идентификатор поддержки, null до первого ответа"`
	Status    string     `gorm:"type:varchar(20);comment:Статус сессии (активна, завершена и т.п.)"`
	StartedAt time.Time  `gorm:"autoCreateTime;comment:Время начала сессии"`
	ClosedAt  *time.Time `gorm:"default:null;comment:Время закрытия сессии"`
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
	ID        uint      `json:"id" gorm:"comment:Первичный ключ сообщения"`
	SessionID uint      `json:"sessionId" gorm:"comment:Ссылка на сессию"`
	SenderID  uint      `json:"senderId" gorm:"comment:ID отправителя"`
	IsSupport bool      `json:"isSupport" gorm:"comment:Флаг сообщения от поддержки"`
	Body      string    `json:"body" gorm:"type:text;comment:Текст сообщения"`
	Timestamp time.Time `json:"timestamp" gorm:"comment:Время отправки сообщения"`
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
