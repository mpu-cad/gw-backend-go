package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mpu-cad/gw-backend-go/internal/configs"
	"github.com/mpu-cad/gw-backend-go/internal/configure"
	handleCourse "github.com/mpu-cad/gw-backend-go/internal/handlers/course"
	handleUser "github.com/mpu-cad/gw-backend-go/internal/handlers/user"
	"github.com/mpu-cad/gw-backend-go/internal/logger"
	"github.com/mpu-cad/gw-backend-go/internal/middleware/log"
	"github.com/mpu-cad/gw-backend-go/internal/middleware/token"
	"github.com/mpu-cad/gw-backend-go/internal/storage/postgresql"
	"github.com/mpu-cad/gw-backend-go/internal/storage/redis"
	"github.com/mpu-cad/gw-backend-go/internal/usecase/course"
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
	courseRepos := postgresql.NewCourseRepos(dbPool)
	_ = postgresql.NewArticleRepos(dbPool)

	// UseCase
	ucMailer := mailer.New(a.cfg.Mailer)
	ucUser := user.NewUCUser(userRepos, ucMailer, redisRepos)
	ucRedis := redisUC.NewUCRepos(redisRepos, userRepos)
	ucCourse := course.NewUCCourse(courseRepos)

	// Handler
	userHandler := handleUser.NewHandleUser(ucUser, ucRedis)
	courseHandler := handleCourse.NewHandleCourse(ucCourse)

	// endpoint
	api := app.Group("/api")

	// эндпоинты для юзеров
	users := api.Group("/user")

	users.Post("/registration", userHandler.Registration)
	users.Post("/login", userHandler.Login, token.SignedToken)
	users.Post("/email/confirm", userHandler.ConfirmEmail)

	// эндпоинты для курсов
	courseEndPoints := api.Group("/course")

	// создать курс
	courseEndPoints.Post("", courseHandler.CreateCourse)

	// получить все курсы
	courseEndPoints.Get("/", courseHandler.GetAllCourses)

	// получить курс по id
	courseEndPoints.Get("/:id", courseHandler.GetCourseByID)

	// удалить курс
	courseEndPoints.Delete("/:id", courseHandler.DeleteCourse)

	// обновить курс
	courseEndPoints.Put("/:id", courseHandler.UpdateCourse)

	// эндпоинты для статьей
	articleEndPoints := courseEndPoints.Group("/:id/article")

	articleEndPoints.Post("", nil)

	err := app.Listen(a.cfg.Server.String())
	if err != nil {
		logger.Log.Errorf("err Listen: %v", err)

		return
	}
}
