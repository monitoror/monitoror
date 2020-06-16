//+build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	TransactionCheckParams struct {
		params.Default

		ID *int `json:"id" query:"id" validate:"required"`
	}
)
