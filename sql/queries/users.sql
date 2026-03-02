-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: ListUsersWithLimit :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateUser :one
INSERT INTO users (
    email, password, name
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = ?,
    name = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUserName :one
UPDATE users
SET name = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET email = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: DeleteUserByEmail :exec
DELETE FROM users
WHERE email = ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = ?);

-- name: UserExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ?);

-- name: SearchUsers :many
SELECT * FROM users
WHERE name LIKE ? OR email LIKE ?
ORDER BY name;

-- name: GetUsersCreatedAfter :many
SELECT * FROM users
WHERE created_at > ?
ORDER BY created_at DESC;

-- name: GetUsersUpdatedAfter :many
SELECT * FROM users
WHERE updated_at > ?
ORDER BY updated_at DESC;
