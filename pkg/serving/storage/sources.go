package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RSSSource represents an RSS source
type RSSSource struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	LastFetchedAt *time.Time `json:"lastFetchedAt,omitempty"`
}

// CreateSourceInput represents the input for creating an RSS source
type CreateSourceInput struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// UpdateSourceInput represents the input for updating an RSS source
type UpdateSourceInput struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// BatchCreateSourcesInput represents the input for batch creating RSS sources
type BatchCreateSourcesInput struct {
	Sources []CreateSourceInput `json:"sources"`
}

// BatchCreateSourcesResult represents the result of batch creating RSS sources
type BatchCreateSourcesResult struct {
	Sources []RSSSource  `json:"sources"`
	Errors  []BatchError `json:"errors"`
}

// BatchDeleteSourcesInput represents the input for batch deleting RSS sources
type BatchDeleteSourcesInput struct {
	SourceIDs []string `json:"sourceIds"`
}

// BatchDeleteSourcesResult represents the result of batch deleting RSS sources
type BatchDeleteSourcesResult struct {
	DeletedCount int          `json:"deletedCount"`
	Errors       []BatchError `json:"errors"`
}

// BatchError represents an error in a batch operation
type BatchError struct {
	Index     int    `json:"index,omitempty"`
	SourceID  string `json:"sourceId,omitempty"`
	ErrorType string `json:"errorType"`
	Message   string `json:"message"`
}

// CreateSource creates a new RSS source
func (s *SQLiteDB) CreateSource(input CreateSourceInput) (*RSSSource, error) {
	// Generate a new UUID for the source
	id := uuid.New().String()
	now := time.Now().UTC()

	// Insert the source into the database
	_, err := s.db.Exec(
		`INSERT INTO rss_sources (id, name, url, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		id, input.Name, input.URL, input.Description, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create RSS source: %w", err)
	}

	// Return the created source
	return &RSSSource{
		ID:          id,
		Name:        input.Name,
		URL:         input.URL,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetSource retrieves an RSS source by ID
func (s *SQLiteDB) GetSource(id string) (*RSSSource, error) {
	var source RSSSource
	var lastFetchedAt sql.NullTime

	err := s.db.QueryRow(
		`SELECT id, name, url, description, created_at, updated_at, last_fetched_at
		FROM rss_sources WHERE id = ?`,
		id,
	).Scan(
		&source.ID,
		&source.Name,
		&source.URL,
		&source.Description,
		&source.CreatedAt,
		&source.UpdatedAt,
		&lastFetchedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Source not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get RSS source: %w", err)
	}

	if lastFetchedAt.Valid {
		source.LastFetchedAt = &lastFetchedAt.Time
	}

	return &source, nil
}

// UpdateSource updates an RSS source
func (s *SQLiteDB) UpdateSource(id string, input UpdateSourceInput) (*RSSSource, error) {
	// Check if the source exists
	source, err := s.GetSource(id)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, nil // Source not found
	}

	// Update the source
	now := time.Now().UTC()
	_, err = s.db.Exec(
		`UPDATE rss_sources
		SET name = ?, url = ?, description = ?, updated_at = ?
		WHERE id = ?`,
		input.Name, input.URL, input.Description, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update RSS source: %w", err)
	}

	// Return the updated source
	source.Name = input.Name
	source.URL = input.URL
	source.Description = input.Description
	source.UpdatedAt = now

	return source, nil
}

// DeleteSource deletes an RSS source
func (s *SQLiteDB) DeleteSource(id string) error {
	// Check if the source exists
	source, err := s.GetSource(id)
	if err != nil {
		return err
	}
	if source == nil {
		return nil // Source not found
	}

	// Delete the source
	_, err = s.db.Exec("DELETE FROM rss_sources WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete RSS source: %w", err)
	}

	return nil
}

// ListSources lists RSS sources with pagination
func (s *SQLiteDB) ListSources(limit int, nextToken string) ([]RSSSource, string, error) {
	// Default limit if not specified
	if limit <= 0 {
		limit = 50
	}

	// Build the query
	query := `
		SELECT id, name, url, description, created_at, updated_at, last_fetched_at
		FROM rss_sources
	`
	args := []interface{}{}

	// Add pagination if nextToken is provided
	if nextToken != "" {
		query += " WHERE id > ?"
		args = append(args, nextToken)
	}

	// Add ordering and limit
	query += " ORDER BY id ASC LIMIT ?"
	args = append(args, limit+1) // Fetch one extra to determine if there are more results

	// Execute the query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list RSS sources: %w", err)
	}
	defer rows.Close()

	// Process the results
	var sources []RSSSource
	for rows.Next() {
		var source RSSSource
		var lastFetchedAt sql.NullTime

		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.URL,
			&source.Description,
			&source.CreatedAt,
			&source.UpdatedAt,
			&lastFetchedAt,
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan RSS source: %w", err)
		}

		if lastFetchedAt.Valid {
			source.LastFetchedAt = &lastFetchedAt.Time
		}

		sources = append(sources, source)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating over RSS sources: %w", err)
	}

	// Determine if there are more results and set the next token
	var newNextToken string
	if len(sources) > limit {
		newNextToken = sources[limit-1].ID
		sources = sources[:limit] // Remove the extra item
	}

	return sources, newNextToken, nil
}

// BatchCreateSources creates multiple RSS sources
func (s *SQLiteDB) BatchCreateSources(input BatchCreateSourcesInput) (*BatchCreateSourcesResult, error) {
	result := &BatchCreateSourcesResult{
		Sources: []RSSSource{},
		Errors:  []BatchError{},
	}

	// Process each source
	for i, sourceInput := range input.Sources {
		source, err := s.CreateSource(sourceInput)
		if err != nil {
			result.Errors = append(result.Errors, BatchError{
				Index:     i,
				ErrorType: "CreateError",
				Message:   err.Error(),
			})
		} else {
			result.Sources = append(result.Sources, *source)
		}
	}

	return result, nil
}

// BatchDeleteSources deletes multiple RSS sources
func (s *SQLiteDB) BatchDeleteSources(input BatchDeleteSourcesInput) (*BatchDeleteSourcesResult, error) {
	result := &BatchDeleteSourcesResult{
		DeletedCount: 0,
		Errors:       []BatchError{},
	}

	// Process each source ID
	for _, id := range input.SourceIDs {
		err := s.DeleteSource(id)
		if err != nil {
			result.Errors = append(result.Errors, BatchError{
				SourceID:  id,
				ErrorType: "DeleteError",
				Message:   err.Error(),
			})
		} else {
			result.DeletedCount++
		}
	}

	return result, nil
}

// UpdateSourceLastFetchedAt updates the last_fetched_at field of an RSS source
func (s *SQLiteDB) UpdateSourceLastFetchedAt(id string, lastFetchedAt time.Time) error {
	_, err := s.db.Exec(
		"UPDATE rss_sources SET last_fetched_at = ? WHERE id = ?",
		lastFetchedAt, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update last_fetched_at: %w", err)
	}
	return nil
}
