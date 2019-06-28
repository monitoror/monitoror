package jenkins

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"
)

const (
	JenkinsBuildTileType TileType = "JENKINS"
)

// Usecase represent the jenkins's usecases
type (
	Usecase interface {
		Build(params *models.JobParams) (*BuildTile, error)
	}
)
