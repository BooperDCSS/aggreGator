package main

import "errors"

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
