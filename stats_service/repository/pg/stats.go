package pg

import (
	"context"

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
