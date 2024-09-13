package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/services"
)

type Handlers struct {
	Payment      *PaymentHandler
	Sms          *SmsHandler
	Notification *NotificationHandler
}

func NewHandlers(
	cfg *config.Config,
	repo *repositories.Repositories,
	srv *services.Services,
) *Handlers {
	return &Handlers{
		Payment:      NewPaymentHandler(srv),
		Sms:          NewSMSHandler(cfg, srv, repo.Redis),
		Notification: NewNotificationHandler(cfg, srv, repo.Redis),
	}
}
