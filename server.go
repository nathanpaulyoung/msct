package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Server represents a server
type Server struct {
	Name        string
	Description string
	Config      *ServerConfig
}

//New creates a new Server from supplied name, description, and *Config
func newServer(name, desc string, cfg *GlobalConfig) *Server {
	return &Server{
		Name:         name,
		Description:  desc,
		ServerConfig: cfg,
	}
}

// Start starts a server safely
func (s *Server) Start() error {
	if !s.Exists() {
		return errors.New("server does not exist")
	}

	if s.IsRunning() {
		return errors.New("server is already running")
	}

	path := s.Config.RootDir
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	dir := s.Name
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	invocation := fmt.Sprintf("new -d -s msct-%s java -Xmx%dM -Xms%dM %s -jar %s%s%s nogui",
		strings.ToLower(s.Name),
		s.Config.RAMMax,
		s.Config.RAMMin,
		s.Config.JavaFlags,
		path,
		dir,
		s.Config.JarName,
	)

	command := exec.Command("tmux", strings.Fields(invocation)...)
	command.Dir = path + dir
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

// Stop stops a running server
func (s *Server) Stop() error {
	if !s.Exists() {
		return errors.New("server does not exist")
	}

	if !s.IsRunning() {
		return errors.New("server is not running")
	}

	return s.Send("stop")
}

// Send sends a command or message to a running server
func (s *Server) Send(message string) error {
	if !s.Exists() {
		return errors.New("server does not exist")
	}

	if !s.IsRunning() {
		return errors.New("server is not running")
	}

	command := exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("msct-%s:0.0", strings.ToLower(s.Name)), message, "Enter")
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

//Resume reattaches to a backgrounded tmux session
func (s *Server) Resume() error {
	if !s.Exists() {
		return errors.New("server does not exist")
	}

	if !s.IsRunning() {
		return errors.New("server is not running")
	}

	command := exec.Command("tmux", "a", "-t", fmt.Sprintf("msct-%s", strings.ToLower(s.Name)))
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

// IsRunning returns a bool indicating server status
func (s *Server) IsRunning() bool {
	command := exec.Command("tmux", "has-session", "-t", fmt.Sprintf("msct-%s", strings.ToLower(s.Name)))
	if err := command.Run(); err != nil {
		return false
	}

	return true
}

//Exists returns a bool indicating the presence of a server's jar file
func (s *Server) Exists() bool {
	path := s.Config.RootDir
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	dir := s.Name
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	if _, err := os.Stat(path + dir + s.Config.JarName); err != nil {
		return false
	}
	return true
}
