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

func (s *FormService) DeleteForm(ctx context.Context, userID int, formID int) error {
	ok, err := s.repo.IsFormOwner(ctx, userID, formID)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("forbidden")
	}

	return s.repo.DeleteForm(ctx, formID)
}

func (s *FormService) CreateResponse(ctx context.Context, userID int, formID int, req model.CreateResponseRequest) error {
	form, err := s.repo.GetForm(ctx, formID)
	if err != nil {
		return err
	}

	if form == nil {
		return errors.New("form not found")
	}

	questions := make(map[int]string)
	for _, q := range form.Questions {
		questions[q.ID] = q.Type
	}

	options := make(map[int]int)
	for _, q := range form.Questions {
		for _, opt := range q.Options {
			options[opt.ID] = q.ID
		}
	}

	radioCheck := make(map[int]bool)

	for _, ans := range req.Answers {
		qType, exists := questions[ans.QuestionID]
		if !exists {
			return errors.New("question not belong form")
		}

		switch qType {
		case "text":
			if ans.TextValue == nil {
				return errors.New("text answer expected")
			}

			if ans.OptionID != nil {
				return errors.New("text answer cant have options")
			}

		case "radio":
			if ans.OptionID == nil {
				return errors.New("option expected")
			}

			if radioCheck[ans.QuestionID] {
				return errors.New("radio answer cant have more than 1 option picked")
			}

			radioCheck[ans.QuestionID] = true

			qID := options[*ans.OptionID]
			if qID != ans.QuestionID {
				return errors.New("options not belong question")
			}

		case "checkbox":
			if ans.OptionID == nil {
				return errors.New("option expected")
			}

			qID := options[*ans.OptionID]
			if qID != ans.QuestionID {
				return errors.New("options not belong question")
			}

		default:
			return errors.New("unknown type of question")
		}
	}

	return s.repo.CreateResponse(ctx, formID, userID, req)
}
