package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"main/logger"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type pingRequest struct {
	MachineID string   `json:"machineId"`
	UserIDs   []string `json:"userIds"`
}

func Ping(ctx context.Context, rdb *redis.Client, r *http.Request) (*[]string, error) {
	var pr pingRequest

	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		logger.L.Err(err).Send()
		return nil, errors.New("Fail to decode")
	}

	var vulnerableUserIDs []string
	for _, userID := range pr.UserIDs {
		err := rdb.Set(ctx, pr.MachineID, userID, 0).Err()
		if err != nil {
			vulnerableUserIDs = append(vulnerableUserIDs, userID)
			logger.L.Err(err).Send()
			continue
		}
	}

	return &vulnerableUserIDs, nil
}
