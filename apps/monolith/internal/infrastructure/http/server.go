package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Server wraps an HTTP server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
	logger     zerolog.Logger
}

// NewServer creates a new HTTP server listening on the given port.
func NewServer(port string, engine *gin.Engine, logger zerolog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
		logger: logger,
	}
}

// Start begins listening for HTTP requests. This method blocks.
func (s *Server) Start() error {
	s.logger.Info().Str("addr", s.httpServer.Addr).Msg("HTTP server starting")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server with a 10-second timeout.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Info().Msg("shutting down HTTP server")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("HTTP server forced to shutdown")
	}
	s.logger.Info().Msg("HTTP server stopped")
}
