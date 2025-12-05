package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	chat "github.com/IT-RushCode/rush_pkg/handlers/chat"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/services"
)

type Handlers struct {
	Sms             *SmsHandler
	Notification    *NotificationHandler
	WebSocket       *chat.WebSocketHandler
	Chat            *chat.ChatHandler
	Policy          *PolicyHandler
	AppVersion      *AppVersionHandler
	YooKassaSetting *YooKassaSettingHandler
}

func NewHandlers(
	cfg *config.Config,
	repo *repositories.Repositories,
	srv *services.Services,
) *Handlers {
	return &Handlers{
		Sms:             NewSMSHandler(cfg, srv, repo.Redis),
		Notification:    NewNotificationHandler(cfg, srv, repo.Redis, repo),
		WebSocket:       chat.NewWebSocketHandler(repo),
		Chat:            chat.NewChatHandler(srv, repo),
		Policy:          NewPolicyHandler(repo),
		AppVersion:      NewAppVersionHandler(repo),
		YooKassaSetting: NewYooKassaSettingHandler(repo),
	}
}
