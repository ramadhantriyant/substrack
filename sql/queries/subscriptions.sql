-- name: GetSubscription :one
SELECT * FROM subscriptions
WHERE id = ? LIMIT 1;

-- name: GetSubscriptionWithCategory :one
SELECT
    s.*,
    c.name as category_name,
    c.description as category_description
FROM subscriptions s
JOIN categories c ON s.category_id = c.id
WHERE s.id = ? LIMIT 1;

-- name: ListSubscriptions :many
SELECT * FROM subscriptions
ORDER BY next_billing_date ASC;

-- name: ListSubscriptionsWithLimit :many
SELECT * FROM subscriptions
ORDER BY next_billing_date ASC
LIMIT ? OFFSET ?;

-- name: ListSubscriptionsWithCategories :many
SELECT
    s.*,
    c.name as category_name
FROM subscriptions s
JOIN categories c ON s.category_id = c.id
ORDER BY s.next_billing_date ASC;

-- name: ListSubscriptionsByCategory :many
SELECT * FROM subscriptions
WHERE category_id = ?
ORDER BY next_billing_date ASC;

-- name: ListSubscriptionsByStatus :many
SELECT * FROM subscriptions
WHERE status = ?
ORDER BY next_billing_date ASC;

-- name: ListActiveSubscriptions :many
SELECT * FROM subscriptions
WHERE status = 'active'
ORDER BY next_billing_date ASC;

-- name: ListSubscriptionsByBillingCycle :many
SELECT * FROM subscriptions
WHERE billing_cycle = ?
ORDER BY next_billing_date ASC;

-- name: CreateSubscription :one
INSERT INTO subscriptions (
    category_id,
    name,
    description,
    cost,
    currency,
    billing_cycle,
    next_billing_date,
    start_date,
    end_date,
    status,
    auto_renew,
    reminder_enabled,
    reminder_days_before,
    payment_method,
    notes
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET category_id = ?,
    name = ?,
    description = ?,
    cost = ?,
    currency = ?,
    billing_cycle = ?,
    next_billing_date = ?,
    start_date = ?,
    end_date = ?,
    status = ?,
    auto_renew = ?,
    reminder_enabled = ?,
    reminder_days_before = ?,
    payment_method = ?,
    notes = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateSubscriptionStatus :one
UPDATE subscriptions
SET status = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateSubscriptionNextBillingDate :one
UPDATE subscriptions
SET next_billing_date = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateSubscriptionCost :one
UPDATE subscriptions
SET cost = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateSubscriptionReminder :one
UPDATE subscriptions
SET reminder_enabled = ?,
    reminder_days_before = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: PauseSubscription :one
UPDATE subscriptions
SET status = 'paused',
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: CancelSubscription :one
UPDATE subscriptions
SET status = 'cancelled',
    end_date = CURRENT_DATE,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: ActivateSubscription :one
UPDATE subscriptions
SET status = 'active',
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE id = ?;

-- name: CountSubscriptions :one
SELECT COUNT(*) FROM subscriptions;

-- name: CountSubscriptionsByStatus :one
SELECT COUNT(*) FROM subscriptions
WHERE status = ?;

-- name: CountSubscriptionsByCategory :one
SELECT COUNT(*) FROM subscriptions
WHERE category_id = ?;

-- name: SubscriptionExists :one
SELECT EXISTS(SELECT 1 FROM subscriptions WHERE id = ?);

-- name: SearchSubscriptions :many
SELECT * FROM subscriptions
WHERE name LIKE ? OR description LIKE ? OR notes LIKE ?
ORDER BY next_billing_date ASC;

-- name: GetSubscriptionsDueForRenewal :many
SELECT * FROM subscriptions
WHERE status = 'active'
  AND next_billing_date <= ?
ORDER BY next_billing_date ASC;

-- name: GetSubscriptionsWithReminders :many
SELECT * FROM subscriptions
WHERE status = 'active'
  AND reminder_enabled = 1
  AND DATE(next_billing_date, '-' || reminder_days_before || ' days') <= ?
ORDER BY next_billing_date ASC;

-- name: GetSubscriptionsBetweenDates :many
SELECT * FROM subscriptions
WHERE next_billing_date BETWEEN ? AND ?
ORDER BY next_billing_date ASC;

-- name: GetExpiredSubscriptions :many
SELECT * FROM subscriptions
WHERE end_date IS NOT NULL
  AND end_date < CURRENT_DATE
  AND status != 'cancelled'
ORDER BY end_date DESC;

-- name: CalculateTotalMonthlyCost :one
SELECT SUM(
    CASE billing_cycle
        WHEN 'daily' THEN cost * 30
        WHEN 'weekly' THEN cost * 4.33
        WHEN 'monthly' THEN cost
        WHEN 'quarterly' THEN cost / 3
        WHEN 'yearly' THEN cost / 12
    END
) as total_monthly_cost
FROM subscriptions
WHERE status = 'active' AND currency = ?;

-- name: CalculateTotalYearlyCost :one
SELECT SUM(
    CASE billing_cycle
        WHEN 'daily' THEN cost * 365
        WHEN 'weekly' THEN cost * 52
        WHEN 'monthly' THEN cost * 12
        WHEN 'quarterly' THEN cost * 4
        WHEN 'yearly' THEN cost
    END
) as total_yearly_cost
FROM subscriptions
WHERE status = 'active' AND currency = ?;

-- name: GetSubscriptionStatsByCategory :many
SELECT
    c.name as category_name,
    COUNT(s.id) as subscription_count,
    SUM(s.cost) as total_cost,
    AVG(s.cost) as average_cost
FROM categories c
LEFT JOIN subscriptions s ON c.id = s.category_id AND s.status = 'active'
GROUP BY c.id, c.name
ORDER BY subscription_count DESC;

-- name: GetUpcomingRenewals :many
SELECT * FROM subscriptions
WHERE status = 'active'
  AND next_billing_date BETWEEN CURRENT_DATE AND DATE(CURRENT_DATE, '+' || ? || ' days')
ORDER BY next_billing_date ASC;

-- name: GetSubscriptionsByAutoRenew :many
SELECT * FROM subscriptions
WHERE auto_renew = ?
  AND status = 'active'
ORDER BY next_billing_date ASC;

-- name: GetRecentSubscriptions :many
SELECT * FROM subscriptions
WHERE created_at >= ?
ORDER BY created_at DESC;

-- name: GetModifiedSubscriptions :many
SELECT * FROM subscriptions
WHERE updated_at >= ?
ORDER BY updated_at DESC;
