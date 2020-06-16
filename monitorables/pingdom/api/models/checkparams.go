//+build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	CheckParams struct {
		params.Default

		ID *int `json:"id" query:"id" validate:"required"`
	}
)
