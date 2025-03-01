package serving

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flyer103/riffle/pkg/serving/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

// Server represents the HTTP server
type Server struct {
	router        *gin.Engine
	db            *storage.SQLiteDB
	options       *ServerOptions
	metricsRouter *gin.Engine
	httpServer    *http.Server
	metricsServer *http.Server
}

// NewServer creates a new server instance
func NewServer(options *ServerOptions) (*Server, error) {
	// Set Gin mode based on log level
	if options.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize the router
	router := gin.New()

	// Initialize the database
	db, err := storage.NewSQLiteDB(options.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create the server
	server := &Server{
		router:  router,
		db:      db,
		options: options,
	}

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Add rate limiting if enabled
	if options.RateLimit > 0 {
		// Use a simple rate limiter middleware
		router.Use(func(c *gin.Context) {
			// In a real implementation, this would be a proper rate limiter
			c.Next()
		})
	}

	// Add CORS if enabled
	if options.EnableCORS {
		corsConfig := cors.DefaultConfig()
		if len(options.CORSOrigins) > 0 {
			corsConfig.AllowOrigins = options.CORSOrigins
		} else {
			corsConfig.AllowAllOrigins = true
		}
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
		corsConfig.AllowCredentials = true
		router.Use(cors.New(corsConfig))
	}

	// Setup routes
	server.setupRoutes()

	// Initialize metrics router if metrics are enabled
	if options.MetricsPort > 0 {
		server.metricsRouter = gin.New()
		server.metricsRouter.Use(gin.Recovery())
	}

	return server, nil
}

// Run starts the server
func (s *Server) Run() error {
	// Create the HTTP server
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.options.Port),
		Handler:      s.router,
		ReadTimeout:  s.options.ReadTimeout,
		WriteTimeout: s.options.WriteTimeout,
	}

	// Start the metrics server if enabled
	if s.options.MetricsPort > 0 {
		go func() {
			if err := s.startMetricsServer(); err != nil && err != http.ErrServerClosed {
				klog.Errorf("Metrics server error: %v", err)
			}
		}()
	}

	// Start the main server
	klog.Infof("Starting server on port %d", s.options.Port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// Shutdown the main server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			klog.Errorf("Error shutting down HTTP server: %v", err)
		}
	}

	// Shutdown the metrics server
	if s.metricsServer != nil {
		if err := s.metricsServer.Shutdown(ctx); err != nil {
			klog.Errorf("Error shutting down metrics server: %v", err)
		}
	}

	// Close the database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			klog.Errorf("Error closing database: %v", err)
		}
	}

	klog.Info("Server shutdown complete")
	return nil
}

// startMetricsServer starts the metrics server
func (s *Server) startMetricsServer() error {
	// Add the metrics handler
	s.metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Create the metrics server
	s.metricsServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.options.MetricsPort),
		Handler: s.metricsRouter,
	}

	// Start the metrics server
	klog.Infof("Starting metrics server on port %d", s.options.MetricsPort)
	return s.metricsServer.ListenAndServe()
}
