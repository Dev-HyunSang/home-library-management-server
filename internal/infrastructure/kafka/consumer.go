package kafka

import (
	"context"
	"encoding/json"

	"github.com/dev-hyunsang/home-library/internal/infrastructure/fcm"
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader     *kafka.Reader
	fcmService *fcm.FCMService
}

func NewConsumer(brokers []string, topic, groupID string, fcmService *fcm.FCMService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:     reader,
		fcmService: fcmService,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			logger.Init().Sugar().Errorf("Failed to read message from kafka: %v", err)
			break
		}

		logger.Init().Sugar().Infof("Message received from kafka: %s", string(m.Value))

		var event NotificationEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			logger.Init().Sugar().Errorf("Failed to unmarshal message: %v", err)
			continue
		}

		// TODO: 실제로는 DB에서 해당 UserID의 FCM 토큰을 조회해야 합니다.
		// 지금은 예시로 event.UserID를 토큰으로 가정하거나, 테스트용 토큰을 사용하거나 해야 합니다.
		// 여기서는 로직의 흐름만 구현합니다.

		// 예: userToken, err := userRepo.GetFCMToken(event.UserID)
		// if err != nil { ... }

		// For demo, logging instead of actual sending if token logic is missing
		if c.fcmService != nil {
			// token := "DEVICE_TOKEN_HERE"
			// err := c.fcmService.SendPush(ctx, token, event.Title, event.Body)
			logger.Init().Sugar().Infof("Processing notification for user %s: %s - %s", event.UserID, event.Title, event.Body)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
