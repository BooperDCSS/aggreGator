package main

import (
	"context"
	"fmt"
)

func (s *state) scrapeFeeds() error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving next feed in database: %w", err)
	}

	nextFeedID, err := s.db.IDFromURL(context.Background(), nextFeed)
	if err != nil {
		return fmt.Errorf("Error retrieving feed URL during scrape process: %w", err)
	}

	if err = s.db.MarkFeedFetched(context.Background(), nextFeedID); err != nil {
		return fmt.Errorf("Error marking 'last_fetched_at' during scrape: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed)
	if err != nil {
		return fmt.Errorf("Error fetching feed during scrape: %w", err)
	}

	fmt.Printf("Scraping %q from %q\n\n", rssFeed.Channel.Title, rssFeed.Channel.Link)
	for _, post := range rssFeed.Channel.Item {
		fmt.Printf("Title: %v\nDesc: %v\nLink: %v\nPubDate: %v\n\n",
			post.Title, post.Description, post.Link, post.PubDate)
	}

	return nil
}
