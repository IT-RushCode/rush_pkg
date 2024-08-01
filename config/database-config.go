package config

// ------------ DATABASE ------------
type DatabaseConfig struct {
	HOST    string `mapstructure:"DBHOST"`
	PORT    int64  `mapstructure:"DBPORT"`
	USER    string `mapstructure:"DBUSER"`
	PASS    string `mapstructure:"DBPASS"`
	NAME    string `mapstructure:"DBNAME"`
	CHARSET string `mapstructure:"CHARSET"`
}
