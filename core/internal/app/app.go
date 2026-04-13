package app

import (
	"lab-tracker/internal/config"
	"lab-tracker/internal/delivery/http/router"
	"lab-tracker/internal/delivery/http/server"

	"go.uber.org/zap"
)

type App struct {
	Server *server.Server
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	router, err := router.Init()
	if err != nil {
		panic(err)
	}

	httpServer := server.New(cfg, router)

	return &App{
		Server: httpServer,
	}
}
