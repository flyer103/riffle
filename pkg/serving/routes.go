package serving

import (
	"net/http/pprof"

	"github.com/flyer103/riffle/pkg/serving/api/handlers"
	"github.com/gin-gonic/gin"
)

// setupRoutes sets up the API routes
func (s *Server) setupRoutes() {
	// Create the handler factory
	factory := handlers.NewFactory(s.db, "1.0.0") // TODO: Get version from build info

	// RSS Sources routes
	sources := s.router.Group("/sources")
	{
		sources.GET("", factory.Sources.ListSources)
		sources.GET("/:id", factory.Sources.GetSource)
		sources.POST("", factory.Sources.CreateSource)
		sources.PUT("/:id", factory.Sources.UpdateSource)
		sources.DELETE("/:id", factory.Sources.DeleteSource)
		sources.POST("/batch", factory.Sources.BatchCreateSources)
		sources.DELETE("/batch", factory.Sources.BatchDeleteSources)
	}

	// RSS Contents routes
	contents := s.router.Group("/contents")
	{
		contents.GET("", factory.Contents.ListContents)
		contents.GET("/:id", factory.Contents.GetContent)
		contents.PUT("/:id", factory.Contents.UpdateContent)
		contents.DELETE("/:id", factory.Contents.DeleteContent)
		contents.DELETE("/batch", factory.Contents.BatchDeleteContents)
		contents.POST("/fetch", factory.Contents.FetchContents)
		contents.GET("/fetch/:jobId", factory.Contents.GetFetchStatus)
		contents.GET("/search", factory.Contents.SearchContents)
	}

	// Recommendations routes
	recommendations := s.router.Group("/recommendations")
	{
		recommendations.GET("", factory.Recommendations.GetRecommendations)
		recommendations.POST("/feedback", factory.Recommendations.SubmitFeedback)
		recommendations.GET("/feedback/:userId", factory.Recommendations.GetUserFeedback)
	}

	// System routes
	s.router.GET("/health", factory.System.HealthCheck)
	s.router.GET("/system/info", factory.System.GetSystemInfo)

	// Add pprof routes if enabled
	if s.options.EnablePprof {
		pprofGroup := s.router.Group("/debug/pprof")
		{
			pprofGroup.GET("/", gin.WrapF(pprof.Index))
			pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
			pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
			pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
			pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
			pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
			pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
			pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
			pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
			pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
			pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
		}
	}
}
