package main

import (
	"errors"
	"fmt"
)

// this represents the function signature for all command "handlers"
func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Provide a single argument for the login command. Usage: %s <arg>", cmd.name)
	}

	err := s.conf.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set successfully: %v\n", s.conf.Current_Username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	requestedFunction, ok := c.commandList[cmd.name]
	if !ok {
		return errors.New("That function does not exist")
	}

	if err := requestedFunction(s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}
