package dto

type GetChatHistoryDTO struct {
	SessionID uint `json:"sessionId"` // Идентификатор сессии
	Limit     uint `json:"limit"`     // Лимит на количество сообщений
	Offset    uint `json:"offset"`    // Смещение для пагинации
}
