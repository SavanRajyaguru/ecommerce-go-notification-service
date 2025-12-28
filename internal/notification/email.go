package notification

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/config"
	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	dialer *gomail.Dialer
}

func NewEmailSender() *EmailSender {
	cfg := config.AppConfig.SMTP
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} // For dev/test simplifiction
	return &EmailSender{dialer: dialer}
}

func (s *EmailSender) Send(to string, subject, body string) error {
	if config.AppConfig.FeatureFlags["enable_email"] == false {
		log.Println("Email feature disabled, skipping email to:", to)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.SMTP.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	log.Printf("Email sent to %s | Subject: %s", to, subject)
	return nil
}
