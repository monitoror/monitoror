//+build faker

package models

import (
	"github.com/monitoror/monitoror/models/tiles"
)

type (
	FakerParamsProvider interface {
		GetStatus() tiles.TileStatus
		GetMessage() string
	}
)
