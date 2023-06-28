package roomservice

import (
	"chatroom/arrayutils"
	"chatroom/keygen"
	"chatroom/logger"
	"chatroom/middlewares"
	roommodel "chatroom/room/model"
	roomrepository "chatroom/room/repository"
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type RoomService interface {
	Rooms(ctx context.Context, page int, limit int) (*roommodel.PaginateRoomsRes, error)
	JoinRoom(ctx context.Context, roomID int64) (*roommodel.JoinRoomRes, error)
	CreateRoom(name string, memberIDs []string) (*roommodel.RoomsRes, error)
	SendMessage(userID int64, text string, roomID int64) (*roommodel.SendMessageRes, error)
}

type roomService struct {
	rr   roomrepository.RoomRepository
	psql sq.StatementBuilderType
	conn *sqlx.DB
}

func New(
	rr roomrepository.RoomRepository,
	psql sq.StatementBuilderType,
	c *sqlx.DB,
) *roomService {
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
) (*roommodel.PaginateRoomsRes, error) {
	userID, _ := middlewares.GetUserID(ctx)

	offset := uint64((page - 1) * limit)

	sql, args, err :=
		rs.psql.
			Select("cr.id as id, cr.name as name").
			From("chat_rooms cr").
			InnerJoin(
				"users_and_chat_rooms uacr ON uacr.room_id = cr.id",
			).
			Where("uacr.user_id = $1", userID).
			Limit(uint64(limit)).
			Offset(offset).
			ToSql()

	logger.L.Info().Msg("[rooms] - " + sql)

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

	rr := []roommodel.RoomsRes{}
	for _, room := range rooms {
		rr = append(rr, roommodel.RoomsRes(room))
	}

	return &roommodel.PaginateRoomsRes{Page: page, Limit: limit, Data: rr}, nil
}

func (rs *roomService) JoinRoom(
	ctx context.Context,
	roomID int64,
) (*roommodel.JoinRoomRes, error) {
	userID, _ := middlewares.GetUserID(ctx)

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

	return &roommodel.JoinRoomRes{ID: strconv.FormatInt(roomID, 10)}, nil
}

func (rs *roomService) CreateRoom(
	name string,
	memberIDs []string,
) (*roommodel.RoomsRes, error) {
	deDupIDs := arrayutils.Uniq(memberIDs)

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

	return &roommodel.RoomsRes{ID: roomID, Name: name}, nil
}

func (rs *roomService) SendMessage(
	userID int64,
	text string,
	roomID int64,
) (*roommodel.SendMessageRes, error) {
	isExist, err := rs.rr.IsRoomExist(roomID)
	if err != nil || !isExist {
		return nil, err
	}

	isMember, err := rs.rr.IsMemberOfRoom(userID, roomID)
	if err != nil || !isMember {
		return nil, err
	}

	id, err := rs.rr.SendMessage(text, userID, roomID)
	if err != nil {
		return nil, err
	}

	return &roommodel.SendMessageRes{ID: strconv.FormatInt(*id, 10)}, nil
}
