package server

import (
	"context"
	"lab-tracker/internal/config"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		Handler:      handler,
	}

	return &Server{
		server: srv,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
