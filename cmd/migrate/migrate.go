package migrate

import (
	"github.com/spf13/cobra"
)

var (
	StartCmd = &cobra.Command{
		Use:          "migrate",
		Short:        "start migrate",
		Example:      "start migrate",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	return nil
}
