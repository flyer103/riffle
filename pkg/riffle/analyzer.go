package riffle

import (
	"html"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// ArticleScore represents the analyzed value of an article
type ArticleScore struct {
	Article *Article
	Score   float64
}

// ContentAnalyzer analyzes article content and scores it
type ContentAnalyzer struct {
	// Keywords that indicate valuable content
	valueKeywords []string
}

// NewContentAnalyzer creates a new ContentAnalyzer
func NewContentAnalyzer() *ContentAnalyzer {
	return &ContentAnalyzer{
		valueKeywords: []string{
			"research", "study", "analysis", "guide", "tutorial",
			"introduction", "review", "comparison", "best practices",
			"how to", "explained", "deep dive", "architecture",
			"performance", "security", "scalability",
		},
	}
}

// AnalyzeArticle scores an article based on various factors
func (ca *ContentAnalyzer) AnalyzeArticle(article *Article) (ArticleScore, error) {
	// Get the full content if available
	content := article.Summary
	if article.Content != "" {
		content = article.Content
	}

	// Clean and prepare the content
	cleanContent := html.UnescapeString(content)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(cleanContent))
	if err != nil {
		return ArticleScore{Article: article}, err
	}

	// Calculate various scores
	textScore := ca.calculateTextScore(doc)
	keywordScore := ca.calculateKeywordScore(doc.Text())
	linkScore := ca.calculateLinkScore(doc)

	// Combine scores with weights
	totalScore := (textScore * 0.4) + (keywordScore * 0.4) + (linkScore * 0.2)

	return ArticleScore{
		Article: article,
		Score:   totalScore,
	}, nil
}

// calculateTextScore evaluates the quality of the text content
func (ca *ContentAnalyzer) calculateTextScore(doc *goquery.Document) float64 {
	text := doc.Text()
	textLen := utf8.RuneCountInString(text)

	// Score based on content length (0-1)
	lengthScore := float64(0)
	switch {
	case textLen > 2000:
		lengthScore = 1.0
	case textLen > 1000:
		lengthScore = 0.8
	case textLen > 500:
		lengthScore = 0.6
	case textLen > 200:
		lengthScore = 0.4
	default:
		lengthScore = 0.2
	}

	return lengthScore
}

// calculateKeywordScore evaluates the presence of valuable keywords
func (ca *ContentAnalyzer) calculateKeywordScore(text string) float64 {
	text = strings.ToLower(text)
	var keywordCount int

	for _, keyword := range ca.valueKeywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			keywordCount++
		}
	}

	// Score based on keyword presence (0-1)
	return float64(keywordCount) / float64(len(ca.valueKeywords))
}

// calculateLinkScore evaluates the quality of links
func (ca *ContentAnalyzer) calculateLinkScore(doc *goquery.Document) float64 {
	links := doc.Find("a")
	linkCount := links.Length()

	if linkCount == 0 {
		return 0.5 // Neutral score for no links
	}

	// Count external links (usually more valuable)
	var externalLinks int
	links.Each(func(_ int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if exists && strings.HasPrefix(href, "http") {
			externalLinks++
		}
	})

	// Score based on external link ratio (0-1)
	return float64(externalLinks) / float64(linkCount)
}
