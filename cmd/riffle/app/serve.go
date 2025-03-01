package app

import (
	"github.com/flyer103/riffle/pkg/serving"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// NewServeCommand creates a new serve command
func NewServeCommand() *cobra.Command {
	opts := serving.NewServerOptions()

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server for RSS management",
		Long:  "Start an HTTP server that provides REST API for RSS feed management and content recommendations",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Complete and validate options
			if err := opts.Complete(); err != nil {
				return err
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			// Run the server
			klog.InfoS("Starting server", "port", opts.Port)
			server, err := serving.NewServer(opts)
			if err != nil {
				return err
			}

			return server.Run()
		},
	}

	// Add flags
	opts.AddFlags(cmd.Flags())

	return cmd
}
