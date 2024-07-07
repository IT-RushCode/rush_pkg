package sms

import (
	"github.com/go-playground/validator/v10"
)

// SMSRequest представляет запрос от клиента для отправки SMS
type SMSRequestDTO struct {
	Messages []SMSMessage `json:"messages" validate:"required"`

	// Если указан true, то Text из Messages игорируется
	// и отправится сгенерированный OTP код на номер телефона.
	// Иначе отправляется Text из Messages.
	IsOTP bool `json:"isOTP"`
}

// SMSMessage представляет собой структуру для одного SMS-сообщения.
type SMSMessage struct {
	Text    string `json:"text" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Channel string `json:"channel" validate:"required" default:"char"`
	Sender  string `json:"sender"`
}

// Validate проверяет структуру SMSMessage с учетом кастомной валидации.
func (m *SMSMessage) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("required_if_char", func(fl validator.FieldLevel) bool {
		smsMessage := fl.Top().Interface().(SMSMessage)
		if smsMessage.Channel == "char" && smsMessage.Sender == "" {
			return false
		}
		return true
	})

	return validate.Struct(m)
}

type VerifyRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	OTPCode     string `json:"code" validate:"required"`
}

type SmsSenderResponse struct {
	Message SMSRes `json:"message,omitempty"`
	Phone   string `json:"phone,omitempty"`
	OTPCode string `json:"otpCode,omitempty"`
}

// Ответ смс отправителя
type SMSRes struct {
	Status string `json:"status,omitempty"`
	Data   []Data `json:"data,omitempty"`
}

type Data struct {
	ID        int64       `json:"id,omitempty"`
	Message   string      `json:"message,omitempty"`
	DebugInfo string      `json:"debugInfo,omitempty"`
	Exception string      `json:"exception,omitempty"`
	CreatedAt int64       `json:"createdAt,omitempty"`
	Tag       interface{} `json:"tag,omitempty"`
	Status    string      `json:"status,omitempty"`
	Phone     string      `json:"phone,omitempty"`
}
