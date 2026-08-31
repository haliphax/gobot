// Package app provides the application entrypoint
package app

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/haliphax/gobot/internal/harness"
	"github.com/haliphax/gobot/internal/harness/openrouter"
	"github.com/haliphax/gobot/internal/platform"
	"github.com/haliphax/gobot/internal/platform/discord"
)

func Main() {
	flag.Parse()

	stop := make(chan os.Signal, 1)
	shutdownBot := make(chan bool, 1)

	signal.Notify(stop, os.Interrupt)
	log.Println("Ctrl+C to exit")
	client := harness.ModelProviderClient(openrouter.New())
	plat := platform.MessagePlatform(discord.New(client))
	go platform.MessagePlatform(plat).Start(shutdownBot)

	<-stop
	log.Println("Gracefully shutting down.")

	// shutdown bot and wait for termination
	shutdownBot <- true
	<-shutdownBot
}
