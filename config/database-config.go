package config

import "time"

// ------------ DATABASE ------------
type DatabaseConfig struct {
	HOST            string        `mapstructure:"HOST"`             // адрес хоста PostgreSQL
	PORT            int64         `mapstructure:"PORT"`             // порт базы
	USER            string        `mapstructure:"USER"`             // пользователь
	PASS            string        `mapstructure:"PASS"`             // пароль
	NAME            string        `mapstructure:"NAME"`             // имя базы
	MaxOpenConns    int           `mapstructure:"MAX_OPEN_CONNS"`   // максимальное число открытых соединений
	MaxIdleConns    int           `mapstructure:"MAX_IDLE_CONNS"`   // максимальное число простаивающих соединений
	ConnMaxLifetime time.Duration `mapstructure:"CONN_MAX_LIFETIME"`// срок жизни соединения
	ConnMaxIdleTime time.Duration `mapstructure:"CONN_MAX_IDLE_TIME"`// время простоя до закрытия
}
