package main

import (
	"errors"
	"fmt"

	"github.com/BooperDCSS/aggreGator/internal/config"
)

type state struct {
	conf *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

// this represents the function signature for all command "handlers"
func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return errors.New("No argument provided for login.")
	} else if len(cmd.arguments) > 2 {
		return errors.New("Argument for login must be a single word/username.")
	}

	err := s.conf.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set: %v", s.conf.Current_Username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	requestedFunction := c.commandList[cmd.name]
	if requestedFunction == nil {
		return errors.New("That function does not exist")
	}

	err := requestedFunction(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}
