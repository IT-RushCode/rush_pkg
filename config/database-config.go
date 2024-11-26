package config

// ------------ DATABASE ------------
type DatabaseConfig struct {
	HOST    string `mapstructure:"HOST"`
	PORT    int64  `mapstructure:"PORT"`
	USER    string `mapstructure:"USER"`
	PASS    string `mapstructure:"PASS"`
	NAME    string `mapstructure:"NAME"`
	CHARSET string `mapstructure:"CHARSET"`
}
