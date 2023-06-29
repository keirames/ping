package api

import (
	"chatroom/db"
	"chatroom/logger"
	"chatroom/middlewares"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type RoomsResponse struct {
	ID   string
	Name string
}

type PaginateRoomsRes struct {
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Data  []RoomsResponse `json:"data"`
}

func Rooms(ctx context.Context, page int, limit int) (*PaginateRoomsRes, error) {
	userID, _ := middlewares.GetUserID(ctx)

	offset := uint64((page - 1) * limit)

	sql, args, err :=
		psql.
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

	rr := []RoomsResponse{}
	for _, room := range rooms {
		rr = append(rr, RoomsResponse(room))
	}

	return &PaginateRoomsRes{Page: 10, Limit: 10, Data: rr}, nil
}

func uniq(arr []string) []string {
	deDupMap := make(map[string]bool)
	result := []string{}

	for _, i := range arr {
		if !deDupMap[i] {
			continue
		}

		deDupMap[i] = true
		result = append(result, i)
	}

	return result
}

func CreateRoom(name string, memberIDs []string) (*RoomsResponse, error) {
	deDupIDs := uniq(memberIDs)

	sql, args, err :=
		psql.
			Insert("chat_rooms").
			Columns("name").
			Values(name).
			Suffix("RETURNING id").
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Invalid sql query")
		return nil, err
	}

	tx, err := db.Conn.Beginx()
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to open transaction")
		return nil, err
	}

	row := tx.QueryRowx(sql, args...)

	var roomID string
	err = row.Scan(&roomID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	sqlBuilder := psql.Insert("users_and_chat_rooms").Columns("user_id", "room_id")
	for _, id := range deDupIDs {
		sqlBuilder.Values(id, roomID)
	}

	sql, args, err = sqlBuilder.ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Invalid sql query")
		tx.Rollback()
		return nil, err
	}

	row = tx.QueryRowx(sql, args...)
	var id string
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to commit a transaction")
		tx.Rollback()
		return nil, err
	}

	return &RoomsResponse{ID: roomID, Name: name}, nil
}
