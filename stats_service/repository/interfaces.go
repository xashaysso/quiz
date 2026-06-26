package repository

import (
	"context"
	"stats/entities"
)

type StatsRepository interface {
	SaveUserStats(ctx context.Context, userID int64, score int) error
	SaveQuizGlobalStats(ctx context.Context, quizID int64, score int) error
	GetUserStats(ctx context.Context, userID int64) (entities.QuizUserStats, error)
	GetQuizStats(ctx context.Context, quizID int64) (entities.QuizGlobalStats, error)
	GetQuizAnalytics(ctx context.Context, quizID int64) (entities.QuizAnalytics, error)
	GetUserAnalytics(ctx context.Context, userID int64) (entities.UserAnalytics, error)
	GetUserLeaderboard(ctx context.Context) ([]entities.UserStats, error)
}
