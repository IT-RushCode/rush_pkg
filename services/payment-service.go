package services

import (
	"context"

	dto "github.com/IT-RushCode/rush_pkg/dto/payment"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"

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

func (s *PaymentService) CreatePayment(
	ctx context.Context,
	store *models.YooKassaSetting,
	req *dto.PaymentRequest,
) (*yoopayment.Payment, error) {
	client := yookassa.NewClient(store.StoreID, store.SecretKey)
	paymentKassa := yookassa.NewPaymentHandler(client)

	var payment *yoopayment.Payment
	payment, err := paymentKassa.CreatePayment(&yoopayment.Payment{
		Metadata:           req.Metadata,
		Description:        req.Description,
		MerchantCustomerID: req.MerchantCustomerID,
		Capture:            true,
		Test:               *store.IsTest,
		Amount: &yoocommon.Amount{
			Value:    req.Amount,
			Currency: req.Currency,
		},
		Confirmation: yoopayment.Redirect{
			Type:      yoopayment.TypeRedirect,
			ReturnURL: req.ReturnURL,
		},
	})
	if err != nil {
		return nil, err
	}

	return payment, nil
}
