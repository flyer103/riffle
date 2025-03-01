package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// SystemHandler handles API requests for system-related endpoints
type SystemHandler struct {
	startTime time.Time
	version   string
}

// SystemInfo represents information about the system
type SystemInfo struct {
	Version      string    `json:"version"`
	GoVersion    string    `json:"goVersion"`
	StartTime    time.Time `json:"startTime"`
	Uptime       string    `json:"uptime"`
	NumGoroutine int       `json:"numGoroutine"`
}

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(version string) *SystemHandler {
	return &SystemHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// HealthCheck handles GET /health
func (h *SystemHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// GetSystemInfo handles GET /system/info
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	info := SystemInfo{
		Version:      h.version,
		GoVersion:    runtime.Version(),
		StartTime:    h.startTime,
		Uptime:       time.Since(h.startTime).String(),
		NumGoroutine: runtime.NumGoroutine(),
	}

	c.JSON(http.StatusOK, info)
}

// GetMetrics handles GET /metrics
// This is a placeholder for Prometheus metrics
// In a real implementation, this would be handled by the Prometheus client library
func (h *SystemHandler) GetMetrics(c *gin.Context) {
	c.String(http.StatusOK, "# This is a placeholder for Prometheus metrics\n")
}
