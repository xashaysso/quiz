package service

import (
	"context"
	"stats/entities"
)

type StatsServiceInterface interface {
	ProcessQuizPassed(ctx context.Context, event entities.QuizPassedEvent) error
	GetUserStats(ctx context.Context, userID int64) (entities.QuizUserStats, error)
	GetQuizGlobalStats(ctx context.Context, quizID int64) (entities.QuizGlobalStats, error)
}
