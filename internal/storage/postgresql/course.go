package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type CourseRepos struct {
	db *pgxpool.Pool
}

func NewCourseRepos(db *pgxpool.Pool) *CourseRepos {
	return &CourseRepos{
		db: db,
	}
}

func (c *CourseRepos) CreateCourse(ctx context.Context, course models.Course) error {
	// c.db.Exec(ctx, `insert into courses (id, ) values ($1, $2, $3)`)

	return nil
}
