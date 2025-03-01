package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// RecommendationFeedback represents user feedback on a recommended content item
type RecommendationFeedback struct {
	ID        string    `json:"id"`
	ContentID string    `json:"contentId"`
	UserID    string    `json:"userId"`
	Rating    int       `json:"rating"` // 1-5 scale
	Timestamp time.Time `json:"timestamp"`
	Comment   string    `json:"comment,omitempty"`
}

// CreateRecommendationFeedbackInput represents the input for creating recommendation feedback
type CreateRecommendationFeedbackInput struct {
	ContentID string `json:"contentId"`
	UserID    string `json:"userId"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment,omitempty"`
}

// RecommendationResult represents a recommended content item with its score
type RecommendationResult struct {
	Content      RSSContent `json:"content"`
	Score        float64    `json:"score"`
	RecommendFor string     `json:"recommendFor,omitempty"`
}

// GetRecommendationsInput represents the input for getting recommendations
type GetRecommendationsInput struct {
	UserID    string   `json:"userId,omitempty"`
	SourceIDs []string `json:"sourceIds,omitempty"`
	Limit     int      `json:"limit"`
}

// CreateRecommendationFeedback creates a new recommendation feedback entry
func (s *SQLiteDB) CreateRecommendationFeedback(input CreateRecommendationFeedbackInput) (*RecommendationFeedback, error) {
	// Validate rating
	if input.Rating < 1 || input.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Check if the content exists
	content, err := s.GetContent(input.ContentID)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, fmt.Errorf("content with ID %s not found", input.ContentID)
	}

	// Generate a new ID for the feedback
	id := generateUUID()
	now := time.Now().UTC()

	// Insert the feedback into the database
	_, err = s.db.Exec(
		`INSERT INTO recommendation_feedback (id, content_id, user_id, rating, timestamp, comment)
		VALUES (?, ?, ?, ?, ?, ?)`,
		id, input.ContentID, input.UserID, input.Rating, now, input.Comment,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create recommendation feedback: %w", err)
	}

	// Return the created feedback
	return &RecommendationFeedback{
		ID:        id,
		ContentID: input.ContentID,
		UserID:    input.UserID,
		Rating:    input.Rating,
		Timestamp: now,
		Comment:   input.Comment,
	}, nil
}

// GetRecommendations retrieves recommended content items
func (s *SQLiteDB) GetRecommendations(input GetRecommendationsInput) ([]RecommendationResult, error) {
	// Default limit if not specified
	if input.Limit <= 0 {
		input.Limit = 10
	}

	// Build the query
	// This is a simplified recommendation algorithm that:
	// 1. Gets recent content (last 7 days)
	// 2. Prioritizes content from sources with higher average ratings (if user has given feedback)
	// 3. Excludes content the user has already rated
	// 4. Sorts by a combination of recency and source popularity
	query := `
		SELECT 
			c.id, c.source_id, c.title, c.link, c.description, c.published_at, c.fetched_at,
			CASE
				WHEN avg_ratings.avg_rating IS NOT NULL THEN avg_ratings.avg_rating * 0.7 + (1.0 - ((JULIANDAY('now') - JULIANDAY(c.published_at)) / 7.0)) * 0.3
				ELSE (1.0 - ((JULIANDAY('now') - JULIANDAY(c.published_at)) / 7.0))
			END as score
		FROM 
			rss_contents c
		LEFT JOIN (
			SELECT 
				s.id as source_id, 
				AVG(rf.rating) as avg_rating
			FROM 
				rss_sources s
			JOIN 
				rss_contents rc ON s.id = rc.source_id
			JOIN 
				recommendation_feedback rf ON rc.id = rf.content_id
			WHERE 
				rf.user_id = ?
			GROUP BY 
				s.id
		) avg_ratings ON c.source_id = avg_ratings.source_id
		WHERE 
			c.published_at >= datetime('now', '-7 day')
	`
	args := []interface{}{input.UserID}

	// Add filter for specific sources if provided
	if len(input.SourceIDs) > 0 {
		query += " AND c.source_id IN (" + createPlaceholders(len(input.SourceIDs)) + ")"
		for _, sourceID := range input.SourceIDs {
			args = append(args, sourceID)
		}
	}

	// Exclude content the user has already rated
	if input.UserID != "" {
		query += `
			AND c.id NOT IN (
				SELECT content_id FROM recommendation_feedback WHERE user_id = ?
			)
		`
		args = append(args, input.UserID)
	}

	// Add ordering and limit
	query += " ORDER BY score DESC LIMIT ?"
	args = append(args, input.Limit)

	// Execute the query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}
	defer rows.Close()

	// Process the results
	var recommendations []RecommendationResult
	for rows.Next() {
		var content RSSContent
		var score float64
		err := rows.Scan(
			&content.ID,
			&content.SourceID,
			&content.Title,
			&content.Link,
			&content.Description,
			&content.PublishedAt,
			&content.FetchedAt,
			&score,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recommendation: %w", err)
		}

		recommendations = append(recommendations, RecommendationResult{
			Content:      content,
			Score:        score,
			RecommendFor: input.UserID,
		})
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over recommendations: %w", err)
	}

	return recommendations, nil
}

// GetUserFeedback retrieves all feedback given by a user
func (s *SQLiteDB) GetUserFeedback(userID string) ([]RecommendationFeedback, error) {
	// Query the feedback
	rows, err := s.db.Query(
		`SELECT id, content_id, user_id, rating, timestamp, comment
		FROM recommendation_feedback
		WHERE user_id = ?
		ORDER BY timestamp DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user feedback: %w", err)
	}
	defer rows.Close()

	// Process the results
	var feedbacks []RecommendationFeedback
	for rows.Next() {
		var feedback RecommendationFeedback
		var comment sql.NullString
		err := rows.Scan(
			&feedback.ID,
			&feedback.ContentID,
			&feedback.UserID,
			&feedback.Rating,
			&feedback.Timestamp,
			&comment,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan feedback: %w", err)
		}

		if comment.Valid {
			feedback.Comment = comment.String
		}

		feedbacks = append(feedbacks, feedback)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over feedback: %w", err)
	}

	return feedbacks, nil
}

// Helper function to create a string of SQL placeholders
func createPlaceholders(count int) string {
	if count <= 0 {
		return ""
	}

	placeholders := "?"
	for i := 1; i < count; i++ {
		placeholders += ",?"
	}
	return placeholders
}

// Helper function to generate a UUID
func generateUUID() string {
	// Using a simple implementation since we can't import uuid package yet
	// This should be replaced with a proper UUID implementation
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
