package pg

import (
	"context"
	"fmt"
	entities "quiz/entities/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgUserRepo struct {
	Pool *pgxpool.Pool
}

func NewUserRepo(p *pgxpool.Pool) *PgUserRepo {
	return &PgUserRepo{
		Pool: p,
	}
}

func (r *PgUserRepo) CreateUser(ctx context.Context, username string, password_hash string) (entities.User, error) {
	var User entities.User
	err := r.Pool.QueryRow(ctx, `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, created_at`, username, password_hash).Scan(&User.ID, &User.Username, &User.CreatedAt)
	if err != nil {
		return entities.User{}, err
	}

	return User, nil
}

func (r *PgUserRepo) GetByUsername(ctx context.Context, username string) (entities.User, error) {
	var User entities.User
	err := r.Pool.QueryRow(ctx, "SELECT id, username, password_hash FROM users WHERE username = $1", username).Scan(&User.ID, &User.Username, &User.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.User{}, fmt.Errorf("user not found")
		}
		return entities.User{}, err
	}

	return User, nil
}
