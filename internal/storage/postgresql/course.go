package postgresql

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/mpu-cad/gw-backend-go/internal/logger"
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

func (r *CourseRepos) InsertCourse(ctx context.Context, course models.Course) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	row := tx.QueryRow(
		ctx,
		`insert into courses (
                     title,
                     poster, 
                     description) 
		values ($1, $2, $3)
		returning id`,
		course.Title,
		course.Poster,
		course.Description,
	)

	var id int
	err = row.Scan(&id)
	if err != nil {
		return errors.Wrap(err, "scan id")
	}

	for _, tag := range course.Tags {
		cmd, err := tx.Exec(
			ctx,
			`insert into course_tags (
                         course_id,
                         tag_name) 
		values ($1,  $2)`,
			id,
			tag,
		)

		if err != nil {
			logger.Log.Error("can not insert tags")
		}

		if !cmd.Insert() {
			return errors.New("no insert to course_tags")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit")
	}

	return nil
}

func (r *CourseRepos) SelectAllCourses(ctx context.Context, limit, offset int) ([]models.Course, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.title, r.poster, r.description, array_agg(ct.tag_name) AS tags 
		FROM courses r 
		LEFT JOIN course_tags ct ON r.id = ct.course_id 
		GROUP BY r.id 
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "select curses")
	}

	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err = rows.Scan(&course.ID, &course.Title, &course.Poster, &course.Description, &course.Tags); err != nil {
			return nil, errors.Wrap(err, "scan course")
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (r *CourseRepos) SelectCourseByID(ctx context.Context, id int) (*models.Course, error) {
	var course models.Course

	err := r.db.QueryRow(ctx, `
		SELECT r.id, r.title, r.poster, r.description, array_agg(ct.tag_name) AS tags 
		FROM courses r 
		LEFT JOIN course_tags ct ON r.id = ct.course_id 
		WHERE r.id = $1 
		GROUP BY r.id
	`, id).Scan(&course.ID, &course.Title, &course.Poster, &course.Description, &course.Tags)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}

		return nil, errors.Wrap(err, "delete courses")
	}

	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.title, a.text, array_agg(at.tag_name) AS tags 
		FROM course_articles ca 
		JOIN articles a ON a.id = ca.article_id 
		LEFT JOIN article_tags at ON a.id = at.article_id 
		WHERE ca.course_id = $1 
		GROUP BY a.id
	`, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return &course, nil
		}

		return nil, errors.Wrap(err, "select courses article")
	}
	defer rows.Close()

	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Text, &article.Tags); err != nil {
			return nil, errors.Wrap(err, "scan course")
		}

		course.Articles = append(course.Articles, article)
	}

	return &course, nil
}

func (r *CourseRepos) UpdateCourse(ctx context.Context, course models.Course) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			logger.Log.Errorf("tx rollback, err: %v", err)
		}
	}()

	_, err = tx.Exec(ctx, `
		UPDATE courses 
		SET title = $1, description = $2, poster = $3 
		WHERE id = $4
	`, course.Title, course.Description, course.Poster, course.ID)
	if err != nil {
		return errors.Wrap(err, "update courses")
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM course_tags WHERE course_id = $1
	`, course.ID)
	if err != nil {
		return errors.Wrap(err, "delete tags courses")
	}

	for _, tag := range course.Tags {
		_, err = tx.Exec(ctx, `
			INSERT INTO course_tags (course_id, tag_name) 
			VALUES ($1, $2)
		`, course.ID, tag)
		if err != nil {
			return errors.Wrap(err, "insert tags courses")
		}
	}

	return tx.Commit(ctx)
}

func (r *CourseRepos) DeleteCourse(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM courses WHERE id = $1`, id)
	return errors.Wrap(err, "delete courses")
}
