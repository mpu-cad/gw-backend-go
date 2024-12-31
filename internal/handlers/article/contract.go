package course

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type articleUC interface {
	CreateArticle(ctx context.Context, courseId int, article models.Article) error
	GetAllArticlesByCourseID(ctx context.Context, courseId int) ([]models.Article, error)
	GetOneArticleByID(ctx context.Context, articleID int) (*models.Article, error)
	UpdateArticle(ctx context.Context, article models.Article) error
	DeleteArticle(ctx context.Context, articleId int) error
}
