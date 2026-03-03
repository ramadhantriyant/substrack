package database

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	Revoked   int64      `json:"revoked"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CreateRefreshTokenParams struct {
	UserID    int64     `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	const query = `
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES (?, ?, ?)
RETURNING id, user_id, token_hash, expires_at, revoked, created_at, updated_at`

	row := q.db.QueryRowContext(ctx, query, arg.UserID, arg.TokenHash, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TokenHash,
		&i.ExpiresAt,
		&i.Revoked,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func (q *Queries) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error) {
	const query = `
SELECT id, user_id, token_hash, expires_at, revoked, created_at, updated_at
FROM refresh_tokens
WHERE token_hash = ? LIMIT 1`

	row := q.db.QueryRowContext(ctx, query, tokenHash)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TokenHash,
		&i.ExpiresAt,
		&i.Revoked,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func (q *Queries) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	const query = `
UPDATE refresh_tokens
SET revoked = 1, updated_at = CURRENT_TIMESTAMP
WHERE token_hash = ?`

	_, err := q.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (q *Queries) RevokeAllUserRefreshTokens(ctx context.Context, userID int64) error {
	const query = `
UPDATE refresh_tokens
SET revoked = 1, updated_at = CURRENT_TIMESTAMP
WHERE user_id = ?`

	_, err := q.db.ExecContext(ctx, query, userID)
	return err
}

func (q *Queries) DeleteExpiredRefreshTokens(ctx context.Context) error {
	const query = `DELETE FROM refresh_tokens WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := q.db.ExecContext(ctx, query)
	return err
}
