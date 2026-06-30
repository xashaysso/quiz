package service

import (
	entities "auth/entities"
	"auth/repository"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
}

func NewAuthService(uRepo repository.UserRepository, sRepo repository.SessionRepository) AuthServiceInterface {
	return &AuthService{
		UserRepo:    uRepo,
		SessionRepo: sRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, username, password string) (entities.User, string, error) {
	if len(username) < 3 {
		return entities.User{}, "", ErrInvalidUsername
	}
	if len(password) < 5 {
		return entities.User{}, "", ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, "", err
	}

	user, err := s.UserRepo.CreateUser(ctx, username, string(hash))
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return entities.User{}, "", ErrUserAlreadyExists
		}
		return entities.User{}, "", err
	}

	slog.Info("new user registered successfully", slog.Int("user_id", user.ID), slog.String("username", user.Username))

	token, err := s.Login(ctx, username, password)
	if err != nil {
		return entities.User{}, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		slog.Warn("login failed: user not found", slog.String("username", username))
		return "", ErrWrongCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		slog.Warn("login failed: wrong password", slog.Int("user_id", user.ID), slog.String("username", user.Username))
		return "", ErrWrongCredentials
	}

	token := uuid.New().String()
	ttl := 24 * time.Hour

	err = s.SessionRepo.Set(ctx, token, user.ID, ttl)
	if err != nil {
		return "", err
	}

	slog.Info("user logged in successfully", slog.Int("user_id", user.ID), slog.String("username", user.Username))

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	err := s.SessionRepo.Delete(ctx, token)
	if err != nil {
		return err
	}

	slog.Info("user logged out successfully")

	return nil
}

func (s *AuthService) CheckSession(ctx context.Context, token string) (int, error) {
	userID, err := s.SessionRepo.Get(ctx, token)
	if err != nil {
		slog.Warn("session expired", slog.Any("err", err))
		return -1, repository.ErrSessionExpired
	}
	return userID, nil
}
