package api

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rhnauf/recipe-api/internal/repository"
	"net/http"
)

type api struct {
	recipeRepository repository.RecipeRepository
}

func NewAPI(pool *sql.DB) *api {

	recipeRepository := repository.NewRecipeRepository(pool)

	return &api{
		recipeRepository: recipeRepository,
	}
}

func (a *api) Server(port string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: a.Routes(),
	}
}

/*
	for simplicity, the architecture only consist of handler -> repository layer (skipping business logic layer)
*/

func (a *api) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Post("/recipe", a.insertRecipe)
	r.Put("/recipe/{id}", a.updateRecipe)
	r.Get("/recipe/{id}", a.getRecipeById)
	r.Delete("/recipe/{id}", a.deleteRecipeById)
	r.Get("/recipe-list", a.getListRecipe)

	return r
}
