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

type UpdateFormRequest struct {
	Title     *string                  `json:"title"`
	Questions *[]UpdateQuestionRequest `json:"questions"`
}

type UpdateQuestionRequest struct {
	ID       *int                   `json:"id"`
	Type     *string                `json:"type"`
	Title    *string                `json:"title"`
	Options  *[]UpdateOptionRequest `json:"options"`
	Position *int                   `json:"position"`
}

type UpdateOptionRequest struct {
	ID       *int    `json:"id"`
	Value    *string `json:"value"`
	Position *int    `json:"position"`
}

type CreateResponseRequest struct {
	Answers []CreateAnswerDTO `json:"answers"`
}

type CreateAnswerDTO struct {
	QuestionID int     `json:"question_id"`
	OptionID   *int    `json:"option_id,omitempty"`
	TextValue  *string `json:"text_value,omitempty"`
}
