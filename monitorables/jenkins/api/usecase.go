package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	JenkinsBuildTileType       coreModels.TileType = "JENKINS-BUILD"
	JenkinsMultiBranchTileType coreModels.TileType = "JENKINS-MULTIBRANCH"
)

type (
	Usecase interface {
		Build(params *models.BuildParams) (*coreModels.Tile, error)
		MultiBranch(params interface{}) ([]builder.Result, error)
	}
)
