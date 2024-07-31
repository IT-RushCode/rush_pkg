package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/repositories"
)

type Handlers struct {
	PaymentHandler *PaymentHandler
	SmsHandler     *SmsHandler
}

func NewHandlers(repo *repositories.Repositories, cfg *config.Config) *Handlers {
	return &Handlers{
		PaymentHandler: NewPaymentHandler(repo),
		SmsHandler:     NewSMSHandler(cfg, repo.Redis),
	}
}
