package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/store"
)

func TestNewMonitororCli(t *testing.T) {
	s := &store.Store{}
	cli := NewMonitororCli(s)

	assert.NotNil(t, cli)
	assert.Equal(t, os.Stdout, cli.Output)
}
