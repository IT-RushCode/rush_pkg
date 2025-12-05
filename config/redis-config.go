package config

import "time"

// ------------ REDIS ------------
type RedisConfig struct {
	HOST               string        `mapstructure:"HOST"`                 // адрес хоста
	PORT               string        `mapstructure:"PORT"`                 // порт Redis
	PASS               string        `mapstructure:"PASS"`                 // пароль
	DB                 int           `mapstructure:"DB"`                   // индекс базы 1/16
	MaxRetries         int           `mapstructure:"MAX_RETRIES"`          // максимум повторов в случае ошибки
	MinRetryBackoff    time.Duration `mapstructure:"MIN_RETRY_BACKOFF"`    // минимальная пауза между повторами
	MaxRetryBackoff    time.Duration `mapstructure:"MAX_RETRY_BACKOFF"`    // максимальная пауза между повторами
	DialTimeout        time.Duration `mapstructure:"DIAL_TIMEOUT"`         // таймаут установления соединения
	DialerRetries      int           `mapstructure:"DIALER_RETRIES"`       // число попыток подключения
	DialerRetryTime    time.Duration `mapstructure:"DIALER_RETRY_TIME"`    // пауза между попытками подключения
	ReadTimeout        time.Duration `mapstructure:"READ_TIMEOUT"`         // таймаут чтения
	WriteTimeout       time.Duration `mapstructure:"WRITE_TIMEOUT"`        // таймаут записи
	ContextTimeout     bool          `mapstructure:"CONTEXT_TIMEOUT"`      // учитывать таймаут контекста
	PoolSize           int           `mapstructure:"POOL_SIZE"`            // размер пула соединений
	MaxActiveConns     int           `mapstructure:"MAX_ACTIVE_CONNS"`     // максимум активных соединений
	MaxConcurrentDials int           `mapstructure:"MAX_CONCURRENT_DIALS"` // максимум одновременных установок соединения
	PoolTimeout        time.Duration `mapstructure:"POOL_TIMEOUT"`         // таймаут ожидания свободного соединения
	MinIdleConns       int           `mapstructure:"MIN_IDLE_CONNS"`       // минимум простаивающих соединений
	MaxIdleConns       int           `mapstructure:"MAX_IDLE_CONNS"`       // максимум простаивающих соединений
	ConnMaxIdleTime    time.Duration `mapstructure:"CONN_MAX_IDLE_TIME"`   // максимум времени простоя
	ConnMaxLifetime    time.Duration `mapstructure:"CONN_MAX_LIFETIME"`    // максимум времени жизни соединения
}
