package config

// ------------ API SMS ------------
type APISMS struct {
	API    string `mapstructure:"API"`
	SENDER string `mapstructure:"SENDER"`
}
