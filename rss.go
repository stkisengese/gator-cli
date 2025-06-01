package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
)

// RSSFeed represents the top-level structure of an RSS feed
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents an individual item in an RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// fetchFeed fetches and parses an RSS feed from the given URL
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return parseFeed(data)
}

// parseFeed parses XML data into an RSSFeed struct
func parseFeed(data []byte) (*RSSFeed, error) {
	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("error unmarshaling XML: %w", err)
	}

	// Unescape HTML entities in all text fields
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Items {
		item := &feed.Channel.Items[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

// String provides a human-readable representation of the RSSFeed
func (f *RSSFeed) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Channel: %s\n", f.Channel.Title))
	sb.WriteString(fmt.Sprintf("Link: %s\n", f.Channel.Link))
	sb.WriteString(fmt.Sprintf("Description: %s\n", f.Channel.Description))
	sb.WriteString("\nItems:\n")

	for i, item := range f.Channel.Items {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Title))
		sb.WriteString(fmt.Sprintf("   Link: %s\n", item.Link))
		sb.WriteString(fmt.Sprintf("   Published: %s\n", item.PubDate))
		sb.WriteString(fmt.Sprintf("   Description: %s\n\n", item.Description))
	}

	return sb.String()
}
