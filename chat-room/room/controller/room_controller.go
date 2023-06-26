package roomcontroller

import (
	"chatroom/logger"
	roommodel "chatroom/room/model"
	roomservice "chatroom/room/service"
	"encoding/json"
	"net/http"

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

func (rc *roomController) JoinRoom(r *http.Request) (
	res *roommodel.JoinRoomRes,
	statusCode int,
	err error,
) {
	var jrr roommodel.JoinRoomReq
	err = json.NewDecoder(r.Body).Decode(&jrr)
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to decode")
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusBadRequest, err
}

func (rc *roomController) Room() {}
