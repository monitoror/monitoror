//+build faker

package models

import coreModels "github.com/monitoror/monitoror/models"

type (
	FakerParamsProvider interface {
		GetStatus() coreModels.TileStatus
		GetMessage() string
		GetValueValues() []string
		GetValueUnit() coreModels.TileValuesUnit
	}
)
