package notification

import "context"

type NotificationService interface {
	Send(ctx context.Context, id int64, cityName string) error
}
