package entities

import "time"

type QuizPassedEvent struct {
	QuizID   int64     `json:"quiz_id"`
	UserID   int64     `json:"user_id"`
	Score    int       `json:"score"`
	PassedAt time.Time `json:"passed_at"`
}
