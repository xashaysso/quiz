package dto

import entities "quiz/entities/api"

type CreateQuestionDTO struct {
	Text      string               `json:"text"`
	Answers   []entities.AnswerAPI `json:"answers"`
}

type UpdateQuestionDTO struct {
	Text *string `json:"text"`
	NewCorrectID *int `json:"correct_answer_id"`
}