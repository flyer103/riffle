package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/flyer103/riffle/pkg/serving/storage"
	"github.com/gin-gonic/gin"
)

// ContentsHandler handles API requests for RSS contents
type ContentsHandler struct {
	db *storage.SQLiteDB
}

// NewContentsHandler creates a new ContentsHandler
func NewContentsHandler(db *storage.SQLiteDB) *ContentsHandler {
	return &ContentsHandler{
		db: db,
	}
}

// ListContents handles GET /contents
func (h *ContentsHandler) ListContents(c *gin.Context) {
	// Parse query parameters
	sourceID := c.Query("sourceId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	nextToken := c.Query("nextToken")

	// Parse date filters if provided
	var startDate, endDate time.Time
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid startDate format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
			})
			return
		}
	}
	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid endDate format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
			})
			return
		}
	}

	// Get contents from the database
	contents, newNextToken, err := h.db.ListContents(sourceID, startDate, endDate, limit, nextToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list contents: " + err.Error(),
		})
		return
	}

	// Return the contents
	c.JSON(http.StatusOK, gin.H{
		"contents":  contents,
		"nextToken": newNextToken,
	})
}

// GetContent handles GET /contents/:id
func (h *ContentsHandler) GetContent(c *gin.Context) {
	// Get the content ID from the URL
	id := c.Param("id")

	// Get the content from the database
	content, err := h.db.GetContent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get content: " + err.Error(),
		})
		return
	}

	// Check if the content exists
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Content not found",
		})
		return
	}

	// Return the content
	c.JSON(http.StatusOK, content)
}

// UpdateContent handles PUT /contents/:id
func (h *ContentsHandler) UpdateContent(c *gin.Context) {
	// Get the content ID from the URL
	id := c.Param("id")

	// Parse the request body
	var input storage.UpdateContentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update the content
	content, err := h.db.UpdateContent(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update content: " + err.Error(),
		})
		return
	}

	// Check if the content exists
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Content not found",
		})
		return
	}

	// Return the updated content
	c.JSON(http.StatusOK, content)
}

// DeleteContent handles DELETE /contents/:id
func (h *ContentsHandler) DeleteContent(c *gin.Context) {
	// Get the content ID from the URL
	id := c.Param("id")

	// Delete the content
	err := h.db.DeleteContent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete content: " + err.Error(),
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Content deleted successfully",
	})
}

// BatchDeleteContents handles DELETE /contents/batch
func (h *ContentsHandler) BatchDeleteContents(c *gin.Context) {
	// Parse the request body
	var input storage.BatchDeleteContentsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Delete the contents
	result, err := h.db.BatchDeleteContents(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to batch delete contents: " + err.Error(),
		})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, result)
}

// FetchContents handles POST /contents/fetch
func (h *ContentsHandler) FetchContents(c *gin.Context) {
	// Parse the request body
	type FetchRequest struct {
		SourceID *string `json:"sourceId"`
		Days     int     `json:"days"`
	}
	var req FetchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate days
	if req.Days <= 0 {
		req.Days = 7 // Default to 7 days if not specified or invalid
	}

	// Create a fetch job
	job, err := h.db.CreateFetchJob(req.SourceID, req.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create fetch job: " + err.Error(),
		})
		return
	}

	// Return the job
	c.JSON(http.StatusAccepted, gin.H{
		"jobId":  job.ID,
		"status": job.Status,
	})

	// Start the fetch process asynchronously
	// In a real implementation, this would be handled by a background worker
	// For now, we'll just update the job status to simulate the process
	go func() {
		// Simulate processing
		time.Sleep(2 * time.Second)

		// Update job status to completed
		err := h.db.UpdateFetchJobStatus(job.ID, "completed", 10, "")
		if err != nil {
			// Log the error but don't return it to the client
			// since this is an asynchronous operation
		}
	}()
}

// GetFetchStatus handles GET /contents/fetch/:jobId
func (h *ContentsHandler) GetFetchStatus(c *gin.Context) {
	// Get the job ID from the URL
	jobID := c.Param("jobId")

	// Get the job from the database
	job, err := h.db.GetFetchJob(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get fetch job: " + err.Error(),
		})
		return
	}

	// Check if the job exists
	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Fetch job not found",
		})
		return
	}

	// Return the job status
	c.JSON(http.StatusOK, job)
}

// SearchContents handles GET /contents/search
func (h *ContentsHandler) SearchContents(c *gin.Context) {
	// Parse query parameters
	keywords := c.Query("keywords")
	sourceID := c.Query("sourceId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Validate keywords
	if keywords == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Keywords parameter is required",
		})
		return
	}

	// Search contents
	contents, err := h.db.SearchContents(keywords, sourceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search contents: " + err.Error(),
		})
		return
	}

	// Return the search results
	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"count":    len(contents),
	})
}
