package course

import (
	"github.com/mpu-cad/gw-backend-go/internal/entity"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type Handle struct {
	courseUC
}

func NewHandleCourse(courseUC courseUC) *Handle {
	return &Handle{
		courseUC: courseUC,
	}
}

func (h *Handle) CreateCourse(ctx *fiber.Ctx) error {
	var course models.Course
	if err := ctx.BodyParser(&course); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: entity.ErrorParseBody,
				Status:  http.StatusBadRequest,
			})
	}

	if err := h.courseUC.CreateCourse(ctx.Context(), course); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: "internal error",
				Status:  http.StatusBadRequest,
			})
	}

	return ctx.SendStatus(http.StatusCreated)
}

func (h *Handle) GetAllCourses(ctx *fiber.Ctx) error {
	limitValue := ctx.Query("limit")
	pageValue := ctx.Query("page")

	var intLimit int
	var err error
	if limitValue != "" {
		intLimit, err = strconv.Atoi(limitValue)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(
				entity.ErrorsRequest{
					Error:   err.Error(),
					Message: "limit is not int",
					Status:  http.StatusBadRequest,
				})
		}

	}

	var intPage int
	if pageValue != "" {
		intPage, err = strconv.Atoi(pageValue)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(
				entity.ErrorsRequest{
					Error:   err.Error(),
					Message: "page is not int",
					Status:  http.StatusBadRequest,
				})
		}
	}

	courses, err := h.courseUC.GetAllCourses(ctx.Context(), intLimit, intPage)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: "",
				Status:  http.StatusInternalServerError,
			})
	}

	if courses == nil {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.Status(http.StatusOK).JSON(courses)
}

func (h *Handle) GetCourseByID(ctx *fiber.Ctx) error {
	paramsInt, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: "can not get id by path, id is need integer",
				Status:  http.StatusBadRequest,
			})
	}

	course, err := h.courseUC.GetCourseByID(ctx.Context(), paramsInt)
	if err != nil {

		return ctx.Status(http.StatusInternalServerError).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: "internal err",
				Status:  http.StatusInternalServerError,
			})
	}

	if course == nil {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.Status(http.StatusOK).JSON(course)
}

func (h *Handle) DeleteCourse(ctx *fiber.Ctx) error {
	paramsInt, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: "can not get id by path, id is need integer",
				Status:  http.StatusBadRequest,
			})
	}

	err = h.courseUC.DeleteCourse(ctx.Context(), paramsInt)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: "internal err",
				Status:  http.StatusInternalServerError,
			})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (h *Handle) UpdateCourse(ctx *fiber.Ctx) error {
	var course models.Course

	paramsInt, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	course.ID = paramsInt

	if err := ctx.BodyParser(&course); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: entity.ErrorParseBody,
				Status:  http.StatusBadRequest,
			})
	}

	if err := h.courseUC.UpdateCourse(ctx.Context(), course); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: "internal error",
				Status:  http.StatusInternalServerError,
			})
	}

	return ctx.SendStatus(http.StatusOK)
}
