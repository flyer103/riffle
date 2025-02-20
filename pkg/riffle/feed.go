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

		articles = append(articles, Article{
			Title:       item.Title,
			Summary:     item.Description,
			PublishedAt: pubDate,
		})
	}

	return articles, nil
}
