package riffle

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Run starts the riffle application
func Run() error {
	opmlFile := flag.String("opml", "", "Path to OPML file")
	flag.Parse()

	if *opmlFile == "" {
		return fmt.Errorf("OPML file path is required")
	}

	feeds, err := ParseOPML(*opmlFile)
	if err != nil {
		return fmt.Errorf("failed to parse OPML file: %w", err)
	}

	ctx := context.Background()
	for _, feed := range feeds {
		fmt.Printf("\nFeed: %s (%s)\n", feed.Title, feed.URL)
		fmt.Println(strings.Repeat("-", 80))

		articles, err := FetchLatestArticles(ctx, feed.URL, 3)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching articles from %s: %v\n", feed.URL, err)
			continue
		}

		for _, article := range articles {
			fmt.Printf("\nTitle: %s\n", article.Title)
			fmt.Printf("Published: %s\n", article.PublishedAt.Format(time.RFC3339))
			fmt.Printf("Summary: %s\n", article.Summary)
		}
		fmt.Println(strings.Repeat("-", 80))
	}

	return nil
}
