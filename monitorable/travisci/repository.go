package travisci

import (
	"context"

	"github.com/monitoror/monitoror/monitorable/travisci/model"
)

// Repository represent the travisci's repository contract
type (
	Repository interface {
		Build(ctx context.Context, group, repository, branch string) (*model.Build, error)
	}
)
