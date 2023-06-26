package roomrepository

import (
	"chatroom/db"
	"chatroom/keygen"
	"chatroom/logger"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type roomsRepository struct {
	Psql sq.StatementBuilderType
	Conn *sqlx.DB
}

type RoomRepository interface {
	IsRoomExist(roomID int64) (bool, error)
	IsMemberOfRoom(userID int64, roomID int64) (bool, error)
	SendMessage(text string, userID int64, roomID int64) (*int64, error)
}

func New(psql sq.StatementBuilderType, conn *sqlx.DB) *roomsRepository {
	return &roomsRepository{Psql: psql, Conn: conn}
}

func (rr *roomsRepository) IsRoomExist(roomID int64) (bool, error) {
	q, args, err :=
		rr.Psql.
			Select("id").
			From("chat_rooms as cr").
			Where(sq.Eq{"cr.id": roomID}).
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to create sql")
		return false, err
	}

	var r int64
	err = db.Conn.Get(&r, q, args...)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to query")
		return false, err
	}

	return true, nil
}

func (rr *roomsRepository) IsMemberOfRoom(userID int64, roomID int64) (bool, error) {
	q, args, err :=
		db.Psql.
			Select("id").
			From("users_and_chat_rooms as uacr").
			Where(
				sq.And{
					sq.Eq{"uacr.user_id": userID},
					sq.Eq{"uacr.room_id": roomID},
				}).
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to create sql")
		return false, err
	}

	var r int64
	err = db.Conn.Get(&r, q, args...)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to query")
		return false, err
	}

	return true, nil
}

func (rr *roomsRepository) SendMessage(
	text string,
	userID int64,
	roomID int64,
) (*int64, error) {
	q, args, err :=
		rr.Psql.
			Insert("messages").
			Columns("id", "content", "user_id", "room_id").
			Values(keygen.Snowflake(), text, userID, roomID).
			Suffix("RETURNING id").
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("Cannot create query")
		return nil, err
	}

	var msgID int64
	row := rr.Conn.QueryRowx(q, args...)
	err = row.Scan(&msgID)
	if err != nil {
		logger.L.Error().Err(err).Msg("Cannot send message")
		return nil, err
	}

	return &msgID, nil
}
