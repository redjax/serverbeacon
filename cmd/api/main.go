package main

import (
	"fmt"
	"log"
	"os"

	"github.com/redjax/serverbeacon/api"
	"github.com/redjax/serverbeacon/internal/auth"
	"github.com/redjax/serverbeacon/internal/config"
	"github.com/redjax/serverbeacon/internal/db"
	"github.com/spf13/pflag"
)

var configFile string

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// Entrypoint/init function
func run() error {
	// Accept -c/--config-file input
	pflag.StringVarP(&configFile, "config-file", "c", "", "config file path (default: ~/.local/share/serverbeacon, $env:LOCALAPPDATA\\serverbeacon)")

	// Parse CLI args
	pflag.Parse()

	if err := config.Init(nil, configFile); err != nil {
		log.Fatal(err)
	}

	cfg := config.GetConfig()

	if cfg.Debug {
		log.Printf("DB path: %s", cfg.DB.Path)
	}

	// Database initialization
	gdb, err := db.Open(cfg.DB.Path)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(gdb); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	// Do migrations
	if err := auth.Migrate(gdb); err != nil {
		log.Fatal(err)
	}

	// Create auth store for users & tokens
	store := auth.NewStore(gdb)

	// Bootstrap inintial users & secrets
	bootstrapSecrets, err := auth.EnsureBootstrap(store)
	if err != nil {
		log.Fatal(err)
	}

	printBootstrapSecrets(bootstrapSecrets)

	// Initialize server
	server := api.NewHttpServer(fmt.Sprintf("%d", cfg.APISettings.Port))

	// Start server
	if err := server.ListenMulti("0.0.0.0"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func printBootstrapSecrets(secrets map[string]string) {
	if len(secrets) == 0 {
		return
	}

	for k, v := range secrets {
		fmt.Fprintf(os.Stderr, "%s: %s\n", k, v)
	}
	fmt.Fprintln(os.Stderr)
}
