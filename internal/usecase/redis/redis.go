package redis

import (
	"context"
	"fmt"

	"github.com/mpu-cad/gw-backend-go/internal/entity"
	"github.com/mpu-cad/gw-backend-go/internal/middleware/token"
)

type UCRedis struct {
	repos repos
	user  userRepos
}

func NewUCRepos(repository repos) *UCRedis {
	return &UCRedis{
		repos: repository,
	}
}

func (r *UCRedis) CreateRefreshToken(ctx context.Context, id int) string {
	refreshToken := token.GenerateRefreshToken()

	r.repos.SaveRefreshToken(ctx, id, refreshToken)

	return refreshToken
}

func (r *UCRedis) CompareRefreshToken(ctx context.Context, token string) error {
	id := r.repos.GetIDByRefreshToken(ctx, token)
	if id == nil {
		return entity.ErrRefreshTokenExpire
	}

	_, err := r.user.SelectUserByID(ctx, *id)
	if err != nil {
		return fmt.Errorf("user not found, err: %w", err)
	}

	return nil
}
