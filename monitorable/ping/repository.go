package ping

import (
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/model"
)

// Repository represent the ping's repository contract
type (
	Repository interface {
		Ping(host string) (*model.Ping, error)
	}
)
