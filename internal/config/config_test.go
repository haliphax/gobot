package config

import (
	"io"
	"log"
	"testing"

	"github.com/spf13/afero"

	"github.com/haliphax/gobot/internal"
)

var testFn = new("test.toml")

// config.Load should parse base configuration from file
func TestConfigLoadParsesBaseConfiguration(t *testing.T) {
	// create test config file in memory-mapped file system
	fs := afero.NewMemMapFs()
	ConfigFilename = testFn
	internal.SetFs(fs)
	file, err := fs.Create(*testFn)
	if err != nil {
		panic(err.Error())
	}

	content := `
[Base]
ModelProviderType = "test-provider"
MessagePlatformType = "test-platform"

[Agent]
Model = "test-model"
`
	_, err = io.WriteString(file, content)
	if err != nil {
		panic(err.Error())
	}

	// parse configuration
	conf := Load()

	if conf.Base.MessagePlatformType != "test-platform" {
		t.Errorf("test-platform != %v", conf.Base.MessagePlatformType)
	}

	if conf.Base.ModelProviderType != "test-provider" {
		t.Errorf("test-provider != %v", conf.Base.ModelProviderType)
	}
}

// config.Load should panic on malformed configuration file
func TestConfigLoadPanicsOnMalformedConfig(t *testing.T) {
	// squash logs during test
	previous := log.Writer()
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(previous) })

	// create test config file in memory-mapped file system
	fs := afero.NewMemMapFs()
	internal.SetFs(fs)
	ConfigFilename = testFn
	file, err := fs.Create(*testFn)
	if err != nil {
		panic(err.Error())
	}

	content := "["
	_, err = io.WriteString(file, content)
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected failure")
		}
	}()

	// parse configuration
	_ = Load()
}
