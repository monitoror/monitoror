package travisci

import (
	"context"

	"github.com/monitoror/monitoror/monitorable/travisci/models"
)

// Repository represent the travisci's repository contract
type (
	Repository interface {
		GetBuildStatus(ctx context.Context, group, repository, branch string) (*models.Build, error)
	}
)
