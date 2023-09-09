package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Aleksao998/LightingUserVault/core/cache"
	"github.com/Aleksao998/LightingUserVault/core/server/routers"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Server is the central manager of the LightingUserVault
type Server struct {
	config     *Config
	httpServer *http.Server
	logger     *zap.Logger
	storage    storage.Storage
}

// NewServer creates a new LightingUserVault server, using the passed in configuration
func NewServer(config *Config) (*Server, error) {
	// Get a production config
	cfg := zap.NewProductionConfig()

	// Set the desired log level
	cfg.Level.SetLevel(zapcore.DebugLevel)

	// Build the logger
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	vault, err := storage.GetStorage(logger)
	if err != nil {
		logger.Error("Failed to get storage", zap.Error(err))

		return nil, err
	}

	cacheMechanism, err := cache.GetCache(logger, "127.0.0.1:11211")
	if err != nil {
		logger.Error("Failed to get cache", zap.Error(err))

		return nil, err
	}

	router := routers.InitRouter(logger, vault, cacheMechanism)

	// create http server instance
	httpServer := &http.Server{
		Addr:              ":9097",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// initialize server
	server := &Server{
		config:     config,
		httpServer: httpServer,
		logger:     logger,
		storage:    vault,
	}

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error while listening and serving", zap.Error(err))
		}
	}()

	logger.Info("Server initialized and listening on :9097")

	return server, nil
}

// Close gracefully shuts down the LightingUserVault server
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.storage.Close(); err != nil {
		s.logger.Error("Storage shutdown failed", zap.Error(err))

		return err
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Server shutdown failed", zap.Error(err))

		return err
	}

	s.logger.Info("Server gracefully stopped")

	return nil
}
