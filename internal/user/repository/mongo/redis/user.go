package redis

import (
	"context"
	"github.com/pt010104/api-golang/internal/models"
)

func (r implRedis) SetSecretKey(ctx context.Context, sessionID string, secretKey string) error {
	return r.redis.Set(ctx, sessionID, secretKey, 0)
}
func (r implRedis) GetSecretKey(ctx context.Context, sessionID string) ([]byte, error) {
	result, err := r.redis.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}
func (redis implRedis) StoreSecretKey(sc models.Scope, secretKey string, ctx context.Context) error {
	return redis.SetSecretKey(ctx, sc.SessionID, secretKey)
}
