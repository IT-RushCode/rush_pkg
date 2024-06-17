package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ------------ GLOBAL CONFIG ------------
type Config struct {
	APP      AppConfig      `mapstructure:"APP"`
	DB       DatabaseConfig `mapstructure:"DB"`
	JWT      JwtConfig      `mapstructure:"JWT"`
	DATETIME DateTimeConfig `mapstructure:"DATETIME"`
	REDIS    RedisConfig    `mapstructure:"REDIS"`
	KAFKA    KafKaConfig    `mapstructure:"KAFKA"`
	RABBITMQ RabbitMQConfig `mapstructure:"RABBITMQ"`
	MONGODB  MongoDBConfig  `mapstructure:"MONGODB"`
	DOCS     DocsAuthConfig `mapstructure:"DOCSAUTH"`
}

// ------------ SERVICES ------------
type AppConfig struct {
	ENV            string `mapstructure:"APP_ENV"`
	DEBUG          bool   `mapstructure:"APP_DEBUG"`
	NAME           string `mapstructure:"APP_NAME"`
	HOST           string `mapstructure:"APP_HOST"`
	PORT           string `mapstructure:"APP_PORT"`
	MAX_CONNECTION int    `mapstructure:"APP_MAX_CONNECTION"`
}

// ------------ JWT ------------
type JwtConfig struct {
	JWT_SECRET  string `mapstructure:"JWT_SECRET"`
	JWT_TTL     int64  `mapstructure:"JWT_TTL"`
	REFRESH_TTL int64  `mapstructure:"REFRESH_TTL"`
}

// ------------ KAFKA ------------
type KafKaConfig struct {
	HOST1 string `mapstructure:"HOST1"`
	HOST2 string `mapstructure:"HOST2"`
	HOST3 string `mapstructure:"HOST3"`
}

// ------------ DATETIME ------------
type DateTimeConfig struct {
	Datetime string `mapstructure:"DATETIME"`
	Date     string `mapstructure:"DATE"`
	Time     string `mapstructure:"TIME"`
}

// ------------ REDIS ------------
type RedisConfig struct {
	HOST string `mapstructure:"RD_HOST"`
	PORT string `mapstructure:"RD_PORT"`
	PASS string `mapstructure:"RD_PASS"`
	DB   int    `mapstructure:"RD_DB"`
}

// ------------ RABBITMQ ------------
type RabbitMQConfig struct {
	URI string `mapstructure:"URI"`
}

// ------------ MONGO DB ------------
type MongoDBConfig struct {
	DB  string `mapstructure:"DB"`
	URI string `mapstructure:"URI"`
}

// ------------ DATABASE ------------
type DatabaseConfig struct {
	HOST    string `mapstructure:"DBHOST"`
	PORT    int64  `mapstructure:"DBPORT"`
	USER    string `mapstructure:"DBUSER"`
	PASS    string `mapstructure:"DBPASS"`
	NAME    string `mapstructure:"DBNAME"`
	CHARSET string `mapstructure:"CHARSET"`
}

type DocsAuthConfig struct {
	LOGIN    string `mapstructure:"LOGIN"`
	PASSWORD string `mapstructure:"PASSWORD"`
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
