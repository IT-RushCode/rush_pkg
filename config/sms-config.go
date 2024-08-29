package config

// ------------ API SMS ------------
type SMS struct {
	URL      string `mapstructure:"URL"`
	METHOD   string `mapstructure:"METHOD"`
	TOKEN    string `mapstructure:"TOKEN"`
	SENDER   string `mapstructure:"SENDER"`
	TEMPLATE string `mapstructure:"TEMPLATE"`
}
