package rooms

import (
	"chatroom/db"
	"chatroom/keygen"
	"chatroom/logger"
	"chatroom/middlewares"
	roomsmodel "chatroom/modules/rooms/model"
	"context"
	"fmt"
	"strconv"
)

func JoinRoom(ctx context.Context, roomID int64) (*roomsmodel.JoinRoomRes, error) {
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

	return &roomsmodel.JoinRoomRes{ID: strconv.FormatInt(roomID, 10)}, nil
}

func Rooms(ctx context.Context, page int, limit int) (*roomsmodel.PaginateRoomsRes, error) {
	userID := middlewares.GetUserID(ctx)

	offset := uint64((page - 1) * limit)

	sql, args, err :=
		db.Psql.
			Select("cr.*").
			From("chat_rooms cr").
			InnerJoin(
				"users_and_chat_rooms uacr ON uacr.room_id = cr.id",
			).
			Where("uacr.user_id = $1", userID).
			Limit(uint64(limit)).
			Offset(offset).
			ToSql()

	if err != nil {
		fmt.Println("Fail to create sql")
		return nil, err
	}

	type chatRoom struct {
		ID   string
		Name string
	}

	rooms := []chatRoom{}

	err = db.Conn.Select(&rooms, sql, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rr := []roomsmodel.RoomsRes{}
	for _, room := range rooms {
		rr = append(rr, roomsmodel.RoomsRes(room))
	}

	return &roomsmodel.PaginateRoomsRes{Page: 10, Limit: 10, Data: rr}, nil
}
