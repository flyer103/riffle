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

// FetchLatestArticles fetches the latest n articles from a feed URL
func FetchLatestArticles(ctx context.Context, feedURL string, n int) ([]Article, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	articles := make([]Article, 0, n)
	for i := 0; i < len(feed.Items) && i < n; i++ {
		item := feed.Items[i]

		// Use published date if available, otherwise use updated date
		var pubDate time.Time
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			pubDate = *item.UpdatedParsed
		} else {
			pubDate = time.Now() // fallback to current time if no date available
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
	}

	return articles, nil
}
