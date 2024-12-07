//go:generate mockgen -source=contract.go -destination mock_test.go -package $GOPACKAGE
package user

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type userUC interface {
	Registration(ctx context.Context, request models.User) (*int, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
}

type redisUC interface {
	CreateRefreshToken(ctx context.Context, id int) string
}
