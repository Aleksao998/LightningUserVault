package server

import (
	"github.com/Aleksao998/LightingUserVault/core/server/handlers"
	"github.com/Aleksao998/LightingUserVault/core/server/routers"
	"github.com/Aleksao998/LightingUserVault/core/storage"
	"log"
	"net/http"
)

// Server is the central manager of the LightingUserVault
type Server struct {
	// config server config
	config *Config
}

// NewServer creates a new LightingUserVault server, using the passed in configuration
func NewServer(config *Config) (*Server, error) {
	// initialize server
	server := &Server{
		config: config,
	}

	vault, _ := storage.GetStorage()

	hendlers := handlers.NewUserHandler(vault)

	// initialize server
	router := routers.InitRouter(hendlers)

	// create http server instance
	srv := &http.Server{
		Addr:    ":9097",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return server, nil
}

// Close closes the LightingUserVault server
func (s *Server) Close() {
}
