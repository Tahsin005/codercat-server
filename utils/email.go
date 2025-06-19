package utils

import (
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
}

func SendEmail(config EmailConfig, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)

	msg := []byte("Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" + body + "\r\n")

	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	return smtp.SendMail(addr, auth, config.From, to, msg)
}
