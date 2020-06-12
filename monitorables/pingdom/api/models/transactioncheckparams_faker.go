//+build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	TransactionCheckParams struct {
		params.Default

		ID *int `json:"id" query:"id" validate:"required"`

		Status coreModels.TileStatus `json:"status" query:"status"`
	}
)
