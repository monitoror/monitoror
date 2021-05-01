//+build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	PingParams struct {
		params.Default

		Query string `json:"query" query:"query" validate:"required"`
	}
)
