package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

func New(addr string, handler http.Handler, logger logger.Logger) *Server {
	server := &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
	return server
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Alert(
				fmt.Sprintf("Failed to start server at %s", s.httpServer.Addr),
				err,
			)
			os.Exit(1)
		}
	}()

	s.logger.Info(fmt.Sprintf("Server started at %s", s.httpServer.Addr))

	<-ctx.Done()

	s.logger.Info("Shutting down server..")

	if err := s.Shutdown(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
