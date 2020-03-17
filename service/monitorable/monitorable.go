package monitorable

import (
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/service/router"
)

type Monitorable interface {
	GetHelp() string
	GetVariants() []string
	Register(variant string, router router.MonitorableRouter, manager config.Manager) bool
}
