package yookassa

import (
	"fmt"
	"log"

	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type PaymentKassa struct {
	paymentHandler *yookassa.PaymentHandler
}

// NewPaymentKassa создает новый PaymentKassa хендлер.
func NewPaymentKassa(client *KassaClient) *PaymentKassa {
	return &PaymentKassa{
		paymentHandler: yookassa.NewPaymentHandler(client.Client),
	}
}

// Создание платежа
func (k *PaymentKassa) CreatePayment(paymentMethod, paymentType, amount, currency, description, returnURL string) {
	var confirmation yoopayment.Confirmer

	switch paymentType {
	case "embedded":
		confirmation = &yoopayment.Embedded{
			Type:              yoopayment.TypeEmbedded,
			ConfirmationToken: "confirmation_token_here",
		}
	case "redirect":
		confirmation = &yoopayment.Redirect{
			Type:      yoopayment.TypeRedirect,
			ReturnURL: returnURL,
		}
	case "mobile_application":
		confirmation = &yoopayment.MobileApplication{
			Type:            yoopayment.TypeMobileApplication,
			ConfirmationURL: "confirmation_url_here",
		}
	case "qr":
		confirmation = &yoopayment.QR{
			Type:             yoopayment.TypeQR,
			ConfirmationData: "confirmation_data_here",
		}
	default:
		log.Fatalf("Неподдерживаемый метод оплаты: %s", paymentType)
	}

	payment, err := k.paymentHandler.CreatePayment(&yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    amount,
			Currency: currency,
		},
		PaymentMethod: yoopayment.PaymentMethodType(paymentMethod),
		Confirmation:  confirmation,
		Description:   description,
	})
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}

	fmt.Println("Успешно: ", payment)
}

// Получение информации о платеже
func (k *PaymentKassa) GetPayment(paymentId string) {
	payment, err := k.paymentHandler.FindPayment(paymentId)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	fmt.Println("Успешно: ", payment)
}

// Получение 'succeeded' платежей по 5 штук за запрос
func (k *PaymentKassa) GetPayments() {
	payments, err := k.paymentHandler.FindPayments(&yoopayment.PaymentListFilter{
		Limit:  5,
		Status: yoopayment.Succeeded,
	})
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	fmt.Println("Успешно: ", payments)
}
