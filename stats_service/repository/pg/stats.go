package pg

import (
	"context"
	"errors"
	"stats/entities"
	"stats/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStatsRepo struct {
	Pool *pgxpool.Pool
}

func NewStatsRepo(p *pgxpool.Pool) *PgStatsRepo {
	return &PgStatsRepo{Pool: p}
}

func (r *PgStatsRepo) SaveUserStats(ctx context.Context, userID int64, score int) error {
	query := `INSERT INTO user_stats (user_id, total_quizzes_passed, total_score, updated_at) VALUES ($1, 1, $2, NOW())
			ON CONFLICT (user_id) DO UPDATE
			SET
				total_quizzes_passed = user_stats.total_quizzes_passed + 1,
				total_score = user_stats.total_score + EXCLUDED.total_score,
				updated_at = NOW()
			`
	_, err := r.Pool.Exec(ctx, query, userID, score)
	return err
}

func (r *PgStatsRepo) SaveQuizGlobalStats(ctx context.Context, quizID int64, score int) error {
	query := `INSERT INTO quiz_global_stats (quiz_id, total_attempts, accumulated_score, updated_at) VALUES ($1, 1, $2, NOW())
			ON CONFLICT (quiz_id) DO UPDATE
			SET
				total_attempts = quiz_global_stats.total_attempts + 1,
				accumulated_score = quiz_global_stats.accumulated_score + EXCLUDED.accumulated_score,
				updated_at = NOW()
			`
	_, err := r.Pool.Exec(ctx, query, quizID, score)
	return err
}

func (r *PgStatsRepo) GetUserStats(ctx context.Context, userID int64) (entities.QuizUserStats, error) {
	var quizUserStats entities.QuizUserStats
	err := r.Pool.QueryRow(ctx, `SELECT user_id, total_quizzes_passed, total_score, updated_at FROM user_stats WHERE user_id = $1`, userID).Scan(&quizUserStats.UserID, &quizUserStats.TotalQuizzesPassed, &quizUserStats.TotalScore, &quizUserStats.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.QuizUserStats{}, repository.ErrRecordNotFound
		}
		return entities.QuizUserStats{}, err
	}

	return quizUserStats, nil
}

func (r *PgStatsRepo) GetQuizStats(ctx context.Context, quizID int64) (entities.QuizGlobalStats, error) {
	var quizGlobalStats entities.QuizGlobalStats
	err := r.Pool.QueryRow(ctx, `SELECT quiz_id, total_attempts, accumulated_score, updated_at FROM quiz_global_stats WHERE quiz_id = $1`, quizID).Scan(&quizGlobalStats.QuizID, &quizGlobalStats.TotalAtts, &quizGlobalStats.AccScore, &quizGlobalStats.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.QuizGlobalStats{}, repository.ErrRecordNotFound
		}
		return entities.QuizGlobalStats{}, err
	}

	return quizGlobalStats, nil
}

func (r *PgStatsRepo) GetQuizAnalytics(ctx context.Context, quizID int64) (entities.QuizAnalytics, error) {
	var quizAnalytics entities.QuizAnalytics
	query := `SELECT quiz_id, total_attempts, accumulated_score, updated_at,
								CASE
									WHEN total_attempts > 0 THEN ROUND(accumulated_score::NUMERIC / total_attempts, 2)
									ELSE 0
								END as average_score
				FROM quiz_global_stats WHERE quiz_id = $1`
	err := r.Pool.QueryRow(ctx, query, quizID).Scan(&quizAnalytics.QuizID, &quizAnalytics.TotalAtts, &quizAnalytics.AccScore, &quizAnalytics.UpdatedAt, &quizAnalytics.AvgScore)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.QuizAnalytics{}, repository.ErrRecordNotFound
		}
		return entities.QuizAnalytics{}, err
	}

	return quizAnalytics, nil
}

func (r *PgStatsRepo) GetUserAnalytics(ctx context.Context, userID int64) (entities.UserAnalytics, error) {
	var userAnalytics entities.UserAnalytics
	query := `SELECT user_id, total_quizzes_passed, total_score, updated_at,
								CASE
									WHEN total_quizzes_passed > 0 THEN ROUND(total_score::NUMERIC / total_quizzes_passed, 2)
									ELSE 0
								END as average_score
				FROM user_stats WHERE user_id = $1`
	err := r.Pool.QueryRow(ctx, query, userID).Scan(&userAnalytics.UserID, &userAnalytics.TotalQuizzesPassed, &userAnalytics.TotalScore, &userAnalytics.UpdatedAt, &userAnalytics.AvgScore)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.UserAnalytics{}, repository.ErrRecordNotFound
		}
		return entities.UserAnalytics{}, err
	}

	return userAnalytics, nil
}

func (r *PgStatsRepo) GetUserLeaderboard(ctx context.Context) ([]entities.UserStats, error) {
	rows, err := r.Pool.Query(ctx, `SELECT user_id, total_score, total_quizzes_passed FROM user_stats ORDER BY total_score DESC LIMIT 10`)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	var leaderboard []entities.UserStats
	for rows.Next() {
		var userStats entities.UserStats
		err := rows.Scan(&userStats.UserID, &userStats.TotalScore, &userStats.TotalQuizzesPassed)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, userStats)
	}
	return leaderboard, nil
}
