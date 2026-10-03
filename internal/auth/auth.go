package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

// Claims é o payload do JWT. Contém os dados do usuário e os
// registered claims (exp, iat, etc).
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken cria um JWT assinado com HS256.
func GenerateToken(userID int, email, secret string) (string, error) {
	// 1. Cria as claims
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "key-ring",
		},
	}

	// 2. Cria o token com as claims e o método de assinatura
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Assina com o secret e retorna a string
	return token.SignedString([]byte(secret))
}

// ValidateToken parseia e valida um JWT, retornando as claims se válido.
func ValidateToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		// SEGURANÇA: garante que o algoritmo é HMAC (HS256).
		// Sem isso, um atacante pode forjar um token com alg: none
		// ou RS256 e enganar a validação.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
