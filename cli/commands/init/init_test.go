package init

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/GeertJohan/go.rice/embedded"
	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/cli"
)

func TestVersionCommand(t *testing.T) {
	delete(embedded.EmbeddedBoxes, "default-files")

	output := &bytes.Buffer{}
	monitororCli := &cli.MonitororCli{Output: output}

	cmd := NewInitCommand(monitororCli)
	assert.Panics(t, func() { _ = cmd.RunE(cmd, []string{}) })
}

func TestRunInit(t *testing.T) {
	delete(embedded.EmbeddedBoxes, "default-files")
	embedded.RegisterEmbeddedBox("default-files", &embedded.EmbeddedBox{
		Name: "default-files",
		Files: map[string]*embedded.EmbeddedFile{
			".env.example":        {Filename: ".env.example", FileModTime: time.Now(), Content: "test"},
			"config-example.json": {Filename: "config-example.json", FileModTime: time.Now(), Content: "test"},
		},
	})
	defer delete(embedded.EmbeddedBoxes, "default-files")

	output := &bytes.Buffer{}
	monitororCli := &cli.MonitororCli{Output: output}

	tmpDir, err := ioutil.TempDir("", "initCommand")
	if assert.NoError(t, err) {
		defer os.RemoveAll(tmpDir)

		assert.NoError(t, runInit(monitororCli, tmpDir))
		files, err := ioutil.ReadDir(tmpDir)
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Equal(t, ".env", files[0].Name())
		assert.Equal(t, "config.json", files[1].Name())
	}
}
