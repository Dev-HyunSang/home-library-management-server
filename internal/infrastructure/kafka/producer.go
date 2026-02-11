package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dev-hyunsang/home-library-backend/logger"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

type NotificationEvent struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Type   string `json:"type"` // e.g., "review", "system"
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

func (p *Producer) ProduceNotification(ctx context.Context, userID, title, body, eventType string) error {
	event := NotificationEvent{
		UserID: userID,
		Title:  title,
		Body:   body,
		Type:   eventType,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(userID),
		Value: payload,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		logger.Init().Sugar().Errorf("Failed to write message to kafka: %v", err)
		return err
	}

	logger.Init().Sugar().Infof("Message produced to kafka: %s", string(payload))
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
