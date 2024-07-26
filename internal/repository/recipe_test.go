package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rhnauf/recipe-api/internal/entity"
	"reflect"
	"testing"
	"time"
)

func assertErr(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertRecipeEqual(t *testing.T, got, want *entity.Recipe) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertRecipesEqual(t *testing.T, got, want []*entity.Recipe) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// should've used suite
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

func TestGetRecipeById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var idSuccess int64 = 1
	var idNotFound int64 = 2

	repo := NewRecipeRepository(db)

	qry := "SELECT * FROM recipes WHERE id = $1"

	t.Run("should return success on get by id query", func(t *testing.T) {
		now := time.Now()

		recipeRow := sqlmock.
			NewRows([]string{"id", "created_at", "title", "description", "instruction", "publish"}).
			AddRow(1, now, "nasi goreng", "nasi goreng desc", "nasi goreng instruction", true)

		mock.
			ExpectQuery(qry).
			WithArgs(idSuccess).
			WillReturnRows(recipeRow)

		got, err := repo.GetRecipeById(idSuccess)

		var p = true
		want := &entity.Recipe{
			Id:          1,
			Title:       "nasi goreng",
			Description: "nasi goreng desc",
			Instruction: "nasi goreng instruction",
			Publish:     &p,
			CreatedAt:   now,
		}

		assertErr(t, err, nil)
		assertRecipeEqual(t, got, want)
	})

	t.Run("should return error not found on get by id query", func(t *testing.T) {
		mock.
			ExpectQuery(qry).
			WithArgs(idNotFound).
			WillReturnError(sql.ErrNoRows)

		got, err := repo.GetRecipeById(idNotFound)

		assertErr(t, err, sql.ErrNoRows)
		assertRecipeEqual(t, got, nil)
	})
}

func TestDeleteRecipeById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id int64 = 1

	repo := NewRecipeRepository(db)

	qry := "DELETE FROM recipes WHERE id = $1"

	t.Run("should return success on delete by id query", func(t *testing.T) {
		mock.
			ExpectExec(qry).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = repo.DeleteRecipeById(id)

		assertErr(t, err, nil)
	})

	t.Run("should return error on delete by id query", func(t *testing.T) {
		mock.
			ExpectExec(qry).
			WithArgs(id).
			WillReturnError(sql.ErrConnDone)

		err = repo.DeleteRecipeById(id)

		assertErr(t, err, sql.ErrConnDone)
	})
}

func TestGetListRecipe(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var offset int64 = 1
	var limit int64 = 10

	repo := NewRecipeRepository(db)

	qry := "SELECT id, title FROM recipes LIMIT $1 OFFSET $2"

	t.Run("should return success on get list query", func(t *testing.T) {
		recipeRows := sqlmock.
			NewRows([]string{"id", "title"}).
			AddRow(1, "nasi goreng")

		mock.
			ExpectQuery(qry).
			WithArgs(limit, offset).
			WillReturnRows(recipeRows)

		got, err := repo.GetListRecipe(limit, offset)

		want := []*entity.Recipe{
			{
				Id:    1,
				Title: "nasi goreng",
			},
		}

		assertErr(t, err, nil)
		assertRecipesEqual(t, got, want)
	})

	t.Run("should return empty slice error on scanning rows", func(t *testing.T) {
		recipeRows := sqlmock.
			NewRows([]string{"id", "title"}).
			AddRow("invalid", "nasi goreng")

		mock.
			ExpectQuery(qry).
			WithArgs(limit, offset).
			WillReturnRows(recipeRows)

		got, err := repo.GetListRecipe(limit, offset)

		var want []*entity.Recipe

		assertErr(t, err, nil)
		assertRecipesEqual(t, got, want)
	})

	t.Run("should return error getting list", func(t *testing.T) {
		mock.
			ExpectQuery(qry).
			WithArgs(limit, offset).
			WillReturnError(sql.ErrConnDone)

		got, err := repo.GetListRecipe(limit, offset)

		assertErr(t, err, sql.ErrConnDone)
		assertRecipesEqual(t, got, nil)
	})
}
