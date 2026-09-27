package repository

import (
	"context"
	"errors"
	"fmt"

	"ghbn/gin-study/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (*models.User, error) {
	query := `
        INSERT INTO users (name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `
	err := r.pool.QueryRow(ctx, query, u.Name, u.Email, u.Password).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`
	u := &models.User{}
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // não achou, não é erro
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`
	u := &models.User{}
	err := r.pool.QueryRow(ctx, query, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, u.Name, u.Email, u.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
