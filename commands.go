package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stkisengese/gator-cli/internal/config"
	"github.com/stkisengese/gator-cli/internal/database"
)

// State of the application
type state struct {
	db  *database.Queries
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
	ctx := context.Background()

	// Check if the user exists
	exists, err := s.db.UserExists(ctx, username)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	if !exists {
		return fmt.Errorf("user %s does not exist", username)
	}

	// Set the user in the config
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("Logging in as: %s\n", username)
	fmt.Println("User switched successfully!")
	return nil
}

// handleRegister handles the register command
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required")
	}

	username := cmd.args[0]
	ctx := context.Background()

	// Check if the user exists
	exists, err := s.db.UserExists(ctx, username)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	if exists {
		return fmt.Errorf("user %s already exists", username)
	}

	// Create the user
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	// Set the user in the config
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("User %s created successfully!\n", user.Name)
	return nil
}
