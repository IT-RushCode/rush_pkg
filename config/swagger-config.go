package config

// ------------ SWAGGER CONFIG ------------
type SwaggerConfig struct {
	BasePath         string `mapstructure:"BASE_PATH"`
	FilePath         string `mapstructure:"FILE_PATH"`
	Path             string `mapstructure:"PATH"`
	Title            string `mapstructure:"TITLE"`
	CacheAge         int    `mapstructure:"CACHE_AGE"`
	BasicAuthEnabled bool   `mapstructure:"BASIC_AUTH_ENABLED"`
	AuthUser         string `mapstructure:"AUTH_USER"`
	AuthPassword     string `mapstructure:"AUTH_PASSWORD"`
}
