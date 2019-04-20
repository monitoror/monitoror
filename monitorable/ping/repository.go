package ping

import (
	"context"

	"github.com/monitoror/monitoror/monitorable/ping/model"
)

// Repository represent the ping's repository contract
type (
	Repository interface {
		Ping(ctx context.Context, hostname string) (*model.Ping, error)
	}
)
