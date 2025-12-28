package repository

import (
	"context"
	"time"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/internal/database"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/models"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	// Add other methods as needed, e.g., GetByReference, etc.
}

type notificationRepository struct{}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	_, err := database.Collection.InsertOne(ctx, notification)
	return err
}
