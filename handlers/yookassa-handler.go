package handler

// import (
// 	"github.com/IT-RushCode/rush_pkg/providers/yookassa"
// 	"github.com/IT-RushCode/rush_pkg/utils"

// 	"github.com/gofiber/fiber/v2"
// )

// type Handlers struct {
// 	settingsKassa *yookassa.SettingsKassa
// 	paymentKassa  *yookassa.PaymentKassa
// 	refundKassa   *yookassa.RefundKassa
// }

// func NewHandlers(client *yookassa.KassaClient) *Handlers {
// 	return &Handlers{
// 		settingsKassa: yookassa.NewSettingsKassa(client),
// 		paymentKassa:  yookassa.NewPaymentKassa(client),
// 		refundKassa:   yookassa.NewRefundKassa(client),
// 	}
// }

// // Хендлеры для работы с настройками
// func (h *Handlers) GetSettings(ctx *fiber.Ctx) error {
// 	settings, err := h.settingsKassa.GetStoreSettings()
// 	if err != nil {
// 		return utils.ErrorInternalServerErrorResponse(ctx, err, nil)
// 	}

// 	return utils.SuccessResponse(ctx, "success", settings)
// }

// // Хендлеры для работы с платежами
// func (h *Handlers) CreatePayment(ctx *fiber.Ctx) error {
// 	var req struct {
// 		PaymentMethod string `json:"paymentMethod"`
// 		PaymentType   string `json:"paymentType"`
// 		Amount        string `json:"amount"`
// 		Currency      string `json:"currency"`
// 		ReturnURL     string `json:"returnUrl"`
// 	}
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
// 	}

// 	payment, err := h.paymentKassa.CreatePayment(req.PaymentMethod, req.PaymentType, req.Amount, req.Currency, req.ReturnURL)
// 	if err != nil {
// 		return utils.ErrorInternalServerErrorResponse(ctx, err, nil)
// 	}

// 	return utils.SuccessResponse(ctx, "success", payment)
// }

// func (h *Handlers) GetPayment(ctx *fiber.Ctx) error {
// 	var req struct {
// 		PaymentID string `json:"paymentId"`
// 	}
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
// 	}

// 	payment, err := h.paymentKassa.GetPayment(req.PaymentID)
// 	if err != nil {
// 		return utils.ErrorInternalServerErrorResponse(ctx, err, nil)
// 	}

// 	return utils.SuccessResponse(ctx, "success", payment)
// }

// func (h *Handlers) GetPayments(ctx *fiber.Ctx) error {
// 	payments, err := h.paymentKassa.GetPayments()
// 	if err != nil {
// 		return utils.ErrorInternalServerErrorResponse(ctx, err, nil)
// 	}

// 	return utils.SuccessResponse(ctx, "success", payments)
// }

// // Хендлеры для работы с возвратами
// func (h *Handlers) CreateRefund(ctx *fiber.Ctx) error {
// 	var req struct {
// 		PaymentID string `json:"paymentId"`
// 		Amount    string `json:"amount"`
// 		Currency  string `json:"currency"`
// 	}
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
// 	}

// 	refund, err := h.refundKassa.CreateRefund(req.PaymentID, req.Amount, req.Currency)
// 	if err != nil {
// 		return utils.ErrorInternalServerErrorResponse(ctx, err, nil)
// 	}

// 	return utils.SuccessResponse(ctx, "success", refund)
// }

// func (h *Handlers) GetRefund(ctx *fiber.Ctx) error {
// 	var req struct {
// 		RefundID string `json:"refundId"`
// 	}
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
// 	}

// 	refund, err := h.refundKassa.FindRefund(req.RefundID)
// 	if err != nil {
// 		return utils.ErrorInternalServerErrorResponse(ctx, err, nil)
// 	}

// 	return utils.SuccessResponse(ctx, "success", refund)
// }

// func (h *Handlers) GetRefunds(ctx *fiber.Ctx) error {
// 	refunds, err := h.refundKassa.FindRefunds()
// 	if err != nil {
// 		return utils.ErrorInternalServerErrorResponse(ctx, err, nil)
// 	}

// 	return utils.SuccessResponse(ctx, "success", refunds)
// }
