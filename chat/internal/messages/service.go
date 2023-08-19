package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"main/broker"
	"main/customerror"
	"main/graph/model"
	"main/internal/auth"
	"main/internal/rooms"
	"main/keygen"
	"main/logger"
	"strconv"

	"github.com/segmentio/kafka-go"
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
	user, err := auth.GetUser(ctx)
	if err != nil {
		return nil, customerror.BadRequest()
	}
	userID := (*user).ID

	roomID, err := strconv.ParseInt(smi.RoomID, 10, 64)
	if err != nil {
		fmt.Println("Fail to parse int roomID")
		return nil, customerror.BadRequest()
	}

	ids, err := ms.roomsRepository.GetMembersIDs(ctx, roomID)
	if err != nil {
		fmt.Println("Fail to get MembersIDs")
		return nil, customerror.BadRequest()
	}

	isExist := false
	for _, id := range *ids {
		if id == userID {
			isExist = true
		}
	}
	if !isExist {
		fmt.Println("user is not from room")
		return nil, customerror.BadRequest()
	}

	// TODO: are we rlly need this ?
	isMember, err := ms.roomsRepository.IsMemberOfRoom(ctx, userID, roomID)
	if err != nil {
		fmt.Println("fail to exec query ismemberofroom")
		return nil, customerror.BadRequest()
	}
	if !isMember {
		fmt.Println("user is not from room")
		return nil, customerror.BadRequest()
	}

	msgID := keygen.Snowflake()
	id, err = ms.messagesRepository.CreateMessage(ctx, CreateMessageParams{
		ID:      msgID,
		Content: smi.Content,
		Type:    smi.Type.String(),
		UserID:  userID,
		RoomID:  roomID,
	})
	if err != nil {
		fmt.Println("create message fail", err)
		return nil, customerror.BadRequest()
	}

	p, err := broker.GetPublisher("room")
	if err != nil {
		logger.L.Err(err).Msg("Room topic publisher not exist")
	}
	if p != nil {
		messages := []kafka.Message{}

		for _, id := range *ids {
			if id != userID {
				mValue := broker.TopicRoomMessage{
					UserID:    id,
					RoomID:    roomID,
					MessageID: msgID,
				}
				v, err := json.Marshal(mValue)
				if err != nil {
					logger.L.Err(err).Msg("Cannot marshal json in room topic - send message")
					continue
				}

				messages = append(messages, kafka.Message{
					Value: []byte(v),
				})
				logger.L.Info().Msg("ping this user")
				logger.L.Info().Msg(string(id))
			}
		}

		err = p.WriteMessages(context.Background(), messages...)
		if err != nil {
			logger.L.Err(err).Msg("Cannot send messages into broker - send message")
		} else {
			logger.L.Info().Msg("Successfully write messages into broker")
		}
	}

	return id, nil
}
