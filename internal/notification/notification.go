package notification

import "context"

type NotificationService interface {
	Send(ctx context.Context) error
}
