package server

import (
	"github.com/spf13/cobra"
)

var (
	StartCmd = &cobra.Command{
		Use:          "server",
		Short:        "start server",
		Example:      "start server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	return nil
}
