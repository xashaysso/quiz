package service

import (
	entities "auth/entities"
	"context"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, username, password string) (entities.User, string, error)
	Login(ctx context.Context, username, password string) (string, error)
	Logout(ctx context.Context, token string) error
	CheckSession(ctx context.Context, token string) (int, error)
}
