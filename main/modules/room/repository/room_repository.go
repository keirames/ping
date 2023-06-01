package repository

import (
	"main/domain"
	"main/tools"

	"github.com/jmoiron/sqlx"
)

type roomRepository struct {
	Conn *sqlx.DB
}

func NewRoomRepository(conn *sqlx.DB) *roomRepository {
	return &roomRepository{
		Conn: conn,
	}
}

func (r *roomRepository) FindByID(id int64) (*domain.Room, error) {
	var room domain.Room

	err := r.Conn.Get(&room, "SELECT * FROM rooms WHERE rooms.id=?", id)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) FindAll() (*[]domain.Room, error) {
	rooms := []domain.Room{}

	err := r.Conn.Select(&rooms, "SELECT * FROM rooms WHERE rooms.id = 1")
	if err != nil {
		return nil, err
	}

	return &rooms, nil
}

func (r *roomRepository) CreateRoom(name string, members []int64) (*domain.Room, error) {
	roomID := tools.Snowflake.Generate()

	tx, err := r.Conn.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO rooms(id, name) VALUES($1, $2)", roomID, name)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	for _, userID := range members {
		_, err = tx.Exec("INSERT INTO users_rooms(user_id, room_id) VALUES($1, $2)", userID, roomID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.Room{ID: roomID.Int64(), Name: name}, nil
}
