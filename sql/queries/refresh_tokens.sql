-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = ? LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = 1, updated_at = CURRENT_TIMESTAMP
WHERE token_hash = ?;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked = 1, updated_at = CURRENT_TIMESTAMP
WHERE user_id = ?;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens WHERE expires_at < CURRENT_TIMESTAMP;
