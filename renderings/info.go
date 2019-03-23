package renderings

import (
	"github.com/jsdidierlaurent/monitowall/config"
)

type (
	//InfoResponse response for info route
	InfoResponse struct {
		BuildInfo config.BuildInfo `json:"build-info"`
		Config    config.Config    `json:"configuration"`
	}
)
