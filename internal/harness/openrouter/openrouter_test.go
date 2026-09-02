package openrouter

import (
	"io"
	"testing"

	"github.com/spf13/afero"

	"github.com/haliphax/gobot/internal"
)

// LoadConfiguration should load the OpenRouter configuration block
func TestLoadConfigurationParsesConfiguration(t *testing.T) {
	fs := afero.NewMemMapFs()
	internal.SetFs(fs)
	file, err := fs.Create("test.toml")
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

	conf := LoadConfiguration("test.toml")

	if conf.Token != "test-token" {
		t.Errorf("test-token != %v", conf.Token)
	}
}
