package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/repositories"
)

type Handlers struct {
	Payment *PaymentHandler
	Sms     *SmsHandler
}

func NewHandlers(
	repo *repositories.Repositories,
	cfg *config.Config,
) *Handlers {
	return &Handlers{
		Payment: NewPaymentHandler(repo),
		Sms:     NewSMSHandler(cfg, repo.Redis),
	}
}
