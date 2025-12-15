package events

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/hawful70/platform-events/pkg/events"
	"github.com/hawful70/shop-identity-service/internal/identity"
)

type KafkaNotifier struct {
	writer *kafka.Writer
}

func NewKafkaNotifier(brokers []string, topic string) *KafkaNotifier {
	return &KafkaNotifier{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			RequiredAcks: kafka.RequireAll,
			Balancer:     &kafka.Hash{},
		},
	}
}

func (n *KafkaNotifier) Close() error {
	if n == nil || n.writer == nil {
		return nil
	}
	return n.writer.Close()
}

func (n *KafkaNotifier) UserCreated(ctx context.Context, user identity.User) error {
	evt := events.NewUserCreated(string(user.ID), user.Email, user.Username)
	payload, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(strings.ToLower(user.Email)),
		Value: payload,
		Time:  time.Now().UTC(),
	}

	return n.writer.WriteMessages(ctx, msg)
}

var _ identity.UserNotifier = (*KafkaNotifier)(nil)
