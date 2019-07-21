package travisci

import (
	"github.com/monitoror/monitoror/monitorable/travisci/models"
)

// Repository represent the travisci's repository contract
type (
	Repository interface {
		GetLastBuildStatus(group, repository, branch string) (*models.Build, error)
	}
)
