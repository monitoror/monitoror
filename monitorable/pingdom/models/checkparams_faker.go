//+build faker

package models

import "github.com/monitoror/monitoror/models"

type (
	CheckParams struct {
		Id *int `json:"id" query:"id"`

		Status models.TileStatus `json:"status" query:"status"`
	}
)

func (p *CheckParams) IsValid() bool {
	return p.Id != nil
}
