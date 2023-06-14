package ws

import (
	"chatroom/logger"
	"chatroom/room/model"
	"fmt"
)

type message struct {
	UserID int64
	Data   string
}

type hub struct {
	subscribe   chan *client
	unsubscribe chan *client
	clients     []*client
	events      chan *message
	service     Service
}

type Service interface {
	SendMessage(
		userID int64,
		text string,
		roomID int64,
	) (*model.SendMessageRes, error)
}

func New(s Service) *hub {
	return &hub{
		subscribe:   make(chan *client),
		unsubscribe: make(chan *client),
		clients:     []*client{},
		events:      make(chan *message, 1),
		service:     s,
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.subscribe:
			fmt.Println("new client", client)

		case client := <-h.unsubscribe:
			fmt.Println("client out", client)

		case data := <-h.events:
			err := eventsHandler(data, h.service)
			if err != nil {
				logger.L.Error().Err(err).Msg("Fail to handler events")
			}
		}
	}
}
