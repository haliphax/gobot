// Package config provides configuration functionality
package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
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

func Load(fn string) *Config {
	bytes, err := os.ReadFile(fn)
	if err != nil {
		log.Panicf("☠️ Unable to load configuration file: %v", err.Error())
	}

	c := &Config{}
	_, err = toml.Decode(string(bytes), &c)
	if err != nil {
		log.Panicf("☠️ Unable to decode TOML: %v", err.Error())
	}

	return c
}
