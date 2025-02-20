package app

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/flyer103/riffle/pkg/riffle"
	"github.com/spf13/cobra"
)

func runRiffle(cmd *cobra.Command, args []string) error {
	feeds, err := riffle.ParseOPML(opmlFile)
	if err != nil {
		return fmt.Errorf("failed to parse OPML file: %w", err)
	}

	analyzer, err := riffle.NewContentAnalyzer(interestsFile)
	if err != nil {
		return fmt.Errorf("failed to initialize content analyzer: %w", err)
	}

	ctx := context.Background()

	// Store all article scores for final recommendation
	var allScores []riffle.ArticleScore

	for _, feed := range feeds {
		fmt.Printf("\nFeed: %s (%s)\n", feed.Title, feed.URL)
		fmt.Println(strings.Repeat("-", 80))

		articles, err := riffle.FetchLatestArticles(ctx, feed.URL, 3)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching articles from %s: %v\n", feed.URL, err)
			continue
		}

		// Analyze and score each article
		var feedScores []riffle.ArticleScore
		for _, article := range articles {
			score, err := analyzer.AnalyzeArticle(&article)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing article '%s': %v\n", article.Title, err)
				continue
			}
			feedScores = append(feedScores, score)
			allScores = append(allScores, score)

			fmt.Printf("\nTitle: %s\n", article.Title)
			fmt.Printf("Published: %s\n", article.PublishedAt.Format(time.RFC3339))
			fmt.Printf("Summary: %s\n", article.Summary)
			fmt.Printf("Scores:\n")
			fmt.Printf("  - Interest Match: %.2f\n", score.InterestScore)
			fmt.Printf("  - Content Quality: %.2f\n", score.ContentScore)
			fmt.Printf("  - Overall: %.2f\n", score.Score)
		}

		// Print the highest-value article for this feed
		if len(feedScores) > 0 {
			sort.Slice(feedScores, func(i, j int) bool {
				return feedScores[i].Score > feedScores[j].Score
			})
			fmt.Printf("\nHighest Value Article in this feed: %s\n", feedScores[0].Article.Title)
			fmt.Printf("Scores:\n")
			fmt.Printf("  - Interest Match: %.2f\n", feedScores[0].InterestScore)
			fmt.Printf("  - Content Quality: %.2f\n", feedScores[0].ContentScore)
			fmt.Printf("  - Overall: %.2f\n", feedScores[0].Score)
		}

		fmt.Println(strings.Repeat("-", 80))
	}

	// Print overall highest-value article recommendation
	if len(allScores) > 0 {
		sort.Slice(allScores, func(i, j int) bool {
			return allScores[i].Score > allScores[j].Score
		})

		fmt.Printf("\n🌟 OVERALL HIGHEST VALUE ARTICLE RECOMMENDATION 🌟\n")
		fmt.Printf("Title: %s\n", allScores[0].Article.Title)
		fmt.Printf("Published: %s\n", allScores[0].Article.PublishedAt.Format(time.RFC3339))
		fmt.Printf("Summary: %s\n", allScores[0].Article.Summary)
		fmt.Printf("Scores:\n")
		fmt.Printf("  - Interest Match: %.2f\n", allScores[0].InterestScore)
		fmt.Printf("  - Content Quality: %.2f\n", allScores[0].ContentScore)
		fmt.Printf("  - Overall: %.2f\n", allScores[0].Score)
	}

	return nil
}
