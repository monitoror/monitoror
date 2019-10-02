//+build faker

package models

import "github.com/monitoror/monitoror/models"

type (
	FakerParamsProvider interface {
		GetStatus() models.TileStatus
		GetMessage() string
	}
)
