package main

import (
	"context"
	"encoding/json"
	"log"
	"main/broker"
	"main/config"
	"main/database"
	"main/graph"
	"main/internal/auth"
	"main/internal/messages"
	"main/internal/rooms"
	"main/keygen"
	"main/logger"
	"main/validator"
	"main/ws"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/segmentio/kafka-go"
)

const defaultPort = "8080"

func main() {
	// TODO: these init func is very dangerous, maybe we will forget to call it
	validator.New()
	config.Load()
	database.Connect()
	keygen.New()
	logger.New()

	var s any
	hub := ws.New(s)
	go hub.Run()

	type RoomTopicMessage struct {
		RoomID    int64 `json:"roomId"`
		UserID    int64 `json:"userId"`
		MessageID int64 `json:"messageId"`
	}

	type NewMessage struct {
		RoomID    int64 `json:"roomId"`
		MessageID int64 `json:"messageId"`
	}

	go broker.CreateConsumer("room", func(m kafka.Message) {
		data := RoomTopicMessage{}
		err := json.Unmarshal(m.Value, &data)
		if err != nil {
			logger.L.Err(err).Msg("Fail to parse data from topic 'room'")
			return
		}

		nm := NewMessage{
			RoomID:    data.RoomID,
			MessageID: data.MessageID,
		}

		dataToClient, err := json.Marshal(&nm)
		if err != nil {
			logger.L.Err(err).Msg("fail to marshal data")
			return
		}

		hub.SendMessageToClient(data.UserID, dataToClient)
	})
	go broker.CreatePublisher("room")

	port := config.C.Port

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(auth.Middleware())

	roomsRepository := rooms.NewRoomsRepository()
	messagesRepository := messages.NewMessagesRepository()
	roomsService := rooms.NewRoomsService(&rooms.NewRoomsServiceParams{
		RR: roomsRepository,
	})
	messagesService := messages.NewMessagesService(&messages.NewMessagesServiceParams{
		MessagesRepository: messagesRepository,
		RoomsRepository:    roomsRepository,
	})

	c := graph.Config{Resolvers: &graph.Resolver{
		RoomsService:    roomsService,
		MessagesService: messagesService,
	}}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Handle("/socket", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ws.Serve(context.Background(), hub, w, r)
		if err != nil {
			http.Error(w, "BadRequest", http.StatusBadRequest)
		}
	}))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
