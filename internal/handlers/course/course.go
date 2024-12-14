package course

import "github.com/gofiber/fiber/v2"

type Handle struct {
	courseUC
}

func NewHandleCourse(courseUC courseUC) *Handle {
	return &Handle{
		courseUC: courseUC,
	}
}

func (h *Handle) CreateCourse(ctx *fiber.Ctx) error {
	return nil
}
