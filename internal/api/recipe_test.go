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
	"reflect"
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

func (m *mockRecipeRepository) GetRecipeById(id int64) (*entity.Recipe, error) {
	if id == 1 {
		return &entity.Recipe{
			Id:    1,
			Title: "nasi goreng",
		}, nil
	} else if id == 0 {
		return nil, sql.ErrNoRows
	}
	return nil, sql.ErrConnDone
}

func (m *mockRecipeRepository) DeleteRecipeById(id int64) error {
	if id == 1 {
		return nil
	}
	return sql.ErrConnDone
}

func (m *mockRecipeRepository) GetListRecipe(limit, offset int64) ([]*entity.Recipe, error) {
	if limit == 10 && offset == 0 {
		return []*entity.Recipe{
			{
				Id:    1,
				Title: "nasi goreng",
			},
		}, nil
	}
	return nil, sql.ErrConnDone
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

func assertNotNil(t *testing.T, got any) {
	t.Helper()
	if got == nil {
		t.Errorf("got %v, want %v", got, nil)
	}
}

func assertRecipesEqual(t *testing.T, got, want []*entity.RecipeDTO) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
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

	t.Run("should return 400 error validating request payload", func(t *testing.T) {
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

	idSuccess := map[string]string{
		"id": "1",
	}
	idInvalid := map[string]string{
		"id": "asdf",
	}

	t.Run("should return 200 success update recipe", func(t *testing.T) {
		body, _ := json.Marshal(recipeSuccess)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), idSuccess)
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

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), idSuccess)
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

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), idSuccess)
		rec := httptest.NewRecorder()

		a.updateRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error decoding request payload")
	})

	t.Run("should return 400 error validating request payload", func(t *testing.T) {
		err := recipeFailedValidation.InsertValidate()
		body, _ := json.Marshal(recipeFailedValidation)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), idSuccess)
		rec := httptest.NewRecorder()

		a.updateRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, err.Error())
	})

	t.Run("should return 400 error validating request path id", func(t *testing.T) {
		body, _ := json.Marshal(recipeSuccess)

		req := AddChiURLParams(httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body)), idInvalid)
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

func TestGetRecipeById(t *testing.T) {

	url := "/recipe"

	idSuccess := map[string]string{
		"id": "1",
	}
	idInvalid := map[string]string{
		"id": "asdf",
	}
	idNotFound := map[string]string{
		"id": "0",
	}
	idError := map[string]string{
		"id": "2",
	}

	t.Run("should return 200 success get detail recipe", func(t *testing.T) {
		req := AddChiURLParams(httptest.NewRequest(http.MethodGet, url, nil), idSuccess)
		rec := httptest.NewRecorder()

		a.getRecipeById(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusOK)
		assertMessage(t, res.Message, "success get detail recipe")
		assertNotNil(t, res.Data)
	})

	t.Run("should return 400 error validating request path id", func(t *testing.T) {
		req := AddChiURLParams(httptest.NewRequest(http.MethodGet, url, nil), idInvalid)
		rec := httptest.NewRecorder()

		a.getRecipeById(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "id must be numeric")
	})

	t.Run("should return 400 error recipe not found", func(t *testing.T) {
		req := AddChiURLParams(httptest.NewRequest(http.MethodGet, url, nil), idNotFound)
		rec := httptest.NewRecorder()

		a.getRecipeById(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "recipe not found")
	})

	t.Run("should return 400 error getting recipe", func(t *testing.T) {
		req := AddChiURLParams(httptest.NewRequest(http.MethodGet, url, nil), idError)
		rec := httptest.NewRecorder()

		a.getRecipeById(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error getting recipe")
	})
}

func TestDeleteRecipeById(t *testing.T) {

	url := "/recipe"

	idSuccess := map[string]string{
		"id": "1",
	}
	idInvalid := map[string]string{
		"id": "asdf",
	}
	idError := map[string]string{
		"id": "2",
	}

	t.Run("should return 200 success delete recipe by id", func(t *testing.T) {
		req := AddChiURLParams(httptest.NewRequest(http.MethodGet, url, nil), idSuccess)
		rec := httptest.NewRecorder()

		a.deleteRecipeById(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusOK)
		assertMessage(t, res.Message, "success delete recipe")
	})

	t.Run("should return 400 error validating request path id", func(t *testing.T) {
		req := AddChiURLParams(httptest.NewRequest(http.MethodGet, url, nil), idInvalid)
		rec := httptest.NewRecorder()

		a.deleteRecipeById(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "id must be numeric")
	})

	t.Run("should return 400 error deleting recipe", func(t *testing.T) {
		req := AddChiURLParams(httptest.NewRequest(http.MethodGet, url, nil), idError)
		rec := httptest.NewRecorder()

		a.deleteRecipeById(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error delete recipe")
	})
}

func TestGetListRecipe(t *testing.T) {
	urlSuccess := "/recipe-list?page=1&limit=10"
	urlPageInvalid := "/recipe-list?page=asdf&limit=10"
	urlLimitInvalid := "/recipe-list?page=1&limit=asdf"
	urlFailed := "/recipe-list?page=2&limit=2"

	t.Run("should return 200 get list recipe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, urlSuccess, nil)
		rec := httptest.NewRecorder()

		a.getListRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		byteData, _ := json.Marshal(res.Data)

		var got []*entity.RecipeDTO
		_ = json.Unmarshal(byteData, &got)

		want := []*entity.RecipeDTO{
			{
				Id:    1,
				Title: "nasi goreng",
			},
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusOK)
		assertMessage(t, res.Message, "success get list recipe")
		assertRecipesEqual(t, got, want)
	})

	t.Run("should return 400 error validating query param page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, urlPageInvalid, nil)
		rec := httptest.NewRecorder()

		a.getListRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "page must be numeric")
	})

	t.Run("should return 400 error validating query param limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, urlLimitInvalid, nil)
		rec := httptest.NewRecorder()

		a.getListRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "limit must be numeric")
	})

	t.Run("should return 400 error get list recipe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, urlFailed, nil)
		rec := httptest.NewRecorder()

		a.getListRecipe(rec, req)

		var res helper.Response
		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatalf("error decoding response body, %v", err.Error())
		}

		assertStatusCode(t, int32(res.StatusCode), http.StatusBadRequest)
		assertMessage(t, res.Message, "error get list recipe")
	})
}

func AddChiURLParams(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}
