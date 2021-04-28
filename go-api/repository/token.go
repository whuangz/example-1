package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/whuangz/go-example/go-api/domain"
)

type tokenRepo struct {
	redis *redis.Client
}

func NewTokenRepo(redisClient *redis.Client) domain.TokenRepository {
	return &tokenRepo{redis: redisClient}
}

func (r *tokenRepo) SetRefreshToken(ctx context.Context, accID string, accessToken string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", accID, accessToken)
	if err := r.redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for accID/accToken: %s/%s: %v\n", accID, accessToken, err)
		return domain.NewInternal()
	}
	return nil
}

func (r *tokenRepo) DeleteRefreshToken(ctx context.Context, accID string, prevAccessToken string) error {
	key := fmt.Sprintf("%s:%s", accID, prevAccessToken)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for accID/accToken: %s/%s: %v\n", accID, prevAccessToken, err)
		return domain.NewInternal()
	}

	return nil
}
