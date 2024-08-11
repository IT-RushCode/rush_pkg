package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ------------ GLOBAL CONFIGS ------------
type Config struct {
	APP      AppConfig      `mapstructure:"APP"`
	DB       DatabaseConfig `mapstructure:"DB"`
	JWT      JwtConfig      `mapstructure:"JWT"`
	REDIS    RedisConfig    `mapstructure:"REDIS"`
	KAFKA    KafKaConfig    `mapstructure:"KAFKA"`
	RABBITMQ RabbitMQConfig `mapstructure:"RABBITMQ"`
	MONGODB  MongoDBConfig  `mapstructure:"MONGODB"`
	APISMS   APISMS         `mapstructure:"APISMS"`
	MAIL     MailConfig     `mapstructure:"MAIL"`
	FIREBASE FirebaseConfig `mapstructure:"FIREBASE"`
	FIBER    FiberConfig    `mapstructure:"FIBER"`
}

func InitConfig(path string) *Config {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		panic("failed to read config file")
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic("failed to unmarshal config")
	}

	return &cfg
}
