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

func (r *FormRepository) GetForm(ctx context.Context, id int) (*model.GetFormResponse, error) {
	query := `
		SELECT
			f.id AS form_id,
			f.title AS form_title,

			q.id AS question_id,
			q.type,
			q.title AS question_title,

			o.id AS option_id,
			o.value AS option_value

		FROM forms f
		LEFT JOIN questions q ON q.form_id = f.id
		LEFT JOIN options o ON o.question_id = q.id

		WHERE f.id = $1

		ORDER BY
			q.position,
			o.position;
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	form := &model.GetFormResponse{}
	qMap := map[int]*model.QuestionDTO{}

	for rows.Next() {
		var (
			formID int
			formTitle string

			qID *int
			qType *string
			qTitle *string

			optID *int
			optValue *string
		)
		
		err := rows.Scan(
			&formID,
			&formTitle,
			&qID,
			&qType,
			&qTitle,
			&optID,
			&optValue,
		)
		if err != nil {
			return nil, err
		}

		form.ID = formID
		form.Title = formTitle

		q, exists := qMap[*qID]
		if !exists {
			q = &model.QuestionDTO{
				ID: *qID,
				Type: *qType,
				Title: *qTitle,
				Options: []model.OptionDTO{},
			}

			qMap[*qID] = q
			form.Questions = append(form.Questions, q)
		}
	
		if optID != nil {
			q.Options = append(q.Options, model.OptionDTO{
				ID: *optID,
				Value: *optValue,
			})
		}
	}

	return form, nil
}