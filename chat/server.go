package main

import (
	"fmt"
	"log"
	"main/broker"
	"main/config"
	"main/database"
	"main/graph"
	"main/internal/auth"
	"main/internal/messages"
	"main/internal/rooms"
	"main/keygen"
	"main/validator"
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

	go broker.CreateConsumer("room", func(m kafka.Message) {
		fmt.Println("consumer topic room run & need to find user inside socket data")
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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
