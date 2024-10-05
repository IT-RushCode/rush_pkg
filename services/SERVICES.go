package services

import (
	"log"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/repositories"
)

type Services struct {
	Firebase *FirebaseService
	Iiko     *IikoService
	Payment  *PaymentService
	Sms      *SmsService
}

func NewServices(cfg *config.Config, repo *repositories.Repositories) *Services {
	fbSrv, err := NewFirebaseService(repo, &cfg.FIREBASE)
	if err != nil {
		log.Println(err)
	}

	return &Services{
		Firebase: fbSrv,
		Iiko:     NewIikoService("", "", "", ""),
		Payment:  NewPaymentService(repo),
		Sms:      NewSmsService(&cfg.SMS),
	}
}
