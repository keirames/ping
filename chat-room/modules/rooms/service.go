package rooms

import (
	"chatroom/db"
	"chatroom/keygen"
	"chatroom/logger"
	"chatroom/middlewares"
	"context"
	"fmt"
	"strconv"
)

type JoinRoomRes struct {
	ID string `json:"id"`
}

func JoinRoom(ctx context.Context, roomID int64) (*JoinRoomRes, error) {
	userID := middlewares.GetUserID(ctx)

	isRoomExist, err := IsRoomExist(roomID)
	if err != nil || !isRoomExist {
		return nil, err
	}

	isJoined, err := IsMemberOfRoom(userID, roomID)
	if err != nil {
		return nil, err
	}
	if isJoined {
		logger.L.Info().Msg("user is already inside room")
		return nil, fmt.Errorf("user is already inside room")
	}

	q, args, err :=
		db.Psql.
			Insert("users_and_chat_rooms").
			Columns("id", "user_id", "room_id").
			Values(keygen.Snowflake(), userID, roomID).
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to create sql")
		return nil, err
	}

	_, err = db.Conn.Exec(q, args...)
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to join room")
		return nil, err
	}

	return &JoinRoomRes{ID: strconv.FormatInt(roomID, 10)}, nil
}
