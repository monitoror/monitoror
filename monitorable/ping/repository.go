package ping

import (
	"github.com/monitoror/monitoror/monitorable/ping/models"
)

// Repository represent the ping's repository contract
type (
	Repository interface {
		ExecutePing(hostname string) (*models.Ping, error)
	}
)
