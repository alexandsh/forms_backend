package repository

import (
	"context"
	"forms/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FormRepository struct {
	pool *pgxpool.Pool
}

func NewFormRepository(pool *pgxpool.Pool) *FormRepository {
	return &FormRepository{pool: pool}
}

func (r *FormRepository) CreateForm(ctx context.Context, userID int, title string) (int, error) {
	query := `
		INSERT INTO forms (user_id, title)
		VALUES ($1, $2)
		RETURNING id
	`

	var formID int

	err := r.pool.QueryRow(ctx, query, userID, title).Scan(&formID)
	if err != nil {
		return -1, err
	}

	return formID, err
}

func (r *FormRepository) CreateQuestion(ctx context.Context, formID int, q model.CreateQuestionDTO, pos int) (int, error) {
	query := `
		INSERT INTO questions (form_id, type, title, position)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var qID int

	err := r.pool.QueryRow(ctx, query, formID, q.Type, q.Title, pos).Scan(&qID)
	if err != nil {
		return -1, err
	}

	return qID, err
}

func (r *FormRepository) CreateOption(ctx context.Context, qID int, opt string, pos int) error {
	query := `
		INSERT INTO options (question_id, value, position)
		VALUES ($1, $2, $3)
	`
	
	_, err := r.pool.Exec(ctx, query, qID, opt, pos)
	
	return err
}