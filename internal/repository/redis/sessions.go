package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sagemyrage/code-quality-expert-system/internal/repository"
)

const sessionKeyPrefix = "session:"

type SessionRepository struct {
	client *goredis.Client
	ttl    time.Duration
}

func NewSessionRepository(client *goredis.Client, ttl time.Duration) *SessionRepository {
	return &SessionRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *SessionRepository) Create(ctx context.Context, userID int64) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}
	key := sessionKey(sessionID)
	value := strconv.FormatInt(userID, 10)

	err = r.client.Set(ctx, key, value, r.ttl).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (r *SessionRepository) GetUserID(ctx context.Context, sessionID string) (int64, error) {
	key := sessionKey(sessionID)

	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, goredis.Nil) {
		return 0, repository.ErrSessionNotFound
	}
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	key := sessionKey(sessionID)

	return r.client.Del(ctx, key).Err()
}

func generateSessionID() (string, error) {
	s := make([]byte, 32)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	sessionID := hex.EncodeToString(s)
	return sessionID, nil
}

func sessionKey(sessionID string) string {
	return sessionKeyPrefix + sessionID
}
