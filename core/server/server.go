package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Aleksao998/LightningUserVault/core/cache"
	"github.com/Aleksao998/LightningUserVault/core/server/routers"
	"github.com/Aleksao998/LightningUserVault/core/storage"
	"go.uber.org/zap"
)

// Server is the central manager of the LightningUserVault
type Server struct {
	config     *Config
	httpServer *http.Server
	logger     *zap.Logger
	storage    storage.Storage
}

// NewServer creates a new LightningUserVault server, using the passed in configuration
func NewServer(config *Config) (*Server, error) {
	// Get a production config
	cfg := zap.NewProductionConfig()

	// Set the desired log level
	cfg.Level.SetLevel(config.LogLevel)

	// Build the logger
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	defer logger.Sync()

	// Create storage config
	storageConfig := storage.Config{
		StorageType: config.StorageType,
		DBHost:      config.DBHost.IP.String(),
		DBPort:      strconv.Itoa(config.DBHost.Port),
		DBPass:      config.DBPass,
		DBName:      config.DBName,
		DBUser:      config.DBUser,
	}

	// Initialize storage
	vault, err := storage.GetStorage(logger, storageConfig)
	if err != nil {
		logger.Error("Failed to get storage", zap.Error(err))

		return nil, err
	}

	// Create cache config
	cacheConfig := cache.Config{
		CacheType:       config.CacheType,
		MemcacheAddress: config.MemcacheAddress,
		Enabled:         config.EnableCache,
	}

	// Initialize cache
	cacheMechanism, err := cache.GetCache(logger, cacheConfig)
	if err != nil {
		logger.Error("Failed to get cache", zap.Error(err))

		return nil, err
	}

	routerConfig := routers.Config{
		CacheEnabled: config.EnableCache,
	}

	router := routers.InitRouter(logger, vault, cacheMechanism, routerConfig)

	// Create http server instance
	httpServer := &http.Server{
		Addr:              config.ServerAddress.String(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Initialize server
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

	logger.Info(fmt.Sprintf("Server initialized and listening on %s", config.ServerAddress.String()))

	return server, nil
}

// Close gracefully shuts down the LightningUserVault server
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
