-- name: GetUserCategory :one
SELECT * FROM users_categories
WHERE id = ? LIMIT 1;

-- name: ListUserCategories :many
SELECT * FROM users_categories
ORDER BY created_at DESC;

-- name: ListCategoriesByUser :many
SELECT c.*
FROM categories c
JOIN users_categories uc ON c.id = uc.category_id
WHERE uc.user_id = ?
ORDER BY c.name;

-- name: ListUsersByCategory :many
SELECT u.*
FROM users u
JOIN users_categories uc ON u.id = uc.user_id
WHERE uc.category_id = ?
ORDER BY u.name;

-- name: GetUserCategoriesWithDetails :many
SELECT
    uc.*,
    u.name as user_name,
    u.email as user_email,
    c.name as category_name,
    c.description as category_description
FROM users_categories uc
JOIN users u ON uc.user_id = u.id
JOIN categories c ON uc.category_id = c.id
WHERE uc.user_id = ?
ORDER BY c.name;

-- name: AddCategoryToUser :one
INSERT INTO users_categories (
    user_id, category_id
) VALUES (
    ?, ?
)
RETURNING *;

-- name: LinkAllCategoriesToUser :exec
INSERT OR IGNORE INTO users_categories (user_id, category_id)
SELECT ?, id FROM categories;

-- name: RemoveCategoryFromUser :exec
DELETE FROM users_categories
WHERE user_id = ? AND category_id = ?;

-- name: RemoveUserCategory :exec
DELETE FROM users_categories
WHERE id = ?;

-- name: RemoveAllCategoriesFromUser :exec
DELETE FROM users_categories
WHERE user_id = ?;

-- name: RemoveAllUsersFromCategory :exec
DELETE FROM users_categories
WHERE category_id = ?;

-- name: UserHasCategory :one
SELECT EXISTS(
    SELECT 1 FROM users_categories
    WHERE user_id = ? AND category_id = ?
);

-- name: CountCategoriesByUser :one
SELECT COUNT(*) FROM users_categories
WHERE user_id = ?;

-- name: CountUsersByCategory :one
SELECT COUNT(*) FROM users_categories
WHERE category_id = ?;

-- name: CountUserCategories :one
SELECT COUNT(*) FROM users_categories;

-- name: GetUserCategoriesCreatedAfter :many
SELECT * FROM users_categories
WHERE created_at > ?
ORDER BY created_at DESC;
