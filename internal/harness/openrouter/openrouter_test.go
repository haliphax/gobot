package openrouter

import (
	"io"
	"testing"

	"github.com/spf13/afero"

	"github.com/haliphax/gobot/internal"
	"github.com/haliphax/gobot/internal/config"
)

var testFn = new("test.toml")

// LoadConfiguration should load the OpenRouter configuration block
func TestLoadConfigurationParsesConfiguration(t *testing.T) {
	// write test configuration file in memory-mapped file system
	fs := afero.NewMemMapFs()
	internal.SetFs(fs)
	config.ConfigFilename = testFn
	file, err := fs.Create(*testFn)
	if err != nil {
		panic(err.Error())
	}

	content := `
[OpenRouter]
Token = "test-token"
`
	_, err = io.WriteString(file, content)
	if err != nil {
		panic(err.Error())
	}

	// parse configuration
	config.Load()
	conf := LoadConfiguration()

	if conf.Token != "test-token" {
		t.Errorf("test-token != %v", conf.Token)
	}
}
