package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/cli/commands"
	"github.com/monitoror/monitoror/cli/printer"
	"github.com/monitoror/monitoror/cli/version"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/service"
	"github.com/monitoror/monitoror/store"

	"github.com/joho/godotenv"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMonitororRootCommand(monitororCli *cli.MonitororCli) {
	cmd := &cobra.Command{
		Use:   "monitoror",
		Short: "Unified monitoring wallboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Init Service
			server := service.Init(monitororCli.Store)

			if err := printer.PrintStartupLog(monitororCli); err != nil {
				return err
			}
			return server.Start()
		},
		Version:       fmt.Sprintf("%s, build %s", version.Version, version.GitCommit),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().BoolP("version", "v", false, "Print version information and quit")
	cmd.PersistentFlags().BoolP("debug", "d", false, "Start monitoror in debug mode")
	_ = viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))

	// Setup this command as root command in cli
	monitororCli.RootCmd = cmd

	// add other command
	commands.AddCommands(monitororCli)
}

func main() {
	// Setup logger
	log.SetPrefix("")
	log.SetHeader("[${level}]")
	log.SetLevel(log.INFO)

	// Load .env
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	_ = godotenv.Load(".env")
	_ = godotenv.Load(filepath.Join(dir, ".env"))

	// Setup Store
	store := &store.Store{
		CoreConfig: config.InitConfig(),
		Registry:   registry.NewRegistry(),
		CacheStore: cache.NewGoCacheStore(time.Minute*5, time.Second), // Default value, Always override
	}

	// Init CLI
	monitororCli := cli.NewMonitororCli(store)
	newMonitororRootCommand(monitororCli)

	// Start CLI
	if err := monitororCli.RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
