package main

import (
	"fmt"
	"main/config"
	eventbus "main/event-bus"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := config.Load()
	if err != nil {
		panic("Load config fail")
	}

	_, err = eventbus.New()
	if err != nil {
		panic("Fail to register event topic")
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	fmt.Println("Listening port: " + config.C.PORT)
	http.ListenAndServe(":"+config.C.PORT, r)
}
