package messageservice

import (
	"chatroom/logger"
	messagemodel "chatroom/message/model"
	messagerepository "chatroom/message/repository"
	"chatroom/middlewares"
	"chatroom/room/repository"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type messageService struct {
	messageRepository messagerepository.MessageRepository
	roomRepository    repository.RoomRepository
	psql              squirrel.StatementBuilderType
	conn              *sqlx.DB
}

func New(
	mr messagerepository.MessageRepository,
	rr repository.RoomRepository,
	p squirrel.StatementBuilderType,
	c *sqlx.DB,
) *messageService {
	return &messageService{
		messageRepository: mr,
		roomRepository:    rr,
		psql:              p,
		conn:              c,
	}
}

func (ms *messageService) Messages(
	ctx context.Context,
	page int,
	roomID int64,
	limit int,
) (*[]messagemodel.MessageEntity, error) {
	userID := middlewares.GetUserID(ctx)
	isExist, err := ms.roomRepository.IsRoomExist(roomID)
	if err != nil {
		logger.L.Error().Err(err).Msg("check room exist error")
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf("room is not exist: %v", roomID)
	}

	isMember, err := ms.roomRepository.IsMemberOfRoom(userID, roomID)
	if err != nil {
		logger.L.Error().Err(err).Msg("check member error")
		return nil, err
	}
	if !isMember {
		logger.L.Info().Msg("user is not a member")
		return nil, fmt.Errorf("user is not a member: %v", userID)
	}

	messages, err := ms.messageRepository.FindByRoomID(roomID, page)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
