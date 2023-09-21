package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"main/config"
	"main/logger"
	"main/messaging"
	"main/pkg/router/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

			cmd := rdb.Get(context.Background(), msg.UserID)
			err := cmd.Err()
			if err == redis.Nil {
				// user not connected in any machine
				logger.L.Err(err).Send()
				return fmt.Errorf("user is offline")
			}
			if err != nil {
				logger.L.Err(err).Send()
				return fmt.Errorf("redis fail")
			}

			machineID := cmd.Val()
			url := machineID
			logger.L.Info().Msg(
				fmt.Sprintf("machine with url %v for user %v", url, msg.UserID),
			)

			body := []byte(`{
				"title": "Post title",
				"body": "Post description",
				"userId": 1
			}`)
			_, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
			if err != nil {
				logger.L.Err(err).Send()
				return errors.New("fail to contact desirer user")
			}

			return nil
		})

	r := chi.NewRouter()

	r.Post("/user-subscribe", func(w http.ResponseWriter, r *http.Request) {
		err := handlers.Subscribe(context.Background(), rdb, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/user-unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		err := handlers.Subscribe(context.Background(), rdb, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/ping", func(w http.ResponseWriter, r *http.Request) {
		vulnerableUserIDs, err := handlers.Ping(context.Background(), rdb, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		render.JSON(w, r, vulnerableUserIDs)
	})

	addr := fmt.Sprintf("%v:%v", config.C.Host, config.C.Port)
	fmt.Println("address", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
