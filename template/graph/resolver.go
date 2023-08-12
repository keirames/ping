package graph

import (
	"main/internal/messages"
	"main/internal/rooms"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	rooms.RoomsService
	messages.MessagesService
}
