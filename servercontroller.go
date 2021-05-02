package main

// ServerController manages and orchestrates all servers for a single MSCT instance
type ServerController struct {
	Servers []Server
	Config  *ServerConfig
}

// initializes a nil ServerController to default values
func newServerController(cfg *ServerConfig) *ServerController {
	return &ServerController{
		Servers: []Server{},
		Config:  cfg,
	}
}
