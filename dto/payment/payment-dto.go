package payment

type PaymentRequest struct {
	PointID     uint   `json:"pointId" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
	OrderNumber string `json:"orderNumber" validate:"required"`
	ReturnURL   string `json:"returnUrl" validate:"required"`
	Description string `json:"description" validate:"required"`
}
