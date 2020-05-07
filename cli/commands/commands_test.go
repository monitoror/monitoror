package commands

import (
	"testing"

	"github.com/monitoror/monitoror/cli"

	"github.com/spf13/cobra"
)

func TestAddCommands(t *testing.T) {
	command := &cobra.Command{Use: "test"}
	cli := &cli.MonitororCli{RootCmd: command}

	AddCommands(cli)
}
