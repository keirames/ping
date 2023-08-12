package messages

import (
	"context"
	"main/customerror"
	"main/graph/model"
	"main/internal/rooms"
	"main/keygen"
	"strconv"
)

type MessagesService interface {
	SendMessage(
		ctx context.Context,
		smi model.SendMessageInput,
	) (id *int64, err error)
}

type messagesService struct {
	messagesRepository MessagesRepository
	roomsRepository    rooms.RoomsRepository
}

type NewMessagesServiceParams struct {
	MessagesRepository *messagesRepository
	rooms.RoomsRepository
}

func NewMessagesService(p *NewMessagesServiceParams) *messagesService {
	return &messagesService{
		messagesRepository: p.MessagesRepository,
		roomsRepository:    p.RoomsRepository,
	}
}

func (ms *messagesService) SendMessage(
	ctx context.Context,
	smi model.SendMessageInput,
) (id *int64, err error) {
	var userID int64

	roomID, err := strconv.ParseInt(smi.RoomID, 64, 10)
	if err != nil {
		return nil, customerror.BadRequest()
	}

	isExist, err := ms.roomsRepository.IsRoomExist(ctx, roomID)
	if err != nil {
		return nil, customerror.BadRequest()
	}
	if !isExist {
		return nil, customerror.BadRequest()
	}

	isMember, err := ms.roomsRepository.IsMemberOfRoom(ctx, userID, roomID)
	if err != nil {
		return nil, customerror.BadRequest()
	}
	if !isMember {
		return nil, customerror.BadRequest()
	}

	id, err = ms.messagesRepository.CreateMessage(ctx, CreateMessageParams{
		ID: keygen.Snowflake(),
	})
	if err != nil {
		return nil, customerror.BadRequest()
	}

	return id, nil
}
