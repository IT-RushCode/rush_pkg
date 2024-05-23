package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ------------ GLOBAL CONFIG ------------
type Config struct {
	APP      AppConfig      `mapstructure:"app"`
	DB       DBConfig       `mapstructure:"db"`
	JWT      JwtConfig      `mapstructure:"jwt"`
	DATETIME DateTimeConfig `mapstructure:"datetime"`
	REDIS    RedisConfig    `mapstructure:"redis"`
	KAFKA    KafKaConfig    `mapstructure:"kafka"`
}

// ------------ SERVICES ------------
type AppConfig struct {
	ENV            string `mapstructure:"env"`
	DEBUG          bool   `mapstructure:"debug"`
	NAME           string `mapstructure:"name"`
	HOST           string `mapstructure:"host"`
	PORT           string `mapstructure:"port"`
	MAX_CONNECTION string `mapstructure:"max_connection"`
}


// ------------ JWT ------------
type JwtConfig struct {
	JWT_SECRET  string `mapstructure:"jwt_secret"`
	JWT_TTL     int64  `mapstructure:"jwt_ttl"`
	REFRESH_TTL int64  `mapstructure:"refresh_ttl"`
}

// ------------ KAFKA ------------
type KafKaConfig struct {
	HOST1 string `mapstructure:"host1"`
	HOST2 string `mapstructure:"host2"`
	HOST3 string `mapstructure:"host3"`
}

// ------------ DATETIME ------------
type DateTimeConfig struct {
	Datetime string `mapstructure:"datetime"`
	Date     string `mapstructure:"date"`
	Time     string `mapstructure:"time"`
}

// ------------ REDIS ------------
type RedisConfig struct {
	HOST string `mapstructure:"host"`
	PORT string `mapstructure:"port"`
	PASS string `mapstructure:"password"`
	DB   int    `mapstructure:"db"`
}

// ------------ MONGO DB ------------
type MongoDBConfig struct {
	URI string `mapstructure:"uri"`
}

// ------------ DATABASE ------------
type DatabaseConfig struct {
	Host    string `mapstructure:"dbhost"`
	Port    int64  `mapstructure:"dbport"`
	User    string `mapstructure:"dbuser"`
	Pass    string `mapstructure:"dbpass"`
	Name    string `mapstructure:"dbname"`
	CHARSET string `mapstructure:"charset"`
}

type DBConfig struct {
	MYSQL DatabaseConfig `mapstructure:"MYSQL"`
	PSQL  DatabaseConfig `mapstructure:"PSQL"`
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
