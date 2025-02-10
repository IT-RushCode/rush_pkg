package services

import (
	"context"

	dto "github.com/IT-RushCode/rush_pkg/dto/payment"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	yookassa "github.com/IT-RushCode/rush_pkg/services/yookassa"
	yoocommon "github.com/IT-RushCode/rush_pkg/services/yookassa/common"
	yoopayment "github.com/IT-RushCode/rush_pkg/services/yookassa/payment"
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

	query := &yoopayment.Payment{
		Metadata:           req.Metadata,
		Description:        req.Description,
		MerchantCustomerID: req.MerchantCustomerID,
		Capture:            true,
		Test:               utils.OrDefault(store.IsTest, false),
		Amount: &yoocommon.Amount{
			Value:    req.Amount,
			Currency: req.Currency,
		},
		Confirmation: yoopayment.Redirect{
			Type:      yoopayment.TypeRedirect,
			ReturnURL: req.ReturnURL,
		},
	}

	// Конвертация OrderDataDTO в yoopayment.Receipt
	if req.ReceiptData != nil {
		receipt := &yoopayment.Receipt{
			Customer: &yoocommon.Customer{
				Email: req.ReceiptData.CustomerEmail,
			},
			Items: []*yoocommon.Item{},
		}
		for _, item := range req.ReceiptData.Items {
			receipt.Items = append(receipt.Items, &yoocommon.Item{
				Description: item.Description,
				Quantity:    item.Quantity,
				Amount: &yoocommon.Amount{
					Value:    item.Amount.Value,
					Currency: item.Amount.Currency,
				},
				VatCode: item.VatCode,
			})
		}

		query.Receipt = receipt
	}

	payment, err := paymentKassa.CreatePayment(query)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
