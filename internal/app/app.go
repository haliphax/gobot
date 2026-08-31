// Package app
package app

import (
	"log"
	"os"
	"os/signal"

	"github.com/haliphax/gobot/internal/bot"
	_ "github.com/haliphax/gobot/internal/harness"
)

var stop = make(chan os.Signal, 1)

func Main() {
	signal.Notify(stop, os.Interrupt)
	shutdownBot := make(chan bool, 1)
	log.Println("Ctrl+C to exit")
	go bot.Main(shutdownBot)
	<-stop
	shutdownBot <- true
	log.Println("Gracefully shutting down.")
}
