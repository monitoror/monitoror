//+build faker

package models

import "github.com/monitoror/monitoror/models"

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`

		Status models.TileStatus `json:"status" query:"status"`
		Values []float64         `json:"value" query:"value"`
	}
)

func (p *PingParams) IsValid() bool {
	return p.Hostname != ""
}
