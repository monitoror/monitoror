//go:build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	MonitorParams struct {
		params.Default

		ID string `json:"id" query:"id" validate:"required"`

		Status coreModels.TileStatus `json:"status" query:"status"`
	}
)
