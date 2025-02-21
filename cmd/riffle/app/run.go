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

	// Print feed articles
	fmt.Printf("\n📰 ARTICLES FROM THE LAST 2 DAYS\n")
	fmt.Println(strings.Repeat("=", 80))

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

		fmt.Printf("\nFeed: %s\n", feed.Title)
		fmt.Println(strings.Repeat("-", 80))

		// Analyze and score each article
		for _, article := range articles {
			score, err := analyzer.AnalyzeArticle(&article)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing article '%s': %v\n", article.Title, err)
				continue
			}
			allScores = append(allScores, score)

			fmt.Printf("\nTitle: %s\n", article.Title)
			fmt.Printf("URL: %s\n", article.URL)
			fmt.Printf("Published: %s\n", article.PublishedAt.Format(time.RFC3339))
			if article.Summary != "" {
				fmt.Printf("Summary: %s\n", article.Summary)
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

		fmt.Printf("\n🌟 MOST VALUABLE ARTICLES ACROSS ALL SOURCES 🌟\n")
		fmt.Println(strings.Repeat("=", 80))
		for i := 0; i < len(allScores) && i < topCount; i++ {
			score := allScores[i]
			fmt.Printf("\n%d. %s\n", i+1, score.Article.Title)
			fmt.Printf("Source: %s\n", getFeedTitleByURL(feeds, score.Article.URL))
			fmt.Printf("URL: %s\n", score.Article.URL)
			fmt.Printf("Published: %s\n", score.Article.PublishedAt.Format(time.RFC3339))
			if score.Article.Summary != "" {
				fmt.Printf("Summary: %s\n", score.Article.Summary)
			}
			fmt.Printf("Scores:\n")
			fmt.Printf("  - Interest Match: %.2f\n", score.InterestScore)
			fmt.Printf("  - Content Quality: %.2f\n", score.ContentScore)
			fmt.Printf("  - Overall: %.2f\n", score.Score)
			fmt.Printf("Why recommended: %s\n", generateRecommendationReason(score))

			// Add Perplexity analysis
			analysis, err := riffle.AnalyzeWithPerplexity(score.Article, modelName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not get AI analysis: %v\n", err)
				continue
			}

			fmt.Printf("\n📊 AI Analysis:\n%s\n", analysis.Content)
		}
		fmt.Println(strings.Repeat("=", 80))
	} else {
		fmt.Printf("\n⚠️ No articles found from the last 2 days in any feed.\n")
	}

	return nil
}

// getFeedTitleByURL returns the feed title for a given article URL
func getFeedTitleByURL(feeds []riffle.Feed, articleURL string) string {
	for _, feed := range feeds {
		baseURL := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(feed.URL, "/feed.xml"), "/feed.atom"), "/rss")
		if strings.Contains(articleURL, baseURL) {
			return feed.Title
		}
	}
	return "Unknown Source"
}
