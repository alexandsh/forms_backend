package service

import (
	"context"

	"forms/internal/model"
	"forms/internal/repository"
)

type FormService struct {
	repo *repository.FormRepository
}

func NewFormService(repo *repository.FormRepository) *FormService {
	return &FormService{repo: repo}
}

func (s *FormService) CreateForm(ctx context.Context, userID int, req model.CreateFormRequest) (int, error) {
	formID, err := s.repo.CreateForm(ctx, userID, req.Title)
	if err != nil {
		return -1, err
	}

	for i, q := range req.Questions {
		qID, err := s.repo.CreateQuestion(ctx, formID, q, i)
		if err != nil {
			return -1, err
		}

		if q.Type == "radio" || q.Type == "checkbox" {
			for j, opt := range q.Options {
				err := s.repo.CreateOption(ctx, qID, opt, j)
				if err != nil {
					return -1, err
				}
			}
		}
	}

	return formID, err
}
