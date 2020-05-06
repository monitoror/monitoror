package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	HTTPProxyParams struct {
		params.Default

		URL string `json:"url" query:"url" validate:"required,url,http"`
	}
)
