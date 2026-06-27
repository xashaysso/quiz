package repository

import (
	"auth/entities"
	"context"
	"time"
)

type SessionRepository interface {
	Get(ctx context.Context, token string) (int, error)
	Set(ctx context.Context, token string, userID int, ttl time.Duration) error
	Delete(ctx context.Context, token string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password_hash string) (entities.User, error)
	GetByUsername(ctx context.Context, username string) (entities.User, error)
}
