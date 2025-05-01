package main

import (
	"errors"
	"fmt"

	"github.com/stkisengese/gator-cli/internal/config"
)

// State of the application
type state struct {
	cfg *config.Config
}

// command represents a CLI command
type command struct {
	name string
	args []string
}

// commands represents all registered set of CLI commands
type commands struct {
	handlers map[string]func(*state, command) error
}

// NewCommands creates a new commands instance
func newCommands() *commands {
	return &commands{
		handlers: make(map[string]func(*state, command) error),
	}
}

// run executes a command
func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.handlers[cmd.name]; ok {
		return f(s, cmd)
	}

	return fmt.Errorf("unknown command: %s", cmd.name)
}

// register adds a new command handler
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

// handlerLogin handles the login command
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required")
	}

	username := cmd.args[0]
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("Logging in as: %s\n", username)
	fmt.Println("User switched successfully!")
	return nil
}
