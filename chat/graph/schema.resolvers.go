package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"
	"main/customerror"
	"main/graph/model"
	"main/internal/auth"
	"main/validator"
	"strconv"
	"time"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, sendMessageInput model.SendMessageInput) (*model.Message, error) {
	type ValidateParams struct {
		Content string `validator:"required"`
		RoomID  string `validator:"required"`
	}

	err := validator.V.Struct(&ValidateParams{
		Content: sendMessageInput.Content,
		RoomID:  sendMessageInput.RoomID,
	})
	if err != nil {
		return nil, customerror.BadRequest()
	}

	id, err := r.MessagesService.SendMessage(ctx, sendMessageInput)
	if err != nil {
		return nil, err
	}

	// TODO: wrong return value
	return &model.Message{
		ID:        strconv.FormatInt(*id, 10),
		Content:   sendMessageInput.Content,
		Type:      sendMessageInput.Type,
		IsDelete:  false,
		ParentID:  nil,
		CreatedAt: time.Now().String(),
		UserID:    strconv.FormatInt(*id, 10),
		RoomID:    "11",
	}, nil
}

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, signInInput *model.SignInInput) (string, error) {
	panic(fmt.Errorf("not implemented: SignIn - signIn"))
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}

// Rooms is the resolver for the rooms field.
func (r *queryResolver) Rooms(ctx context.Context, page int) (*model.PagedRooms, error) {
	user, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("got user", *user)

	rooms, err := r.RoomsService.Rooms(ctx, 67309604664680448, page)
	if err != nil {
		return nil, err
	}

	var roomsModel []*model.Room
	for _, r := range *rooms {
		roomsModel = append(roomsModel, &model.Room{
			ID:   strconv.Itoa(int(r.ID)),
			Name: r.Name,
		})
	}

	return &model.PagedRooms{
		Page:  page,
		Items: roomsModel,
	}, nil
}

// Room is the resolver for the room field.
func (r *queryResolver) Room(ctx context.Context, id string) (*model.Room, error) {
	panic(fmt.Errorf("not implemented: Room - room"))
}

// Messages is the resolver for the messages field.
func (r *queryResolver) Messages(ctx context.Context, messagesInput model.MessagesInput) (*model.PagedMessages, error) {
	panic(fmt.Errorf("not implemented: Messages - messages"))
}

// Message is the resolver for the message field.
func (r *queryResolver) Message(ctx context.Context, getMessageInput *model.GetMessageInput) (*model.Message, error) {
	uc, err := auth.GetUser(ctx)
	if err != nil {
		return nil, customerror.BadRequest()
	}

	messageID, err := strconv.ParseInt(getMessageInput.MessageID, 10, 64)
	if err != nil {
		return nil, customerror.BadRequest()
	}

	roomID, err := strconv.ParseInt(getMessageInput.RoomID, 10, 64)
	if err != nil {
		return nil, customerror.BadRequest()
	}

	m, err := r.MessagesService.GetMessage(ctx, messageID, uc.ID, roomID)
	if err != nil {
		return nil, customerror.BadRequest()
	}

	parentID := strconv.FormatInt(m.ParentID.Int64, 10)

	result := &model.Message{
		ID:        getMessageInput.MessageID,
		Content:   m.Content,
		Type:      model.MessageType(m.Type.String),
		IsDelete:  m.IsDelete.Bool,
		ParentID:  &parentID,
		CreatedAt: m.CreatedAt.Time.String(),
		UserID:    strconv.FormatInt(m.UserID, 10),
		RoomID:    strconv.FormatInt(m.RoomID, 10),
	}

	return result, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
