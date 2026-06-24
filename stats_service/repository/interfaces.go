package repository

import "context"

type StatsRepository interface {
	SaveUserStats(ctx context.Context, userID int64, score int) error
	SaveQuizGlobalStats(ctx context.Context, quizID int64, score int) error
}
