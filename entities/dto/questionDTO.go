package dto

import entities "quiz/entities/db"

type CreateQuestionDTO struct {
	Text    string            `json:"text"`
	Answers []CreateAnswerDTO `json:"answers"`
}

type UpdateQuestionDTO struct {
	Text         *string `json:"text"`
	NewCorrectID *int    `json:"correct_answer_id"`
}

type QuestionResponse struct {
	ID      int                    `json:"id"`
	Text    string                 `json:"text"`
	Answers []AnswerPublicResponse `json:"answers"`
}

func NewQuestionResponse(q entities.Question, answers []entities.Answer) QuestionResponse {
	answersDTO := make([]AnswerPublicResponse, len(answers))
	for i, a := range answers {
		answersDTO[i] = AnswerPublicResponse{
			ID:   a.ID,
			Text: a.Text,
		}
	}
	return QuestionResponse{
		ID:      q.ID,
		Text:    q.Text,
		Answers: answersDTO,
	}
}
