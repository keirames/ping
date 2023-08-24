package ws

import (
	"fmt"
	"main/logger"
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
			logger.L.Info().Msg(fmt.Sprintf("user subscribe: %v", client.id))
			h.clients = append(h.clients, client)

		case client := <-h.unsubscribe:
			logger.L.Info().Msg(fmt.Sprintf("user unsubscribe: %v", client.id))

		case <-h.events:
			// err := eventsHandler(data, h.service)
			// if err != nil {
			// 	logger.L.Error().Err(err).Msg("Fail to handler events")
			// }
			fmt.Println("got event")
		}
	}
}

func (h *hub) SendMessageToClient(clientID int64, message []byte) {
	fmt.Println("send message to client: ", clientID, h.clients)
	for _, c := range h.clients {
		fmt.Println(c.id)
		if c.id == clientID {
			fmt.Println("found client ", clientID, "subscribe to server")
			c.send <- message
		}
	}
}
