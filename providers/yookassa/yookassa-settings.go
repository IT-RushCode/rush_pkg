package yookassa

import (
	"fmt"

	"github.com/rvinnie/yookassa-sdk-go/yookassa"
)

type SettingsKassa struct {
	settingsHandler *yookassa.SettingsHandler
}

// NewSettingsKassa создает новый SettingsKassa хендлер.
func NewSettingsKassa(client *KassaClient) *SettingsKassa {
	return &SettingsKassa{
		settingsHandler: yookassa.NewSettingsHandler(client.Client),
	}
}

// Получение информации о настройках магазина или шлюза
func (k *SettingsKassa) GetStoreSettings() {
	settings, err := k.settingsHandler.GetAccountSettings(nil)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	fmt.Println("Успешно: ", settings)
}
