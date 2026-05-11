package app

import (
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"lab-tracker/internal/config"
	"lab-tracker/internal/delivery/http/router"
	"lab-tracker/internal/delivery/http/server"
	"lab-tracker/internal/repository"
	"lab-tracker/internal/repository/postgres"
	"lab-tracker/internal/service"
)

type App struct {
	Server *server.Server
}

func NewApp(cfg *config.Config) *App {
	log, err := buildLogger(cfg.Env)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
	)

	conn, err := postgres.New(dsn)
	if err != nil {
		panic(fmt.Sprintf("app: connect to db: %v", err))
	}

	repo := repository.New(conn.Db)
	svc := service.New(
		log,
		repo.UserRepo,
		repo.GroupRepo,
		repo.LabWorkRepo,
		repo.AssignmentRepo,
		repo.SubmissionRepo,
		repo.GradeRepo,
		cfg.JWTSecret,
	)

	r, err := router.Init(svc, log, cfg.JWTSecret)
	if err != nil {
		panic(fmt.Sprintf("app: init router: %v", err))
	}

	return &App{
		Server: server.New(cfg, r),
	}
}

func buildLogger(env string) (*zap.Logger, error) {
	if env == "prod" || env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
