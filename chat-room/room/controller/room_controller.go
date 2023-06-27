package roomcontroller

import (
	"chatroom/logger"
	roommodel "chatroom/room/model"
	roomservice "chatroom/room/service"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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
	*roommodel.JoinRoomRes,
	int,
	error,
) {
	var jrr roommodel.JoinRoomReq
	err := json.NewDecoder(r.Body).Decode(&jrr)
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to decode")
		return nil, http.StatusBadRequest, err
	}

	err = rc.validate.Struct(jrr)
	if err != nil {
		logger.L.Error().Err(err).Msg("Fail to validate")
		return nil, http.StatusBadRequest, err
	}

	roomID, err := strconv.ParseInt(jrr.RoomID, 10, 64)
	if err != nil {
		logger.L.Error().Err(err).Msg("Cannot parse into uint")
		return nil, http.StatusBadRequest, err
	}

	result, err := rc.roomService.JoinRoom(context.Background(), roomID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return result, http.StatusOK, err
}

func (rc *roomController) Rooms(r *http.Request) (
	*roommodel.PaginateRoomsRes,
	int,
	error,
) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		logger.L.Error().Err(err).Msg("Invalid params")
		return nil, http.StatusBadRequest, err
	}

	err = rc.validate.Var(page, "gt=0")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	rooms, err := rc.roomService.Rooms(r.Context(), page, 10)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return rooms, http.StatusOK, nil
}
