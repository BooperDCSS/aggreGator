package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "aggreGator")

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var rssFeedData RSSFeed

	if err := xml.Unmarshal(data, &rssFeedData); err != nil {
		return &RSSFeed{}, err
	}

	//clean up the strings

	rssFeedData.Channel.Title = html.UnescapeString(rssFeedData.Channel.Title)
	rssFeedData.Channel.Description = html.UnescapeString(rssFeedData.Channel.Description)
	for i, item := range rssFeedData.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeedData.Channel.Item[i] = item
	}

	return &rssFeedData, nil
}
