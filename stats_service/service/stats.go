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

func (s *StatsService) GetQuizStats(ctx context.Context, quizID int64) (entities.QuizGlobalStats, error) {
	quizStats, err := s.StatsRepo.GetQuizStats(ctx, quizID)
	if err != nil {
		slog.Error("failed to get quiz stats", slog.Any("err", err))
		return entities.QuizGlobalStats{}, err
	}

	return quizStats, nil
}

func (s *StatsService) GetUserLeaderboard(ctx context.Context) ([]entities.UserStats, error) {
	leaderboard, err := s.StatsRepo.GetUserLeaderboard(ctx)
	if err != nil {
		slog.Error("failed to get user leaderboard", slog.Any("err", err))
		return nil, err
	}
	return leaderboard, nil
}

func (s *StatsService) GetQuizAnalytics(ctx context.Context, quizID int64) (entities.QuizAnalytics, error) {
	qAnalytics, err := s.StatsRepo.GetQuizAnalytics(ctx, quizID)
	if err != nil {
		slog.Error("failed to get quiz analytics", slog.Any("err", err))
		return entities.QuizAnalytics{}, err
	}
	return qAnalytics, nil
}

func (s *StatsService) GetUserAnalytics(ctx context.Context, userID int64) (entities.UserAnalytics, error) {
	uAnalytics, err := s.StatsRepo.GetUserAnalytics(ctx, userID)
	if err != nil {
		slog.Error("failed to get user analytics", slog.Any("err", err))
		return entities.UserAnalytics{}, err
	}
	return uAnalytics, nil
}
