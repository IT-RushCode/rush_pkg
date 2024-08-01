package config

// ------------ MONGO DB ------------
type MongoDBConfig struct {
	DB  string `mapstructure:"DB"`
	URI string `mapstructure:"URI"`
}
