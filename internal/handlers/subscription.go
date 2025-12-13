package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func (h *Handler) ListSubscription(w http.ResponseWriter, r *http.Request) {
	var subscriptions []database.Subscription
	var err error

	categoryIDString := r.URL.Query().Get("category_id")

	if categoryIDString == "" {
		subscriptions, err = h.config.Queries.ListSubscriptions(r.Context())
	} else {
		categoryID, err := strconv.Atoi(categoryIDString)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid category id")
			return
		}
		subscriptions, err = h.config.Queries.ListSubscriptionsByCategory(r.Context(), int64(categoryID))
	}

	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	subscriptionList := models.SubscriptionList{
		Total:         len(subscriptions),
		Subscriptions: subscriptions,
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscriptionList); err != nil {
		return
	}
}

func (h *Handler) ListActiveSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.config.Queries.ListActiveSubscriptions(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	subscriptionList := models.SubscriptionList{
		Total:         len(subscriptions),
		Subscriptions: subscriptions,
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscriptionList); err != nil {
		return
	}
}

func (h *Handler) ListExpiredSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.config.Queries.GetExpiredSubscriptions(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	subscriptionList := models.SubscriptionList{
		Total:         len(subscriptions),
		Subscriptions: subscriptions,
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscriptionList); err != nil {
		return
	}
}

func (h *Handler) ListSubscriptionsByBillingCycle(w http.ResponseWriter, r *http.Request) {
	billingCycle := r.PathValue("billCycle")

	subscriptions, err := h.config.Queries.ListSubscriptionsByBillingCycle(r.Context(), billingCycle)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid billing cycle")
		return
	}

	subscriptionList := models.SubscriptionList{
		Total:         len(subscriptions),
		Subscriptions: subscriptions,
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscriptionList); err != nil {
		return
	}
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	subscription, err := h.config.Queries.GetSubscription(r.Context(), int64(id))
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscription); err != nil {
		return
	}
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var subscriptionRequest models.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subscriptionRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	subscription, err := h.config.Queries.CreateSubscription(r.Context(), database.CreateSubscriptionParams{
		CategoryID:         subscriptionRequest.CategoryID,
		Name:               subscriptionRequest.Name,
		Description:        subscriptionRequest.Description,
		Cost:               subscriptionRequest.Cost,
		Currency:           subscriptionRequest.Currency,
		BillingCycle:       subscriptionRequest.BillingCycle,
		NextBillingDate:    subscriptionRequest.NextBillingDate,
		StartDate:          subscriptionRequest.StartDate,
		EndDate:            subscriptionRequest.EndDate,
		Status:             subscriptionRequest.Status,
		AutoRenew:          subscriptionRequest.AutoRenew,
		ReminderEnabled:    subscriptionRequest.ReminderEnabled,
		ReminderDaysBefore: subscriptionRequest.ReminderDaysBefore,
		PaymentMethod:      subscriptionRequest.PaymentMethod,
		Notes:              subscriptionRequest.Notes,
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, subscription); err != nil {
		return
	}
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	var subscriptionRequest models.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subscriptionRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	subscription, err := h.config.Queries.UpdateSubscription(r.Context(), database.UpdateSubscriptionParams{
		CategoryID:         subscriptionRequest.CategoryID,
		Name:               subscriptionRequest.Name,
		Description:        subscriptionRequest.Description,
		Cost:               subscriptionRequest.Cost,
		Currency:           subscriptionRequest.Currency,
		BillingCycle:       subscriptionRequest.BillingCycle,
		NextBillingDate:    subscriptionRequest.NextBillingDate,
		StartDate:          subscriptionRequest.StartDate,
		EndDate:            subscriptionRequest.EndDate,
		Status:             subscriptionRequest.Status,
		AutoRenew:          subscriptionRequest.AutoRenew,
		ReminderEnabled:    subscriptionRequest.ReminderEnabled,
		ReminderDaysBefore: subscriptionRequest.ReminderDaysBefore,
		PaymentMethod:      subscriptionRequest.PaymentMethod,
		Notes:              subscriptionRequest.Notes,
		ID:                 int64(id),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscription); err != nil {
		return
	}
}

func (h *Handler) UpdateSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	var subscriptionRequest models.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subscriptionRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	subscription, err := h.config.Queries.UpdateSubscriptionStatus(r.Context(), database.UpdateSubscriptionStatusParams{
		Status: subscriptionRequest.Status,
		ID:     int64(id),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusOK, "subscription not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, subscription); err != nil {
		return
	}
}

func (h *Handler) UpdateSubscriptionCost(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	var subscriptionRequest models.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subscriptionRequest); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	subscription, err := h.config.Queries.UpdateSubscriptionCost(r.Context(), database.UpdateSubscriptionCostParams{
		Cost: subscriptionRequest.Cost,
		ID:   int64(id),
	})
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscription); err != nil {
		return
	}
}

func (h *Handler) PauseSubscription(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	subscription, err := h.config.Queries.PauseSubscription(r.Context(), int64(id))
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, subscription); err != nil {
		return
	}
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	if err := h.config.Queries.DeleteSubscription(r.Context(), int64(id)); err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "subscription not found")
		return
	}

	if err := utils.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return
	}
}
