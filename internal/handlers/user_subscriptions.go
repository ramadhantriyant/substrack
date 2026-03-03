package handlers

import (
	"net/http"
	"strconv"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/middlewares"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func (h *Handler) ListUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	subscriptions, err := h.config.Queries.ListSubscriptionsByUser(r.Context(), userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, models.SubscriptionList{
		Total:         len(subscriptions),
		Subscriptions: subscriptions,
	}); err != nil {
		return
	}
}

func (h *Handler) AddUserSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idPath := r.PathValue("id")
	subscriptionID, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	exists, err := h.config.Queries.UserHasSubscription(r.Context(), database.UserHasSubscriptionParams{
		UserID:         userID,
		SubscriptionID: int64(subscriptionID),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if exists > 0 {
		utils.WriteJSONError(w, http.StatusConflict, "subscription already added")
		return
	}

	record, err := h.config.Queries.AddSubscriptionToUser(r.Context(), database.AddSubscriptionToUserParams{
		UserID:         userID,
		SubscriptionID: int64(subscriptionID),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, record); err != nil {
		return
	}
}

func (h *Handler) RemoveUserSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idPath := r.PathValue("id")
	subscriptionID, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	if err := h.config.Queries.RemoveSubscriptionFromUser(r.Context(), database.RemoveSubscriptionFromUserParams{
		UserID:         userID,
		SubscriptionID: int64(subscriptionID),
	}); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return
	}
}
