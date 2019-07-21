package jenkins

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"
)

const (
	JenkinsBuildTileType TileType = "JENKINS-BUILD"
)

// Usecase represent the jenkins's usecases
type (
	Usecase interface {
		Build(params *models.BuildParams) (*BuildTile, error)
	}
)
