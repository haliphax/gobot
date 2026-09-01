// Package config provides configuration functionality
package config

import (
	"io"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
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

func Load(fs afero.Fs, fn string) *Config {
	file, err := fs.Open(fn)
	if err != nil {
		panic(err.Error())
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	c := &Config{}
	_, err = toml.Decode(string(bytes), &c)
	if err != nil {
		log.Panicf("☠️ Unable to decode TOML: %v", err.Error())
	}

	return c
}
