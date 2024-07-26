package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rhnauf/recipe-api/internal/entity"
	"testing"
)

func assertErr(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestInsertRecipe(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var b = true
	recipe := entity.Recipe{
		Title:       "nasi goreng",
		Description: "desc nasi goreng",
		Instruction: "instruction nasi goreng",
		Publish:     &b,
	}

	repo := NewRecipeRepository(db)

	qry := "INSERT INTO recipes(title, description, instruction, publish) VALUES ($1, $2, $3, $4)"

	t.Run("should return success on insert query", func(t *testing.T) {
		mock.
			ExpectExec(qry).
			WithArgs(recipe.Title, recipe.Description, recipe.Instruction, recipe.Publish).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.InsertRecipe(recipe)
		assertErr(t, err, nil)
	})

	t.Run("should return error on insert query", func(t *testing.T) {
		mock.
			ExpectExec(qry).
			WithArgs(recipe.Title, recipe.Description, recipe.Instruction, recipe.Publish).
			WillReturnError(sql.ErrConnDone)

		err = repo.InsertRecipe(recipe)
		assertErr(t, err, sql.ErrConnDone)
	})
}

func TestUpdateRecipe(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var b = true
	recipe := entity.Recipe{
		Id:          1,
		Title:       "nasi goreng",
		Description: "desc nasi goreng",
		Instruction: "instruction nasi goreng",
		Publish:     &b,
	}

	repo := NewRecipeRepository(db)

	qry := "UPDATE recipes SET title = $1, description = $2, instruction = $3, publish = $4 WHERE id = $5"

	t.Run("should return success on update query", func(t *testing.T) {
		mock.
			ExpectExec(qry).
			WithArgs(recipe.Title, recipe.Description, recipe.Instruction, recipe.Publish, recipe.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.UpdateRecipe(recipe)
		assertErr(t, err, nil)
	})

	t.Run("should return error on update query", func(t *testing.T) {
		mock.
			ExpectExec(qry).
			WithArgs(recipe.Title, recipe.Description, recipe.Instruction, recipe.Publish, recipe.Id).
			WillReturnError(sql.ErrConnDone)

		err = repo.UpdateRecipe(recipe)
		assertErr(t, err, sql.ErrConnDone)
	})
}
