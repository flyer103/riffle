package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	opmlFile      string
	interestsFile string
)

// NewRiffleCommand creates the root command for riffle
func NewRiffleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "riffle",
		Short: "Riffle is an RSS feed analyzer and content recommender",
		Long: `Riffle analyzes RSS feeds from an OPML file and recommends articles
based on content quality and user interests. It helps you find the most
valuable content from your RSS subscriptions.`,
		RunE: runRiffle,
	}

	// Add flags
	cmd.Flags().StringVarP(&opmlFile, "opml", "o", "", "Path to OPML file (required)")
	cmd.Flags().StringVarP(&interestsFile, "interests", "i", "", "Path to file containing interests (one per line)")

	// Mark required flags
	cmd.MarkFlagRequired("opml")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cmd := NewRiffleCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
