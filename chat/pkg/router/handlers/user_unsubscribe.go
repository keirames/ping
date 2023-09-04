package handlers

import (
	"context"
	"encoding/json"
	"main/logger"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type unsubscribeInput struct {
	UserID    string `json:"userId"`
	MachineID string `json:"machineId"`
}

func Unsubscribe(ctx context.Context, rdb *redis.Client, r *http.Request) error {
	var ui unsubscribeInput

	if err := json.NewDecoder(r.Body).Decode(&ui); err != nil {
		logger.L.Err(err).Msg("fail to decode unsubscribe input")
		return err
	}

	if err := rdb.Del(ctx, ui.UserID).Err(); err != nil {
		logger.L.Err(err).Msg("delete fail")
		return err
	}

	return nil
}
