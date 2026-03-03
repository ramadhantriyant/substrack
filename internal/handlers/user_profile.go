package handlers

import (
	"encoding/json"
	"net/http"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/middlewares"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.config.Queries.GetUser(r.Context(), userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, models.UserToResponse(user)); err != nil {
		return
	}
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.Email != "" && !utils.ValidateEmail(req.Email) {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid email address")
		return
	}

	// Fetch current values to fill in any omitted fields
	current, err := h.config.Queries.GetUser(r.Context(), userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	email := current.Email
	if req.Email != "" {
		email = req.Email
	}

	name := current.Name
	if req.Name != "" {
		name = req.Name
	}

	updated, err := h.config.Queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: email,
		Name:  name,
		ID:    userID,
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, models.UserToResponse(updated)); err != nil {
		return
	}
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	if !utils.ValidatePassword(req.NewPassword) {
		utils.WriteJSONError(w, http.StatusBadRequest, "new password must be at least 8 characters")
		return
	}

	user, err := h.config.Queries.GetUser(r.Context(), userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	match, err := utils.VerifyPassword(req.OldPassword, user.Password)
	if err != nil || !match {
		utils.WriteJSONError(w, http.StatusUnauthorized, "old password is incorrect")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if _, err := h.config.Queries.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		Password: hashedPassword,
		ID:       userID,
	}); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated"}); err != nil {
		return
	}
}

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.config.Queries.DeleteUser(r.Context(), userID); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return
	}
}
