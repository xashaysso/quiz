package redis

import (
	"context"
	"fmt"
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
