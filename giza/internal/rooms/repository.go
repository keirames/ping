package rooms

import (
	"context"
	"main/database"
	"main/logger"
	"main/query"

	"github.com/jackc/pgx/v5"
)

type RoomsRepository interface {
	IsRoomExist(ctx context.Context, id int64) (bool, error)
	IsMemberOfRoom(
		ctx context.Context, userID int64, roomID int64,
	) (bool, error)
	GetRooms(
		ctx context.Context,
		userID int64,
		page int,
	) (*[]query.ChatRoom, error)
}

type roomsRepository struct{}

func NewRoomsRepository() *roomsRepository {
	return &roomsRepository{}
}

func (rr *roomsRepository) GetRooms(
	ctx context.Context,
	userID int64,
	page int,
) (*[]query.ChatRoom, error) {
	rooms, err := database.Queries.GetRooms(
		ctx, query.GetRoomsParams{UserID: userID, Offset: int32(page * 10)},
	)
	if err != nil {
		return nil, err
	}

	return &rooms, nil
}

func (rr *roomsRepository) IsRoomExist(ctx context.Context, id int64) (bool, error) {
	_, err := database.Queries.IsRoomExist(ctx, id)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logger.ExecQueryError(err)
		return false, err
	}

	return true, nil
}

func (rr *roomsRepository) IsMemberOfRoom(
	ctx context.Context, userID int64, roomID int64,
) (bool, error) {
	_, err := database.Queries.IsMemberOfRoom(
		ctx, query.IsMemberOfRoomParams{
			UserID: userID, RoomID: roomID,
		},
	)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logger.ExecQueryError(err)
		return false, err
	}

	return true, nil
}
