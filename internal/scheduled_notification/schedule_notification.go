package scheduled_notification

import "time"

const SCHEDULED_NOTIFICATION_TABLE = "scheduled_notification"

type ScheduledNotification struct {
	Id               int64              `json:"id" db:"id"`
	Status           NotificationStatus `json:"status" db:"status"`
	Date             time.Time          `json:"date" db:"date"`
	CityName         string             `json:"city_name" db:"city_name"`
	UserId           int64              `json:"user_id" db:"user_id"`
	NotificationType NotificationType   `json:"notification_type" db:"notification_type"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`
}

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)

type NotificationType string

const (
	TypeWebhook NotificationType = "webhook"
	TypeEmail   NotificationType = "email"
	TypeSMS     NotificationType = "sms"
	TypePush    NotificationType = "push"
)
