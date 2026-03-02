-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_subscriptions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    subscription_id INTEGER NOT NULL REFERENCES subscriptions(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_subscriptions;
-- +goose StatementEnd
