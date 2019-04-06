//+build faker

package model

import "github.com/monitoror/monitoror/models/tiles"

type (
	PingParams struct {
		Hostname string           `json:"hostname" query:"hostname"`
		Status   tiles.TileStatus `json:"status" query:"status"`
		Message  string           `json:"message" query:"message"`
	}
)

func (p *PingParams) Validate() bool {
	return p.Hostname != ""
}
