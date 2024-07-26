package repository

import (
	"database/sql"
	"github.com/rhnauf/recipe-api/internal/entity"
)

type recipeRepository struct {
	db *sql.DB
}

type RecipeRepository interface {
	InsertRecipe(recipe entity.Recipe) error
	UpdateRecipe(recipe entity.Recipe) error
}

func NewRecipeRepository(db *sql.DB) *recipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) InsertRecipe(recipe entity.Recipe) error {
	_, err := r.db.Exec(`
		INSERT INTO recipes(title, description, instruction, publish)
		VALUES ($1, $2, $3, $4)`,
		recipe.Title,
		recipe.Description,
		recipe.Instruction,
		*recipe.Publish,
	)

	return err
}
func (r *recipeRepository) UpdateRecipe(recipe entity.Recipe) error {
	_, err := r.db.Exec(`
		UPDATE recipes
		SET title = $1, description = $2, instruction = $3, publish = $4
		WHERE id = $5`,
		recipe.Title,
		recipe.Description,
		recipe.Instruction,
		*recipe.Publish,
		recipe.Id,
	)

	return err
}
