package dto

type CreateAnswerDTO struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type UpdateAnswerDTO struct {
	Text         *string `json:"text"`
	NewCorrectID *int    `json:"correct_id"`
}

type CheckAnswerDTO struct {
	AnswerID int `json:"answer_id"`
}

type AnswerResponse struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}