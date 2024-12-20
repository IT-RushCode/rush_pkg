package config

// ------------ API SMS ------------
type SmsConfig struct {
	URL              string `mapstructure:"URL"`              // URL для отправки SMS, указанный в API провайдера.
	METHOD           string `mapstructure:"METHOD"`           // HTTP-метод для запроса (обычно POST).
	TOKEN            string `mapstructure:"TOKEN"`            // Токен авторизации для доступа к API.
	SENDER           string `mapstructure:"SENDER"`           // Имя отправителя (Sender ID), отображаемое в SMS.
	TEMPLATE         string `mapstructure:"TEMPLATE"`         // Шаблон текста SMS с поддержкой форматирования (например, для кода OTP).
	ACTIVE_SEND      bool   `mapstructure:"ACTIVE_SEND"`      // Флаг, указывающий, активна ли отправка SMS (false для тестирования).
	IGNORING_NUMBERS string `mapstructure:"IGNORING_NUMBERS"` // Список номеров, для которых отправка SMS игнорируется (через запятую).
	IGNORING_OTP     string `mapstructure:"IGNORING_OTP"`     // Тестовый OTP-код, для которого проверка игнорируется.
}
