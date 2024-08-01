package config

// ------------ JWT ------------
type JwtConfig struct {
	JWT_SECRET  string `mapstructure:"JWT_SECRET"`
	JWT_TTL     int64  `mapstructure:"JWT_TTL"`
	REFRESH_TTL int64  `mapstructure:"REFRESH_TTL"`
}
