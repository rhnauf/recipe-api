package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rhnauf/recipe-api/internal/entity"
	"github.com/rhnauf/recipe-api/internal/helper"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRecipeRepository struct{}

func (m *mockRecipeRepository) InsertRecipe(recipe entity.Recipe) error {
	if recipe.Title == "failed" {
		return sql.ErrConnDone
	}
	return nil
}

func (m *mockRecipeRepository) UpdateRecipe(recipe entity.Recipe) error {
	if recipe.Title == "failed" {
		return sql.ErrConnDone
	}
	return nil
}

func assertStatusCode(t *testing.T, got, want int32) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertMessage(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

var (
	mockRepo = &mockRecipeRepository{}
	a        = api{
		recipeRepository: mockRepo,
	}
)

func TestInsertRecipe(t *testing.T) {

	var p = true
	recipeSuccess := entity.RecipeDTO{
		Title:       "test",
		Description: "test",
		Instruction: "test",
		Publish:     &p,
	}

	recipeFailed := entity.RecipeDTO{
		Title:       "failed",
		Description: "failed",
		Instruction: "failed",
		Publish:     &p,
	}

	recipeFailedValidation := entity.RecipeDTO{
		Title: "",
	}

	url := "/recipe"

	t.Run("should return 200 success insert recipe", func(t *testing.T) {
		body, _ := json.Marshal(recipeSuccess)

		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		a.insertRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusOK)
		assertMessage(t, res.Message, "success insert recipe")
	})

	t.Run("should return 400 error insert recipe", func(t *testing.T) {
		body, _ := json.Marshal(recipeFailed)

		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		a.insertRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error insert recipe")
	})

	t.Run("should return 400 error decode payload", func(t *testing.T) {
		body := []byte(`invalid body`)

		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		a.insertRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error decoding request payload")
	})

	t.Run("should return 400 error error validating request payload", func(t *testing.T) {
		err := recipeFailedValidation.InsertValidate()
		body, _ := json.Marshal(recipeFailedValidation)

		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		a.insertRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, err.Error())
	})
}

func TestUpdateRecipe(t *testing.T) {

	var p = true
	recipeSuccess := entity.RecipeDTO{
		Title:       "test",
		Description: "test",
		Instruction: "test",
		Publish:     &p,
	}

	recipeFailed := entity.RecipeDTO{
		Title:       "failed",
		Description: "failed",
		Instruction: "failed",
		Publish:     &p,
	}

	recipeFailedValidation := entity.RecipeDTO{
		Id: 0,
	}

	url := "/recipe"

	path := map[string]string{
		"id": "1",
	}
	pathInvalid := map[string]string{
		"id": "asdf",
	}

	t.Run("should return 200 success update recipe", func(t *testing.T) {
		body, _ := json.Marshal(recipeSuccess)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), path)
		rec := httptest.NewRecorder()

		a.updateRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusOK)
		assertMessage(t, res.Message, "success update recipe")
	})

	t.Run("should return 400 error update recipe", func(t *testing.T) {
		body, _ := json.Marshal(recipeFailed)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), path)
		rec := httptest.NewRecorder()

		a.updateRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error update recipe")
	})

	t.Run("should return 400 error decode payload", func(t *testing.T) {
		body := []byte(`invalid body`)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), path)
		rec := httptest.NewRecorder()

		a.updateRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error decoding request payload")
	})

	t.Run("should return 400 error error validating request payload", func(t *testing.T) {
		err := recipeFailedValidation.InsertValidate()
		body, _ := json.Marshal(recipeFailedValidation)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), path)
		rec := httptest.NewRecorder()

		a.updateRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, err.Error())
	})

	t.Run("should return 400 error error validating request path id", func(t *testing.T) {
		body, _ := json.Marshal(recipeSuccess)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), pathInvalid)
		rec := httptest.NewRecorder()

		a.updateRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "id must be numeric")
	})
}

func AddChiURLParams(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}
