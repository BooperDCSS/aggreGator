package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/BooperDCSS/aggreGator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

func (s *state) scrapeFeeds() error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving next feed in database: %w", err)
	}

	feedID, err := s.db.IDFromURL(context.Background(), nextFeed)
	if err != nil {
		return fmt.Errorf("Error retrieving feed URL during scrape process: %w", err)
	}

	if err = s.db.MarkFeedFetched(context.Background(), feedID); err != nil {
		return fmt.Errorf("Error marking 'last_fetched_at' during scrape: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed)
	if err != nil {
		return fmt.Errorf("Error fetching feed during scrape: %w", err)
	}

	fmt.Printf("Storing posts from %q in database\n", rssFeed.Channel.Title)
	err = s.storePostInDB(rssFeed, feedID)
	if err != nil {
		return fmt.Errorf("Error storing posts in database: %w", err)
	}

	return nil
}

func (s *state) storePostInDB(rssFeed *RSSFeed, feedID uuid.UUID) error {
	for _, post := range rssFeed.Channel.Item {
		publishedDescription := rssDescriptionCheck(post.Description)
		publishedDate, err := convertTimeString(post.PubDate)
		if err != nil {
			fmt.Printf("Invalid PublishedAt date: %q. Using null value.\n", post.PubDate)
		}

		postData := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PublishedAt: publishedDate,
			Title:       post.Title,
			Url:         post.Link,
			Description: publishedDescription,
			FeedID:      feedID,
		}

		var pqError *pq.Error

		_, err = s.db.CreatePost(context.Background(), postData)
		if errors.As(err, &pqError) {
			if pqError.Code == pqerror.UniqueViolation {
				fmt.Printf("%q already in database, skipping post.\n\n", postData.Url)
				continue
			} else {
				return fmt.Errorf("pq error encountered: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("Error creating post in database: %w", err)
		}
		fmt.Printf("Success: %q stored in database\nURL: %q\n\n", postData.Title, postData.Url)
	}

	return nil
}

func rssDescriptionCheck(rssDescription string) sql.NullString {
	if rssDescription == "" {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: rssDescription, Valid: true}
}

func convertTimeString(rssTimeString string) (sql.NullTime, error) {

	if rssTimeString == "" {
		return sql.NullTime{Valid: false}, errors.New("Null PublishedAt string from feed. Using null value.")
	}

	rfcformat1, err := time.Parse(time.RFC1123Z, rssTimeString)
	if err == nil {
		return sql.NullTime{Time: rfcformat1, Valid: true}, nil
	}

	rfcformat2, err := time.Parse(time.RFC1123, rssTimeString)
	if err == nil {
		return sql.NullTime{Time: rfcformat2, Valid: true}, nil
	}

	rfcformat3, err := time.Parse(time.RFC3339, rssTimeString)
	if err == nil {
		return sql.NullTime{Time: rfcformat3, Valid: true}, nil
	}

	unixformat, err := time.Parse(time.UnixDate, rssTimeString)
	if err == nil {
		return sql.NullTime{Time: unixformat, Valid: true}, nil
	}

	styles := []string{
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"2 Jan 2006 15:04:05 -0700",
		"2 Jan 2006 15:04:05 MST",
		"2006-01-02T15:04:05Z0700",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
	}

	var finalErr error
	for _, style := range styles {
		validTime, err := time.Parse(style, rssTimeString)
		if err == nil {
			return sql.NullTime{Time: validTime, Valid: true}, nil
		}

		finalErr = err
	}

	return sql.NullTime{Valid: false}, finalErr
}
