package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/middlewares"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	if !utils.ValidateEmail(req.Email) {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid email address")
		return
	}

	if !utils.ValidatePassword(req.Password) {
		utils.WriteJSONError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	exists, err := h.config.Queries.UserExistsByEmail(r.Context(), req.Email)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if exists > 0 {
		utils.WriteJSONError(w, http.StatusConflict, "email already registered")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	user, err := h.config.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	_ = h.config.Queries.LinkAllCategoriesToUser(r.Context(), user.ID)

	if err := utils.WriteJSON(w, http.StatusCreated, models.UserToResponse(user)); err != nil {
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := h.config.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	match, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil || !match {
		utils.WriteJSONError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	accessToken, err := utils.MakeJWT(user.ID, h.config.JWTSecret)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	rawRefreshToken, err := utils.MakeRefreshToken()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	tokenHash := utils.HashRefreshToken(rawRefreshToken)
	_, err = h.config.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(utils.RefreshTokenDuration),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
	}); err != nil {
		return
	}
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.RefreshToken == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokenHash := utils.HashRefreshToken(req.RefreshToken)
	stored, err := h.config.Queries.GetRefreshTokenByHash(r.Context(), tokenHash)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	if stored.Revoked != 0 {
		utils.WriteJSONError(w, http.StatusUnauthorized, "refresh token has been revoked")
		return
	}

	if stored.ExpiresAt.Before(time.Now()) {
		utils.WriteJSONError(w, http.StatusUnauthorized, "refresh token has expired")
		return
	}

	accessToken, err := utils.MakeJWT(stored.UserID, h.config.JWTSecret)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		TokenType:    "Bearer",
	}); err != nil {
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.RefreshToken == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokenHash := utils.HashRefreshToken(req.RefreshToken)
	// Ignore error — idempotent: revoking a non-existent token is fine
	_ = h.config.Queries.RevokeRefreshToken(r.Context(), tokenHash)

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"}); err != nil {
		return
	}
}
