package models

import (
	"time"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
)

type SubscriptionRequest struct {
	CategoryID         int64      `json:"category_id"`
	Name               string     `json:"name"`
	Description        *string    `json:"description"`
	Cost               float64    `json:"cost"`
	Currency           *string    `json:"currency"`
	BillingCycle       string     `json:"billing_cycle"`
	NextBillingDate    time.Time  `json:"next_billing_date"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            *time.Time `json:"end_date"`
	Status             *string    `json:"status"`
	AutoRenew          *bool      `json:"auto_renew"`
	ReminderEnabled    *bool      `json:"reminder_enabled"`
	ReminderDaysBefore *int64     `json:"reminder_days_before"`
	PaymentMethod      *string    `json:"payment_method"`
	Notes              *string    `json:"notes"`
}

type SubscriptionList struct {
	Total         int                     `json:"total"`
	Subscriptions []database.Subscription `json:"subscriptions"`
}
