package user

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/mpu-cad/gw-backend-go/internal/entity"
	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type Handle struct {
	userUC  userUC
	redisUC redisUC
}

func NewHandleUser(userUC userUC, redisUC redisUC) *Handle {
	return &Handle{
		userUC:  userUC,
		redisUC: redisUC,
	}
}

func (r *Handle) Registration(ctx *fiber.Ctx) error {
	var data registrationUserRequest
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			entity.ErrorsRequest{
				Error:   err.Error(),
				Message: entity.ErrorParseBody,
				Status:  http.StatusBadRequest,
			})
	}

	_, err = r.userUC.Registration(ctx.Context(), models.User{
		Phone:    data.Phone,
		Email:    data.Email,
		HashPass: data.Password,
		Name:     data.Name,
		Surname:  data.Surname,
	})

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: "can not create user",
			Status:  http.StatusBadRequest,
		})
	}

	return ctx.Status(http.StatusCreated).Send(nil)
}

func (r *Handle) Login(ctx *fiber.Ctx) error {
	var user loginUserRequest
	err := ctx.BodyParser(&user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: entity.ErrorParseBody,
			Status:  http.StatusBadRequest,
		})
	}

	getUser, err := r.userUC.Login(ctx.Context(), user.Email, user.Password)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: "can not login user",
			Status:  http.StatusBadRequest,
		})
	}

	refreshToken := r.redisUC.CreateRefreshToken(ctx.Context(), getUser.ID)

	ctx.Locals("UserID", getUser.ID)
	ctx.Locals("RefreshToken", refreshToken)

	return ctx.Next()
}
