package payment

type PaymentRequest struct {
	PointID     uint                   `json:"pointId" validate:"required"`
	Amount      string                 `json:"amount" validate:"required"`
	Currency    string                 `json:"currency" validate:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
	ReturnURL   string                 `json:"returnUrl" validate:"required"`
	Description string                 `json:"description" validate:"required"`
}

type WebhookRequest struct {
	Event  string        `json:"event"`  // Тип события, например: "payment.succeeded" или "payment.canceled"
	Object WebhookObject `json:"object"` // Детали объекта (платежа) из вебхука
}

type WebhookObject struct {
	ID       string                 `json:"id"`       // ID платежа в YooKassa
	Status   string                 `json:"status"`   // Статус платежа, например: "succeeded", "canceled"
	Amount   WebhookAmount          `json:"amount"`   // Сумма платежа
	Metadata map[string]interface{} `json:"metadata"` // Дополнительные данные, например, orderNumber
}

type WebhookAmount struct {
	Value    string `json:"value"`    // Сумма платежа
	Currency string `json:"currency"` // Валюта платежа, например, "RUB"
}
