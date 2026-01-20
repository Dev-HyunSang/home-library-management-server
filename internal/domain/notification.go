package domain

import "context"

type NotificationProducer interface {
	ProduceNotification(ctx context.Context, userID, title, body, eventType string) error
}
