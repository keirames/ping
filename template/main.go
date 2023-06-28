package main

import (
	"fmt"
	"main/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := config.Load()
	if err != nil {
		panic("Load config fail")
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	fmt.Println("Listening port: " + config.C.PORT)
	http.ListenAndServe(":"+config.C.PORT, r)
}
