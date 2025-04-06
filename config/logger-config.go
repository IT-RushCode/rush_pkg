package config

// LogConfig содержит конфигурацию логирования
type LogConfig struct {
	Level            string         `mapstructure:"LEVEL"`              // Уровень логирования (info, debug, error, etc.)
	AccessLog        RotationConfig `mapstructure:"ACCESS_LOG"`         // Настройки для access логов
	ErrorLog         RotationConfig `mapstructure:"ERROR_LOG"`          // Настройки для error логов
	FileLog          bool           `mapstructure:"FILE_LOG"`           // Включить запись логов в файлы
	FileJsonFormat   bool           `mapstructure:"FILE_JSON_FORMAT"`   // Тип Включить форматирования в json
	LogRequestBody   bool           `mapstructure:"LOG_REQUEST_BODY"`   // Логировать тело запроса
	EnableStackTrace bool           `mapstructure:"ENABLE_STACK_TRACE"` // Включить трассировку стека вызовов

}

// RotationConfig содержит настройки ротации логов
type RotationConfig struct {
	Filename   string `mapstructure:"FILENAME"`    // Имя файла лога
	MaxSize    int    `mapstructure:"MAX_SIZE"`    // Максимальный размер в МБ до ротации
	MaxBackups int    `mapstructure:"MAX_BACKUPS"` // Максимальное количество файлов
	MaxAge     int    `mapstructure:"MAX_AGE"`     // Максимальный возраст файлов в днях
	Compress   bool   `mapstructure:"COMPRESS"`    // Сжимать ротированные файлы
}
