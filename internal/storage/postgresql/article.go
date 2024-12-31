package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/mpu-cad/gw-backend-go/internal/logger"
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

func (r *ArticleRepos) InsertArticle(ctx context.Context, courseId int, article models.Article) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			logger.Log.Errorf("tx rollback, err: %v", err)
		}
	}()

	if _, err = tx.Exec(ctx, `
		INSERT INTO course_articles (course_id, article_id) 
		VALUES ($1, $2)
	`, courseId, article.ID); err != nil {
		return errors.Wrap(err, "insert article")
	}

	if _, err = tx.Exec(ctx, `
		INSERT INTO articles (id, title, text) 
		VALUES ($1, $2, $3)`,
		article.ID,
		article.Title,
		article.Text,
	); err != nil {
		return errors.Wrap(err, "insert article")
	}

	for _, tag := range article.Tags {
		if _, err := tx.Exec(ctx, `
				INSERT INTO article_tags (article_id, tag_name)
				VALUES ($1,  $2)
			`,
			article.ID,
			tag,
		); err != nil {
			return errors.Wrap(err, "insert articles tags")
		}
	}

	//nolint:wrapcheck
	return tx.Commit(ctx)
}

func (r *ArticleRepos) SelectAllArticlesByCourseID(ctx context.Context, courseId int) ([]models.Article, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.title, a.text, array_agg(at.tag_name) AS tags 
		FROM course_articles ca 
		JOIN articles a ON a.id = ca.article_id 
		LEFT JOIN article_tags at ON a.id = at.article_id 
		WHERE ca.course_id = $1 
		GROUP BY a.id
	`, courseId)
	if err != nil {
		return nil, errors.Wrap(err, "select article")
	}

	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Text,
			&article.Tags); err != nil {
			return nil, errors.Wrap(err, "scan article")
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *ArticleRepos) SelectOneArticleByID(ctx context.Context, articleID int) (*models.Article, error) {
	row := r.db.QueryRow(ctx, `
			SELECT a.id, a.title, a.text, array_agg(at.tag_name) AS tags
			FROM articles a
			LEFT JOIN article_tags at ON a.id = at.article_id
			WHERE a.id = $1 
			GROUP BY a.id`,
		articleID,
	)

	var article models.Article
	if err := row.Scan(&article); err != nil {
		return nil, errors.Wrap(err, "scan article")
	}

	return &article, nil
}

func (r *ArticleRepos) UpdateArticle(ctx context.Context, article models.Article) error {
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
		UPDATE articles 
		SET title = $1, text = $2 
		WHERE id = $3
	`, article.Title, article.Text, article.ID)
	if err != nil {
		return errors.Wrap(err, "update article")
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM article_tags WHERE article_id = $1
	`, article.ID)
	if err != nil {
		return errors.Wrap(err, "delete tags article")
	}

	for _, tag := range article.Tags {
		_, err = tx.Exec(ctx, `
			INSERT INTO article_tags (article_id, tag_name) 
			VALUES ($1, $2)
		`, article.ID, tag)
		if err != nil {
			return errors.Wrap(err, "insert tags article")
		}
	}

	//nolint:wrapcheck
	return tx.Commit(ctx)
}

func (r *ArticleRepos) DeleteArticle(ctx context.Context, articleId int) error {
	if _, err := r.db.Exec(ctx, `DELETE FROM articles WHERE id = $1`, articleId); err != nil {
		return errors.Wrap(err, "delete article")
	}

	return nil
}
