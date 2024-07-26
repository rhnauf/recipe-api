package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rhnauf/recipe-api/internal/entity"
	"github.com/rhnauf/recipe-api/internal/helper"
	"net/http"
	"strconv"
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
