//go:generate mockgen -source=contract.go -destination mock_test.go -package $GOPACKAGE
package user

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type userRepos interface {
	InsertUser(ctx context.Context, user models.User) (*int, error)
	SelectUserByID(ctx context.Context, id int) (*models.User, error)
	SelectUserByLogin(ctx context.Context, email string) (*models.User, error)
	ConfirmEmail(ctx context.Context, userID int) error
}

type mailer interface {
	SendEmail(gmail models.Gmail) error
}

type redis interface {
	SaveUsersRegistrationCode(ctx context.Context, code string, userID int)
	GetUsersRegistrationCode(ctx context.Context, userID int) string
}
