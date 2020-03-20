//+build faker

package models

import coreModels "github.com/monitoror/monitoror/models"

type (
	PortParams struct {
		Hostname string `json:"hostname" query:"hostname"`
		Port     int    `json:"port" query:"port"`

		Status coreModels.TileStatus `json:"status" query:"status"`
	}
)

func (p *PortParams) IsValid() bool {
	return p.Hostname != "" && p.Port != 0
}
