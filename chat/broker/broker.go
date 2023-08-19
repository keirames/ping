package broker

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type TopicRoomMessage struct {
	UserID    int64
	RoomID    int64
	MessageID int64
}

var publishers []*kafka.Writer

func CreateConsumer(topic string, handler func(m kafka.Message)) {
	fmt.Println("create consumer for topic", topic)
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
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		handler(m)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func CreatePublisher(topic string) {
	fmt.Println("create publisher for topic", topic)
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	isExist := false
	for _, p := range publishers {
		if p.Topic == topic {
			isExist = true
			break
		}
	}

	if isExist {
		return
	} else {
		publishers = append(publishers, w)
	}
}

func GetPublisher(topic string) (*kafka.Writer, error) {
	for _, p := range publishers {
		if p.Topic == topic {
			return p, nil
		}
	}

	return nil, fmt.Errorf("Publisher not found")
}
