package ping

import (
	"github.com/monitoror/monitoror/monitorable/ping/models"
)

type (
	Repository interface {
		ExecutePing(hostname string) (*models.Ping, error)
	}
)
