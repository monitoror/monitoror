//+build faker

package models

import "github.com/monitoror/monitoror/models"

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`

		Status  models.TileStatus `json:"status" query:"status"`
		Message string            `json:"message" query:"message"`
	}
)

func (p *PingParams) IsValid() bool {
	return p.Hostname != ""
}
