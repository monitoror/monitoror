package ping

import (
	"context"

	"github.com/monitoror/monitoror/monitorable/ping/models"
)

// Repository represent the ping's repository contract
type (
	Repository interface {
		Ping(ctx context.Context, hostname string) (*models.Ping, error)
	}
)
