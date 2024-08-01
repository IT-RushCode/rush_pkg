package config

// ------------ REDIS ------------
type RedisConfig struct {
	HOST string `mapstructure:"RD_HOST"`
	PORT string `mapstructure:"RD_PORT"`
	PASS string `mapstructure:"RD_PASS"`
	DB   int    `mapstructure:"RD_DB"`
}
