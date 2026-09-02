// Package app provides the application entrypoint
package app

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/afero"

	"github.com/haliphax/gobot/internal"
	"github.com/haliphax/gobot/internal/config"
	"github.com/haliphax/gobot/internal/harness"
	"github.com/haliphax/gobot/internal/harness/openrouter"
	"github.com/haliphax/gobot/internal/platform"
	"github.com/haliphax/gobot/internal/platform/discord"
)

var ConfigFilename = new("gobot.toml")

func Main() {
	internal.SetFs(afero.NewOsFs())

	// load base configuration
	conf := config.Load(*ConfigFilename)

	var (
		modelProviderClient harness.ModelProviderClient
		messagePlatform     platform.MessagePlatform
	)

	flag.Parse()

	// load model provider
	switch conf.Base.ModelProviderType {
	case "openrouter":
		modelProviderClient = harness.ModelProviderClient(openrouter.New(*ConfigFilename))
	default:
		log.Fatalf("☠️ Unsupported model provider: %v", conf.Base.ModelProviderType)
	}

	// set default model
	modelProviderClient.SetModel(conf.Agent.Model)

	// load message platform
	switch conf.Base.MessagePlatformType {
	case "discord":
		messagePlatform = platform.MessagePlatform(discord.New(*ConfigFilename, modelProviderClient))
	default:
		log.Fatalf("☠️ Unsupported message platform: %v", conf.Base.MessagePlatformType)
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
