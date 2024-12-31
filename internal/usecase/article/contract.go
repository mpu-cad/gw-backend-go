package article

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type articleRepos interface {
	InsertArticle(ctx context.Context, courseId int, article models.Article) error
	SelectAllArticlesByCourseID(ctx context.Context, courseId int) ([]models.Article, error)
	SelectOneArticleByID(ctx context.Context, articleID int) (*models.Article, error)
	UpdateArticle(ctx context.Context, article models.Article) error
	DeleteArticle(ctx context.Context, articleId int) error
}
