package course

import (
	"context"

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
	return nil
}
