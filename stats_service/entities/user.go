package entities

import "time"

type UserStats struct {
	UserID             int64 `json:"user_id"`
	TotalScore         int   `json:"total_score"`
	TotalQuizzesPassed int   `json:"total_quizzes_passed"`
}

type UserAnalytics struct {
	UserID             int64     `json:"user_id"`
	TotalScore         int       `json:"total_score"`
	TotalQuizzesPassed int       `json:"total_quizzes_passed"`
	UpdatedAt          time.Time `json:"updated_at"`
	AvgScore           int       `json:"average_score"`
}
