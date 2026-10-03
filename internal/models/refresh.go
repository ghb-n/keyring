package models

import "time"

type RefreshToken struct {
	ID        int
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}
