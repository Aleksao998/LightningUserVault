package server

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

	return server, nil
}

// Close closes the LightingUserVault server
func (s *Server) Close() {
}
