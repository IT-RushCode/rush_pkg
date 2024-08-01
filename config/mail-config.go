package config

// ------------ MAIL CONFIG ------------
type MailConfig struct {
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SenderEmail  string `mapstructure:"SENDER_EMAIL"`
	SenderName   string `mapstructure:"SENDER_NAME"`
}
