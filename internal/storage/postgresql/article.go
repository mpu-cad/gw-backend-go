package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type ArticleRepos struct {
	db *pgxpool.Pool
}

func NewArticleRepos(db *pgxpool.Pool) *ArticleRepos {
	return &ArticleRepos{
		db: db,
	}
}

func (r *CourseRepos) InsertArticle(ctx context.Context, articleId, courseId int) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO course_articles (course_id, article_id) 
		VALUES ($1, $2)
	`, courseId, articleId)
	return err
}

func (r *CourseRepos) SelectAllArticlesByCourseID(ctx context.Context, courseId int) ([]models.Article, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.title, a.text, array_agg(at.tag_name) AS tags 
		FROM course_articles ca 
		JOIN articles a ON a.id = ca.article_id 
		LEFT JOIN article_tags at ON a.id = at.article_id 
		WHERE ca.course_id = $1 
		GROUP BY a.id
	`, courseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Text, &article.Tags); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (r *CourseRepos) UpdateArticle(ctx context.Context, article models.Article) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE articles 
		SET title = $1, text = $2 
		WHERE id = $3
	`, article.Title, article.Text, article.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM article_tags WHERE article_id = $1
	`, article.ID)
	if err != nil {
		return err
	}

	for _, tag := range article.Tags {
		_, err = tx.Exec(ctx, `
			INSERT INTO article_tags (article_id, tag_name) 
			VALUES ($1, $2)
		`, article.ID, tag)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *CourseRepos) DeleteArticle(ctx context.Context, articleId int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM articles WHERE id = $1`, articleId)
	return err
}
