package roomcontroller

import (
	roomservice "chatroom/room/service"

	"github.com/go-playground/validator/v10"
)

type roomController struct {
	validate    *validator.Validate
	roomService roomservice.RoomService
}

type Options struct {
	Validate    *validator.Validate
	RoomService roomservice.RoomService
}

func New(o *Options) *roomController {
	return &roomController{
		validate:    o.Validate,
		roomService: o.RoomService,
	}
}

func (rc *roomController) Rooms() {}

func (rc *roomController) Room() {}
