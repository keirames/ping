package messaging

import (
	"context"
	"fmt"
	"main/logger"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Register()
}

type consumer struct {
}

func CreateConsumer() *consumer {
	return &consumer{}
}

func (c *consumer) Register(
	ctx context.Context, topic string, handler func(*kafka.Message) error,
) {
	logger.L.Info().Msg("create consumer for topic " + topic)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9094"},
		Topic:    topic,
		GroupID:  "consumer-group-router",
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			logger.L.Err(err).Msg("fail to fetch msg from topic")
			break
		}
		fmt.Printf(
			"message at offset %d: %s = %s\n",
			m.Offset, string(m.Key),
			string(m.Value),
		)

		if err := handler(&m); err != nil {
			// Skip commit msg when execute fail.
			logger.L.Err(err).Msg("skip commit msg because handler execute fail")
			continue
		}

		fmt.Println("commit msg")
		if err := r.CommitMessages(ctx, m); err != nil {
			logger.L.Err(err).Msg("failed to commit messages")
			break
		}
	}

	if err := r.Close(); err != nil {
		logger.L.Err(err).Msg("failed to close reader")
	}

	panic("unable to connect to broker")
}
