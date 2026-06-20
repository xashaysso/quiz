package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"quiz/db/repositories"
	entities "quiz/entities/db"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepo struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *SessionRepo {
	return &SessionRepo{
		client: client,
	}
}

func (r *SessionRepo) Get(ctx context.Context, token string) (int, error) {
	key := fmt.Sprintf("session:%s", token)
	userID, err := r.client.Get(ctx, key).Int()
	if err != nil {
		return -1, err
	}
	return userID, nil
}

func (r *SessionRepo) Set(ctx context.Context, token string, userID int, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	err := r.client.Set(ctx, key, userID, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepo) Delete(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepo) SaveQuizSession(ctx context.Context, session entities.QuizSession, ttl time.Duration) error {
	jsonData, err := json.Marshal(session)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("quiz_session:%s", session.SessionID)
	err = r.client.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) GetQuizSession(ctx context.Context, sessionID string) (*entities.QuizSession, error) {
	key := fmt.Sprintf("quiz_session:%s", sessionID)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, repositories.ErrSessionExpired
		}
		return nil, err
	}

	var session entities.QuizSession
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepo) DeleteQuizSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("quiz_session:%s", sessionID)
	return r.client.Del(ctx, key).Err()
}
