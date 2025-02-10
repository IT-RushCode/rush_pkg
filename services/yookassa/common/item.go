package yoocommon

type Item struct {
	// parameter with the name of the product or service
	Description string `json:"description"`

	// parameter with the amount per unit of product
	Quantity uint `json:"quantity"`

	// parameter specifying the quantity of goods (only integers, for example 1)
	Amount *Amount `json:"amount"`

	// parameter with the fixed value 1 (price without VAT)
	VatCode uint `json:"vat_code"`

	// default: full_payment
	PaymentMode string `json:"payment_mode"`

	// default: commodity
	PaymentSubject string `json:"payment_subject"`
}
