package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/handler"
)

type server struct {
	httpServer *http.Server
}

func NewServer(port int, db database.DB) *server {
	mux := http.NewServeMux()
	mux.Handle("/", handler.NewHandler(db))

	return &server{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

func (s *server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}
	return nil
}
