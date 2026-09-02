package discord

import (
	"io"
	"testing"

	"github.com/spf13/afero"

	"github.com/haliphax/gobot/internal"
	"github.com/haliphax/gobot/internal/config"
)

var testFn = new("test.toml")

// LoadConfiguration should load the Discord configuration block
func TestLoadConfigurationParsesConfiguration(t *testing.T) {
	fs := afero.NewMemMapFs()
	config.ConfigFilename = testFn
	internal.SetFs(fs)
	file, err := fs.Create(*testFn)
	if err != nil {
		panic(err.Error())
	}

	content := `
[Discord]
Token = "test-token"
`
	_, err = io.WriteString(file, content)
	if err != nil {
		panic(err.Error())
	}

	conf := LoadConfiguration()

	if conf.Token != "test-token" {
		t.Errorf("test-token != %v", conf.Token)
	}
}
