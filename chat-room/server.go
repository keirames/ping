package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"chatroom/config"
	"chatroom/db"
	"chatroom/keygen"
	"chatroom/logger"
	"chatroom/middlewares"
	"chatroom/room/repository"
	"chatroom/room/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func main() {
	logger.New()
	keygen.New()

	err := config.Load()
	if err != nil {
		panic("cannot load config file!")
	}

	err = db.Connect()
	if err != nil {
		fmt.Println(err)
		logger.L.Error().Err(err).Msg("DB connection error")
		panic("cannot connect to the database")
	}
	logger.L.Info().Msg("DB connected")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middlewares.Auth)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	roomRepository := repository.New(db.Psql, db.Conn)
	roomService := service.New(roomRepository, db.Psql, db.Conn)

	type joinRoomReq struct {
		RoomID string `json:"roomId" validate:"required"`
	}

	r.Post("/v1/join-room", func(w http.ResponseWriter, r *http.Request) {
		var jrr joinRoomReq
		err := json.NewDecoder(r.Body).Decode(&jrr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.L.Error().Err(err).Msg("Fail to decode")
			return
		}

		validate := validator.New()
		err = validate.Struct(jrr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.L.Error().Err(err).Msg("Fail to validate")
			return
		}

		roomID, err := strconv.ParseInt(jrr.RoomID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.L.Error().Err(err).Msg("Cannot parse into uint")
			return
		}

		res, err := roomService.JoinRoom(r.Context(), roomID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		render.JSON(w, r, res)
	})

	r.Get("/v1/rooms", func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			logger.L.Error().Err(err).Msg("Invalid params")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err = validate.Var(page, "gt=0")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rooms, err := roomService.Rooms(r.Context(), page, 10)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// w.WriteHeader(http.StatusOK)
		render.JSON(w, r, rooms)
	})

	type createRoomRequest struct {
		Name      string   `json:"name" validate:"required,max=50"`
		MemberIDs []string `json:"memberIds" validate:"required,min=1"`
	}

	r.Post("/v1/create-room", func(w http.ResponseWriter, r *http.Request) {

		validate := validator.New()

		var createRoomRequest createRoomRequest

		err = json.NewDecoder(r.Body).Decode(&createRoomRequest)
		if err != nil {
			logger.L.Error().Err(err).Msg("[API create-room] Cannot decode request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := validate.Struct(createRoomRequest)
		if err != nil {
			logger.L.Error().Err(err).Msg("[API create-room] Invalidate request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		room, err := roomService.CreateRoom(createRoomRequest.Name, createRoomRequest.MemberIDs)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, room)
	})

	http.ListenAndServe(":3000", r)
}
