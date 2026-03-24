package cmd

import (
	"fmt"
	"log"

	"github.com/redjax/serverbeacon/internal/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
	debug      bool
	cfg        *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "serverbeacon",
	Short: "Create beacons for various protocols",
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no command is provided
		cmd.Help()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Parse CLI flags

	// Config file
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "config file path (default: ~/.local/share/serverbeacon/config.yml)")
	// Debug flag
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug logging")

	// Load configuration
	rootCmd.PersistentPreRunE = loadConfig

	// Handle persistent flgs
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Handle --debug flag
		if d, _ := cmd.Flags().GetBool("debug"); d {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Println("DEBUG logging enabled")
		}
	}

	// Register subcommands
	// rootCmd.AddCommand(cmdPkg.ExampleCmd)
}

func loadConfig(cmd *cobra.Command, args []string) error {
	c, err := config.LoadConfig(cmd.PersistentFlags(), configFile)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	cfg = c

	// Debug print config
	// fmt.Printf("Config: %+v\n", cfg)

	return nil
}
