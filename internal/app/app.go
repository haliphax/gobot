// Package app
package app

import (
	"log"
	"os"
	"os/signal"

	"github.com/haliphax/gobot/internal/harness"
	"github.com/haliphax/gobot/internal/platform"
)

var (
	stop   = make(chan os.Signal, 1)
	client = &harness.OpenRouterClient{}
)

func Main() {
	signal.Notify(stop, os.Interrupt)
	shutdownBot := make(chan bool, 1)
	log.Println("Ctrl+C to exit")
	go platform.Discord(client, shutdownBot)
	<-stop
	shutdownBot <- true
	log.Println("Gracefully shutting down.")
}
