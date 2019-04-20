package port

import "context"

// Repository represent the port's repository contract
type (
	Repository interface {
		OpenSocket(ctx context.Context, hostname string, port int) error
	}
)
