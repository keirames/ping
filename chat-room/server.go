package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"chatroom/config"
	"chatroom/db"
	"chatroom/jwt"
	"chatroom/keygen"
	"chatroom/logger"
	messagemodel "chatroom/message/model"
	messagerepository "chatroom/message/repository"
	messageservice "chatroom/message/service"
	"chatroom/middlewares"
	"chatroom/room/repository"
	"chatroom/room/service"
	"chatroom/ws"

	"github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)

	if config.C.ENV == "DEV" {
		logger.L.Info().Msg("Allow all cors in DEV mode")
		r.Use(cors.New(
			cors.Options{
				AllowedOrigins: []string{"localhost:3000"},
				AllowedMethods: []string{
					http.MethodHead,
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
					http.MethodOptions,
				},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
				Debug:            true,
			},
		).Handler)
	}

	r.Group(func(r chi.Router) {
		type signInReq struct {
			ID string `json:"id" validate:"required"`
		}
		r.Post("/v1/sign-in", func(w http.ResponseWriter, r *http.Request) {
			var sir signInReq
			err := json.NewDecoder(r.Body).Decode(&sir)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logger.L.Error().Err(err).Msg("Bad request sign-in")
				return
			}

			validate := validator.New()
			err = validate.Struct(sir)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logger.L.Error().Err(err).Msg("Bad request sign-in validate fail")
				return
			}

			id, err := strconv.ParseInt(sir.ID, 10, 64)
			if err != nil {
				logger.L.Error().Err(err).Msg("sign-in parse fail")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			q, args, err := db.Psql.Select("1 as flag").
				From("users u").
				Where(squirrel.Eq{"u.id": id}).
				ToSql()
			if err != nil {
				logger.L.Error().Err(err).Msg("Fail to prepare query")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var flag int
			err = db.Conn.Get(&flag, q, args...)
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err != nil {
				logger.L.Error().Err(err).Msg("sql exec error")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			jwt, err := jwt.GenerateJwt(
				context.Background(),
				id,
			)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logger.L.Error().Err(err).Msg("Bad request sign-in fail to gen jwt")
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "x-token",
				Value:    *jwt,
				Expires:  time.Now().Add(time.Hour * 2400),
				Secure:   false,
				HttpOnly: true,
			})
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, *jwt)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		roomRepository := repository.New(db.Psql, db.Conn)
		roomService := service.New(roomRepository, db.Psql, db.Conn)

		messageRepository := messagerepository.New(db.Psql, db.Conn)
		messageService := messageservice.New(messageRepository, roomRepository, db.Psql, db.Conn)

		type joinRoomReq struct {
			RoomID string `json:"roomId" validate:"required"`
		}

		hub := ws.New(roomService)
		go hub.Run()

		r.Get("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
			ws.Serve(hub, w, r)
		})

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

		r.Get("/v1/messages", func(w http.ResponseWriter, r *http.Request) {
			page, err := strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				logger.L.Error().Err(err).Msg("Invalid params")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			roomID, err := strconv.ParseInt(r.URL.Query().Get("roomId"), 10, 64)
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

			messages, err := messageService.Messages(r.Context(), page, roomID, 10)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			res := messagemodel.MapMessagesEntityToModel(*messages)
			render.JSON(w, r, res)
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

		type sendMessageRes struct {
			Text   string `json:"text" validate:"required,max=255"`
			RoomID string `json:"roomId" validate:"required"`
		}

		r.Post("/v1/send-message", func(w http.ResponseWriter, r *http.Request) {
			userID := middlewares.GetUserID(r.Context())

			validate := validator.New()

			var smr sendMessageRes

			err = json.NewDecoder(r.Body).Decode(&smr)
			if err != nil {
				logger.L.Error().Err(err).Msg("[API send-message] Cannot decode request body")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := validate.Struct(smr)
			if err != nil {
				logger.L.Error().Err(err).Msg("[API send-message] Invalidate request body")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			roomID, err := strconv.ParseInt(smr.RoomID, 10, 64)
			if err != nil {
				logger.L.Error().Err(err).Msg("[API send-message] Invalid roomID")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			room, err := roomService.SendMessage(userID, smr.Text, roomID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			render.Status(r, http.StatusCreated)
			render.JSON(w, r, room)
		})
	})

	http.ListenAndServe(":8080", r)
}
