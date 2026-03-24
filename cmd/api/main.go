package main

import (
	"log"

	"github.com/redjax/serverbeacon/api"
)

func main() {
	server := api.NewHttpServer("18080")

	if err := server.ListenMulti("0.0.0.0"); err != nil {
		log.Fatal(err)
	}
}
