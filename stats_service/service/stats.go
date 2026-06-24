package service

import (
	"context"
	"log/slog"
	"stats/entities"
	"stats/repository"
)

type StatsService struct {
	StatsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) *StatsService {
	return &StatsService{StatsRepo: statsRepo}
}

func (s *StatsService) ProcessQuizPassed(ctx context.Context, event entities.QuizPassedEvent) error {
	err := s.StatsRepo.SaveUserStats(ctx, event.UserID, event.Score)
	if err != nil {
		slog.Error("failed to save to repo", slog.Any("err", err))
		return err
	}

	err = s.StatsRepo.SaveQuizGlobalStats(ctx, event.QuizID, event.Score)
	if err != nil {
		slog.Error("failed to save quiz global stats", slog.Any("err", err))
		return err
	}

	slog.Info("stats saved successfully", slog.Int64("user_id", event.UserID), slog.Int64("quiz_id", event.QuizID), slog.Int("score", event.Score))
	return nil
}

func (s *StatsService) GetUserStats(ctx context.Context, userID int64) (entities.QuizUserStats, error) {
	userStats, err := s.StatsRepo.GetUserStats(ctx, userID)
	if err != nil {
		slog.Error("failed to get user stats", slog.Any("err", err))
		return entities.QuizUserStats{}, err
	}

	return userStats, nil
}

func (s *StatsService) GetQuizGlobalStats(ctx context.Context, quizID int64) (entities.QuizGlobalStats, error) {
	quizStats, err := s.StatsRepo.GetQuizGlobalStats(ctx, quizID)
	if err != nil {
		slog.Error("failed to get quiz stats", slog.Any("err", err))
		return entities.QuizGlobalStats{}, err
	}

	return quizStats, nil
}
