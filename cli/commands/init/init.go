package init

import (
	"io/ioutil"
	"strings"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/internal/pkg/path"

	rice "github.com/GeertJohan/go.rice"
	"github.com/spf13/cobra"
)

func NewInitCommand(monitororCli *cli.MonitororCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init monitoror with default config.json and .env",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(monitororCli, path.MonitororBaseDir)
		},
	}
	return cmd
}

func runInit(_ *cli.MonitororCli, basedir string) error {
	// Remove LocateAppend and LocateFS because we don't use it and it cause tests issue
	riceConfig := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded},
	}
	// Never use constant or variable according to docs : https://github.com/GeertJohan/go.rice#calling-findbox-and-mustfindbox
	defaultFiles, err := riceConfig.FindBox("default-files")
	if err != nil {
		panic("static default config files not found. Build them with `make package-defaultconfig` first.")
	}

	// Create defautl config.json
	_ = ioutil.WriteFile(path.ToAbsolute(basedir, "config.json"), defaultFiles.MustBytes("config-example.json"), 0644)

	// Create default .env
	dotEnv := defaultFiles.MustString(".env.example")
	dotEnv = strings.Replace(dotEnv, "config-example.json", "config.json", 1)
	_ = ioutil.WriteFile(path.ToAbsolute(basedir, ".env"), []byte(dotEnv), 0644)

	return nil
}
