package ws

import "fmt"

type hub struct {
	subscribe   chan *client
	unsubscribe chan *client
	clients     []*client
	events      chan []byte
}

func New() *hub {
	return &hub{
		subscribe:   make(chan *client),
		unsubscribe: make(chan *client),
		clients:     []*client{},
		events:      make(chan []byte),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.subscribe:
			fmt.Println("new client", client)

		case client := <-h.unsubscribe:
			fmt.Println("client out", client)

		case event := <-h.events:
			fmt.Println(event)
		}
	}
}
