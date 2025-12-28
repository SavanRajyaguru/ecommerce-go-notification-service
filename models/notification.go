package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationStatus string
type ChannelType string

const (
	StatusSent   NotificationStatus = "SENT"
	StatusFailed NotificationStatus = "FAILED"

	ChannelEmail ChannelType = "EMAIL"
	ChannelSMS   ChannelType = "SMS"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EventType string             `bson:"event_type"`   // order.created, etc.
	Reference string             `bson:"reference_id"` // order_id or payment_id
	UserID    string             `bson:"user_id,omitempty"`
	Payload   interface{}        `bson:"payload"`
	Channel   ChannelType        `bson:"channel"`
	Status    NotificationStatus `bson:"status"`
	Error     string             `bson:"error,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
