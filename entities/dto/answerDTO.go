package dto

type CreateAnswerDTO struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}