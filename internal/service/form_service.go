package service

import (
	"context"
	"errors"

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

func (s *FormService) GetForm(ctx context.Context, id int) (*model.GetFormResponse, error) {
	return s.repo.GetForm(ctx, id)
}

func (s *FormService) GetForms(ctx context.Context, userID int) (*model.GetFormsResponse, error) {
	return s.repo.GetForms(ctx, userID)
}

func (s *FormService) UpdateForm(ctx context.Context, userID int, formID int, req model.UpdateFormRequest) error {
	ok, err := s.repo.IsFormOwner(ctx, userID, formID)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("forbidden")
	}
	
	return s.repo.UpdateForm(ctx, userID, formID, req)
}