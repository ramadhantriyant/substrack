-- name: GetUserSubscription :one
SELECT * FROM users_subscriptions
WHERE id = ? LIMIT 1;

-- name: ListUserSubscriptions :many
SELECT * FROM users_subscriptions
ORDER BY created_at DESC;

-- name: ListSubscriptionsByUser :many
SELECT s.*
FROM subscriptions s
JOIN users_subscriptions us ON s.id = us.subscription_id
WHERE us.user_id = ?
ORDER BY s.next_billing_date ASC;

-- name: ListUsersBySubscription :many
SELECT u.*
FROM users u
JOIN users_subscriptions us ON u.id = us.user_id
WHERE us.subscription_id = ?
ORDER BY u.name;

-- name: GetUserSubscriptionsWithDetails :many
SELECT
    us.*,
    u.name as user_name,
    u.email as user_email,
    s.name as subscription_name,
    s.cost as subscription_cost,
    s.currency as subscription_currency,
    s.billing_cycle as billing_cycle,
    s.next_billing_date as next_billing_date,
    s.status as subscription_status
FROM users_subscriptions us
JOIN users u ON us.user_id = u.id
JOIN subscriptions s ON us.subscription_id = s.id
WHERE us.user_id = ?
ORDER BY s.next_billing_date ASC;

-- name: GetUserActiveSubscriptions :many
SELECT s.*
FROM subscriptions s
JOIN users_subscriptions us ON s.id = us.subscription_id
WHERE us.user_id = ? AND s.status = 'active'
ORDER BY s.next_billing_date ASC;

-- name: AddSubscriptionToUser :one
INSERT INTO users_subscriptions (
    user_id, subscription_id
) VALUES (
    ?, ?
)
RETURNING *;

-- name: RemoveSubscriptionFromUser :exec
DELETE FROM users_subscriptions
WHERE user_id = ? AND subscription_id = ?;

-- name: RemoveUserSubscription :exec
DELETE FROM users_subscriptions
WHERE id = ?;

-- name: RemoveAllSubscriptionsFromUser :exec
DELETE FROM users_subscriptions
WHERE user_id = ?;

-- name: RemoveAllUsersFromSubscription :exec
DELETE FROM users_subscriptions
WHERE subscription_id = ?;

-- name: UserHasSubscription :one
SELECT EXISTS(
    SELECT 1 FROM users_subscriptions
    WHERE user_id = ? AND subscription_id = ?
);

-- name: CountSubscriptionsByUser :one
SELECT COUNT(*) FROM users_subscriptions
WHERE user_id = ?;

-- name: CountUsersBySubscription :one
SELECT COUNT(*) FROM users_subscriptions
WHERE subscription_id = ?;

-- name: CountUserSubscriptions :one
SELECT COUNT(*) FROM users_subscriptions;

-- name: GetUserSubscriptionsCreatedAfter :many
SELECT * FROM users_subscriptions
WHERE created_at > ?
ORDER BY created_at DESC;

-- name: CalculateUserTotalMonthlyCost :one
SELECT SUM(
    CASE s.billing_cycle
        WHEN 'daily' THEN s.cost * 30
        WHEN 'weekly' THEN s.cost * 4.33
        WHEN 'monthly' THEN s.cost
        WHEN 'quarterly' THEN s.cost / 3
        WHEN 'yearly' THEN s.cost / 12
    END
) as total_monthly_cost
FROM subscriptions s
JOIN users_subscriptions us ON s.id = us.subscription_id
WHERE us.user_id = ? AND s.status = 'active' AND s.currency = ?;

-- name: CalculateUserTotalYearlyCost :one
SELECT SUM(
    CASE s.billing_cycle
        WHEN 'daily' THEN s.cost * 365
        WHEN 'weekly' THEN s.cost * 52
        WHEN 'monthly' THEN s.cost * 12
        WHEN 'quarterly' THEN s.cost * 4
        WHEN 'yearly' THEN s.cost
    END
) as total_yearly_cost
FROM subscriptions s
JOIN users_subscriptions us ON s.id = us.subscription_id
WHERE us.user_id = ? AND s.status = 'active' AND s.currency = ?;

-- name: GetUserUpcomingRenewals :many
SELECT s.*
FROM subscriptions s
JOIN users_subscriptions us ON s.id = us.subscription_id
WHERE us.user_id = ?
  AND s.status = 'active'
  AND s.next_billing_date BETWEEN CURRENT_DATE AND DATE(CURRENT_DATE, '+' || ? || ' days')
ORDER BY s.next_billing_date ASC;

-- name: GetUserSubscriptionsByCategory :many
SELECT s.*
FROM subscriptions s
JOIN users_subscriptions us ON s.id = us.subscription_id
WHERE us.user_id = ? AND s.category_id = ?
ORDER BY s.next_billing_date ASC;
