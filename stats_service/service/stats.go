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
