package main

import (
	"context"
	"encoding/json"
	"fmt"
	"main/config"
	"main/logger"
	"main/messaging"
	"main/pkg/router/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	logger.New()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.C.RedisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go messaging.CreateConsumer().Register(
		context.Background(), "room", func(m *kafka.Message) error {
			type KafkaMessage struct {
				UserID string `json:"userId"`
			}

			var msg KafkaMessage
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				fmt.Println("fail to unmarshal msg from broker", err)
				return fmt.Errorf("fail to unmarshal msg")
			}

			fmt.Println("got u", msg)

			cmd := rdb.Get(context.Background(), msg.UserID)
			err := cmd.Err()
			if err == redis.Nil {
				// user not connected in any machine
				logger.L.Err(err).Msg("user offline")
				return fmt.Errorf("user is offline")
			}
			if err != nil {
				logger.L.Err(err).Msg("redis fail")
				return fmt.Errorf("redis fail")
			}

			machineID := cmd.Val()
			// TODO: call desirer endpoint
			// ???.???.???/receive-message
			fmt.Println(machineID, "machineID")

			return fmt.Errorf("fail")
		})

	r := chi.NewRouter()

	r.Post("/user-subscribe", func(w http.ResponseWriter, r *http.Request) {
		err := handlers.Subscribe(context.Background(), rdb, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Post("/user-unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		err := handlers.Subscribe(context.Background(), rdb, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	addr := fmt.Sprintf("%v:%v", config.C.Host, config.C.Port)
	fmt.Println("address", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
