package services

import (
	"context"
	"errors"
	"quiz/db/repositories"
	entities "quiz/entities/db"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    repositories.UserRepository
	SessionRepo repositories.SessionRepository
}

func NewAuthService(uRepo repositories.UserRepository, sRepo repositories.SessionRepository) AuthServiceInterface {
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
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			return entities.User{}, "", ErrUserAlreadyExists
		}
		return entities.User{}, "", err
	}

	token, err := s.Login(ctx, username, password)
	if err != nil {
		return entities.User{}, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", ErrWrongCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrWrongCredentials
	}

	token := uuid.New().String()
	ttl := 24 * time.Hour

	err = s.SessionRepo.Set(ctx, token, user.ID, ttl)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	err := s.SessionRepo.Delete(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
