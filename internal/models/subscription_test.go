package models

import (
	"encoding/json"
	"testing"
	"time"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
)

func TestSubscriptionRequest(t *testing.T) {
	now := time.Now()
	description := "Test subscription"
	currency := "USD"
	status := "active"
	autoRenew := true
	reminderEnabled := true
	reminderDays := int64(7)
	paymentMethod := "credit_card"
	notes := "Some notes"
	endDate := now.Add(24 * time.Hour * 30)

	request := SubscriptionRequest{
		CategoryID:         1,
		Name:               "Netflix",
		Description:        &description,
		Cost:               15.99,
		Currency:           &currency,
		BillingCycle:       "monthly",
		NextBillingDate:    now.Add(30 * 24 * time.Hour),
		StartDate:          now,
		EndDate:            &endDate,
		Status:             &status,
		AutoRenew:          &autoRenew,
		ReminderEnabled:    &reminderEnabled,
		ReminderDaysBefore: &reminderDays,
		PaymentMethod:      &paymentMethod,
		Notes:              &notes,
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal SubscriptionRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["category_id"] != float64(request.CategoryID) {
		t.Errorf("SubscriptionRequest category_id = %v, want %v", result["category_id"], request.CategoryID)
	}

	if result["name"] != request.Name {
		t.Errorf("SubscriptionRequest name = %v, want %v", result["name"], request.Name)
	}

	if result["cost"] != request.Cost {
		t.Errorf("SubscriptionRequest cost = %v, want %v", result["cost"], request.Cost)
	}

	if result["billing_cycle"] != request.BillingCycle {
		t.Errorf("SubscriptionRequest billing_cycle = %v, want %v", result["billing_cycle"], request.BillingCycle)
	}
}

func TestSubscriptionRequestWithNilPointers(t *testing.T) {
	now := time.Now()

	request := SubscriptionRequest{
		CategoryID:         1,
		Name:               "Basic Subscription",
		Description:        nil,
		Cost:               9.99,
		Currency:           nil,
		BillingCycle:       "monthly",
		NextBillingDate:    now,
		StartDate:          now,
		EndDate:            nil,
		Status:             nil,
		AutoRenew:          nil,
		ReminderEnabled:    nil,
		ReminderDaysBefore: nil,
		PaymentMethod:      nil,
		Notes:              nil,
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal SubscriptionRequest with nil pointers: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check that required fields are present
	if result["category_id"] == nil {
		t.Error("SubscriptionRequest category_id should not be nil")
	}

	if result["name"] == nil {
		t.Error("SubscriptionRequest name should not be nil")
	}

	// Nil pointer fields should be null in JSON
	if result["description"] != nil {
		t.Error("SubscriptionRequest description should be null when nil pointer")
	}
}

func TestSubscriptionList(t *testing.T) {
	now := time.Now()
	subscriptions := []database.Subscription{
		{
			ID:              1,
			CategoryID:      1,
			Name:            "Netflix",
			Description:     nil,
			Cost:            15.99,
			Currency:        nil,
			BillingCycle:    "monthly",
			NextBillingDate: now,
			StartDate:       now,
			EndDate:         nil,
			Status:          nil,
			AutoRenew:       nil,
			CreatedAt:       &now,
			UpdatedAt:       &now,
		},
	}

	list := SubscriptionList{
		Total:         len(subscriptions),
		Subscriptions: subscriptions,
	}

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal SubscriptionList: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["total"] != float64(list.Total) {
		t.Errorf("SubscriptionList total = %v, want %v", result["total"], list.Total)
	}

	if subs, ok := result["subscriptions"].([]interface{}); !ok {
		t.Error("SubscriptionList subscriptions should be an array")
	} else if len(subs) != list.Total {
		t.Errorf("SubscriptionList subscriptions length = %v, want %v", len(subs), list.Total)
	}
}

func TestSubscriptionListEmpty(t *testing.T) {
	list := SubscriptionList{
		Total:         0,
		Subscriptions: []database.Subscription{},
	}

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal empty SubscriptionList: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["total"] != float64(0) {
		t.Errorf("Empty SubscriptionList total = %v, want 0", result["total"])
	}

	if subs, ok := result["subscriptions"].([]interface{}); !ok {
		t.Error("Empty SubscriptionList subscriptions should be an array")
	} else if len(subs) != 0 {
		t.Errorf("Empty SubscriptionList subscriptions length = %v, want 0", len(subs))
	}
}
