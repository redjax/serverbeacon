package startcommand

import "github.com/spf13/cobra"

func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start serverbeacon services",
	}

	// Add subcommands
	cmd.AddCommand(NewRestApiCommand())

	return cmd
}
