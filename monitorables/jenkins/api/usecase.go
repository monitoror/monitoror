//go:generate mockery -name Usecase

package api

import (
	models2 "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"
)

const (
	JenkinsBuildTileType       coreModels.TileType = "JENKINS-BUILD"
	JenkinsMultiBranchTileType coreModels.TileType = "JENKINS-MULTIBRANCH"
)

type (
	Usecase interface {
		Build(params *models.BuildParams) (*coreModels.Tile, error)
		MultiBranch(params interface{}) ([]models2.DynamicTileResult, error)
	}
)
