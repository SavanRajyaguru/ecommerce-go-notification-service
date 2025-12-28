package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/internal/notification"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/models"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/repository"
)

type NotificationProcessor struct {
	repo        repository.NotificationRepository
	emailSender *notification.EmailSender
	smsSender   *notification.SMSSender
}

func NewNotificationProcessor(repo repository.NotificationRepository, email *notification.EmailSender, sms *notification.SMSSender) *NotificationProcessor {
	return &NotificationProcessor{
		repo:        repo,
		emailSender: email,
		smsSender:   sms,
	}
}

// Event Payload Structures (simplified)
type OrderCreatedEvent struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	UserEmail string  `json:"user_email"` // Assuming event has email, or we fetch user details
}

type PaymentEvent struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

func (p *NotificationProcessor) ProcessEvent(ctx context.Context, topic string, payload []byte) error {
	log.Printf("Processing event from topic: %s", topic)

	// Idempotency check could be here (check if notification for this ref exists)

	var err error
	switch topic {
	case "order.created":
		err = p.handleOrderCreated(ctx, payload)
	case "payment.success":
		err = p.handlePaymentSuccess(ctx, payload)
	case "order.cancelled":
		// err = p.handleOrderCancelled(ctx, payload)
		log.Println("Order Cancelled event received (TODO)")
	case "payment.failed":
		// err = p.handlePaymentFailed(ctx, payload)
		log.Println("Payment Failed event received (TODO)")
	default:
		log.Printf("Unknown topic: %s", topic)
	}

	return err
}

func (p *NotificationProcessor) handleOrderCreated(ctx context.Context, payload []byte) error {
	var event OrderCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	// In a real system, we might need to fetch User Email from User Service if not in event.
	// For now, assume event has it or use a default.
	emailTo := event.UserEmail
	if emailTo == "" {
		emailTo = "customer@example.com" // Fallback
	}

	// 1. Send Email
	subject := fmt.Sprintf(notification.TemplateOrderCreatedSubject, event.OrderID)
	body := fmt.Sprintf(notification.TemplateOrderCreatedBody, event.OrderID)
	err := p.emailSender.Send(emailTo, subject, body)

	status := models.StatusSent
	errorMsg := ""
	if err != nil {
		status = models.StatusFailed
		errorMsg = err.Error()
		log.Printf("Failed to send email: %v", err)
	}

	// 2. Log to DB
	notif := &models.Notification{
		EventType: "order.created",
		Reference: event.OrderID,
		UserID:    event.UserID,
		Payload:   event,
		Channel:   models.ChannelEmail,
		Status:    status,
		Error:     errorMsg,
	}
	return p.repo.Create(ctx, notif)
}

func (p *NotificationProcessor) handlePaymentSuccess(ctx context.Context, payload []byte) error {
	var event PaymentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	// Mock Email
	subject := fmt.Sprintf(notification.TemplatePaymentSuccessSubject, event.OrderID)
	body := fmt.Sprintf(notification.TemplatePaymentSuccessBody, event.OrderID, "PAID") // Amount missing in simple struct

	err := p.emailSender.Send("customer@example.com", subject, body)

	status := models.StatusSent
	errorMsg := ""
	if err != nil {
		status = models.StatusFailed
		errorMsg = err.Error()
	}

	notif := &models.Notification{
		EventType: "payment.success",
		Reference: event.OrderID, // or PaymentID
		Payload:   event,
		Channel:   models.ChannelEmail,
		Status:    status,
		Error:     errorMsg,
	}
	return p.repo.Create(ctx, notif)
}
