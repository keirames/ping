package service

import (
	"chatroom/keygen"
	"chatroom/logger"
	"chatroom/middlewares"
	"chatroom/room/model"
	"chatroom/room/repository"
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type roomService struct {
	rr   repository.RoomRepository
	psql sq.StatementBuilderType
	conn *sqlx.DB
}

func New(rr repository.RoomRepository, psql sq.StatementBuilderType, c *sqlx.DB) *roomService {
	return &roomService{
		rr:   rr,
		psql: psql,
		conn: c,
	}
}

func (rs *roomService) Rooms(
	ctx context.Context,
	page int,
	limit int,
) (*model.PaginateRoomsRes, error) {
	userID := middlewares.GetUserID(ctx)

	offset := uint64((page - 1) * limit)

	sql, args, err :=
		rs.psql.
			Select("cr.*").
			From("chat_rooms cr").
			InnerJoin(
				"users_and_chat_rooms uacr ON uacr.room_id = cr.id",
			).
			Where("uacr.user_id = $1", userID).
			Limit(uint64(limit)).
			Offset(offset).
			ToSql()

	fmt.Println(sql, args)

	if err != nil {
		fmt.Println("Fail to create sql")
		return nil, err
	}

	type chatRoom struct {
		ID   string
		Name string
	}

	rooms := []chatRoom{}

	err = rs.conn.Select(&rooms, sql, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rr := []model.RoomsRes{}
	for _, room := range rooms {
		rr = append(rr, model.RoomsRes(room))
	}

	return &model.PaginateRoomsRes{Page: page, Limit: limit, Data: rr}, nil
}

func (rs *roomService) JoinRoom(
	ctx context.Context,
	roomID int64,
) (*model.JoinRoomRes, error) {
	userID := middlewares.GetUserID(ctx)

	isRoomExist, err := rs.rr.IsRoomExist(roomID)
	if err != nil || !isRoomExist {
		return nil, err
	}

	isJoined, err := rs.rr.IsMemberOfRoom(userID, roomID)
	if err != nil {
		return nil, err
	}
	if isJoined {
		logger.L.Info().Msg("user is already inside room")
		return nil, fmt.Errorf("user is already inside room")
	}

	q, args, err :=
		rs.psql.
			Insert("users_and_chat_rooms").
			Columns("id", "user_id", "room_id").
			Values(keygen.Snowflake(), userID, roomID).
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to create sql")
		return nil, err
	}

	_, err = rs.conn.Exec(q, args...)
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to join room")
		return nil, err
	}

	return &model.JoinRoomRes{ID: strconv.FormatInt(roomID, 10)}, nil
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

func (rs *roomService) CreateRoom(
	name string,
	memberIDs []string,
) (*model.RoomsRes, error) {
	deDupIDs := uniq(memberIDs)

	sql, args, err :=
		rs.psql.
			Insert("chat_rooms").
			Columns("name").
			Values(name).
			Suffix("RETURNING id").
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Invalid sql query")
		return nil, err
	}

	tx, err := rs.conn.Beginx()
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

	sqlBuilder := rs.psql.Insert("users_and_chat_rooms").Columns("user_id", "room_id")
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

	return &model.RoomsRes{ID: roomID, Name: name}, nil
}