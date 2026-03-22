package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func (h *Handler) ListCategory(w http.ResponseWriter, r *http.Request) {
	categories, err := h.config.Queries.ListCategories(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	categoryList := models.CategoryList{
		Total:      len(categories),
		Categories: categories,
	}

	if err := utils.WriteJSON(w, http.StatusOK, &categoryList); err != nil {
		return
	}
}

func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	category, err := h.config.Queries.GetCategory(r.Context(), int64(id))
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, category); err != nil {
		return
	}
}

func (h *Handler) GetCategoryByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	category, err := h.config.Queries.GetCategoryByName(r.Context(), name)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, category); err != nil {
		return
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryRequest models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	category, err := h.config.Queries.CreateCategory(r.Context(), database.CreateCategoryParams{
		Name:        categoryRequest.Name,
		Description: &categoryRequest.Description,
	})
	if err != nil {
		log.Printf("CreateCategory error: %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.WriteJSONError(w, http.StatusConflict, "category name already exists")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, category); err != nil {
		return
	}
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	var categoryRequest models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	category, err := h.config.Queries.UpdateCategory(r.Context(), database.UpdateCategoryParams{
		Name:        categoryRequest.Name,
		Description: &categoryRequest.Description,
		ID:          int64(id),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, category); err != nil {
		return
	}
}

func (h *Handler) UpdateCategoryName(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	var categoryRequest models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	category, err := h.config.Queries.UpdateCategoryName(r.Context(), database.UpdateCategoryNameParams{
		Name: categoryRequest.Name,
		ID:   int64(id),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, category); err != nil {
		return
	}
}

func (h *Handler) UpdateCategoryDescription(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	var categoryRequest models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	category, err := h.config.Queries.UpdateCategoryDescription(r.Context(), database.UpdateCategoryDescriptionParams{
		Description: &categoryRequest.Description,
		ID:          int64(id),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, category); err != nil {
		return
	}
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")

	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	if err := h.config.Queries.DeleteCategory(r.Context(), int64(id)); err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return
	}
}
