package course

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type courseUC interface {
	CreateCourse(ctx context.Context, course models.Course) error
}
