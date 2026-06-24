package repository

import (
	"context"
	"stats/entities"
)

type StatsRepository interface {
	SaveUserStats(ctx context.Context, userID int64, score int) error
	SaveQuizGlobalStats(ctx context.Context, quizID int64, score int) error
	GetUserStats(ctx context.Context, userID int64) (entities.QuizUserStats, error)
	GetQuizGlobalStats(ctx context.Context, quizID int64) (entities.QuizGlobalStats, error)
	GetQuizAnalytics(ctx context.Context, quizID int64) (entities.QuizAnalytics, error)
}
