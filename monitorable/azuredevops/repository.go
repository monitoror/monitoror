package azuredevops

import (
	"github.com/monitoror/monitoror/monitorable/azuredevops/models"

	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops/build"
	"github.com/jsdidierlaurent/azure-devops-go-api/azuredevops/release"
)

type (
	Connection interface {
		GetBuildConnection() (build.Client, error)
		GetReleaseConnection() (release.Client, error)
	}

	Repository interface {
		GetBuild(project string, definition int, branch *string) (*models.Build, error)
		GetRelease(project string, definition int) (*models.Release, error)
	}
)
