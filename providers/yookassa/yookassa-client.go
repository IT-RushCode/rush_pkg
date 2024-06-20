package providers

import (
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
)

// KassaClient представляет собой структуру, содержащую клиента ЮKassa.
type KassaClient struct {
	Client *yookassa.Client
}

// NewKassaClient создает нового клиента ЮKassa.
func NewKassaClient(storeId, secretKey string) *KassaClient {
	client := yookassa.NewClient(storeId, secretKey)
	return &KassaClient{Client: client}
}
