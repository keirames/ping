package eventbus

import (
	"context"
	"fmt"
	"main/config"

	"github.com/segmentio/kafka-go"
)

var Conn *kafka.Conn

func New() (*kafka.Conn, error) {
	topic := "subscribe-session"
	partition := 5

	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		config.C.KAFKA_HOST,
		topic,
		partition,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("Subscribe to kafka topic: " + topic)
	Conn = conn
	return conn, nil
}
