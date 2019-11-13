package travisci

import (
	"github.com/monitoror/monitoror/monitorable/travisci/models"
)

type (
	Repository interface {
		GetLastBuildStatus(group, repository, branch string) (*models.Build, error)
	}
)
