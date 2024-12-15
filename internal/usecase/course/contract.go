package course

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type courseRepos interface {
	InsertCourse(ctx context.Context, course models.Course) error
	SelectAllCourses(ctx context.Context, limit, offset int) ([]models.Course, error)
	SelectCourseByID(ctx context.Context, id int) (*models.Course, error)
	UpdateCourse(ctx context.Context, course models.Course) error
	DeleteCourse(ctx context.Context, id int) error
}
