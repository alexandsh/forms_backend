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

type GetFormsResponse struct {
	Forms []*GetFormResponse `json:"forms"`
}

type GetFormResponse struct {
	ID        int            `json:"id"`
	Title     string         `json:"title"`
	Questions []*QuestionDTO `json:"questions"`
}

type QuestionDTO struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Options []OptionDTO `json:"options,omitempty"`
}

type OptionDTO struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}
