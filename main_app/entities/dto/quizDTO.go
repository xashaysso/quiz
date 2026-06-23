package dto

import (
	"time"
)

type CreateQuizDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateQuizDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type QuizResponse struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Questions   []QuestionResponse `json:"questions"`
}

type QuizPassedEvent struct {
	QuizID   int64     `json:"quiz_id"`
	UserID   int64     `json:"user_id"`
	Score    int       `json:"score"`
	PassedAt time.Time `json:"passed_at"`
}
