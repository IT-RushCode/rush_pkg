package handlers

import (
	"context"

	dto "github.com/IT-RushCode/rush_pkg/dto/payment"
	rpYKassa "github.com/IT-RushCode/rush_pkg/models/yookassa"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	yookassa "github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	repo repositories.Repositories
}

func NewPaymentHandler(repo *repositories.Repositories) *PaymentHandler {
	return &PaymentHandler{
		repo: *repo,
	}
}

// TODO: ДОРАБОТАТЬ МЕТОД
func (h *PaymentHandler) CreatePayment(ctx *fiber.Ctx) error {
	var req dto.PaymentRequest
	var err error
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	store := &rpYKassa.YooKassaSetting{}
	if err := h.repo.YooKassaSetting.Filter(
		context.Background(),
		map[string]interface{}{"point_id": req.PointID},
		store,
	); err != nil {
		return utils.ErrorNotFoundResponse(ctx, "настройки YooKassa не найдены", nil)
	}

	client := yookassa.NewClient(store.StoreID, store.SecretKey)
	paymentKassa := yookassa.NewPaymentHandler(client)

	var payment *yoopayment.Payment
	switch yoopayment.PaymentMethodType(req.PaymentType) {
	case yoopayment.PaymentTypeBankCard:
		payment, err = paymentKassa.CreatePayment(&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value:    req.Amount,
				Currency: req.Currency,
			},
			PaymentMethod: yoopayment.PaymentMethodType(yoopayment.PaymentTypeBankCard),
			// Card: &yoopayment.PaymentMethodDataCard{
			// 	Number:      req.CardNumber,
			// 	ExpiryMonth: req.ExpiryMonth,
			// 	ExpiryYear:  req.ExpiryYear,
			// 	Cvc:         req.Cvc,
			// },
			Confirmation: &yoopayment.Redirect{
				Type:      yoopayment.TypeRedirect,
				ReturnURL: req.ReturnURL,
			},
			Description: req.Description,
		})
	case yoopayment.PaymentTypeTinkoffBank, yoopayment.PaymentTypeSberbank, yoopayment.PaymentTypeYooMoney:
		payment, err = paymentKassa.CreatePayment(&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value:    req.Amount,
				Currency: req.Currency,
			},
			PaymentMethod: yoopayment.PaymentMethodType(req.PaymentType),
			Confirmation: &yoopayment.Redirect{
				Type:      yoopayment.TypeRedirect,
				ReturnURL: req.ReturnURL,
			},
			Description: req.Description,
		})
	case yoopayment.PaymentTypeCash:
		payment, err = paymentKassa.CreatePayment(&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value:    req.Amount,
				Currency: req.Currency,
			},
			PaymentMethod: yoopayment.PaymentMethodType(yoopayment.PaymentTypeCash),
			Confirmation: &yoopayment.Redirect{
				Type:      yoopayment.TypeRedirect,
				ReturnURL: req.ReturnURL,
			},
			Description: req.Description,
		})
	default:
		return utils.ErrorBadRequestResponse(ctx, "не поддерживаемый тип оплаты", nil)
	}

	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, utils.Success, payment)
}
