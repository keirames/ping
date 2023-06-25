package messagerepository

import (
	"chatroom/logger"
	messagemodel "chatroom/message/model"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	IsMessageExist(id int64, roomID int64) (bool, error)
	FindByRoomID(id int64, page int) (*[]messagemodel.MessageEntity, error)
}

type messageRepository struct {
	Psql squirrel.StatementBuilderType
	Conn *sqlx.DB
}

func New(p squirrel.StatementBuilderType, c *sqlx.DB) *messageRepository {
	return &messageRepository{
		Psql: p, Conn: c,
	}
}

func (mr *messageRepository) IsMessageExist(id int64, roomID int64) (bool, error) {
	q, args, err :=
		mr.Psql.Select("1 as flag").
			From("messages m").
			Where(squirrel.And{
				squirrel.Eq{"m.room_id": roomID},
				squirrel.Eq{"m.id": id},
			}).
			ToSql()

	fmt.Println(q, args, err)

	var result int
	err = mr.Conn.Get(&result, q, args...)
	if err != nil && err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logger.L.Error().Err(err).Msg("Query give error" + q)
		return false, err
	}

	return true, nil
}

func (mr *messageRepository) FindByRoomID(
	roomID int64,
	page int,
) (*[]messagemodel.MessageEntity, error) {
	offset := uint64((page - 1) * 10)
	q, args, err :=
		mr.Psql.Select("*").
			From("messages m").
			Where(squirrel.Eq{"m.room_id": roomID}).
			OrderBy("m.created_at ASC").
			Limit(10).
			Offset(offset).
			ToSql()
	if err != nil {
		logger.L.Error().Err(err).Msg("fail to prepare query")
		return nil, err
	}

	var messages []messagemodel.MessageEntity
	err = mr.Conn.Select(&messages, q, args...)
	if err != nil {
		logger.L.Error().Err(err).Msg("fail to exec query")
		return nil, err
	}

	return &messages, nil
}
