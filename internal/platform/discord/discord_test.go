package discord

import (
	"io"
	"testing"

	"github.com/spf13/afero"
)

// LoadConfiguration should load the Discord configuration block
func TestLoadConfigurationParsesConfiguration(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("test.toml")
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

	conf := LoadConfiguration(fs, "test.toml")

	if conf.Token != "test-token" {
		t.Errorf("test-token != %v", conf.Token)
	}
}
