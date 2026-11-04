//go:generate mockery --name Repository

package api

import (
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"
)

type (
	Repository interface {
		GetMonitors() ([]models.Monitor, error)
	}
)
