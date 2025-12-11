package repositories

import (
	"context"

	dto "github.com/IT-RushCode/rush_pkg/dto/chat"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/utils"

	"gorm.io/gorm"
)

// ChatRepository представляет методы для работы с чатом
type ChatRepository interface {
	CreateSession(ctx context.Context, session *models.ChatSession) error
	CloseSession(ctx context.Context, sessionID string) error
	CreateMessage(ctx context.Context, message *models.ChatMessage) error
	GetMessages(ctx context.Context, dto *dto.GetChatHistoryDTO) ([]models.ChatMessage, error)
	GetActiveSession(ctx context.Context, clientID uint) (*models.ChatSession, error)
	GetSessionByID(ctx context.Context, sessionID uint) (*models.ChatSession, error)
}

// chatRepository представляет собой реализацию репозитория для работы с чатом
type chatRepository struct {
	db *gorm.DB
}

// NewChatRepository создает новый экземпляр репозитория чата
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

// CreateSession создает новую сессию чата
func (r *chatRepository) CreateSession(ctx context.Context, session *models.ChatSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// CloseSession закрывает сессию чата
func (r *chatRepository) CloseSession(ctx context.Context, sessionID string) error {
	var session models.ChatSession
	if err := r.db.WithContext(ctx).
		Where("id = ? AND status = 'active'", sessionID).
		First(&session).Error; err != nil {
		return err
	}

	// Закрываем сессию, если она активна
	return r.db.WithContext(ctx).Model(&models.ChatSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":    "closed",
			"closed_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// CreateMessage создает новое сообщение в чате
func (r *chatRepository) CreateMessage(ctx context.Context, message *models.ChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetMessages получает все сообщения для сессии чата
func (r *chatRepository) GetMessages(ctx context.Context, dto *dto.GetChatHistoryDTO) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Session(&gorm.Session{Context: ctx}).
		Where("session_id = ?", dto.SessionID).
		Order("timestamp DESC").
		Scopes(utils.Paginate(dto.Offset, dto.Limit)).
		Find(&messages).Error
	return messages, err
}

// GetActiveSession получает активную сессию для конкретного клиента
func (r *chatRepository) GetActiveSession(ctx context.Context, clientID uint) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.WithContext(ctx).
		Where("client_id = ? AND status = ?", clientID, "active").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionByID получает сессию по её ID
func (r *chatRepository) GetSessionByID(ctx context.Context, sessionID uint) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.WithContext(ctx).
		Where("id = ?", sessionID).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}
