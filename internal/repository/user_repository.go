package repository

import (
	"context"

	"forms/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	query := `
		INSERT INTO users(email, password_hash)
		VALUES ($1, $2)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserRepository) Get(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email=$1
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		email,
	)

	var user model.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
