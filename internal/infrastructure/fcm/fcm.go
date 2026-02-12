package fcm

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/dev-hyunsang/home-library-backend/logger"
	"google.golang.org/api/option"
)

type FCMService struct {
	client *messaging.Client
}

func NewFCMService(serviceAccountPath string) (*FCMService, error) {
	opt := option.WithCredentialsFile(serviceAccountPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &FCMService{
		client: client,
	}, nil
}

func (s *FCMService) SendPush(ctx context.Context, token, title, body string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := s.client.Send(ctx, message)
	if err != nil {
		logger.Sugar().Errorf("Failed to send FCM message: %v", err)
		return err
	}

	logger.Sugar().Infof("Successfully sent message: %s", response)
	return nil
}
