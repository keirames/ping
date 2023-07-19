package rooms

import (
	"context"
	"fmt"
	"main/customerror"
	"main/query"
)

type RoomsService interface {
	Rooms(
		ctx context.Context, userID int64, page int,
	) (*[]query.ChatRoom, error)
}

type roomsService struct {
	roomsRepository
}

type NewRoomsServiceParams struct {
	RR *roomsRepository
}

func NewRoomsService(p *NewRoomsServiceParams) *roomsService {
	return &roomsService{
		roomsRepository: *p.RR,
	}
}

func (rs *roomsService) Rooms(
	ctx context.Context, userID int64, page int,
) (*[]query.ChatRoom, error) {
	rooms, err := rs.roomsRepository.GetRooms(ctx, userID, page)
	if err != nil {
		fmt.Println(err)
		return nil, customerror.BadRequest()
	}

	return rooms, nil
}
