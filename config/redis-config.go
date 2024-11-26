package config

// ------------ REDIS ------------
type RedisConfig struct {
	HOST string `mapstructure:"HOST"`
	PORT string `mapstructure:"PORT"`
	PASS string `mapstructure:"PASS"`
	DB   int    `mapstructure:"DB"`
}
