package handlers

import (
	"net/http"
	"strconv"

	"github.com/flyer103/riffle/pkg/serving/storage"
	"github.com/gin-gonic/gin"
)

// SourcesHandler handles API requests for RSS sources
type SourcesHandler struct {
	db *storage.SQLiteDB
}

// NewSourcesHandler creates a new SourcesHandler
func NewSourcesHandler(db *storage.SQLiteDB) *SourcesHandler {
	return &SourcesHandler{
		db: db,
	}
}

// ListSources handles GET /sources
func (h *SourcesHandler) ListSources(c *gin.Context) {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	nextToken := c.Query("nextToken")

	// Get sources from the database
	sources, newNextToken, err := h.db.ListSources(limit, nextToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list sources: " + err.Error(),
		})
		return
	}

	// Return the sources
	c.JSON(http.StatusOK, gin.H{
		"sources":   sources,
		"nextToken": newNextToken,
	})
}

// GetSource handles GET /sources/:id
func (h *SourcesHandler) GetSource(c *gin.Context) {
	// Get the source ID from the URL
	id := c.Param("id")

	// Get the source from the database
	source, err := h.db.GetSource(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get source: " + err.Error(),
		})
		return
	}

	// Check if the source exists
	if source == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Source not found",
		})
		return
	}

	// Return the source
	c.JSON(http.StatusOK, source)
}

// CreateSource handles POST /sources
func (h *SourcesHandler) CreateSource(c *gin.Context) {
	// Parse the request body
	var input storage.CreateSourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Create the source
	source, err := h.db.CreateSource(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create source: " + err.Error(),
		})
		return
	}

	// Return the created source
	c.JSON(http.StatusCreated, source)
}

// UpdateSource handles PUT /sources/:id
func (h *SourcesHandler) UpdateSource(c *gin.Context) {
	// Get the source ID from the URL
	id := c.Param("id")

	// Parse the request body
	var input storage.UpdateSourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update the source
	source, err := h.db.UpdateSource(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update source: " + err.Error(),
		})
		return
	}

	// Check if the source exists
	if source == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Source not found",
		})
		return
	}

	// Return the updated source
	c.JSON(http.StatusOK, source)
}

// DeleteSource handles DELETE /sources/:id
func (h *SourcesHandler) DeleteSource(c *gin.Context) {
	// Get the source ID from the URL
	id := c.Param("id")

	// Delete the source
	err := h.db.DeleteSource(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete source: " + err.Error(),
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Source deleted successfully",
	})
}

// BatchCreateSources handles POST /sources/batch
func (h *SourcesHandler) BatchCreateSources(c *gin.Context) {
	// Parse the request body
	var input storage.BatchCreateSourcesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Create the sources
	result, err := h.db.BatchCreateSources(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to batch create sources: " + err.Error(),
		})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, result)
}

// BatchDeleteSources handles DELETE /sources/batch
func (h *SourcesHandler) BatchDeleteSources(c *gin.Context) {
	// Parse the request body
	var input storage.BatchDeleteSourcesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Delete the sources
	result, err := h.db.BatchDeleteSources(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to batch delete sources: " + err.Error(),
		})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, result)
}
