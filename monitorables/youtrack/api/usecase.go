package api

import (
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	Usecase interface {
		CountIssues(params *models.) (*coreModels.Tile, error)
	}
)
