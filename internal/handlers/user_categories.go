package handlers

import (
	"net/http"
	"strconv"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/middlewares"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func (h *Handler) ListUserCategories(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	categories, err := h.config.Queries.ListCategoriesByUser(r.Context(), userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, models.CategoryList{
		Total:      len(categories),
		Categories: categories,
	}); err != nil {
		return
	}
}

func (h *Handler) AddUserCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idPath := r.PathValue("id")
	categoryID, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	exists, err := h.config.Queries.UserHasCategory(r.Context(), database.UserHasCategoryParams{
		UserID:     userID,
		CategoryID: int64(categoryID),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if exists > 0 {
		utils.WriteJSONError(w, http.StatusConflict, "category already added")
		return
	}

	record, err := h.config.Queries.AddCategoryToUser(r.Context(), database.AddCategoryToUserParams{
		UserID:     userID,
		CategoryID: int64(categoryID),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, record); err != nil {
		return
	}
}

func (h *Handler) RemoveUserCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idPath := r.PathValue("id")
	categoryID, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	if err := h.config.Queries.RemoveCategoryFromUser(r.Context(), database.RemoveCategoryFromUserParams{
		UserID:     userID,
		CategoryID: int64(categoryID),
	}); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return
	}
}
