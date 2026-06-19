package service

import (
	"context"
	"stats/entities"
)

type StatsServiceInterface interface {
	ProcessQuizPassed(ctx context.Context, event entities.QuizPassedEvent) error
}
