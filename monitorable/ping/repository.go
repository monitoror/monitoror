package ping

import (
	"github.com/monitoror/monitoror/monitorable/ping/model"
)

// Repository represent the ping's repository contract
type (
	Repository interface {
		CheckPing(hostname string) (*model.Ping, error)
	}
)
