package riffle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// ArticleScore represents the analyzed value of an article
type ArticleScore struct {
	Article       *Article
	Score         float64
	InterestScore float64 // Score based on user interests
	ContentScore  float64 // Score based on content quality
}

// ContentAnalyzer analyzes article content and scores it
type ContentAnalyzer struct {
	// Keywords that indicate valuable content
	valueKeywords []string
	// User's current interests
	interests []string
}

// NewContentAnalyzer creates a new ContentAnalyzer
func NewContentAnalyzer(interestsFile string) (*ContentAnalyzer, error) {
	analyzer := &ContentAnalyzer{
		valueKeywords: []string{
			"research", "study", "analysis", "guide", "tutorial",
			"introduction", "review", "comparison", "best practices",
			"how to", "explained", "deep dive", "architecture",
			"performance", "security", "scalability",
		},
	}

	// Load interests if file is provided
	if interestsFile != "" {
		file, err := os.Open(interestsFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			interest := strings.TrimSpace(scanner.Text())
			if interest != "" {
				analyzer.interests = append(analyzer.interests, interest)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return analyzer, nil
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

	// Calculate content quality scores
	textScore := ca.calculateTextScore(doc)
	keywordScore := ca.calculateKeywordScore(doc.Text())
	linkScore := ca.calculateLinkScore(doc)

	// Calculate interest relevance score
	interestScore := ca.calculateInterestScore(article.Title + " " + doc.Text())

	// Combine content quality scores (50% weight)
	contentScore := (textScore * 0.4) + (keywordScore * 0.4) + (linkScore * 0.2)

	// Final score is 50% interest relevance and 50% content quality
	totalScore := (interestScore * 0.5) + (contentScore * 0.5)

	return ArticleScore{
		Article:       article,
		Score:         totalScore,
		InterestScore: interestScore,
		ContentScore:  contentScore,
	}, nil
}

// calculateInterestScore evaluates how well the content matches user interests
func (ca *ContentAnalyzer) calculateInterestScore(text string) float64 {
	if len(ca.interests) == 0 {
		return 0.5 // Neutral score if no interests defined
	}

	text = strings.ToLower(text)
	var matchCount int

	for _, interest := range ca.interests {
		// Split interest into words for more flexible matching
		interestWords := strings.Fields(strings.ToLower(interest))

		// Count how many words from this interest appear in the text
		wordMatches := 0
		for _, word := range interestWords {
			if strings.Contains(text, word) {
				wordMatches++
			}
		}

		// Consider it a match if more than half of the words match
		if float64(wordMatches) >= float64(len(interestWords))*0.5 {
			matchCount++
		}
	}

	// Score based on the proportion of matching interests
	return float64(matchCount) / float64(len(ca.interests))
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

// PerplexityAnalysis represents the analysis result from Perplexity API
type PerplexityAnalysis struct {
	Content string
}

// AnalyzeWithPerplexity uses Perplexity API to analyze an article
func AnalyzeWithPerplexity(article *Article, model string) (*PerplexityAnalysis, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.perplexity.ai"
	}

	client := &http.Client{}
	prompt := fmt.Sprintf("Analyze this article and provide: 1) A concise summary 2) Key points 3) Why it's significant\n\nTitle: %s\nContent: %s",
		article.Title,
		article.Content)

	data := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert at analyzing articles and providing insightful summaries.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the raw response body for error logging
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(respBody))
	}

	// Check for API error response
	if errMsg, ok := result["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("API error: %v", errMsg)
	}

	// Extract the analysis from the response
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("invalid response format: missing or empty choices array, response: %s", string(respBody))
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message format in response: %s", string(respBody))
	}

	content, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content format in response: %s", string(respBody))
	}

	if content == "" {
		return nil, fmt.Errorf("received empty content from API")
	}

	return &PerplexityAnalysis{Content: content}, nil
}
