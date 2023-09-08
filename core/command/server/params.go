package server

import "github.com/Aleksao998/LightingUserVault/core/server"

var (
	params = &serverParams{}
)

type serverParams struct {
}

func (p *serverParams) initRawParams() error {
	return nil
}

func (p *serverParams) generateConfig() *server.Config {
	return &server.Config{}
}
