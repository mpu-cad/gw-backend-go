package course

import (
	"github.com/gofiber/fiber/v2"
)

type Handle struct {
	articleUC
}

func NewHandleArticle(repos articleUC) *Handle {
	return &Handle{
		repos,
	}
}

func (h *Handle) CreateArticle(ctx *fiber.Ctx) error {
	return nil
}

func (h *Handle) GetAllArticleByCourseID(ctx *fiber.Ctx) error {
	return nil
}

func (h *Handle) GetOneArticleByID(ctx *fiber.Ctx) error {
	return nil
}

func (h *Handle) DeleteArticle(ctx *fiber.Ctx) error {
	return nil
}

func (h *Handle) UpdateArticle(ctx *fiber.Ctx) error {
	return nil
}
