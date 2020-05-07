package init

import (
	"io/ioutil"
	"path"

	"github.com/monitoror/monitoror/cli"

	rice "github.com/GeertJohan/go.rice"
	"github.com/spf13/cobra"
)

const DefaultFilePath = "."

func NewInitCommand(monitororCli *cli.MonitororCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init monitoror with default config.json and .env",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(DefaultFilePath, monitororCli)
		},
	}
	return cmd
}

func runInit(defaultFilePath string, _ *cli.MonitororCli) error {
	// Remove LocateAppend and LocateFS because we don't use it and it cause tests issue
	riceConfig := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded},
	}
	// Never use constant or variable according to docs : https://github.com/GeertJohan/go.rice#calling-findbox-and-mustfindbox
	defaultFiles, err := riceConfig.FindBox("default-files")
	if err != nil {
		panic("static default config files not found. Build them with `make package-defaultconfig` first.")
	}

	_ = ioutil.WriteFile(path.Join(defaultFilePath, ".env"), defaultFiles.MustBytes(".env.example"), 0644)
	_ = ioutil.WriteFile(path.Join(defaultFilePath, "config.json"), defaultFiles.MustBytes("config-example.json"), 0644)

	return nil
}
