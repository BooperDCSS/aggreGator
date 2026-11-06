package main

import (
	"context"
	"fmt"
	"time"

	"github.com/BooperDCSS/aggreGator/internal/database"
	"github.com/google/uuid"
)

// uses s.conf.SetUser to change current user
func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Provide a single argument for the login command. Usage: %s <arg>", cmd.name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("User %q does not exist in database", cmd.arguments[0])
	}

	err = s.conf.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set successfully: %v\n", s.conf.Current_Username)
	return nil
}

// creates and sets a new user as the current user
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

// resets database, removing all users
func handlerReset(s *state, cmd command) error {
	if len(cmd.arguments) > 0 {
		return fmt.Errorf("No argument required for resetting the user database")
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("'users' database successfully reset. All users removed.")
	return nil
}

// lists all users currently in the database; indicates which is the current user in the config
func handlerListUsers(s *state, cmd command) error {
	if len(cmd.arguments) > 0 {
		return fmt.Errorf("No argument required for listing all users.")
	}

	userList, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving user from database: %w", err)
	}

	for _, user := range userList {
		if user == s.conf.Current_Username {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}

	return nil
}

func handlerFetch(s *state, cmd command) error {

	if len(cmd.arguments) > 0 {
		return fmt.Errorf("No argument required for the aggregator command.")
	}

	feed, err := fetchFeed(context.Background(), s.target)
	if err != nil {
		return fmt.Errorf("Error fetching XML from target: %w", err)
	}

	cleanStringsRSS(feed)

	return nil
}
