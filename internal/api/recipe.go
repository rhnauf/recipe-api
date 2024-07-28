package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/rhnauf/recipe-api/internal/entity"
	"github.com/rhnauf/recipe-api/internal/helper"
)

func (a *api) insertRecipe(w http.ResponseWriter, r *http.Request) {
	var requestRecipe entity.RecipeDTO
	if err := json.NewDecoder(r.Body).Decode(&requestRecipe); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "error decoding request payload", nil)
		return
	}

	if err := requestRecipe.InsertValidate(); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	recipe := entity.Recipe{
		Title:       requestRecipe.Title,
		Description: requestRecipe.Description,
		Instruction: requestRecipe.Instruction,
		Publish:     requestRecipe.Publish,
	}

	if err := a.recipeRepository.InsertRecipe(recipe); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "error insert recipe", nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, "success insert recipe", nil)
}

func (a *api) updateRecipe(w http.ResponseWriter, r *http.Request) {
	var requestRecipe entity.RecipeDTO
	if err := json.NewDecoder(r.Body).Decode(&requestRecipe); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "error decoding request payload", nil)
		return
	}

	pathParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(pathParam, 0, 64)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "id must be numeric", nil)
		return
	}
	requestRecipe.SetId(id)

	if err := requestRecipe.UpdateValidate(); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	recipe := entity.Recipe{
		Id:          requestRecipe.Id,
		Title:       requestRecipe.Title,
		Description: requestRecipe.Description,
		Instruction: requestRecipe.Instruction,
		Publish:     requestRecipe.Publish,
	}

	if err := a.recipeRepository.UpdateRecipe(recipe); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "error update recipe", nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, "success update recipe", nil)
}

func (a *api) getRecipeById(w http.ResponseWriter, r *http.Request) {
	pathParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(pathParam, 0, 64)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "id must be numeric", nil)
		return
	}

	recipe, err := a.recipeRepository.GetRecipeById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.HandleResponse(w, http.StatusBadRequest, "recipe not found", nil)
			return
		}
		helper.HandleResponse(w, http.StatusBadRequest, "error getting recipe", nil)
		return
	}

	recipeDto := entity.RecipeDTO{
		Id:          recipe.Id,
		Title:       recipe.Title,
		Description: recipe.Description,
		Instruction: recipe.Instruction,
		Publish:     recipe.Publish,
		CreatedAt:   recipe.CreatedAt.Format("02-01-2006"),
	}

	helper.HandleResponse(w, http.StatusOK, "success get detail recipe", recipeDto)
}

func (a *api) deleteRecipeById(w http.ResponseWriter, r *http.Request) {
	pathParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(pathParam, 0, 64)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "id must be numeric", nil)
		return
	}

	err = a.recipeRepository.DeleteRecipeById(id)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "error delete recipe", nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, "success delete recipe", nil)
}

const (
	DefaultPage = 1
	MinPage     = 1

	DefaultLimit = 10
	MinLimit     = 1
	MaxLimit     = 100
)

func (a *api) getListRecipe(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	var page int64 = DefaultPage
	var limit int64 = DefaultLimit

	if pageParam != "" {
		p, err := strconv.ParseInt(pageParam, 0, 64)
		if err != nil {
			helper.HandleResponse(w, http.StatusBadRequest, "page must be numeric", nil)
			return
		}
		if p >= MinPage {
			page = p
		}
	}

	if limitParam != "" {
		l, err := strconv.ParseInt(limitParam, 0, 64)
		if err != nil {
			helper.HandleResponse(w, http.StatusBadRequest, "limit must be numeric", nil)
			return
		}
		if l >= MinLimit && l <= MaxLimit {
			limit = l
		}
	}

	offset := limit * (page - 1)

	recipes, err := a.recipeRepository.GetListRecipe(limit, offset)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, "error get list recipe", nil)
		return
	}

	res := make([]*entity.RecipeDTO, len(recipes))
	for idx, recipe := range recipes {
		res[idx] = &entity.RecipeDTO{
			Id:    recipe.Id,
			Title: recipe.Title,
		}
	}

	helper.HandleResponse(w, http.StatusOK, "success get list recipe", res)
}
