package event

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func Subscribe(topic string, handler func(m kafka.Message)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
	})
	fmt.Println("Subscribe to " + topic)

	for {
		m, err := r.ReadMessage(context.Background())
		fmt.Println("Got message ", string(m.Value))
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		handler(m)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
