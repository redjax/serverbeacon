package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "serverbeacon",
	Short: "Create beacons for various protocols",
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no command is provided
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Register subcommands
	// rootCmd.AddCommand(cmdPkg.ExampleCmd)

	// Global persistent flgs
	// Empty default means 'use default with .local fallback'
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "config file path (default: ~/.local/share/serverbeacon/config.yml)")
}
