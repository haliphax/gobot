package config

import (
	"io"
	"testing"

	"github.com/spf13/afero"
)

// config.Load should parse base configuration from file
func TestConfigLoadParsesBaseConfiguration(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("test.toml")
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

	conf := Load(fs, "test.toml")

	if conf.Base.MessagePlatformType != "test-platform" {
		t.Errorf("%v != %v", "test-platform", conf.Base.MessagePlatformType)
	}

	if conf.Base.ModelProviderType != "test-provider" {
		t.Errorf("%v != %v", "test-provider", conf.Base.ModelProviderType)
	}
}

// config.Load should panic on malformed configuration file
func TestConfigLoadPanicsOnMalformedConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("test.toml")
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

	_ = Load(fs, "test.toml")
}
