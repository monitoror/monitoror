//+build faker

package models

import coreModels "github.com/monitoror/monitoror/models"

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`

		Status      coreModels.TileStatus `json:"status" query:"status"`
		ValueValues []string              `json:"valueValues" query:"valueValues"`
	}
)

func (p *PingParams) IsValid() bool {
	return p.Hostname != ""
}
