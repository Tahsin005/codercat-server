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

// SendEmail sends a plain text email
func SendEmail(config EmailConfig, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)

	msg := []byte("Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" + body + "\r\n")

	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	return smtp.SendMail(addr, auth, config.From, to, msg)
}

// SendHTMLEmail sends an HTML email
func SendHTMLEmail(config EmailConfig, to []string, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)

	msg := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + htmlBody + "\r\n")

	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	return smtp.SendMail(addr, auth, config.From, to, msg)
}
