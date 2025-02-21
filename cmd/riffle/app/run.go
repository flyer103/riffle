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

// generateRecommendationReason generates a detailed explanation of why an article is recommended
func generateRecommendationReason(score riffle.ArticleScore) string {
	var reasons []string

	// Analyze interest match score
	if score.InterestScore >= 0.7 {
		reasons = append(reasons, "Strongly matches your interests")
	} else if score.InterestScore >= 0.5 {
		reasons = append(reasons, "Moderately aligns with your interests")
	}

	// Analyze content quality components
	if score.ContentScore >= 0.7 {
		reasons = append(reasons, "High-quality content with detailed information")
	} else if score.ContentScore >= 0.5 {
		reasons = append(reasons, "Good content quality")
	}

	// If no specific reasons found, provide a general reason
	if len(reasons) == 0 {
		reasons = append(reasons, "Balanced combination of relevance and quality")
	}

	return strings.Join(reasons, ", ")
}

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
	// Track feeds without recent updates
	var noUpdateFeeds []string

	for _, feed := range feeds {
		articles, err := riffle.FetchLatestArticles(ctx, feed.URL, articleCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching articles from %s: %v\n", feed.URL, err)
			continue
		}

		// Skip printing feed details if no recent articles
		if len(articles) == 0 {
			noUpdateFeeds = append(noUpdateFeeds, feed.Title)
			continue
		}

		fmt.Printf("\nFeed: %s (%s)\n", feed.Title, feed.URL)
		fmt.Println(strings.Repeat("-", 80))

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
			fmt.Printf("URL: %s\n", article.URL)
			fmt.Printf("Published: %s\n", article.PublishedAt.Format(time.RFC3339))
			fmt.Printf("Summary: %s\n", article.Summary)
			fmt.Printf("Scores:\n")
			fmt.Printf("  - Interest Match: %.2f\n", score.InterestScore)
			fmt.Printf("  - Content Quality: %.2f\n", score.ContentScore)
			fmt.Printf("  - Overall: %.2f\n", score.Score)
		}

		// Print the highest-value articles for this feed
		if len(feedScores) > 0 {
			sort.Slice(feedScores, func(i, j int) bool {
				return feedScores[i].Score > feedScores[j].Score
			})

			fmt.Printf("\nTop %d Articles in this feed:\n", topCount)
			for i := 0; i < len(feedScores) && i < topCount; i++ {
				score := feedScores[i]
				fmt.Printf("%d. %s\n", i+1, score.Article.Title)
				fmt.Printf("   URL: %s\n", score.Article.URL)
				fmt.Printf("   Scores:\n")
				fmt.Printf("   - Interest Match: %.2f\n", score.InterestScore)
				fmt.Printf("   - Content Quality: %.2f\n", score.ContentScore)
				fmt.Printf("   - Overall: %.2f\n", score.Score)
				fmt.Printf("   Why recommended: %s\n", generateRecommendationReason(score))
			}
		}

		fmt.Println(strings.Repeat("-", 80))
	}

	// Print feeds without recent updates
	if len(noUpdateFeeds) > 0 {
		fmt.Printf("\n📅 RSS Sources Without Recent Updates (Last 2 Days):\n")
		fmt.Println(strings.Repeat("-", 50))
		for i, feedTitle := range noUpdateFeeds {
			fmt.Printf("%d. %s\n", i+1, feedTitle)
		}
		fmt.Println(strings.Repeat("-", 50))
	}

	// Print overall highest-value article recommendations
	if len(allScores) > 0 {
		sort.Slice(allScores, func(i, j int) bool {
			return allScores[i].Score > allScores[j].Score
		})

		fmt.Printf("\n🌟 OVERALL TOP %d ARTICLE RECOMMENDATIONS 🌟\n", topCount)
		for i := 0; i < len(allScores) && i < topCount; i++ {
			score := allScores[i]
			fmt.Printf("\n%d. %s\n", i+1, score.Article.Title)
			fmt.Printf("URL: %s\n", score.Article.URL)
			fmt.Printf("Published: %s\n", score.Article.PublishedAt.Format(time.RFC3339))
			fmt.Printf("Summary: %s\n", score.Article.Summary)
			fmt.Printf("Scores:\n")
			fmt.Printf("  - Interest Match: %.2f\n", score.InterestScore)
			fmt.Printf("  - Content Quality: %.2f\n", score.ContentScore)
			fmt.Printf("  - Overall: %.2f\n", score.Score)
			fmt.Printf("Why recommended: %s\n", generateRecommendationReason(score))
		}
	} else {
		fmt.Printf("\n⚠️ No articles found from the last 2 days in any feed.\n")
	}

	return nil
}
