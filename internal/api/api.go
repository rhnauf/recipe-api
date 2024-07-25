package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type api struct {
}

func NewAPI() *api {
	return &api{}
}

func (a *api) Server(port string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return r
}
