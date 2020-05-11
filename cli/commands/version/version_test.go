package version

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/cli/version"
)

func TestVersionCommand(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := &cli.MonitororCli{Output: output}
	version.Version = "1.0.0"
	version.GitCommit = "HASH"
	version.BuildTime = "Now"
	version.BuildTags = "test"

	expected := ` Version:    1.0.0 (test)
 Git commit: HASH
 Built:      Now
`

	cmd := NewVersionCommand(monitororCli)
	assert.NoError(t, cmd.RunE(cmd, []string{}))
	assert.True(t, strings.HasPrefix(output.String(), expected))
}
