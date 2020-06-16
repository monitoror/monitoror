package cli

import (
	"io"
	"os"

	"github.com/monitoror/monitoror/store"

	"github.com/spf13/cobra"
)

type MonitororCli struct {
	RootCmd *cobra.Command
	Store   *store.Store
	Output  io.Writer
}

func NewMonitororCli(store *store.Store) *MonitororCli {
	return &MonitororCli{
		Store:  store,
		Output: os.Stdout,
	}
}
