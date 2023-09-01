package bus

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

func (c *consumer) Register(topic string, handler func(*kafka.Message)) {
	logger.L.Info().Msg("create consumer for topic " + topic)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf(
			"message at offset %d: %s = %s\n",
			m.Offset, string(m.Key),
			string(m.Value),
		)
		if err := r.CommitMessages(context.Background(), m); err != nil {
			logger.L.Err(err).Msg("fail to commit message")
			continue
		}

		handler(&m)
	}

	if err := r.Close(); err != nil {
		logger.L.Err(err).Msg("failed to close reader")
	}
}
