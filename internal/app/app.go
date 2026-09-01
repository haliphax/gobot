// Package app provides the application entrypoint
package app

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/haliphax/gobot/internal"
	"github.com/haliphax/gobot/internal/harness"
	"github.com/haliphax/gobot/internal/harness/openrouter"
	"github.com/haliphax/gobot/internal/platform"
	"github.com/haliphax/gobot/internal/platform/discord"
)

func Main() {
	// hard-coded config (for now)
	conf := internal.Config{
		ModelProviderType:   "openrouter",
		MessagePlatformType: "discord",
	}

	var (
		modelProviderClient harness.ModelProviderClient
		messagePlatform     platform.MessagePlatform
	)

	flag.Parse()

	// single provider type (for now)
	switch conf.ModelProviderType {
	default:
		modelProviderClient = harness.ModelProviderClient(openrouter.New())
	}

	// single platform type (for now)
	switch conf.MessagePlatformType {
	default:
		messagePlatform = platform.MessagePlatform(discord.New(modelProviderClient))
	}

	// coordination channels
	stop := make(chan os.Signal, 1)
	shutdownPlatform := make(chan bool, 1)

	// start message platform
	signal.Notify(stop, os.Interrupt)
	log.Println("🚪 Ctrl+C to exit")
	go messagePlatform.Start(shutdownPlatform)

	// wait for interrupt
	<-stop
	log.Println("🪦 Gracefully shutting down.")

	// shutdown platform and wait for termination
	shutdownPlatform <- true
	<-shutdownPlatform
}
