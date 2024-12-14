package user

import (
	"net/http"
	"strconv"

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

	id, err := r.userUC.Registration(ctx.Context(), models.User{
		Phone:    data.Phone,
		Email:    data.Email,
		LastName: data.LastName,
		Login:    data.Login,
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

	return ctx.Status(http.StatusCreated).SendString(strconv.Itoa(*id))
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

	getUser, err := r.userUC.Login(ctx.Context(), user.Login, user.Password)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: "can not login user",
			Status:  http.StatusBadRequest,
		})
	}

	refreshToken := r.redisUC.CreateRefreshToken(ctx.Context(), getUser.ID)

	ctx.Locals("UserID", &getUser.ID)
	ctx.Locals("RefreshToken", refreshToken)

	return ctx.Next()
}

func (r *Handle) ConfirmEmail(ctx *fiber.Ctx) error {
	var confirm confirmEmail
	if err := ctx.BodyParser(&confirm); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: entity.ErrorParseBody,
			Status:  http.StatusBadRequest,
		})
	}

	if err := r.userUC.ConfirmMail(ctx.Context(), confirm.ID, confirm.Code); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: "can not confirm email",
			Status:  http.StatusBadRequest,
		})
	}

	return ctx.SendStatus(http.StatusNoContent)
}
