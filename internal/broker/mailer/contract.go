package mailer

import "github.com/mpu-cad/gw-backend-go/internal/models"

type mailer interface {
	SendEmail(models.Gmail) error
}
