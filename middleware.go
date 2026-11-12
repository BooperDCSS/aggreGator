package main

import (
	"context"
	"fmt"

	"github.com/BooperDCSS/aggreGator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.conf.Current_Username)
		if err != nil {
			return fmt.Errorf("Cannot retrieve current user in middleware process: %w", err)
		}

		err = handler(s, cmd, user)
		if err != nil {
			return fmt.Errorf("Cannot pass user to handler in middlware process: %w", err)
		}

		return nil
	}

}
