package azuredevops

import (
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops/models"
)

const (
	AzureDevOpsBuildTileType   TileType = "AZUREDEVOPS-BUILD"
	AzureDevOpsReleaseTileType TileType = "AZUREDEVOPS-RELEASE"
)

// Usecase represent the Azure devops's usecases
type (
	Usecase interface {
		Build(params *models.BuildParams) (*Tile, error)
		Release(params *models.ReleaseParams) (*Tile, error)
	}
)
