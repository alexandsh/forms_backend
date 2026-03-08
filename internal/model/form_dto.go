package model

type CreateFormRequest struct {
	Title     string              `json:"title"`
	Questions []CreateQuestionDTO `json:"questions"`
}

type CreateQuestionDTO struct {
	Type    string   `json:"type"`
	Title   string   `json:"title"`
	Options []string `json:"options,omitempty"`
}
