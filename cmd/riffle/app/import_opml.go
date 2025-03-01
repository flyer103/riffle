package app

import (
	"fmt"

	"github.com/flyer103/riffle/pkg/riffle"
	"github.com/flyer103/riffle/pkg/serving/storage"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// NewImportOPMLCommand creates a new import-opml command
func NewImportOPMLCommand() *cobra.Command {
	var (
		opmlFile string
		dbPath   string
	)

	cmd := &cobra.Command{
		Use:   "import-opml",
		Short: "Import RSS sources from an OPML file into the database",
		Long:  "Parse an OPML file and import the RSS sources into the SQLite database for use by the serve command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return importOPML(opmlFile, dbPath)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&opmlFile, "opml", "o", "", "Path to OPML file (required)")
	cmd.Flags().StringVar(&dbPath, "db-path", "./riffle.db", "Path to the SQLite database file")

	// Mark required flags
	cmd.MarkFlagRequired("opml")

	return cmd
}

// importOPML imports RSS sources from an OPML file into the database
func importOPML(opmlFile, dbPath string) error {
	// Parse the OPML file
	feeds, err := riffle.ParseOPML(opmlFile)
	if err != nil {
		return fmt.Errorf("failed to parse OPML file: %w", err)
	}

	// Connect to the database
	db, err := storage.NewSQLiteDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Prepare batch input
	var sourcesInput storage.BatchCreateSourcesInput
	for _, feed := range feeds {
		sourcesInput.Sources = append(sourcesInput.Sources, storage.CreateSourceInput{
			Name:        feed.Title,
			URL:         feed.URL,
			Description: fmt.Sprintf("Imported from OPML file: %s", opmlFile),
		})
	}

	// Import the sources
	result, err := db.BatchCreateSources(sourcesInput)
	if err != nil {
		return fmt.Errorf("failed to import RSS sources: %w", err)
	}

	// Print results
	klog.InfoS("OPML import completed",
		"totalFeeds", len(feeds),
		"importedSources", len(result.Sources),
		"errors", len(result.Errors))

	if len(result.Sources) > 0 {
		fmt.Println("\nSuccessfully imported RSS sources:")
		for _, source := range result.Sources {
			fmt.Printf("- %s (%s)\n", source.Name, source.URL)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors during import:")
		for _, err := range result.Errors {
			fmt.Printf("- Index %d: %s - %s\n", err.Index, err.ErrorType, err.Message)
		}
	}

	// Schedule a fetch job for the newly imported sources
	if len(result.Sources) > 0 {
		// Create a fetch job for all sources (null sourceID means all sources)
		days := 7 // Fetch articles from the last 7 days
		job, err := db.CreateFetchJob(nil, days)
		if err != nil {
			return fmt.Errorf("failed to create fetch job: %w", err)
		}

		fmt.Printf("\nCreated fetch job %s to retrieve content from all sources\n", job.ID)
		fmt.Println("You can check the job status using the API when running the serve command")
	}

	return nil
}
