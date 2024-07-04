package payment

type PaymentRequest struct {
	PaymentMethod string `json:"paymentMethod,omitempty"`
	PointID       uint   `json:"pointId"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
	CardNumber    string `json:"cardNumber,omitempty"`
	ExpiryMonth   string `json:"expiryMonth,omitempty"`
	ExpiryYear    string `json:"expiryYear,omitempty"`
	Cvc           string `json:"cvc,omitempty"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	ReturnURL     string `json:"returnUrl"`
}
