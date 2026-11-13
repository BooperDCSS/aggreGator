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

// fetches feeds using the scrapeFeeds function and a time interval defined by the user
func handlerFetch(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Enter a time duration to schedule regular RSS scraping. Usage %s <1m or 30s, etc.>", cmd.name)
	}

	time_between_reqs, _ := time.ParseDuration(cmd.arguments[0])
	fmt.Printf("Collecting feeds every %v\n\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)

	for ; ; <-ticker.C {
		err := s.scrapeFeeds()
		if err != nil {
			fmt.Printf("Error encountered during 'agg' process: %v", err)
			break
		}
	}

	return nil
}

// adds a feed to the database, associated with the current user using feed_follows table
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("Add a feed by supplying its name and URL: Usage %s <name> <url>", cmd.name)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("Error generating feed follow during AddFeed operation: %w", err)
	}

	fmt.Printf("Feed added: %+v", feed)
	return nil

}

// lists all feeds, their urls, and who created them
func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("No arguments required when retrieving list of feeds.")
	}

	feedList, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds in database: %w", err)
	}

	for i, feed := range feedList {
		creator, err := s.db.NameFromID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error identifying feed creator: %w", err)
		}
		fmt.Printf("Feed #%v\nFeed Name: %s\nFeed URL: %s\nFeed creator: %s\n\n", i+1, feed.Name, feed.Url, creator)
	}

	return nil
}

// allows current user to follow a feed that already exists in the db
func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Must provide URL feed to follow. Usaage: %s <url>", cmd.name)
	}

	feedID, err := s.db.IDFromURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Error retrieving Feed ID from URL: %w", err)
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	feedData, err := s.db.CreateFeedFollow(context.Background(), followParams)

	fmt.Printf("Feed Followed: %s, Current User: %s", feedData.FeedName, feedData.UserName)
	return nil
}

// finds the feeds associated with the current user
func handlerCurrentUserFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("No arguments are required when listing the current users feeds.")
	}

	feedsList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error retrieving list of feeds for current user: %w", err)
	}

	fmt.Printf("Feed follows for %s:\n\n", s.conf.Current_Username)
	for _, feed := range feedsList {
		fmt.Printf("%s\n", feed)
	}

	return nil
}

// deletes a followed feed using the current user and the feed ID
func handlerDeleteFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Provide the URL follow you wish to delete for current user. Usage %s <url>", cmd.name)
	}

	urlFeedID, err := s.db.IDFromURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Error retrieving feed ID for delete: %w", err)
	}

	deleteParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: urlFeedID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteParams)
	if err != nil {
		return fmt.Errorf("Error deleting feed follow for current user: %w", err)
	}

	fmt.Printf("User %q successfully removed %q from follows", s.conf.Current_Username, cmd.arguments[0])

	return nil
}
