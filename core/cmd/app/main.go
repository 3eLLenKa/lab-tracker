package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lab-tracker/internal/app"
	"lab-tracker/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application := app.NewApp(cfg)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err := application.Server.Run(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Server.Stop(shutdownCtx); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}

	log.Println("server stopped cleanly")
}
