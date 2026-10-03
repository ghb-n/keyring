package repository

import (
	"context"
	"errors"
	"fmt"

	"ghbn/gin-study/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshRepository(pool *pgxpool.Pool) *RefreshRepository {
	return &RefreshRepository{pool: pool}
}

func (r *RefreshRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	return r.pool.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *RefreshRepository) FindByHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	query := `
        SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
        FROM refresh_tokens
        WHERE token_hash = $1
    `
	t := &models.RefreshToken{}
	err := r.pool.QueryRow(ctx, query, hash).Scan(
		&t.ID, &t.UserID, &t.TokenHash,
		&t.ExpiresAt, &t.CreatedAt, &t.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	return t, nil
}

func (r *RefreshRepository) Revoke(ctx context.Context, id int) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (r *RefreshRepository) RevokeAllForUser(ctx context.Context, userID int) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("revoke all refresh tokens: %w", err)
	}
	return nil
}
