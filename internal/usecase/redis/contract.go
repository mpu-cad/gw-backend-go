package redis

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type repos interface {
	SaveRefreshToken(ctx context.Context, userID int, refreshToken string)
	GetIDByRefreshToken(ctx context.Context, refreshToken string) *int
	DeleteRefreshToken(ctx context.Context, refreshToken string)
}

type userRepos interface {
	InsertUser(ctx context.Context, user models.User) (*int, error)
	SelectUserByID(ctx context.Context, id int) (*models.User, error)
}
