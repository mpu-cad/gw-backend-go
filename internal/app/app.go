package app

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mpu-cad/gw-backend-go/internal/usecase/article"

	"github.com/mpu-cad/gw-backend-go/internal/configs"
	"github.com/mpu-cad/gw-backend-go/internal/configure"
	handleArticle "github.com/mpu-cad/gw-backend-go/internal/handlers/article"
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
	app.Use(recover.New())

	logger.InitLogger(a.cfg.Logger)

	// Подключение к разным БД
	dbPool := configure.Postgres(ctx, a.cfg.Postgres)
	defer dbPool.Close()

	if err := a.cfg.Postgres.MigrationsUp(); err != nil {
		logger.Log.Errorf("can not up migration in postgres, err: %v", err)
	}

	redisDB := configure.Redis(a.cfg.Redis)

	// Слой работы с базой данных
	userRepos := postgresql.NewUserRepos(dbPool)
	redisRepos := redis.NewTokenRepos(redisDB)
	courseRepos := postgresql.NewCourseRepos(dbPool)
	articleRepos := postgresql.NewArticleRepos(dbPool)

	// Слой логики приложения
	ucMailer := mailer.New(a.cfg.Mailer)
	ucUser := user.NewUCUser(userRepos, ucMailer, redisRepos)
	ucRedis := redisUC.NewUCRepos(redisRepos, userRepos)
	ucCourse := course.NewUCCourse(courseRepos)
	ucArticle := article.NewArticleUC(articleRepos)

	// Слой обработчиков запросов
	userHandler := handleUser.NewHandleUser(ucUser, ucRedis)
	courseHandler := handleCourse.NewHandleCourse(ucCourse)
	articleHandler := handleArticle.NewHandleArticle(ucArticle)

	// endpoint
	api := app.Group("/api")

	// эндпоинты для юзеров /api/user
	users := api.Group("/user")

	users.Post("/registration", userHandler.Registration)
	users.Post("/login", userHandler.Login, token.SignedToken)
	users.Post("/email/confirm", userHandler.ConfirmEmail)

	// эндпоинты для курсов /api/course
	courseEndPoints := api.Group("/course")

	// создать курс
	courseEndPoints.Post("", courseHandler.CreateCourse)

	// получить все курсы
	courseEndPoints.Get("/", courseHandler.GetAllCourses)

	// получить курс по id
	courseEndPoints.Get("/:course_id", courseHandler.GetCourseByID)

	// удалить курс
	courseEndPoints.Delete("/:course_id", courseHandler.DeleteCourse)

	// обновить курс
	courseEndPoints.Put("/:course_id", courseHandler.UpdateCourse)

	// эндпоинты для статьей /api/course/:course_id/article
	articleEndPoints := courseEndPoints.Group("/:id/article")

	// создать статью
	articleEndPoints.Post("/", articleHandler.CreateArticle)

	// получить все статьи по курсу
	articleEndPoints.Post("/", articleHandler.GetAllArticleByCourseID)

	// обновить статью
	articleEndPoints.Post("/:article_id", articleHandler.UpdateArticle)

	// удалить статью
	articleEndPoints.Post("/:article_id", articleHandler.UpdateArticle)

	err := app.Listen(a.cfg.Server.String())
	if err != nil {
		logger.Log.Errorf("err Listen: %v", err)

		return
	}
}
