package ws

import (
	"chatroom/logger"
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

type Hub interface {
	Run()
	SendMessageToClient(clientID int64, message int64)
}

type Service interface {
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
			h.clients = append(h.clients, client)

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

func (h *hub) SendMessageToClient(clientID int64, message int64) {
	fmt.Println("send message to client: ", clientID, h.clients)
	for _, c := range h.clients {
		fmt.Println(c.id)
		if c.id == clientID {
			fmt.Println("found client ", clientID, "subscribe to server")
			c.send <- []byte(string(message))
		}
	}
}
