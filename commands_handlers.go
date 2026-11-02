package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BooperDCSS/aggreGator/internal/database"
	"github.com/google/uuid"
)

// this represents the function signature for all command "handlers"
func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Provide a single argument for the login command. Usage: %s <arg>", cmd.name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		log.Fatalf("User %q does not exist in database", cmd.arguments[0])
	}

	err = s.conf.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set successfully: %v\n", s.conf.Current_Username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Provide a single argument for the register command. Usage: %s <arg>", cmd.name)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}

	registeredUser, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("couldn't create the user, error: %w", err)
	}

	err = s.conf.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s created.\nNew user: %+v", cmd.arguments[0], registeredUser)
	return nil

}

func (c *commands) run(s *state, cmd command) error {
	requestedFunction, ok := c.commandList[cmd.name]
	if !ok {
		return errors.New("That function does not exist")
	}

	// no need to check err again since this function always checks and returns errors anyway
	// this will either return the error or return nil
	return requestedFunction(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}
