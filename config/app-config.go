package config

// ------------ APP CONFIG ------------
type AppConfig struct {
	ENV           string `mapstructure:"ENV"`
	DEBUG         bool   `mapstructure:"DEBUG"`
	DB_DEBUG      bool   `mapstructure:"DB_DEBUG"`
	NAME          string `mapstructure:"NAME"`
	HOST          string `mapstructure:"HOST"`
	PORT          string `mapstructure:"PORT"`
	CACHE_ACTIVE  bool   `mapstructure:"CACHE_ACTIVE"`
	CACHE_TIMEOUT int64  `mapstructure:"CACHE_TIMEOUT"`
}
