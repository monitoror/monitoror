//+build faker

package models

import "github.com/monitoror/monitoror/models/tiles"

type (
	PortParams struct {
		Hostname string `json:"hostname" query:"hostname"`
		Port     int    `json:"port" query:"port"`

		Status tiles.TileStatus `json:"status" query:"status"`
	}
)

func (p *PortParams) IsValid() bool {
	return p.Hostname != "" && p.Port != 0
}
