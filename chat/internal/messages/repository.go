package messages

import (
	"context"
	"fmt"
	"main/database"
	"main/logger"
	"main/query"

	"github.com/jackc/pgx/v5/pgtype"
)

type MessagesRepository interface {
	CreateMessage(
		ctx context.Context,
		arg CreateMessageParams,
	) (id *int64, err error)
	GetMessage(
		ctx context.Context, id int64,
	) (*query.Message, error)
}

type messagesRepository struct {
}

func NewMessagesRepository() *messagesRepository {
	return &messagesRepository{}
}

type CreateMessageParams struct {
	ID      int64
	Content string
	Type    string
	UserID  int64
	RoomID  int64
}

func (mr *messagesRepository) CreateMessage(
	ctx context.Context,
	arg CreateMessageParams,
) (id *int64, err error) {
	m, err := database.Queries.CreateMessage(ctx, query.CreateMessageParams{
		ID:      arg.ID,
		Content: arg.Content,
		Type:    pgtype.Text{String: arg.Type, Valid: true},
		RoomID:  arg.RoomID,
		UserID:  arg.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &m.ID, nil
}

func (mr *messagesRepository) GetMessage(
	ctx context.Context, id int64,
) (*query.Message, error) {
	m, err := database.Queries.GetMessage(ctx, id)
	if err != nil {
		logger.L.Err(err).Msg("fail to exec query")
		return nil, fmt.Errorf("fail to get message")
	}

	return &m, nil
}
