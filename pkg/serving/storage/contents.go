package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RSSContent represents an RSS content item
type RSSContent struct {
	ID          string     `json:"id"`
	SourceID    string     `json:"sourceId"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Description string     `json:"description"`
	Content     string     `json:"content,omitempty"`
	PublishedAt time.Time  `json:"publishedAt"`
	FetchedAt   time.Time  `json:"fetchedAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	Author      string     `json:"author,omitempty"`
	Categories  []string   `json:"categories,omitempty"`
}

// UpdateContentInput represents the input for updating an RSS content item
type UpdateContentInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Categories  []string `json:"categories"`
}

// BatchDeleteContentsInput represents the input for batch deleting RSS content items
type BatchDeleteContentsInput struct {
	ContentIDs []string `json:"contentIds"`
}

// BatchDeleteContentsResult represents the result of batch deleting RSS content items
type BatchDeleteContentsResult struct {
	DeletedCount int          `json:"deletedCount"`
	Errors       []BatchError `json:"errors"`
}

// FetchJob represents an RSS content fetch job
type FetchJob struct {
	ID             string     `json:"jobId"`
	Status         string     `json:"status"`
	StartedAt      time.Time  `json:"startedAt"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
	ItemsProcessed int        `json:"itemsProcessed"`
	SourceID       *string    `json:"sourceId,omitempty"`
	Days           int        `json:"days"`
	Errors         []string   `json:"errors,omitempty"`
}

// CreateContent creates a new RSS content item
func (s *SQLiteDB) CreateContent(content *RSSContent) error {
	// Generate a new UUID if not provided
	if content.ID == "" {
		content.ID = uuid.New().String()
	}

	// Set fetched time if not provided
	if content.FetchedAt.IsZero() {
		content.FetchedAt = time.Now().UTC()
	}

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the content into the database
	_, err = tx.Exec(
		`INSERT INTO rss_contents (id, source_id, title, link, description, content, published_at, fetched_at, author)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		content.ID, content.SourceID, content.Title, content.Link, content.Description,
		content.Content, content.PublishedAt, content.FetchedAt, content.Author,
	)
	if err != nil {
		return fmt.Errorf("failed to create RSS content: %w", err)
	}

	// Insert categories if provided
	if len(content.Categories) > 0 {
		for _, category := range content.Categories {
			_, err = tx.Exec(
				"INSERT INTO content_categories (content_id, category) VALUES (?, ?)",
				content.ID, category,
			)
			if err != nil {
				return fmt.Errorf("failed to insert category: %w", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetContent retrieves an RSS content item by ID
func (s *SQLiteDB) GetContent(id string) (*RSSContent, error) {
	// Query the content
	var content RSSContent
	var updatedAt sql.NullTime
	var author sql.NullString
	var contentText sql.NullString

	err := s.db.QueryRow(
		`SELECT id, source_id, title, link, description, content, published_at, fetched_at, updated_at, author
		FROM rss_contents WHERE id = ?`,
		id,
	).Scan(
		&content.ID,
		&content.SourceID,
		&content.Title,
		&content.Link,
		&content.Description,
		&contentText,
		&content.PublishedAt,
		&content.FetchedAt,
		&updatedAt,
		&author,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Content not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get RSS content: %w", err)
	}

	// Set optional fields if present
	if updatedAt.Valid {
		content.UpdatedAt = &updatedAt.Time
	}
	if author.Valid {
		content.Author = author.String
	}
	if contentText.Valid {
		content.Content = contentText.String
	}

	// Query categories
	rows, err := s.db.Query(
		"SELECT category FROM content_categories WHERE content_id = ?",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	// Process categories
	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over categories: %w", err)
	}

	content.Categories = categories
	return &content, nil
}

// UpdateContent updates an RSS content item
func (s *SQLiteDB) UpdateContent(id string, input UpdateContentInput) (*RSSContent, error) {
	// Check if the content exists
	content, err := s.GetContent(id)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, nil // Content not found
	}

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update the content
	now := time.Now().UTC()
	_, err = tx.Exec(
		`UPDATE rss_contents
		SET title = ?, description = ?, content = ?, updated_at = ?
		WHERE id = ?`,
		input.Title, input.Description, input.Content, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update RSS content: %w", err)
	}

	// Delete existing categories
	_, err = tx.Exec("DELETE FROM content_categories WHERE content_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing categories: %w", err)
	}

	// Insert new categories
	for _, category := range input.Categories {
		_, err = tx.Exec(
			"INSERT INTO content_categories (content_id, category) VALUES (?, ?)",
			id, category,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert category: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Update the content object
	content.Title = input.Title
	content.Description = input.Description
	content.Content = input.Content
	content.Categories = input.Categories
	content.UpdatedAt = &now

	return content, nil
}

// DeleteContent deletes an RSS content item
func (s *SQLiteDB) DeleteContent(id string) error {
	// Check if the content exists
	content, err := s.GetContent(id)
	if err != nil {
		return err
	}
	if content == nil {
		return nil // Content not found
	}

	// Delete the content (categories will be deleted via CASCADE)
	_, err = s.db.Exec("DELETE FROM rss_contents WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete RSS content: %w", err)
	}

	return nil
}

// ListContents lists RSS content items with filtering and pagination
func (s *SQLiteDB) ListContents(sourceID string, startDate, endDate time.Time, limit int, nextToken string) ([]RSSContent, string, error) {
	// Default limit if not specified
	if limit <= 0 {
		limit = 50
	}

	// Build the query
	query := `
		SELECT c.id, c.source_id, c.title, c.link, c.description, c.published_at, c.fetched_at
		FROM rss_contents c
		WHERE 1=1
	`
	args := []interface{}{}

	// Add filters
	if sourceID != "" {
		query += " AND c.source_id = ?"
		args = append(args, sourceID)
	}
	if !startDate.IsZero() {
		query += " AND c.published_at >= ?"
		args = append(args, startDate)
	}
	if !endDate.IsZero() {
		query += " AND c.published_at <= ?"
		args = append(args, endDate)
	}
	if nextToken != "" {
		query += " AND c.id > ?"
		args = append(args, nextToken)
	}

	// Add ordering and limit
	query += " ORDER BY c.id ASC LIMIT ?"
	args = append(args, limit+1) // Fetch one extra to determine if there are more results

	// Execute the query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list RSS contents: %w", err)
	}
	defer rows.Close()

	// Process the results
	var contents []RSSContent
	for rows.Next() {
		var content RSSContent
		err := rows.Scan(
			&content.ID,
			&content.SourceID,
			&content.Title,
			&content.Link,
			&content.Description,
			&content.PublishedAt,
			&content.FetchedAt,
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan RSS content: %w", err)
		}
		contents = append(contents, content)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating over RSS contents: %w", err)
	}

	// Determine if there are more results and set the next token
	var newNextToken string
	if len(contents) > limit {
		newNextToken = contents[limit-1].ID
		contents = contents[:limit] // Remove the extra item
	}

	return contents, newNextToken, nil
}

// BatchDeleteContents deletes multiple RSS content items
func (s *SQLiteDB) BatchDeleteContents(input BatchDeleteContentsInput) (*BatchDeleteContentsResult, error) {
	result := &BatchDeleteContentsResult{
		DeletedCount: 0,
		Errors:       []BatchError{},
	}

	// Process each content ID
	for _, id := range input.ContentIDs {
		err := s.DeleteContent(id)
		if err != nil {
			result.Errors = append(result.Errors, BatchError{
				SourceID:  id, // Reusing SourceID field for ContentID
				ErrorType: "DeleteError",
				Message:   err.Error(),
			})
		} else {
			result.DeletedCount++
		}
	}

	return result, nil
}

// CreateFetchJob creates a new fetch job
func (s *SQLiteDB) CreateFetchJob(sourceID *string, days int) (*FetchJob, error) {
	// Generate a new UUID for the job
	id := uuid.New().String()
	now := time.Now().UTC()

	// Insert the job into the database
	_, err := s.db.Exec(
		`INSERT INTO fetch_jobs (id, status, started_at, source_id, days)
		VALUES (?, ?, ?, ?, ?)`,
		id, "pending", now, sourceID, days,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch job: %w", err)
	}

	// Return the created job
	return &FetchJob{
		ID:             id,
		Status:         "pending",
		StartedAt:      now,
		ItemsProcessed: 0,
		SourceID:       sourceID,
		Days:           days,
		Errors:         []string{},
	}, nil
}

// GetFetchJob retrieves a fetch job by ID
func (s *SQLiteDB) GetFetchJob(id string) (*FetchJob, error) {
	var job FetchJob
	var sourceID sql.NullString
	var completedAt sql.NullTime

	// Query the job
	err := s.db.QueryRow(
		`SELECT id, status, started_at, completed_at, items_processed, source_id, days
		FROM fetch_jobs WHERE id = ?`,
		id,
	).Scan(
		&job.ID,
		&job.Status,
		&job.StartedAt,
		&completedAt,
		&job.ItemsProcessed,
		&sourceID,
		&job.Days,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Job not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get fetch job: %w", err)
	}

	// Set optional fields if present
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	if sourceID.Valid {
		job.SourceID = &sourceID.String
	}

	// Query job errors
	rows, err := s.db.Query(
		"SELECT error_message FROM job_errors WHERE job_id = ? ORDER BY timestamp",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query job errors: %w", err)
	}
	defer rows.Close()

	// Process errors
	var errors []string
	for rows.Next() {
		var errorMsg string
		if err := rows.Scan(&errorMsg); err != nil {
			return nil, fmt.Errorf("failed to scan job error: %w", err)
		}
		errors = append(errors, errorMsg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over job errors: %w", err)
	}

	job.Errors = errors
	return &job, nil
}

// UpdateFetchJobStatus updates the status of a fetch job
func (s *SQLiteDB) UpdateFetchJobStatus(id string, status string, itemsProcessed int, errorMsg string) error {
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update job status
	var completedAt *time.Time
	if status == "completed" || status == "failed" {
		now := time.Now().UTC()
		completedAt = &now
	}

	if completedAt != nil {
		_, err = tx.Exec(
			"UPDATE fetch_jobs SET status = ?, items_processed = ?, completed_at = ? WHERE id = ?",
			status, itemsProcessed, completedAt, id,
		)
	} else {
		_, err = tx.Exec(
			"UPDATE fetch_jobs SET status = ?, items_processed = ? WHERE id = ?",
			status, itemsProcessed, id,
		)
	}
	if err != nil {
		return fmt.Errorf("failed to update fetch job status: %w", err)
	}

	// Add error message if provided
	if errorMsg != "" {
		_, err = tx.Exec(
			"INSERT INTO job_errors (job_id, error_message, timestamp) VALUES (?, ?, ?)",
			id, errorMsg, time.Now().UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert job error: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetContentByURL retrieves an RSS content item by URL
func (s *SQLiteDB) GetContentByURL(url string) (*RSSContent, error) {
	var id string
	err := s.db.QueryRow("SELECT id FROM rss_contents WHERE link = ?", url).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil // Content not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get RSS content by URL: %w", err)
	}

	return s.GetContent(id)
}

// SearchContents searches for RSS content items by keywords
func (s *SQLiteDB) SearchContents(keywords string, sourceID string, limit int) ([]RSSContent, error) {
	// Default limit if not specified
	if limit <= 0 {
		limit = 50
	}

	// Split keywords
	keywordList := strings.Split(keywords, ",")
	for i, k := range keywordList {
		keywordList[i] = strings.TrimSpace(k)
	}

	// Build the query
	query := `
		SELECT c.id, c.source_id, c.title, c.link, c.description, c.published_at, c.fetched_at
		FROM rss_contents c
		WHERE 1=1
	`
	args := []interface{}{}

	// Add source filter if provided
	if sourceID != "" {
		query += " AND c.source_id = ?"
		args = append(args, sourceID)
	}

	// Add keyword filters
	if len(keywordList) > 0 {
		query += " AND ("
		for i, keyword := range keywordList {
			if i > 0 {
				query += " OR "
			}
			query += "c.title LIKE ? OR c.description LIKE ? OR c.content LIKE ?"
			likePattern := "%" + keyword + "%"
			args = append(args, likePattern, likePattern, likePattern)
		}
		query += ")"
	}

	// Add ordering and limit
	query += " ORDER BY c.published_at DESC LIMIT ?"
	args = append(args, limit)

	// Execute the query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search RSS contents: %w", err)
	}
	defer rows.Close()

	// Process the results
	var contents []RSSContent
	for rows.Next() {
		var content RSSContent
		err := rows.Scan(
			&content.ID,
			&content.SourceID,
			&content.Title,
			&content.Link,
			&content.Description,
			&content.PublishedAt,
			&content.FetchedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan RSS content: %w", err)
		}
		contents = append(contents, content)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over RSS contents: %w", err)
	}

	return contents, nil
}
