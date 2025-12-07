package config

// TusConfig описывает настройки для tus-загрузчика.
type TusConfig struct {
	Enabled   bool   `mapstructure:"ENABLED"`
	BasePath  string `mapstructure:"BASE_PATH"`
	UploadDir string `mapstructure:"UPLOAD_DIR"`
	FinalDir  string `mapstructure:"FINAL_DIR"`
	LockDir   string `mapstructure:"LOCK_DIR"`
	MaxSize   int64  `mapstructure:"MAX_SIZE"`
}
