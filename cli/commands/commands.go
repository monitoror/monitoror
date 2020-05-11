package commands

import (
	"github.com/monitoror/monitoror/cli"
	initCmd "github.com/monitoror/monitoror/cli/commands/init"
	"github.com/monitoror/monitoror/cli/commands/version"
)

func AddCommands(cli *cli.MonitororCli) {
	cli.RootCmd.AddCommand(
		// INIT
		initCmd.NewInitCommand(cli),
		// VERSION
		version.NewVersionCommand(cli),
	)
}
