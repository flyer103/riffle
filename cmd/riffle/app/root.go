package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRiffleCommand creates the root command for riffle
func NewRiffleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "riffle",
		Short: "Riffle is an RSS feed analyzer and content recommender",
		Long: `Riffle analyzes RSS feeds from an OPML file and recommends articles
based on content quality and user interests. It helps you find the most
valuable content from your RSS subscriptions.`,
	}

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cmd := NewRiffleCommand()

	// Add subcommands
	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewServeCommand())
	cmd.AddCommand(NewImportOPMLCommand())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
