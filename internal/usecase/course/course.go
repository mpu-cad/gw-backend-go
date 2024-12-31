package course

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type UCCourse struct {
	courseRepos
}

func NewUCCourse(courseRepos courseRepos) *UCCourse {
	return &UCCourse{
		courseRepos: courseRepos,
	}
}

func (c *UCCourse) CreateCourse(ctx context.Context, course models.Course) error {
	//nolint:wrapcheck
	return c.courseRepos.InsertCourse(ctx, course)
}

func (c *UCCourse) GetAllCourses(ctx context.Context, limit, offset int) ([]models.Course, error) {
	if limit == 0 {
		limit = 20
	}

	courses, err := c.courseRepos.SelectAllCourses(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "select all courses")
	}

	return courses, nil
}

func (c *UCCourse) GetCourseByID(ctx context.Context, id int) (*models.Course, error) {
	//nolint:wrapcheck
	return c.courseRepos.SelectCourseByID(ctx, id)
}

func (c *UCCourse) UpdateCourse(ctx context.Context, course models.Course) error {
	//nolint:wrapcheck
	return c.courseRepos.UpdateCourse(ctx, course)
}

func (c *UCCourse) DeleteCourse(ctx context.Context, id int) error {
	//nolint:wrapcheck
	return c.courseRepos.DeleteCourse(ctx, id)
}
