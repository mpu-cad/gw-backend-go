package course

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type courseUC interface {
	CreateCourse(ctx context.Context, course models.Course) error
	GetAllCourses(ctx context.Context, limit, offset int) ([]models.Course, error)
	GetCourseByID(ctx context.Context, id int) (*models.Course, error)
	UpdateCourse(ctx context.Context, course models.Course) error
	DeleteCourse(ctx context.Context, id int) error
}
