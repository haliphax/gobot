// Package config provides configuration functionality
package config

import (
	"flag"
	"io"
	"log"

	"github.com/BurntSushi/toml"

	"github.com/haliphax/gobot/internal"
)

type BaseConfig struct {
	ModelProviderType   string
	MessagePlatformType string
}

type AgentConfig struct {
	Model string
}

type Config struct {
	Base  BaseConfig
	Agent AgentConfig
}

var (
	ConfigFilename = flag.String("config", "gobot.toml", "Configuration file path")
	Configuration  string
)

func Load() *Config {
	fs := internal.GetFs()

	file, err := fs.Open(*ConfigFilename)
	if err != nil {
		panic(err.Error())
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	Configuration = string(bytes)
	c := &Config{}
	_, err = toml.Decode(Configuration, &c)
	if err != nil {
		log.Panicf("☠️ Unable to decode TOML: %v", err.Error())
	}

	return c
}
