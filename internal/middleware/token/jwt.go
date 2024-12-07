package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/mpu-cad/gw-backend-go/internal/configs"
	"github.com/mpu-cad/gw-backend-go/internal/entity"
)

// SignedToken - подписание JWT для авторизированного пользователя токена.
func SignedToken(ctx *fiber.Ctx) error {
	id, ok := ctx.Context().Value("UserID").(*int)
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   "can not signed token",
			Message: "not header login",
			Status:  http.StatusBadRequest,
		})
	}

	refreshToken := ctx.Context().Value("RefreshToken").(string)

	payload := jwt.MapClaims{
		"ExpiresAt": jwt.NewNumericDate(time.Now().UTC().Add(entity.ExpiresMinuteAccessToken * time.Minute)),
		"UserID":    id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token.Header["id"] = id

	path := "./config.yaml"
	cfg := configs.MustConfig(&path)

	secretKey := []byte(cfg.JWT.SecretAuth)

	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: "can not signed string",
			Status:  http.StatusBadRequest,
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// CompareToken - проверка корректности JWT токена.
func CompareToken(ctx *fiber.Ctx) error {
	headers := ctx.GetReqHeaders()
	tokens, ok := headers[entity.HeaderAccessToken]
	if !ok || len(tokens) < 1 {
		return ctx.Status(http.StatusUnauthorized).JSON(entity.ErrorsRequest{
			Error:   "token is empty",
			Message: "can not get token",
			Status:  http.StatusUnauthorized,
		})
	}

	token := tokens[0]
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok = token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		path := "./config.yaml"
		cfg := configs.MustConfig(&path)

		secretKey := []byte(cfg.JWT.SecretAuth)

		return secretKey, nil
	})
	if err != nil {
		return ctx.Status(http.StatusForbidden).JSON(entity.ErrorsRequest{
			Error:   err.Error(),
			Message: "token token is not valid",
			Status:  http.StatusForbidden,
		})
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(entity.ErrorsRequest{
			Error:   "payload is not a valid",
			Message: "",
			Status:  http.StatusUnauthorized,
		})
	}

	t, ok := payload["ExpiresAt"].(float64)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(entity.ErrorsRequest{
			Error:   "ExpiresAt is empty",
			Message: "",
			Status:  http.StatusUnauthorized,
		})
	}

	now := time.Now().Unix()

	if now > int64(t) {
		return ctx.Status(http.StatusUnauthorized).JSON(entity.ErrorsRequest{
			Error:   "JWT token is expired",
			Message: "",
			Status:  http.StatusUnauthorized,
		})
	}

	ctx.Locals("UserID", jwtToken.Header["id"])

	return ctx.Next()
}
