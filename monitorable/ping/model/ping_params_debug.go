//+build debug

package model

import "github.com/jsdidierlaurent/monitowall/models/tiles"

type (
	PingParamsDebug struct {
		Hostname string           `json:"hostname" query:"hostname"`
		Status   tiles.TileStatus `json:"status" query:"status"`
		Message  string           `json:"message" query:"message"`
	}
)

func (p *PingParamsDebug) Validate() bool {
	return p.Hostname != ""
}
