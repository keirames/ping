package handlers

import (
	"context"
	"encoding/json"
	"main/logger"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type subscribeInput struct {
	UserID    string `json:"userId"`
	MachineID string `json:"machineId"`
}

func Subscribe(ctx context.Context, rdb *redis.Client, r *http.Request) error {
	var si subscribeInput

	if err := json.NewDecoder(r.Body).Decode(&si); err != nil {
		logger.L.Err(err).Msg("fail to decode subscribe input")
		return err
	}

	if err := rdb.Set(ctx, si.UserID, si.MachineID, 0).Err(); err != nil {
		logger.L.Err(err).Msg("write fail to redis")
		return err
	}

	return nil
}
