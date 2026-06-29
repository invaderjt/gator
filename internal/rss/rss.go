package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
	"time"
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := NewClient(10 * time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feedData *RSSFeed
	if err = xml.Unmarshal(body, &feedData); err != nil {
		log.Fatalf("Error unmarshalling XML: %v\n", err)
	}

	feedData.Channel.Title = html.UnescapeString(feedData.Channel.Title)
	feedData.Channel.Description = html.UnescapeString(feedData.Channel.Description)
	for i, item := range feedData.Channel.Item {
		feedData.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feedData.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return feedData, nil
}
