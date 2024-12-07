package app

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mpu-cad/gw-backend-go/internal/configs"
	"github.com/mpu-cad/gw-backend-go/internal/configure"
	handle "github.com/mpu-cad/gw-backend-go/internal/handlers/user"
	"github.com/mpu-cad/gw-backend-go/internal/logger"
	"github.com/mpu-cad/gw-backend-go/internal/middleware/log"
	"github.com/mpu-cad/gw-backend-go/internal/middleware/token"
	"github.com/mpu-cad/gw-backend-go/internal/storage/postgresql"
	"github.com/mpu-cad/gw-backend-go/internal/storage/redis"
	"github.com/mpu-cad/gw-backend-go/internal/usecase/mailer"
	redisUC "github.com/mpu-cad/gw-backend-go/internal/usecase/redis"
	"github.com/mpu-cad/gw-backend-go/internal/usecase/user"
)

type App struct {
	cfg *configs.Config
}

func New(cfg *configs.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) {
	app := fiber.New()
	app.Use(log.New())

	logger.InitLogger(a.cfg.Logger)

	// DB
	dbPool := configure.Postgres(ctx, a.cfg.Postgres)
	defer dbPool.Close()

	if err := a.cfg.Postgres.MigrationsUp(); err != nil {
		logger.Log.Errorf("can not up migration in postgres, err: %v", err)
	}

	redisDB := configure.Redis(a.cfg.Redis)

	// Repos
	userRepos := postgresql.NewUserRepos(dbPool)
	redisRepos := redis.NewTokenRepos(redisDB)

	// UseCase
	ucMailer := mailer.New(a.cfg.Mailer)
	ucUser := user.NewUCUser(userRepos, ucMailer)
	ucRedis := redisUC.NewUCRepos(redisRepos, userRepos)

	// Handler
	userHandler := handle.NewHandleUser(ucUser, ucRedis)

	// endpoint
	api := app.Group("/api")

	// эндпоинты для юзеров
	users := api.Group("/user")

	users.Post("/registration", userHandler.Registration)
	users.Post("/login", userHandler.Login, token.SignedToken)

	// эндпоинты для курсов
	course := api.Group("/course")

	// получить все курсы
	course.Get("", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.Status(http.StatusOK).SendString("<h1>Hello world</h1>")
	})
	// получить курс по id
	course.Get("/:id", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.Status(http.StatusOK).SendString("<h1>Hello world</h1>")
	})

	course.Post("", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.Status(http.StatusOK).SendString("<h1>Hello world</h1>")
	}) // создать курс

	course.Delete("/:id", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.Status(http.StatusOK).SendString("<h1>Hello world</h1>")
	}) // удалить курс
	course.Put("/:id", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.Status(http.StatusOK).SendString("<h1>Hello world</h1>")
	}) // обновить курс

	// эндпоинты для статьей
	api.Group("/article", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.Status(http.StatusOK).SendString("<h1>Hello world</h1>")
	})

	err := app.Listen(a.cfg.Server.String())
	if err != nil {
		logger.Log.Errorf("err Listen: %v", err)

		return
	}
}
