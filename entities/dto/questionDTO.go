package dto

import entities "quiz/entities/api"

type CreateQuestionDTO struct {
	Text      string               `json:"text"`
	CorrectID int                  `json:"correct_id"`
	Answers   []entities.AnswerAPI `json:"answers"`
}