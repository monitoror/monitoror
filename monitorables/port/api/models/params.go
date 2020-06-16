//+build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	PortParams struct {
		params.Default

		Hostname string `json:"hostname" query:"hostname" validate:"required"`
		Port     int    `json:"port" query:"port" validate:"required,gt=0"`
	}
)
