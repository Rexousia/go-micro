package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {

	//creating a new mux
	mux := chi.NewRouter()

	//use the following origins, methods, and headers.
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//checking to see if connection is alive
	mux.Use(middleware.Heartbeat("/ping"))

	//at this path serve app.Authenticate
	mux.Post("/authenticate", app.Authenticate)
	return mux
}
