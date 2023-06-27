package messageservice

import (
	"chatroom/logger"
	messagemodel "chatroom/message/model"
	messagerepository "chatroom/message/repository"
	roomrepository "chatroom/room/repository"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type MessageService interface {
	Messages(
		userID int64,
		page int,
		roomID int64,
		limit int,
	) (*[]messagemodel.MessageEntity, error)
}

type messageService struct {
	messageRepository messagerepository.MessageRepository
	roomRepository    roomrepository.RoomRepository
	psql              squirrel.StatementBuilderType
	conn              *sqlx.DB
}

func New(
	mr messagerepository.MessageRepository,
	rr roomrepository.RoomRepository,
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
	userID int64,
	page int,
	roomID int64,
	limit int,
) (*[]messagemodel.MessageEntity, error) {
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
