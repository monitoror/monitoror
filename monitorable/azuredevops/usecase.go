package azuredevops

import (
	"github.com/monitoror/monitoror/models"
	azureModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
)

const (
	AzureDevOpsBuildTileType   models.TileType = "AZUREDEVOPS-BUILD"
	AzureDevOpsReleaseTileType models.TileType = "AZUREDEVOPS-RELEASE"
)

type (
	Usecase interface {
		Build(params *azureModels.BuildParams) (*models.Tile, error)
		Release(params *azureModels.ReleaseParams) (*models.Tile, error)
	}
)
