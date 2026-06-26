package service

import (
	"context"
	"stats/entities"
)

type StatsServiceInterface interface {
	ProcessQuizPassed(ctx context.Context, event entities.QuizPassedEvent) error
	GetUserStats(ctx context.Context, userID int64) (entities.QuizUserStats, error)
	GetQuizStats(ctx context.Context, quizID int64) (entities.QuizGlobalStats, error)
	GetUserLeaderboard(ctx context.Context) ([]entities.UserStats, error)
	GetQuizAnalytics(ctx context.Context, quizID int64) (entities.QuizAnalytics, error)
	GetUserAnalytics(ctx context.Context, userID int64) (entities.UserAnalytics, error)
}
