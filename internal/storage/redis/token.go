package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/mpu-cad/gw-backend-go/internal/entity"
)

type TokenRepos struct {
	cli *redis.Client
}

func NewTokenRepos(redisCli *redis.Client) *TokenRepos {
	return &TokenRepos{
		cli: redisCli,
	}
}

func (t *TokenRepos) SaveRefreshToken(ctx context.Context, userID int, refreshToken string) {
	t.cli.Set(ctx, strconv.Itoa(userID), refreshToken, entity.ExpiresDayRefreshToken)
}

func (t *TokenRepos) GetIDByRefreshToken(ctx context.Context, refreshToken string) *int {
	get := t.cli.Get(ctx, refreshToken)
	if get == nil {
		return nil
	}

	id, err := get.Int()
	if err != nil {
		return nil
	}

	return &id
}

func (t *TokenRepos) DeleteRefreshToken(ctx context.Context, refreshToken string) {
	t.cli.Del(ctx, refreshToken)
}
