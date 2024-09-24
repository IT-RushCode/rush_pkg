package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/IT-RushCode/rush_pkg/config"
)

// SendEmail отправляет email с заданным субъектом и телом на указанный адрес
//
// to - email получателя
//
// subject - заголовок письма
//
// body - текст письма
func SendEmail(cfg *config.MailConfig, to, subject, body string) error {
	message := fmt.Sprintf(
		"From: %s <%s>\nTo: %s\nSubject: %s\n\n%s",
		cfg.SenderName,
		cfg.SenderEmail,
		to,
		subject,
		body,
	)

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort))
	if err != nil {
		return err
	}
	defer client.Close()

	// Начинаем TLS соединение
	if err = client.StartTLS(&tls.Config{ServerName: cfg.SMTPHost}); err != nil {
		return err
	}

	// Аутентификация
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Установка адреса отправителя
	if err = client.Mail(cfg.SenderEmail); err != nil {
		return err
	}

	// Установка адреса получателя
	if err = client.Rcpt(to); err != nil {
		return err
	}

	// Отправка письма
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
