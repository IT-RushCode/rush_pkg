package handlers

import (
	"fmt"

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
		ctx.Context(),
		map[string]interface{}{"point_id": req.PointID},
		store,
	); err != nil {
		return utils.ErrorNotFoundResponse(ctx, "настройки YooKassa не найдены", nil)
	}

	client := yookassa.NewClient(store.StoreID, store.SecretKey)
	paymentKassa := yookassa.NewPaymentHandler(client)

	var payment *yoopayment.Payment
	switch yoopayment.PaymentMethodType(req.PaymentMethod) {
	case yoopayment.PaymentTypeBankCard:
		payment, err = paymentKassa.CreatePayment(&yoopayment.Payment{
			Metadata: map[string]interface{}{
				"orderID": 1,
			},
			MerchantCustomerID: "felixkot00@gmail.com",
			Capture:            true,
			Amount: &yoocommon.Amount{
				Value:    req.Amount,
				Currency: req.Currency,
			},
			PaymentMethod: yoopayment.BankCard{
				Card: yoopayment.Card{
					First6:        "220000",
					Last4:         "0004",
					ExpiryYear:    "05",
					ExpiryMonth:   "2030",
					CardType:      "MIR",
					IssuerCountry: "RU",
					IssuerName:    "Sberbank",
					Source:        "sber_pay",
				},
			},
			Confirmation: &yoopayment.Redirect{
				Type:      yoopayment.TypeRedirect,
				Enforce:   true,
				ReturnURL: "http://localhost:8000",
			},
			Description: req.Description,
		})
	case yoopayment.PaymentTypeSberbank,
		yoopayment.PaymentTypeYooMoney,
		yoopayment.PaymentTypeSBP:
		payment, err = paymentKassa.CreatePayment(&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value:    req.Amount,
				Currency: req.Currency,
			},
			PaymentMethod: yoopayment.PaymentMethodType(req.PaymentMethod),
			Confirmation:  yoopayment.TypeEmbedded,

			// &yoopayment.Redirect{
			// 	Type:      yoopayment.TypeRedirect,
			// 	ReturnURL: req.ReturnURL,
			// },
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
		return utils.ErrorBadRequestResponse(
			ctx,
			fmt.Sprintf(
				"не поддерживаемый метод оплаты. Используйте один из методов (%s, %s, %s, %s, %s)",
				yoopayment.PaymentTypeCash,
				yoopayment.PaymentTypeBankCard,
				yoopayment.PaymentTypeYooMoney,
				yoopayment.PaymentTypeSBP,
				yoopayment.PaymentTypeSberbank,
			),
			nil,
		)
	}

	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, utils.Success, payment)
}
