package travisci

import (
	"github.com/monitoror/monitoror/monitorable/travisci/models"
)

type (
	Repository interface {
		GetLastBuildStatus(owner, repository, branch string) (*models.Build, error)
	}
)
