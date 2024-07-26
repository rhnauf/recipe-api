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
	GetRecipeById(id int64) (*entity.Recipe, error)
	DeleteRecipeById(id int64) error
	GetListRecipe(limit, offset int64) ([]*entity.Recipe, error)
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

func (r *recipeRepository) GetRecipeById(id int64) (*entity.Recipe, error) {
	var recipe entity.Recipe

	err := r.db.QueryRow("SELECT * FROM recipes WHERE id = $1", id).
		Scan(&recipe.Id, &recipe.CreatedAt, &recipe.Title, &recipe.Description, &recipe.Instruction, &recipe.Publish)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (r *recipeRepository) DeleteRecipeById(id int64) error {
	_, err := r.db.Exec("DELETE FROM recipes WHERE id = $1", id)
	if err != nil {
		return err
	}

	/*
		could get rows affected from result to get specific err no rows, if rows affected == 0
		for simplicity just return success regardless the rows affected
	*/

	return nil
}

func (r *recipeRepository) GetListRecipe(limit, offset int64) ([]*entity.Recipe, error) {
	var recipes []*entity.Recipe

	rows, err := r.db.Query("SELECT id, title FROM recipes LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var recipe entity.Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title); err != nil {
			continue
		}
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}
