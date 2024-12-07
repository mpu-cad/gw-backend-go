package log

import (
	"encoding/json"
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/gofiber/fiber/v2"

	"github.com/sirupsen/logrus"

	"github.com/mpu-cad/gw-backend-go/internal/logger"
)

const (
	reg   = "/api/user/registration"
	login = "/api/user/login"
)

func New() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		msk, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return fmt.Errorf("can not load location, err: %w", err)
		}

		start := time.Now().In(msk)

		err1 := ctx.Next()

		endTime := time.Now().In(msk)
		interval := endTime.Sub(start)
		runTime := interval.Milliseconds()

		resBody := jsonParse(ctx.Context().Response.Body())
		reqBody := jsonParse(ctx.Context().Request.Body())

		path := ctx.Path()

		if path == reg || path == login {
			reqBody = nil
		}

		logger.Log.WithFields(logrus.Fields{
			"api":         path,
			"method":      ctx.Method(),
			"status code": ctx.Context().Response.StatusCode(),
			"run time":    runTime,
			"response":    resBody,
			"request":     reqBody,
		}).Infoln()

		if err1 != nil {
			return err1
		}

		return nil
	}
}

func jsonParse(body []byte) interface{} {
	if len(body) == 0 {
		return nil
	}

	var result interface{}
	_ = json.Unmarshal(body, &result)

	return result
}
