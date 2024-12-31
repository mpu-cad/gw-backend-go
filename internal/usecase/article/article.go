package article

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type UCArticle struct {
	articleRepos
}

func NewArticleUC(repos articleRepos) *UCArticle {
	return &UCArticle{
		articleRepos: repos,
	}
}

func (uc *UCArticle) CreateArticle(ctx context.Context, courseId int, article models.Article) error {
	//nolint:wrapcheck
	return uc.articleRepos.InsertArticle(ctx, courseId, article)
}

func (uc *UCArticle) GetAllArticlesByCourseID(ctx context.Context, courseId int) ([]models.Article, error) {
	//nolint:wrapcheck
	return uc.articleRepos.SelectAllArticlesByCourseID(ctx, courseId)
}

func (uc *UCArticle) GetOneArticleByID(ctx context.Context, articleID int) (*models.Article, error) {
	//nolint:wrapcheck
	return uc.articleRepos.SelectOneArticleByID(ctx, articleID)
}

func (uc *UCArticle) UpdateArticle(ctx context.Context, article models.Article) error {
	//nolint:wrapcheck
	return uc.articleRepos.UpdateArticle(ctx, article)
}

func (uc *UCArticle) DeleteArticle(ctx context.Context, articleId int) error {
	//nolint:wrapcheck
	return uc.articleRepos.DeleteArticle(ctx, articleId)
}
