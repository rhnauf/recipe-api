package repository

import "database/sql"

type recipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *recipeRepository {
	return &recipeRepository{db: db}
}
