package entities

import "time"

type QuizPassedEvent struct {
	QuizID   int64     `json:"quiz_id"`
	UserID   int64     `json:"user_id"`
	Score    int       `json:"score"`
	PassedAt time.Time `json:"passed_at"`
}

type QuizUserStats struct {
	UserID             int64     `json:"user_id"`
	TotalQuizzesPassed int       `json:"total_quizzes_passed"`
	TotalScore         int       `json:"total_score"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type QuizGlobalStats struct {
	QuizID    int64     `json:"quiz_id"`
	TotalAtts int       `json:"total_attempts"`
	AccScore  int64     `json:"accumulated_score"`
	UpdatedAt time.Time `json:"updated_at"`
}

type QuizAnalytics struct {
	QuizID    int64     `json:"quiz_id"`
	TotalAtts int       `json:"total_attempts"`
	AccScore  int64     `json:"accumulated_score"`
	UpdatedAt time.Time `json:"updated_at"`
	AvgScore  int       `json:"average_score"`
}
