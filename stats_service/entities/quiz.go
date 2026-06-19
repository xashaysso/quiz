package entities

import "time"

type QuizPassedEvent struct {
	QuizID         int64     `json:"quiz_id"`
	UserID         int64     `json:"user_id"`
	Score          int       `json:"score"`
	TotalQuestions int       `json:"total_q"`
	PassedAt       time.Time `json:"passed_at"`
}
