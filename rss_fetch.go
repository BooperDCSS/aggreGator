package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/BooperDCSS/aggreGator/internal/rss"
)

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rss.RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "aggreGator")

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return &rss.RSSFeed{}, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return &rss.RSSFeed{}, err
	}

	var rssFeedData rss.RSSFeed

	if err := xml.Unmarshal(data, &rssFeedData); err != nil {
		return &rss.RSSFeed{}, err
	}

	return &rssFeedData, nil
}

func cleanStringsRSS(r *rss.RSSFeed) {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	fmt.Printf("Channel: %+v\nChannel Description: %+v", r.Channel.Title, r.Channel.Description)

	for _, item := range r.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		fmt.Printf("--> %+v", item)
	}

}
