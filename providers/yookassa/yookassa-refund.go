package providers

import (
	"fmt"

	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoorefund "github.com/rvinnie/yookassa-sdk-go/yookassa/refund"
)

type RefundKassa struct {
	refundHandler *yookassa.RefundHandler
}

// NewRefundKassa создает новый RefundKassa хендлер.
func NewRefundKassa(client *KassaClient) *RefundKassa {
	return &RefundKassa{
		refundHandler: yookassa.NewRefundHandler(client.Client),
	}
}

// Создание возврата
func (k *RefundKassa) CreateRefund(paymentId, amountValue, amountCurrency string) {
	refund, err := k.refundHandler.CreateRefund(&yoorefund.Refund{
		PaymentId: paymentId,
		Amount: &yoocommon.Amount{
			Value:    amountValue,
			Currency: amountCurrency,
		},
		Description: "Test refund :)",
	})
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	fmt.Println("Успешно: ", refund)
}

// Получение возврата
func (k *RefundKassa) FindRefund(refundId string) {
	refund, err := k.refundHandler.FindRefund(refundId)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	fmt.Println("Успешно: ", refund)
}

// Получение списка объектов возврата (последние 5 со статусом succeeded)
func (k *RefundKassa) FindRefunds() {
	refunds, err := k.refundHandler.FindRefunds(&yoorefund.RefundListFilter{
		Status: yoorefund.Succeeded,
		Limit:  5,
	})
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	fmt.Println("Успешно: ", refunds)
}
