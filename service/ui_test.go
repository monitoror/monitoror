package service

import (
	"testing"

	cliMocks "github.com/monitoror/monitoror/cli/mocks"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/service/store"
)

func TestInitUI_Dev(t *testing.T) {
	cliMock := new(cliMocks.CLI)
	cliMock.On("PrintDevMode")
	InitUI(&Server{
		Echo: nil,
		store: &store.Store{
			CoreConfig: &config.CoreConfig{Env: "develop"},
			Cli:        cliMock,
		},
	})
	cliMock.AssertNumberOfCalls(t, "PrintDevMode", 1)
}
