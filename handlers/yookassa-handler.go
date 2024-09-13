// handlers/payment_handler.go
package handlers

import (
	"github.com/IT-RushCode/rush_pkg/dto/payment"
	"github.com/IT-RushCode/rush_pkg/services"
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	srv *services.Services
}

func NewPaymentHandler(srv *services.Services) *PaymentHandler {
	return &PaymentHandler{srv: srv}
}

func (h *PaymentHandler) CreatePayment(ctx *fiber.Ctx) error {
	var req payment.PaymentRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	payment, err := h.srv.Payment.CreatePayment(ctx.Context(), &req)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, utils.Success, payment)
}
