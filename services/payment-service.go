package services

import (
	"context"
	"fmt"

	dto "github.com/IT-RushCode/rush_pkg/dto/payment"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	yookassa "github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type PaymentService struct {
	repo *repositories.Repositories
}

func NewPaymentService(repo *repositories.Repositories) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.PaymentRequest) (*yoopayment.Payment, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	store := &models.YooKassaSetting{}
	if err := s.repo.YooKassaSetting.Filter(
		ctx,
		map[string]interface{}{"point_id": req.PointID},
		store,
	); err != nil {
		return nil, fmt.Errorf("настройки YooKassa для PointID = %d не найдены", req.PointID)
	}

	client := yookassa.NewClient(store.StoreID, store.SecretKey)
	paymentKassa := yookassa.NewPaymentHandler(client)

	var payment *yoopayment.Payment
	payment, err := paymentKassa.CreatePayment(&yoopayment.Payment{
		Metadata:    map[string]interface{}{"orderNumber": req.OrderNumber},
		Description: req.Description,
		Capture:     true,
		Amount: &yoocommon.Amount{
			Value:    req.Amount,
			Currency: req.Currency,
		},
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: req.ReturnURL,
		},
	})
	if err != nil {
		return nil, err
	}

	return payment, nil
}
