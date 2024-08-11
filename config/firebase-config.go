package config

// ------------ FIREBASE CONFIG ------------
type FirebaseConfig struct {
	PROJECT_ID                  string `mapstructure:"PROJECT_ID"`
	PRIVATE_KEY_ID              string `mapstructure:"PRIVATE_KEY_ID"`
	PRIVATE_KEY                 string `mapstructure:"PRIVATE_KEY"`
	CLIENT_EMAIL                string `mapstructure:"CLIENT_EMAIL"`
	CLIENT_ID                   string `mapstructure:"CLIENT_ID"`
	AUTH_URI                    string `mapstructure:"AUTH_URI"`
	TOKEN_URI                   string `mapstructure:"TOKEN_URI"`
	AUTH_PROVIDER_X509_CERT_URL string `mapstructure:"AUTH_PROVIDER_X509_CERT_URL"`
	CLIENT_X509_CERT_URL        string `mapstructure:"CLIENT_X509_CERT_URL"`
}
