package handlers

import (
	"net/http"
	"strconv"

	"github.com/flyer103/riffle/pkg/serving/storage"
	"github.com/gin-gonic/gin"
)

// RecommendationsHandler handles API requests for recommendations
type RecommendationsHandler struct {
	db *storage.SQLiteDB
}

// NewRecommendationsHandler creates a new RecommendationsHandler
func NewRecommendationsHandler(db *storage.SQLiteDB) *RecommendationsHandler {
	return &RecommendationsHandler{
		db: db,
	}
}

// GetRecommendations handles GET /recommendations
func (h *RecommendationsHandler) GetRecommendations(c *gin.Context) {
	// Parse query parameters
	userID := c.Query("userId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Parse source IDs if provided
	var sourceIDs []string
	if sourceIDsStr := c.Query("sourceIds"); sourceIDsStr != "" {
		// Split by comma
		for _, id := range c.QueryArray("sourceIds") {
			if id != "" {
				sourceIDs = append(sourceIDs, id)
			}
		}
	}

	// Get recommendations from the database
	input := storage.GetRecommendationsInput{
		UserID:    userID,
		SourceIDs: sourceIDs,
		Limit:     limit,
	}
	recommendations, err := h.db.GetRecommendations(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get recommendations: " + err.Error(),
		})
		return
	}

	// Return the recommendations
	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// SubmitFeedback handles POST /recommendations/feedback
func (h *RecommendationsHandler) SubmitFeedback(c *gin.Context) {
	// Parse the request body
	var input storage.CreateRecommendationFeedbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if input.ContentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "contentId is required",
		})
		return
	}
	if input.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userId is required",
		})
		return
	}
	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "rating must be between 1 and 5",
		})
		return
	}

	// Create the feedback
	feedback, err := h.db.CreateRecommendationFeedback(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to submit feedback: " + err.Error(),
		})
		return
	}

	// Return the created feedback
	c.JSON(http.StatusCreated, feedback)
}

// GetUserFeedback handles GET /recommendations/feedback/:userId
func (h *RecommendationsHandler) GetUserFeedback(c *gin.Context) {
	// Get the user ID from the URL
	userID := c.Param("userId")

	// Get the user's feedback from the database
	feedback, err := h.db.GetUserFeedback(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user feedback: " + err.Error(),
		})
		return
	}

	// Return the feedback
	c.JSON(http.StatusOK, gin.H{
		"feedback": feedback,
		"count":    len(feedback),
	})
}
