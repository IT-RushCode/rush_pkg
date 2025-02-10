package payment

type PaymentRequest struct {
	PointID            uint          `json:"pointId" validate:"required"`
	Amount             string        `json:"amount" validate:"required"`
	Currency           string        `json:"currency" validate:"required"`
	ReturnURL          string        `json:"returnUrl" validate:"required"`
	Description        string        `json:"description" validate:"required"`
	MerchantCustomerID string        `json:"merchantCustomerID"` // email или номер телефона
	ReceiptData        *OrderDataDTO `json:"orderReceipt"`

	// Метаданные должны быть уникальными для каждого платежа
	Metadata map[string]interface{} `json:"metadata"`
}

type OrderDataDTO struct {
	CustomerEmail string         `json:"customerEmail" validate:"required"`
	Items         []OrderItemDTO `json:"items"`
}

type OrderItemDTO struct {
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	Amount      *OrderItemAmountDTO
	VatCode     uint `json:"vat_code"`

	// default: full_payment
	PaymentMode string `json:"payment_mode"`

	// default: commodity
	PaymentSubject string `json:"payment_subject"`
}

type OrderItemAmountDTO struct {
	// Amount in the selected currency, in the form of a string with a dot separator,
	// for example, 10.00.
	Value string `json:"value,omitempty"`

	// Three-letter currency code in the ISO-4217 format. Example: RUB.
	Currency string `json:"currency,omitempty"`
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
