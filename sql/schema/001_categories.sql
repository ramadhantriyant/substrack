-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_name ON categories(name);

-- Insert default categories
INSERT INTO categories (id, name, description) VALUES
    (1, 'Streaming', 'Video and music streaming services'),
    (2, 'Software', 'Software and SaaS subscriptions'),
    (3, 'Cloud Storage', 'Cloud storage and backup services'),
    (4, 'Gaming', 'Gaming platforms and services'),
    (5, 'News & Media', 'News, magazines, and media subscriptions'),
    (6, 'Productivity', 'Productivity and business tools'),
    (7, 'Health & Fitness', 'Health and fitness apps'),
    (8, 'Education', 'Online courses and learning platforms');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_categories_name;
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
