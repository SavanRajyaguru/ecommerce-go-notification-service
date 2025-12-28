package notification

import (
	"log"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/config"
)

type SMSSender struct{}

func NewSMSSender() *SMSSender {
	return &SMSSender{}
}

func (s *SMSSender) Send(to string, message string) error {
	if config.AppConfig.FeatureFlags["enable_sms"] == false {
		log.Println("SMS feature disabled, skipping SMS to:", to)
		return nil
	}

	// Mock implementation
	log.Printf("[MOCK SMS] Sent to %s: %s", to, message)
	return nil
}
