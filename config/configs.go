package config

import (
	"log"

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
	SMS      SMS            `mapstructure:"SMS"`
	MAIL     MailConfig     `mapstructure:"MAIL"`
	FIREBASE FirebaseConfig `mapstructure:"FIREBASE"`
	FIBER    FiberConfig    `mapstructure:"FIBER"`
}

func InitConfig(path string) *Config {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		panic("failed to read config file")
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic("failed to unmarshal config")
	}

	if cfg.APP.ENV != "dev" && cfg.APP.ENV != "test" && cfg.APP.ENV != "prod" {
		log.Fatal("APP_ENV should be one of type: dev/test/prod")
	}

	if cfg.JWT.JWT_SECRET == "" {
		log.Fatal("JWT_SECRET обязателен к заполнению")
	}

	return &cfg
}
