package startcommand

import (
	"fmt"

	"github.com/redjax/serverbeacon/api"
	"github.com/redjax/serverbeacon/internal/config"
	"github.com/redjax/serverbeacon/internal/validators"
	"github.com/spf13/cobra"
)

func NewRestApiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rest-api",
		Aliases: []string{"rest"},
		Short:   "Start REST API server",
		Long:    `Start serverbeacon REST API with configurable proto/host/port`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()

			// Override Host config if passed from CLI
			if h, err := cmd.Flags().GetString("host"); err == nil && h != "" {
				cfg.APISettings.Host = h
			}

			// Override/convert Port config if passed from CLI
			if p, err := cmd.Flags().GetInt16("port"); err == nil && p != 0 {
				cfg.APISettings.Port = int64(p) // int16 -> int64
			}

			if cfg.APISettings.Port < 1 || cfg.APISettings.Port > 65535 {
				return fmt.Errorf("invalid port: %d (must be 1-65535)", cfg.APISettings.Port)
			}

			if err := validators.ValidateTCPPort(cfg.APISettings.Port); err != nil {
				return err
			}

			server := api.NewHttpServer(fmt.Sprintf("%d", cfg.APISettings.Port))

			return server.ListenMulti(cfg.APISettings.Host)
		},
	}

	// Local flags to override config values
	cmd.Flags().StringP("host", "H", "", "API bind host (default: 0.0.0.0)")
	cmd.Flags().Int16P("port", "p", 0, "Port binding for REST API (default: 0, config default: 18080)")
	cmd.Flags().StringP("proto", "", "", "Protocol (default: http)")

	return cmd
}
