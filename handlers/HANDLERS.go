package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/controllers"
	chat "github.com/IT-RushCode/rush_pkg/handlers/chat"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/services"
)

type Handlers struct {
	Sms          *SmsHandler
	Notification *NotificationHandler
	WebSocket    *chat.WebSocketHandler
	Chat         *chat.ChatHandler
	Policy       *PolicyHandler
}

func NewHandlers(
	cfg *config.Config,
	repo *repositories.Repositories,
	srv *services.Services,
	ctrl *controllers.Controllers,
) *Handlers {
	return &Handlers{
		Sms:          NewSMSHandler(cfg, srv, repo.Redis),
		Notification: NewNotificationHandler(cfg, srv, repo.Redis),
		WebSocket:    chat.NewWebSocketHandler(ctrl.Chat),
		Chat:         chat.NewChatHandler(srv, repo),
		Policy:       NewPolicyHandler(repo),
	}
}
