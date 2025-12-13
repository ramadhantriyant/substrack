-- name: GetCategory :one
SELECT * FROM categories
WHERE id = ? LIMIT 1;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = ? LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name;

-- name: ListCategoriesWithLimit :many
SELECT * FROM categories
ORDER BY name
LIMIT ? OFFSET ?;

-- name: CreateCategory :one
INSERT INTO categories (
    name, description
) VALUES (
    ?, ?
)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name = ?,
    description = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateCategoryName :one
UPDATE categories
SET name = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateCategoryDescription :one
UPDATE categories
SET description = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?;

-- name: DeleteCategoryByName :exec
DELETE FROM categories
WHERE name = ?;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories;

-- name: CategoryExists :one
SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?);

-- name: CategoryExistsByName :one
SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?);

-- name: SearchCategories :many
SELECT * FROM categories
WHERE name LIKE ? OR description LIKE ?
ORDER BY name;

-- name: GetCategoriesCreatedAfter :many
SELECT * FROM categories
WHERE created_at > ?
ORDER BY created_at DESC;

-- name: GetCategoriesUpdatedAfter :many
SELECT * FROM categories
WHERE updated_at > ?
ORDER BY updated_at DESC;
