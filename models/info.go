package models

import (
	"github.com/jsdidierlaurent/monitowall/configs"
)

type (
	//InfoResponse response for info route
	InfoResponse struct {
		BuildInfo configs.BuildInfo `json:"build-info"`
		Config    configs.Config    `json:"configuration"`
	}
)
