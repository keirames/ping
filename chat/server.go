package main

import (
	"context"
	"fmt"
	"log"
	"main/config"
	"main/database"
	"main/graph"
	"main/internal/messages"
	"main/internal/rooms"
	"main/validator"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

const defaultPort = "8080"

func main() {
	validator.New()
	config.Load()
	database.Connect()

	port := config.C.Port

	router := chi.NewRouter()
	// router.Use(cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowCredentials: true,
	// 	Debug:            true,
	// }).Handler)

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
	c.Directives.Auth = func(
		ctx context.Context, obj interface{}, next graphql.Resolver,
	) (interface{}, error) {
		fmt.Print(ctx)
		return next(ctx)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
