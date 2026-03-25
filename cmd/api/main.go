package main

import (
	"fmt"
	"log"

	"github.com/redjax/serverbeacon/api"
	"github.com/redjax/serverbeacon/internal/config"
	"github.com/spf13/pflag"
)

var configFile string

func main() {
	// Accept -c/--config-file input
	pflag.StringVarP(&configFile, "config-file", "c", "", "config file path (default: ~/.local/share/serverbeacon, $env:LOCALAPPDATA\\serverbeacon)")

	// Parse CLI args
	pflag.Parse()

	if err := config.Init(nil, configFile); err != nil {
		log.Fatal(err)
	}

	cfg := config.GetConfig()

	server := api.NewHttpServer(fmt.Sprintf("%d", cfg.APiSettings.Port))

	if err := server.ListenMulti("0.0.0.0"); err != nil {
		log.Fatal(err)
	}
}
