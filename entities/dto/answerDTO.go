package dto

import entities "quiz/entities/db"

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

type AnswerPublicResponse struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Text       string `json:"text"`
}

func NewAnswerResponse(a entities.Answer) AnswerPublicResponse {
	return AnswerPublicResponse{
		ID:         a.ID,
		QuestionID: a.QuestionID,
		Text:       a.Text,
	}
}
