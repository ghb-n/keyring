package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"ghbn/gin-study/internal/models"
	"ghbn/gin-study/internal/repository"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredRefreshToken = errors.New("expired refresh token")
	ErrRevokedRefreshToken = errors.New("revoked refresh token")
)

const (
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days
	refreshTokenBytes    = 32                 // 256 bits
)

type RefreshService struct {
	repo *repository.RefreshRepository
}

func NewRefreshService(repo *repository.RefreshRepository) *RefreshService {
	return &RefreshService{repo: repo}
}

// Generate creates a new refresh token for a user.
// Returns the plain token (to send to the client) and stores the hash.
func (s *RefreshService) Generate(ctx context.Context, userID int) (string, error) {
	// 1. Generate random bytes
	raw := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate random: %w", err)
	}

	// 2. Encode as base64url (URL-safe, no padding)
	plain := base64.RawURLEncoding.EncodeToString(raw)

	// 3. Hash it
	hash := sha256.Sum256([]byte(plain))

	// 4. Store the hash + expiry
	token := &models.RefreshToken{
		UserID:    userID,
		TokenHash: base64.RawURLEncoding.EncodeToString(hash[:]),
		ExpiresAt: time.Now().Add(RefreshTokenDuration),
	}
	if err := s.repo.Create(ctx, token); err != nil {
		return "", fmt.Errorf("create refresh token: %w", err)
	}

	// 5. Return the plain token
	return plain, nil
}

// Validate checks if a refresh token is valid. Returns the token record.
func (s *RefreshService) Validate(ctx context.Context, plain string) (*models.RefreshToken, error) {
	// 1. Hash the incoming token
	hash := sha256.Sum256([]byte(plain))
	hashStr := base64.RawURLEncoding.EncodeToString(hash[:])

	// 2. Look up by hash
	token, err := s.repo.FindByHash(ctx, hashStr)
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	if token == nil {
		return nil, ErrInvalidRefreshToken
	}

	// 3. Check expiry
	if time.Now().After(token.ExpiresAt) {
		return nil, ErrExpiredRefreshToken
	}

	// 4. Check revocation
	if token.RevokedAt != nil {
		return nil, ErrRevokedRefreshToken
	}

	return token, nil
}

// Rotate validates an old token, revokes it, and generates a new one.
// Returns the new plain token and the userID.
func (s *RefreshService) Rotate(ctx context.Context, oldPlain string) (string, int, error) {
	// 1. Validate old token
	old, err := s.Validate(ctx, oldPlain)
	if err != nil {
		return "", 0, err
	}

	// 2. Revoke old token
	if err := s.repo.Revoke(ctx, old.ID); err != nil {
		return "", 0, fmt.Errorf("revoke old token: %w", err)
	}

	// 3. Generate new token
	newPlain, err := s.Generate(ctx, old.UserID)
	if err != nil {
		return "", 0, fmt.Errorf("generate new token: %w", err)
	}

	return newPlain, old.UserID, nil
}

// RevokeAll revokes all refresh tokens for a user (e.g. on password change).
func (s *RefreshService) RevokeAll(ctx context.Context, userID int) error {
	return s.repo.RevokeAllForUser(ctx, userID)
}
