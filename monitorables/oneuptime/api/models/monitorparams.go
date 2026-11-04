//go:build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	MonitorParams struct {
		params.Default

		ID string `json:"id" query:"id" validate:"required"`
	}
)
