package riffle

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

// Article represents a single article from a feed
type Article struct {
	Title       string
	Summary     string
	Content     string // Full content of the article
	URL         string // URL of the article
	PublishedAt time.Time
}

// FetchLatestArticles fetches articles from the last 2 days from a feed URL, up to n articles
func FetchLatestArticles(ctx context.Context, feedURL string, n int) ([]Article, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	// Calculate the cutoff time (2 days ago)
	cutoffTime := time.Now().AddDate(0, 0, -2)

	articles := make([]Article, 0, n)
	for _, item := range feed.Items {
		// Get the article's publication time
		var pubDate time.Time
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			pubDate = *item.UpdatedParsed
		} else {
			continue // Skip articles with no date
		}

		// Skip articles older than 2 days
		if pubDate.Before(cutoffTime) {
			continue
		}

		// Get the full content if available, otherwise use description
		content := item.Content
		if content == "" {
			content = item.Description
		}

		// Get the article URL
		url := item.Link
		if url == "" {
			url = item.GUID // fallback to GUID if link is not available
		}

		articles = append(articles, Article{
			Title:       item.Title,
			Summary:     item.Description,
			Content:     content,
			URL:         url,
			PublishedAt: pubDate,
		})

		// Stop if we have enough articles
		if len(articles) >= n {
			break
		}
	}

	return articles, nil
}
