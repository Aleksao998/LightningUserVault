package server

import (
	"context"
	"github.com/Aleksao998/LightingUserVault/core/server/routers"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"log"
	"net/http"
	"time"
)

// Server is the central manager of the LightingUserVault
type Server struct {
	config     *Config
	httpServer *http.Server
}

// NewServer creates a new LightingUserVault server, using the passed in configuration
func NewServer(config *Config) (*Server, error) {
	vault, err := storage.GetStorage()
	if err != nil {
		return nil, err
	}

	router := routers.InitRouter(vault)

	// create http server instance
	httpServer := &http.Server{
		Addr:    ":9097",
		Handler: router,
	}

	// initialize server
	server := &Server{
		config:     config,
		httpServer: httpServer,
	}

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err)
		}
	}()

	return server, nil
}

// Close gracefully shuts down the LightingUserVault server
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Failed:%+v", err)
		return err
	}
	log.Println("Server gracefully stopped")
	return nil
}
