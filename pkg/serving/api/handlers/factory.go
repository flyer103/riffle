package handlers

import (
	"github.com/flyer103/riffle/pkg/serving/storage"
)

// Factory creates and initializes all API handlers
type Factory struct {
	Sources         *SourcesHandler
	Contents        *ContentsHandler
	Recommendations *RecommendationsHandler
	System          *SystemHandler
}

// NewFactory creates a new handler factory
func NewFactory(db *storage.SQLiteDB, version string) *Factory {
	return &Factory{
		Sources:         NewSourcesHandler(db),
		Contents:        NewContentsHandler(db),
		Recommendations: NewRecommendationsHandler(db),
		System:          NewSystemHandler(version),
	}
}
