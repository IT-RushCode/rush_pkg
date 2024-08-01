package config

// ------------ APP CONFIG ------------
type AppConfig struct {
	ENV            string `mapstructure:"APP_ENV"`
	DEBUG          bool   `mapstructure:"APP_DEBUG"`
	NAME           string `mapstructure:"APP_NAME"`
	HOST           string `mapstructure:"APP_HOST"`
	PORT           string `mapstructure:"APP_PORT"`
	MAX_CONNECTION int    `mapstructure:"APP_MAX_CONNECTION"`
	DATETIME       string `mapstructure:"DATETIME"`
	DATE           string `mapstructure:"DATE"`
	TIME           string `mapstructure:"TIME"`
}
